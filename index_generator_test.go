package opensearchutil

import (
	"encoding/json"
	"testing"

	"github.com/onsi/gomega"
	"github.com/pkg/errors"
)

func TestIndexGenerator_GenerateIndexJson_buildsATree(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	mappingProperties := []MappingProperty{
		{
			FieldName: "id",
			FieldType: "integer",
			Children:  nil,
		},
		{
			FieldName: "price",
			FieldType: "float",
		},
		{
			FieldName: "location",
			Children: []MappingProperty{
				{
					FieldName: "full_address",
					FieldType: "text",
				},
				{
					FieldName: "confirmed",
					FieldType: "boolean",
				},
			},
		},
		{
			FieldName: "company",
			Children: []MappingProperty{
				{
					FieldName: "name",
					FieldType: "text",
				},
				{
					FieldName: "parent_company",
					Children: []MappingProperty{
						{
							FieldName: "name",
							FieldType: "text",
						},
					},
				},
			},
		},
	}

	resultJson, err := NewIndexGenerator().GenerateIndexJson(mappingProperties, nil)
	g.Expect(err).To(gomega.BeNil())

	assertJsonsEqual(g, resultJson, []byte(`{
   "mappings": {
      "properties": {
         "company": {
            "properties": {
               "name": {
                  "type": "text"
               },
               "parent_company": {
                  "properties": {
                     "name": {
                        "type": "text"
                     }
                  }
               }
            }
         },
         "id": {
            "type": "integer"
         },
         "location": {
            "properties": {
               "confirmed": {
                  "type": "boolean"
               },
               "full_address": {
                  "type": "text"
               }
            }
         },
         "price": {
            "type": "float"
         }
      }
   }
}`))
}

func TestIndexGenerator_GenerateIndexJson_addsFormatIfSpecified(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	mappingProperties := []MappingProperty{
		{
			FieldName:   "created_at",
			FieldType:   "date",
			FieldFormat: MakePtr("basic_time"),
		},
	}

	resultJson, err := NewIndexGenerator().GenerateIndexJson(mappingProperties, nil)
	g.Expect(err).To(gomega.BeNil())

	g.Expect(string(resultJson)).To(gomega.Equal(`{
   "mappings": {
      "properties": {
         "created_at": {
            "format": "basic_time",
            "type": "date"
         }
      }
   }
}`))
}

func TestIndexGenerator_GenerateIndexJson_addsSettings(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	resultJson, err := NewIndexGenerator().GenerateIndexJson([]MappingProperty{
		{
			FieldName: "id",
			FieldType: "integer",
		},
	}, &IndexSettings{
		NumberOfShards:  MakePtr(uint16(1)),
		Hidden:          MakePtr(true),
		RefreshInterval: MakePtr("-1"),
	})
	g.Expect(err).To(gomega.BeNil())

	assertJsonsEqual(g, resultJson, []byte(`{
   "mappings": {
      "properties": {
         "id": {
            "type": "integer"
         }
      }
   },
   "settings": {
      "hidden": true,
      "number_of_shards": 1,
      "refresh_interval": "-1"
   }
}`))
}

func TestIndexGenerator_GenerateIndexJson_addsDynamic(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	resultJson, err := NewIndexGenerator().GenerateIndexJson([]MappingProperty{
		{
			FieldName: "id",
			FieldType: "integer",
		},
	}, nil, WithStrictMapping(true))
	g.Expect(err).To(gomega.BeNil())

	assertJsonsEqual(g, resultJson, []byte(`{
   "mappings": {
	  "dynamic": "strict",
      "properties": {
         "id": {
            "type": "integer"
         }
      }
   }
}`))
}

func TestIndexGenerator_GenerateIndexJson_addsCustomProps(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	resultJson, err := NewIndexGenerator().GenerateIndexJson([]MappingProperty{
		{
			FieldName: "name",
			FieldType: "text",
			IndexPrefixes: MakePtr(map[string]string{
				"min_chars": "2",
				"max_chars": "10",
			}),
		},
	}, nil, WithStrictMapping(true))
	g.Expect(err).To(gomega.BeNil())

	assertJsonsEqual(g, resultJson, []byte(`{
   "mappings": {
	  "dynamic": "strict",
      "properties": {
         "name": {
            "type": "text",
			"index_prefixes": {
				"min_chars": "2",
				"max_chars": "10"
			}
         }
      }
   }
}`))
}

func TestIndexGenerator_GenerateIndexJson_addsAnalyzer(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	mappingProperties := []MappingProperty{
		{
			FieldName: "company_name",
			FieldType: "text",
			Analyzer:  MakePtr("standard"),
		},
	}

	resultJson, err := NewIndexGenerator().GenerateIndexJson(mappingProperties, nil)
	g.Expect(err).To(gomega.BeNil())

	g.Expect(string(resultJson)).To(gomega.Equal(`{
   "mappings": {
      "properties": {
         "company_name": {
            "analyzer": "standard",
            "type": "text"
         }
      }
   }
}`))
}

func TestIndexGenerator_GenerateIndexJson_addsSearchAnalyzerAndCopyTo(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	mappingProperties := []MappingProperty{
		{
			FieldName:      "title",
			FieldType:      "text",
			Analyzer:       MakePtr("standard"),
			SearchAnalyzer: MakePtr("english"),
			CopyTo:         []string{"all_text", "foo_text"},
		},
	}

	resultJson, err := NewIndexGenerator().GenerateIndexJson(mappingProperties, nil)
	g.Expect(err).To(gomega.BeNil())

	g.Expect(string(resultJson)).To(gomega.Equal(`{
   "mappings": {
      "properties": {
         "title": {
            "analyzer": "standard",
            "copy_to": [
               "all_text",
               "foo_text"
            ],
            "search_analyzer": "english",
            "type": "text"
         }
      }
   }
}`))
}

func makeJsonObj(jsonBytes []byte) (map[string]interface{}, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &m); err != nil {
		return map[string]interface{}{}, errors.Wrapf(err, "json.Unmarshal")
	}
	return m, nil
}

func assertJsonsEqual(g *gomega.WithT, jsonA []byte, jsonB []byte) {
	objA, err := makeJsonObj(jsonA)
	g.Expect(err).To(gomega.BeNil())

	objB, err := makeJsonObj(jsonB)
	g.Expect(err).To(gomega.BeNil())

	g.Expect(objA).To(gomega.Equal(objB))
}
