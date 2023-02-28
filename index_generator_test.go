package opensearchutil

import (
	"github.com/onsi/gomega"
	"testing"
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
	g.Expect(string(resultJson)).To(gomega.Equal(`{
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
	g.Expect(string(resultJson)).To(gomega.Equal(`{
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
