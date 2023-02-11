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

var _ = Describe("GenerateIndexJson", func() {
	It("Generates an index JSON string", func() {
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
		str, err := GenerateIndexJson(person{})
		Expect(err).To(BeNil())
		Expect(str).To(Equal(`{
   "mappings": {
      "properties": {
         "account_balance": {
            "type": "float"
         },
         "age": {
            "type": "integer"
         },
         "home_loc": {
            "properties": {
               "confirmed": {
                  "type": "boolean"
               },
               "full_address": {
                  "type": "text"
               }
            }
         },
         "is_dead": {
            "type": "boolean"
         },
         "name": {
            "type": "text"
         },
         "social_security": {
            "type": "text"
         },
         "work_loc": {
            "properties": {
               "confirmed": {
                  "type": "boolean"
               },
               "full_address": {
                  "type": "text"
               }
            }
         }
      }
   },
   "settings": {
      "number_of_replicas": 2,
      "number_of_shards": 1
   }
}`))
	})
})

var _ = Describe("getMappingProperties", func() {
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
		mp, err := getMappingProperties(person{})
		Expect(err).To(BeNil())
		Expect(mp).To(ConsistOf(
			mappingProperty{
				FieldName: "name",
				FieldType: "text",
			},
			mappingProperty{
				FieldName: "age",
				FieldType: "integer",
			},
			mappingProperty{
				FieldName: "account_balance",
				FieldType: "float",
			},
			mappingProperty{
				FieldName: "is_dead",
				FieldType: "boolean",
			},
			mappingProperty{
				FieldName: "social_security",
				FieldType: "text",
			},
			mappingProperty{
				FieldName: "home_loc",
				Children: []mappingProperty{
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
			mappingProperty{
				FieldName: "work_loc",
				Children: []mappingProperty{
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
