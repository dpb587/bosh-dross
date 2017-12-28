package loaders

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dpb587/bosh-dross/distich/schema"
	yaml "gopkg.in/yaml.v2"
)

type LocalReleaseCache struct {
	DownloadMissing bool
}

type localReleaseMF struct {
	Name               string `yaml:"name"`
	Version            string `yaml:"version"`
	UncommittedChanges bool   `yaml:"uncommitted_changes"`
	CommitHash         string `yaml:"commit_hash"`
}

type localReleaseJobMF struct {
	Name       string                                 `yaml:"name"`
	Properties map[string]localReleaseJobMFProperties `yaml:"properties"`
}

type localReleaseJobMFProperties struct {
	Default     interface{} `yaml:"-"` // @todo yaml map[interface{}]interface{} -> json map[interface{}]interface{} fails?
	Description string      `yaml:"description"`
	// Example     string      `yaml:"example"` // @todo not always a string
	Type string `yaml:"type"`
}

var _ schema.Loader = LocalReleaseCache{}

var localReleaseCacheTarball = regexp.MustCompilePOSIX(`^jobs/.+\.tgz$`)

func (l LocalReleaseCache) IsSupported(uri string) bool {
	if !strings.HasPrefix(uri, "https://dpb587.github.io/bosh-json-schema-release/v0/") {
		return false
	}

	_, cache, err := l.parseURI(uri)
	if err != nil {
		return false
	}

	if cache == "" {
		return false
	} else if _, err := os.Stat(cache); os.IsNotExist(err) {
		return false
	}

	return true
}

func (l LocalReleaseCache) Load(uri string) ([]byte, error) {
	_, cache, err := l.parseURI(uri)
	if err != nil {
		return nil, err
	}

	fileReader, err := os.OpenFile(cache, os.O_RDONLY, 0000)
	if err != nil {
		return nil, err
	}

	defer fileReader.Close()

	gzipReader, err := gzip.NewReader(fileReader)
	if err != nil {
		return nil, err
	}

	tarReader := tar.NewReader(gzipReader)

	schemaNode := schema.Node{
		Definitions: map[string]*schema.Node{},
	}

	var releaseManifest localReleaseMF

	for {
		file, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		// @todo possibly including job tarballs not referenced by release.MF
		if localReleaseCacheTarball.Match([]byte(strings.TrimPrefix(file.Name, "./"))) {
			err = l.appendJobTarball(&schemaNode, tarReader)
			if err != nil {
				return nil, err
			}
		} else if strings.TrimPrefix(file.Name, "./") == "release.MF" {
			releaseBytes, err := ioutil.ReadAll(tarReader)
			if err != nil {
				return nil, err
			}

			err = yaml.Unmarshal(releaseBytes, &releaseManifest)
			if err != nil {
				return nil, err
			}
		}
	}

	for _, jobSchema := range schemaNode.Properties {
		jobSchema.Description = fmt.Sprintf("The %s job from the %s/%s (%s%s) release.", jobSchema.Title, releaseManifest.Name, releaseManifest.Version, releaseManifest.CommitHash, "?") // @todo uncommitted_changes
		jobSchema.Properties["release"] = &schema.Node{
			Enum: []string{releaseManifest.Name},
		}
	}

	marshalled, err := json.Marshal(schemaNode)
	if err != nil {
		return nil, err
	}

	return marshalled, nil
}

func (l LocalReleaseCache) parseURI(uri string) (string, string, error) {
	parsedURI, err := url.Parse(uri)
	if err != nil {
		return "", "", err
	}

	url := parsedURI.Query().Get("url")

	sha1 := parsedURI.Query().Get("sha1")
	cachePath := filepath.Join(os.Getenv("HOME"), ".bosh", "cache", sha1)

	return url, cachePath, nil
}

func (l LocalReleaseCache) appendJobTarball(schemaNode *schema.Node, jobReader io.Reader) error {
	gzipReader, err := gzip.NewReader(jobReader)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(gzipReader)

	for {
		file, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		} else if strings.TrimPrefix(file.Name, "./") != "job.MF" {
			continue
		}

		jobBytes, err := ioutil.ReadAll(tarReader)
		if err != nil {
			return err
		}

		var jobManifest localReleaseJobMF

		err = yaml.Unmarshal(jobBytes, &jobManifest)
		if err != nil {
			return err
		}

		// @todo additional job properties are not validated; consumes/provides
		jobSchemaNode := &schema.Node{
			Title:      jobManifest.Name,
			Properties: map[string]*schema.Node{},
		}

		for propertyName, propertySpec := range jobManifest.Properties {
			err = l.appendJobSchemaProperty(jobSchemaNode, propertyName, propertySpec)
			if err != nil {
				return err
			}
		}

		schemaNode.Definitions[fmt.Sprintf("%s_properties", jobManifest.Name)] = jobSchemaNode
	}

	return nil
}

func (l LocalReleaseCache) appendJobSchemaProperty(node *schema.Node, name string, spec localReleaseJobMFProperties) error {
	nameSplit := strings.SplitN(name, ".", 2)

	if len(nameSplit) == 2 {
		if _, found := node.Properties[nameSplit[0]]; !found {
			node.Properties[nameSplit[0]] = &schema.Node{
				Type:       "object",
				Properties: map[string]*schema.Node{},
			}
		} else if node.Properties[nameSplit[0]].Type != "object" {
			return errors.New("property is trying to overwrite existing multi-level property")
		}

		return l.appendJobSchemaProperty(node.Properties[nameSplit[0]], nameSplit[1], spec)
	}

	node.Properties[nameSplit[0]] = &schema.Node{
		Title:       strings.Title(strings.Replace(name, "_", " ", -1)),
		Default:     spec.Default,
		Description: spec.Description,
	}

	return nil
}
