package schemaguessers

import (
	"github.com/dpb587/bosh-dross/distich/data/visitors"
	"github.com/dpb587/bosh-dross/distich/schema"
	"github.com/dpb587/bosh-dross/manifest"
)

type IntrospectiveBuiltin struct {
	resolver *schema.Resolver
	builtin  *NaiveBuiltin
}

var _ manifest.SchemaGuesser = &IntrospectiveBuiltin{}

type introspectiveBuiltinRelease struct {
	Name    string
	Version string
	URL     string
	SHA1    string
}

func NewIntrospectiveBuiltin(resolver *schema.Resolver, builtin *NaiveBuiltin) IntrospectiveBuiltin {
	return IntrospectiveBuiltin{
		resolver: resolver,
		builtin:  builtin,
	}
}

// @todo this shouldn't modify schema pointer; copy?
func (sg IntrospectiveBuiltin) Guess(m manifest.Manifest) (*schema.Node, error) {
	schemaNode, err := sg.builtin.Guess(m)
	if err != nil {
		return nil, err
	}

	data, err := m.GetAppliedDataNode()
	if err != nil {
		return nil, err
	}

	if _, err := data.Traverse("releases"); err == nil {
		releaseCollector := visitors.ReleaseCollector{}

		err := data.Visit(&releaseCollector)
		if err != nil {
			return nil, err
		}

		// releases := releaseCollector.GetReleases()
		//
		// for _, igNode := range schemaNode.Properties["instance_groups"] {
		// 	igData, err := data.Traverse("instance_groups")
		// 	if err != nil {
		// 		return nil, err
		// 	}
		//
		// 	for jobIdx, jobNode := range igNode.Properties {
		// 		jobData, err := data.Traverse(jobIdx)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		// optimistic
		// 		jobReleaseNameData, _ := jobData.Traverse("release")
		// 		jobReleaseJobData, _ := jobData.Traverse("job")
		// 		props, err := sg.resolver.Load(fmt.Sprintf("%s#/definitions/%s_properties", releases[jobReleaseNameData.Export().(string)].SchemaURL(), jobReleaseJobData.Export().(string)))
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		jobNode.Properties["properties"] = props
		// 	}
		// }

		// @todo addons
	}

	return schemaNode, err
}
