package visitors

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dpb587/bosh-dross/distich/data"
)

type ReleaseCollectorRelease struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	URL     string `json:"url"`
	SHA1    string `json:"sha1"`
}

func (rcr ReleaseCollectorRelease) SchemaURL() string {
	return fmt.Sprintf(
		"https://dpb587.github.io/bosh-json-schema-release/v0/fake-host/fake-owner/fake-repo/%s/%s/jobs?sha1=%s&url=%s",
		rcr.Name,
		rcr.Version,
		rcr.SHA1,
		rcr.URL,
	)
}

type ReleaseCollector struct {
	releases map[string]ReleaseCollectorRelease
}

func NewReleaseCollector() ReleaseCollector {
	return ReleaseCollector{
		releases: map[string]ReleaseCollectorRelease{},
	}
}

var _ data.NodeVisitor = &ReleaseCollector{}

func (t *ReleaseCollector) EnterNode(node data.Node) error {
	if !strings.HasPrefix(node.GetPath(), "/releases") {
		return nil
	}

	// lazy
	marshalled, err := json.Marshal(node.Export())
	if err != nil {
		return err
	}

	var release ReleaseCollectorRelease

	err = json.Unmarshal(marshalled, release)
	if err != nil {
		return err
	}

	t.releases[release.Name] = release

	return nil
}

func (t *ReleaseCollector) LeaveNode(node data.Node) error {
	return nil
}

func (t *ReleaseCollector) GetReleases() map[string]ReleaseCollectorRelease {
	return t.releases
}
