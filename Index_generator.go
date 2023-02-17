package opensearchutil

import (
	"bytes"
	_ "embed"
	"github.com/pkg/errors"
	"text/template"
)

//go:embed index.gotmpl
var indexTmpl string

type IndexGenerator struct {
	optionContainer indexGeneratorOptionContainer
}

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

func (g *IndexGenerator) GenerateIndexJson(mappingProperties []MappingProperty) ([]byte, error) {
	type indexTmplData struct {
		MappingProperties []MappingProperty
	}

	var funcMap template.FuncMap = map[string]interface{}{
		"notLast": func(index int, len int) bool {
			return index+1 < len
		},
	}

	tmpl, err := template.New("IndexTmpl").Funcs(funcMap).Parse(indexTmpl)
	if err != nil {
		return nil, errors.Wrapf(err, "parse template")
	}

	var tmplResult bytes.Buffer
	err = tmpl.Execute(&tmplResult, indexTmplData{
		MappingProperties: mappingProperties,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "tpl.Execute")
	}

	formattedJson, err := g.optionContainer.jsonFormatter.FormatJson(tmplResult.Bytes())
	if err != nil {
		return nil, errors.Wrapf(err, "formatJson")
	}

	return formattedJson, nil
}
