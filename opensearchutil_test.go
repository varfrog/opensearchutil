package opensearchutil

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPrintFieldsForOpensearch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PrintFieldsForOpensearch Suite")
}

var _ = Describe("BuildMappingProperties", func() {
	It("Gets field names and their types", func() {
		type location struct {
			FullAddress string
			Confirmed   bool
		}
		type person struct {
			Name           string
			Age            uint8
			AccountBalance float64
			IsDead         bool
			HomeLoc        location
			WorkLoc        *location
			SocialSecurity *string
		}
		mp, err := BuildMappingProperties(person{})
		Expect(err).To(BeNil())
		Expect(mp).To(ConsistOf(
			MappingProperty{
				FieldName: "name",
				FieldType: "text",
			},
			MappingProperty{
				FieldName: "age",
				FieldType: "integer",
			},
			MappingProperty{
				FieldName: "account_balance",
				FieldType: "float",
			},
			MappingProperty{
				FieldName: "is_dead",
				FieldType: "boolean",
			},
			MappingProperty{
				FieldName: "social_security",
				FieldType: "text",
			},
			MappingProperty{
				FieldName: "home_loc",
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
			MappingProperty{
				FieldName: "work_loc",
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
		))
	})
})

var _ = Describe("GenerateIndexJson", func() {
	It("Generates an index JSON string", func() {
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
		resultJson, err := GenerateIndexJson(mappingProperties)
		Expect(err).To(BeNil())
		Expect(resultJson).To(Equal([]byte(`{
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
}`)))
	})
})

var _ = Describe("toSnakeCase", func() {
	It("Converts strings to snake_case", func() {
		Expect(toSnakeCase("Foo")).To(Equal("foo"))
		Expect(toSnakeCase("FOo")).To(Equal("foo"))
		Expect(toSnakeCase("FOO")).To(Equal("foo"))
		Expect(toSnakeCase("foo")).To(Equal("foo"))
		Expect(toSnakeCase("fooBar")).To(Equal("foo_bar"))
		Expect(toSnakeCase("FooBar")).To(Equal("foo_bar"))
		Expect(toSnakeCase("FooBarBaz123")).To(Equal("foo_bar_baz_123"))
		Expect(toSnakeCase("Foo123BarBaz123")).To(Equal("foo_123_bar_baz_123"))
		Expect(toSnakeCase("123")).To(Equal("123"))
		Expect(toSnakeCase("123foo")).To(Equal("123foo"))
		Expect(toSnakeCase("123Foo")).To(Equal("123_foo"))
	})
})
