package schemaguessers

import (
	"github.com/dpb587/bosh-dross/distich/schema"
	"github.com/dpb587/bosh-dross/manifest"
)

type NaiveBuiltin struct {
	resolver *schema.Resolver
}

var _ manifest.SchemaGuesser = &NaiveBuiltin{}

func NewNaiveBuiltin(resolver *schema.Resolver) NaiveBuiltin {
	return NaiveBuiltin{
		resolver: resolver,
	}
}

func (sg NaiveBuiltin) Guess(m manifest.Manifest) (*schema.Node, error) {
	data, err := m.GetAppliedDataNode()
	if err != nil {
		return nil, err
	}

	if _, err := data.Traverse("cloud_provider"); err == nil {
		return sg.resolver.Load("https://dpb587.github.io/bosh-json-schema/v0/director/v0/deployment-v2.json")
		// return sg.resolver.Load("https://dpb587.github.io/bosh-json-schema/v0/director/v0/create-env.json") // @todo more correct
	} else if _, err := data.Traverse("azs"); err == nil {
		return sg.resolver.Load("https://dpb587.github.io/bosh-json-schema/v0/director/v0/cloud-config.json")
	} else if _, err := data.Traverse("update"); err == nil {
		if _, err := data.Traverse("networks"); err == nil {
			return sg.resolver.Load("https://dpb587.github.io/bosh-json-schema/v0/director/v0/deployment.json")
		} else {
			return sg.resolver.Load("https://dpb587.github.io/bosh-json-schema/v0/director/v0/deployment-v2.json")
		}
	} else if _, err := data.Traverse("addons"); err == nil {
		return sg.resolver.Load("https://dpb587.github.io/bosh-json-schema/v0/director/v0/runtime-config.json")
	}

	return nil, nil
}
