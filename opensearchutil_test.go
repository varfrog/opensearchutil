package opensearchutil

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestPrintFieldsForOpensearch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PrintFieldsForOpensearch Suite")
}

var _ = Describe("BuildMappingProperties", func() {
	It("Works with primitive values and their pointers (test with ints)", func() {
		type person struct {
			Age  uint8
			Age2 *uint8
		}
		mp, err := BuildMappingProperties(person{})
		Expect(err).To(BeNil())
		Expect(mp).To(ConsistOf(
			MappingProperty{
				FieldName: "age",
				FieldType: "integer",
			},
			MappingProperty{
				FieldName: "age_2",
				FieldType: "integer",
			},
		))
	})

	It("Works with primitive values and their pointers (test with strings)", func() {
		type person struct {
			Name  string
			Name2 *string
		}
		mp, err := BuildMappingProperties(person{})
		Expect(err).To(BeNil())
		Expect(mp).To(ConsistOf(
			MappingProperty{
				FieldName: "name",
				FieldType: "text",
			},
			MappingProperty{
				FieldName: "name_2",
				FieldType: "text",
			},
		))
	})

	It("Works with custom structs and their pointers", func() {
		type location struct {
			FullAddress string
		}
		type person struct {
			HomeLoc location
			WorkLoc *location
		}
		mp, err := BuildMappingProperties(person{})
		Expect(err).To(BeNil())
		Expect(mp).To(ConsistOf(
			MappingProperty{
				FieldName: "home_loc",
				Children: []MappingProperty{
					{
						FieldName: "full_address",
						FieldType: "text",
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
				},
			},
		))
	})

	It("Sets the specified type or falls back to default", func() {
		type person struct {
			Name  string
			Email string `opensearch:"type:keyword"`
		}
		mp, err := BuildMappingProperties(person{})
		Expect(err).To(BeNil())
		Expect(mp).To(ConsistOf(
			MappingProperty{
				FieldName: "name",
				FieldType: "text",
			},
			MappingProperty{
				FieldName: "email",
				FieldType: "keyword",
			},
		))
	})

	It("Sets the specified type for time.Time or falls back to the default for time.Time", func() {
		Skip("cunt") // todo unskip
		type person struct {
			Created time.Time
			DOB     time.Time `opensearch:"type:basic_date"`
		}
		mp, err := BuildMappingProperties(person{})
		Expect(err).To(BeNil())
		Expect(mp).To(ConsistOf(
			MappingProperty{
				FieldName: "created",
				FieldType: "basic_date_time",
			},
			MappingProperty{
				FieldName: "dob",
				FieldType: "basic_date",
			},
		))
	})

	It("Does not exceed default MaxDepth when there is recursion", func() {
		type location struct {
			name string
			loc  *location
		}
		mappingProperties, err := BuildMappingProperties(location{})
		Expect(err).To(BeNil())
		Expect(mappingProperties).To(Equal([]MappingProperty{ // Depth Level 1
			{FieldName: "name", FieldType: "text"},
			{
				FieldName: "loc",
				Children: []MappingProperty{ // Depth level 2
					// No field for "loc", as MaxDepth is reached
					{FieldName: "name", FieldType: "text"},
				},
			},
		}))
	})

	It("Does not exceed the given MaxDepth when there is recursion", func() {
		type location struct {
			name string
			loc  *location
		}
		mappingProperties, err := BuildMappingProperties(location{}, WithMaxDepth(3))
		Expect(err).To(BeNil())

		expectedMappingProperties := []MappingProperty{ // Level 1
			{FieldName: "name", FieldType: "text"},
			{
				FieldName: "loc",
				Children: []MappingProperty{ // Level 2
					{FieldName: "name", FieldType: "text"},
					{
						FieldName: "loc",
						Children: []MappingProperty{ // Level 3
							{FieldName: "name", FieldType: "text"},
						},
					},
				},
			},
		}
		Expect(mappingProperties).To(Equal(expectedMappingProperties))
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
