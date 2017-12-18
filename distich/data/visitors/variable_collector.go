package visitors

import (
	"regexp"

	"github.com/dpb587/bosh-dross/distich/data"
)

var interpolationRegex = regexp.MustCompile(`\(\((!?[-/\.\w\pL]+)\)\)`)

type VariableCollector struct {
	references map[string][]string
}

func NewVariableCollector() VariableCollector {
	return VariableCollector{
		references: map[string][]string{},
	}
}

var _ data.NodeVisitor = &VariableCollector{}

func (t *VariableCollector) EnterNode(node data.Node) error {
	if _, ok := node.(*data.ScalarNode); !ok {
		return nil
	}

	export, ok := node.Export().(string)
	if !ok {
		return nil
	}

	varrefs := interpolationRegex.FindAllString(export, -1)
	if len(varrefs) == 0 {
		return nil
	}

	for _, varref := range varrefs {
		varpath := varref[2 : len(varref)-2]

		if _, ok := t.references[varpath]; !ok {
			t.references[varpath] = []string{}
		}

		t.references[varpath] = append(t.references[varpath], node.GetPath())
	}

	return nil
}

func (t *VariableCollector) LeaveNode(node data.Node) error {
	return nil
}

func (t *VariableCollector) GetReferences() map[string][]string {
	return t.references
}
