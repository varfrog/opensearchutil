package opensearchutil

import (
	"errors"
	"github.com/onsi/gomega"
	"testing"
	"time"
)

func TestMappingPropertiesBuilder_BuildMappingProperties_PrimitivesAndTheirPtrs(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	type person struct {
		Age      uint8
		Age2     *uint8
		Name     string
		Name2    *string
		Balance  float64
		Balance2 *float64
	}

	builder := NewMappingPropertiesBuilder()
	mps, err := builder.BuildMappingProperties(person{})
	g.Expect(err).To(gomega.BeNil())
	g.Expect(mps).To(gomega.ConsistOf(
		MappingProperty{
			FieldName: "age",
			FieldType: "integer",
		},
		MappingProperty{
			FieldName: "age_2",
			FieldType: "integer",
		},
		MappingProperty{
			FieldName: "name",
			FieldType: "text",
		},
		MappingProperty{
			FieldName: "name_2",
			FieldType: "text",
		},
		MappingProperty{
			FieldName: "balance",
			FieldType: "float",
		},
		MappingProperty{
			FieldName: "balance_2",
			FieldType: "float",
		},
	))
}

func TestMappingPropertiesBuilder_BuildMappingProperties_StructsAndTheirPtrs(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	type location struct {
		FullAddress string
	}
	type person struct {
		HomeLoc location
		WorkLoc *location
	}

	builder := NewMappingPropertiesBuilder()
	mps, err := builder.BuildMappingProperties(person{})
	g.Expect(err).To(gomega.BeNil())
	g.Expect(mps).To(gomega.ConsistOf(
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
}

func TestMappingPropertiesBuilder_BuildMappingProperties_SetsSpecifiedTypeOrFallsBackToDefault(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	type person struct {
		Name  string
		Email string `opensearch:"type:keyword"`
	}

	builder := NewMappingPropertiesBuilder()
	mps, err := builder.BuildMappingProperties(person{})
	g.Expect(err).To(gomega.BeNil())
	g.Expect(mps).To(gomega.ConsistOf(
		MappingProperty{
			FieldName: "name",
			FieldType: "text",
		},
		MappingProperty{
			FieldName: "email",
			FieldType: "keyword",
		},
	))
}

func TestMappingPropertiesBuilder_BuildMappingProperties_DefaultsToCorrectTimeFormats(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	type foo struct {
		A TimeBasicDateTime
		B TimeBasicDateTimeNoMillis
		C TimeBasicDate
	}

	builder := NewMappingPropertiesBuilder()
	mps, err := builder.BuildMappingProperties(foo{})
	g.Expect(err).To(gomega.BeNil())
	g.Expect(mps).To(gomega.ConsistOf(
		MappingProperty{
			FieldName:   "a",
			FieldType:   "date",
			FieldFormat: MakePtr("basic_date_time"),
		},
		MappingProperty{
			FieldName:   "b",
			FieldType:   "date",
			FieldFormat: MakePtr("basic_date_time_no_millis"),
		},
		MappingProperty{
			FieldName:   "c",
			FieldType:   "date",
			FieldFormat: MakePtr("basic_date"),
		},
	))
}

func TestMappingPropertiesBuilder_BuildMappingProperties_ErrorsWhenFieldIsTime(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	type person struct {
		DOB time.Time `opensearch:"format:basic_date"`
	}

	builder := NewMappingPropertiesBuilder()
	_, err := builder.BuildMappingProperties(person{})
	g.Expect(errors.Is(err, ErrGotBuiltInTimeField)).To(gomega.BeTrue())
}

func TestMappingPropertiesBuilder_BuildMappingProperties_DoesNotExceedDefaultMaxDepthWithRecursiveField(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	type location struct {
		name string
		loc  *location
	}

	builder := NewMappingPropertiesBuilder()
	mps, err := builder.BuildMappingProperties(location{})
	g.Expect(err).To(gomega.BeNil())
	g.Expect(mps).To(gomega.Equal(
		[]MappingProperty{ // Depth Level 1
			{FieldName: "name", FieldType: "text"},
			{
				FieldName: "loc",
				Children: []MappingProperty{ // Depth level 2
					// No field for "loc", as MaxDepth is reached
					{FieldName: "name", FieldType: "text"},
				},
			},
		},
	))
}

func TestMappingPropertiesBuilder_BuildMappingProperties_DoesNotExceedGivenMaxDepthWithRecursiveField(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	type location struct {
		name string
		loc  *location
	}

	builder := NewMappingPropertiesBuilder(WithMaxDepth(3))

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

	mps, err := builder.BuildMappingProperties(location{})
	g.Expect(err).To(gomega.BeNil())
	g.Expect(mps).To(gomega.Equal(expectedMappingProperties))
}
