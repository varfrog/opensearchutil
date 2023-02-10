package main

import (
	"fmt"
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
		type address struct {
			FullAddress string
			HouseNumber uint8
		}
		var person struct {
			Name           string
			Age            uint8
			AccountBalance float64
			IsDead         bool
			Address        address
			SocialSecurity *string
		}

		str, err := GenerateIndexJson(person)
		Expect(err).To(BeNil())
		fmt.Printf("%s\n", str)
	})
})

var _ = Describe("getMappingProperties", func() {
	It("Gets field names and their types", func() {
		type address struct {
			FullAddress string
		}
		var person struct {
			Name           string
			Age            uint8
			AccountBalance float64
			IsDead         bool
			Address        address
			SocialSecurity *string
		}
		mp, err := getMappingProperties(person)
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
