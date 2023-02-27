package opensearchutil

import (
	_ "embed"
	"encoding/json"
	"github.com/pkg/errors"
)

type IndexGenerator struct {
	optionContainer indexGeneratorOptionContainer
}

type (
	index struct {
		Mappings parentNode `json:"mappings"`
	}
	parentNode struct {
		// Property maps from a property name to another parentNode or to a leafNode
		Properties map[string]interface{} `json:"properties"`
	}
	leafNode struct {
		Type   string  `json:"type"`
		Format *string `json:"format,omitempty"`
	}
)

func NewIndexGenerator(options ...IndexGeneratorOption) *IndexGenerator {
	optContainer := indexGeneratorOptionContainer{}
	for _, o := range options {
		o.apply(&optContainer)
	}
	if optContainer.jsonFormatter == nil {
		optContainer.jsonFormatter = NewMarshalIndentJsonFormatter()
	}

	return &IndexGenerator{optionContainer: optContainer}
}

func (g *IndexGenerator) GenerateIndexJson(
	mappingProperties []MappingProperty,
	_ *IndexSettings, // todo
) ([]byte, error) {
	jsonBytes, err := json.Marshal(index{
		Mappings: parentNode{
			Properties: g.buildProperties(mappingProperties),
		},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "json.Marshal")
	}

	formattedJson, err := g.optionContainer.jsonFormatter.FormatJson(jsonBytes)
	if err != nil {
		return nil, errors.Wrapf(err, "formatJson")
	}

	return formattedJson, nil
}

func (g *IndexGenerator) buildProperties(mappingProperties []MappingProperty) map[string]interface{} {
	m := make(map[string]interface{}, len(mappingProperties))
	for _, mp := range mappingProperties {
		if mp.Children == nil {
			m[mp.FieldName] = leafNode{
				Type:   mp.FieldType,
				Format: mp.FieldFormat,
			}
		} else {
			m[mp.FieldName] = parentNode{Properties: g.buildProperties(mp.Children)}
		}
	}
	return m
}
