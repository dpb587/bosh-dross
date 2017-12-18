package manifest

import (
	"fmt"
	"io/ioutil"

	boshtpl "github.com/cloudfoundry/bosh-cli/director/template"
	"github.com/cppforlife/go-patch/patch"
	"github.com/dpb587/bosh-dross/distich/data"
	yaml "gopkg.in/yaml.v2"
)

type Manifest struct {
	manifestPath string
	opsFilePaths []string

	appliedManifest []byte
	appliedDataNode data.Node
}

func NewManifest(manifestPath string, opsFilePaths []string) Manifest {
	return Manifest{
		manifestPath: manifestPath,
		opsFilePaths: opsFilePaths,
		// varFlags:     varFlags, // @todo
	}
}

func (m *Manifest) GetAppliedManifest() ([]byte, error) {
	if m.appliedManifest == nil {
		manifestBytes, err := ioutil.ReadFile(m.manifestPath)
		if err != nil {
			return nil, err
		}

		tpl := boshtpl.NewTemplate(manifestBytes)

		vars := boshtpl.StaticVariables{}
		ops := patch.Ops{}

		bytes, err := tpl.Evaluate(vars, ops, boshtpl.EvaluateOpts{})
		if err != nil {
			return nil, err
		}

		m.appliedManifest = bytes
	}

	return m.appliedManifest, nil
}

func (m *Manifest) GetAppliedDataNode() (data.Node, error) {
	if m.appliedDataNode == nil {
		var source map[interface{}]interface{}

		appliedManifest, err := m.GetAppliedManifest()
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(appliedManifest, &source)
		if err != nil {
			return nil, err
		}

		dataNode, err := data.CreateNode(source)
		if err != nil {
			return nil, err
		}

		m.appliedDataNode = dataNode
	}

	return m.appliedDataNode, nil
}

func (m *Manifest) GetKnownVariables() ([]string, error) {
	appliedDataNode, err := m.GetAppliedDataNode()
	if err != nil {
		return nil, err
	}

	variablesDataNode, err := appliedDataNode.Traverse("variables")
	if err == data.PathNotFound {
		return []string{}, nil
	} else if err != nil {
		return nil, err
	}

	variablesSimple, ok := variablesDataNode.Export().([]interface{})
	if !ok {
		return nil, fmt.Errorf("variables section is unknown type")
	}

	var knownVariables []string

	for variablesIdx, variables := range variablesSimple {
		variables, ok := variables.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("variable %d is unknown type", variablesIdx)
		}

		knownVariables = append(knownVariables, variables["name"].(string))
	}

	return knownVariables, nil
}
