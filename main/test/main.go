package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dpb587/bosh-dross/distich"
	"github.com/dpb587/bosh-dross/distich/data/visitors"
	"github.com/dpb587/bosh-dross/distich/schema"
	"github.com/dpb587/bosh-dross/distich/schema/loaders"
	"github.com/dpb587/bosh-dross/editor/form"
	"github.com/dpb587/bosh-dross/editor/form/fieldfactories"
	"github.com/dpb587/bosh-dross/manifest"
	"github.com/dpb587/bosh-dross/manifest/schemaguessers"
)

func main() {
	// contextual
	schemaLoader := schema.PrioritizedFactoryLoader{}
	schemaLoader.Add(loaders.NewRemappedLoader(loaders.LocalFile{}, "https://dpb587.github.io/bosh-json-schema/", "file://./bosh-json-schema/"), 50)
	schemaLoader.Add(loaders.LocalReleaseCache{}, 10)
	schemaResolver := schema.NewResolver(&schemaLoader)

	schema, err := schemaResolver.Load("https://dpb587.github.io/bosh-json-schema-release/v0/github.com/cloudfoundry/bosh/bosh/v264.5.0/jobs?sha1=c7133d0069b604e96bee203801bb8d686c9a98d9&url=https://bosh.io/d/github.com/cloudfoundry/bosh?v=264.5.0")
	if err != nil {
		panic(err)
	}

	mm1, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", mm1)

	return

	editorFieldFactory := form.PrioritizedFieldFactory{}
	editorFieldFactory.Add(fieldfactories.NewJSONSchema(&schemaResolver), 50)
	// editorFieldFactory.Add(awscpi.NewFieldFactory(), 10)

	// manifest
	manifest := manifest.NewManifest("tmp/bosh-aws.yml", []string{})
	manifestSchemaGuesser := schemaguessers.NewNaiveBuiltin(&schemaResolver)

	// guess schema
	schemaNode, err := manifestSchemaGuesser.Guess(manifest)
	if err != nil {
		panic(err)
	}

	mm, err := json.MarshalIndent(schemaNode, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", mm)

	// usefulness
	manifestNode, err := manifest.GetAppliedDataNode()
	if err != nil {
		panic(err)
	}

	// referenced variables
	vars := visitors.NewVariableCollector()
	manifestNode.Visit(&vars)

	// renderable
	for refVariable, refUsages := range vars.GetReferences() {
		path := refUsages[0]

		fmt.Printf("<h1>%s</h1>\n", refVariable)
		fmt.Printf("<p>Found at %s</p>\n", strings.Join(refUsages, ", "))

		traversedDataNode, traversedSchemaNode, err := distich.Traverse(manifestNode, *schemaNode, path)
		if err != nil {
			fmt.Printf(`  <div class="alert alert-danger" role="alert">
    <div class="grid">
      <div class="col col-middle">traverse: %s</div>
    </div>
  </div>
`, err)

			continue
		}

		fmt.Printf("<!--\n%#+v\n-->\n", traversedDataNode)
		fmt.Printf("<!--\n%#+v\n-->\n", traversedSchemaNode)

		// try
		field, err := editorFieldFactory.Create(traversedSchemaNode.ID, path, form.FieldOptions{})
		if err != nil {
			fmt.Printf(`  <div class="alert alert-danger" role="alert">
    <div class="grid">
      <div class="col col-middle">create: %s</div>
    </div>
  </div>
`, err)

			continue
		}

		err = field.Set(traversedDataNode.Export())
		if err != nil {
			fmt.Printf(`  <div class="alert alert-danger" role="alert">
    <div class="grid">
      <div class="col col-middle">set: %s</div>
    </div>
  </div>
`, err)

			continue
		}

		fmt.Printf("%s\n", field.HTML())
	}
}
