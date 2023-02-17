package opensearchutil

import (
	"github.com/onsi/gomega"
	"testing"
)

func TestIndexGenerator_GenerateIndexJson(t *testing.T) {
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

	resultJson, err := NewIndexGenerator().GenerateIndexJson(mappingProperties)
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
