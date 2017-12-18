package fieldfactories

import (
	"errors"

	"github.com/dpb587/bosh-dross/distich/schema"
	"github.com/dpb587/bosh-dross/editor/form"
	"github.com/dpb587/bosh-dross/editor/form/fields"
)

type JSONSchema struct {
	resolver *schema.Resolver
}

var _ form.FieldFactory = &JSONSchema{}

func NewJSONSchema(resolver *schema.Resolver) JSONSchema {
	return JSONSchema{
		resolver: resolver,
	}
}

func (JSONSchema) IsSupported(uri string) bool {
	return true // @todo not really
}

func (ff JSONSchema) Create(uri, path string, options form.FieldOptions) (form.Field, error) {
	schema, err := ff.resolver.Load(uri)
	if err != nil {
		return nil, err
	}

	baseField := fields.BaseField{
		Path:  path,
		Title: schema.Title,
		Help:  schema.Description,
	}

	baseField.Name = baseField.ID()

	if len(schema.Enum) > 0 {
		options := map[interface{}]string{}

		for _, value := range schema.Enum {
			options[value] = value
		}

		return &fields.Select{
			BaseField: baseField,
			Options:   options,
		}, nil
	} else if schema.Type == "string" || schema.Type == "" {
		return &fields.Select{
			BaseField: baseField,
		}, nil
	} else if schema.Type == "object" {
		// yamlconv
	} else if schema.Type == "integer" || schema.Type == "number" {
		return &fields.Number{
			BaseField: baseField,
		}, nil
	} else if schema.Type == "boolean" {
		return &fields.Select{
			BaseField: baseField,
			Options: map[interface{}]string{
				true:  "Enabled",
				false: "Disabled",
			},
		}, nil
	}

	return nil, errors.New("Unexpected field type")
	// $rawSchema = $schema->getSchema();
	//       } elseif ((!isset($rawSchema->type)) || ('string' == $rawSchema->type)) {
	//           $builder->add($name, 'text', $formOptions);
	//       } elseif ('object' == $rawSchema->type) {
	//           $builder->add($name, 'form', $formOptions);
	//           $subBuilder = $builder->get($name);
	//           foreach ($rawSchema->properties as $propertyName => $propertyRelativeSchema) {
	//               $formOptions = [];
	//               if (empty($rawSchema->required)) {
	//                   $formOptions['required'] = false;
	//               } elseif (in_array($propertyName, $rawSchema->required)) {
	//                   $formOptions['required'] = true;
	//               } else {
	//                   $formOptions['required'] = false;
	//               }
	//               $this->buildForm(
	//                   $subBuilder,
	//                   new ArraySchemaNode($this->jsonSchema->getSchema($this->getSchemaPath($rawSchema->id, '/properties/' . $propertyName))),
	//                   $propertyName,
	//                   $formOptions
	//               );
	//           }
	//       } else {
	//           throw new \LogicException(sprintf('Unsupported field type: %s', $rawSchema->type));
	//       }
	//       } elseif (isset($rawSchema->items)) {
	//           $builder->add($name, 'form');
	//
	//       } elseif (isset($rawSchema->oneOf)) {
	//           $formOptions['forms'] = [];
	//           foreach ($rawSchema->oneOf as $oneOfIdx => $oneOf) {
	//               $oneOfSchema = $this->getResolvedSchema(new ArraySchemaNode($oneOf));
	//               $formOptions['forms'][$oneOfIdx] = $this->buildForm(
	//                   $builder->create('stub', 'form'),
	//                   $oneOfSchema,
	//                   'via_' . $oneOfIdx,
	//                   [
	//                       'label' => isset($oneOfSchema->getSchema()->title) ? $oneOfSchema->getSchema()->title : null,
	//                   ]
	//               );
	//           }
	//           $builder->add(
	//               $name,
	//               'veneer_core_form_picker',
	//               $formOptions
	//           );
}
