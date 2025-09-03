package opensearchutil

import (
	"errors"
	"testing"
	"time"

	"github.com/onsi/gomega"
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
	for _, mp := range mps {
		g.Expect(mp.GetDepth() <= DefaultMaxDepth).To(gomega.BeTrue())
	}
}

func TestMappingPropertiesBuilder_BuildMappingProperties_DoesNotExceedGivenMaxDepthWithRecursiveField(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	type location struct {
		name string
		loc  *location
	}

	const depth = 3

	builder := NewMappingPropertiesBuilder(WithMaxDepth(depth))

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
	for _, mp := range mps {
		g.Expect(mp.GetDepth() <= depth).To(gomega.BeTrue())
	}
}

func TestMappingPropertiesBuilder_BuildMappingProperties_SetsCustomProps(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	type person struct {
		BusinessName string `opensearch:"index_prefixes:min_chars=2;max_chars=10"`
	}

	builder := NewMappingPropertiesBuilder()
	mps, err := builder.BuildMappingProperties(person{})
	g.Expect(err).To(gomega.BeNil())
	g.Expect(mps).To(gomega.HaveLen(1))
	g.Expect(mps[0].IndexPrefixes).ToNot(gomega.BeNil())
	g.Expect(*mps[0].IndexPrefixes).To(gomega.Equal(map[string]string{
		"min_chars": "2",
		"max_chars": "10",
	}))
}

func TestMappingPropertiesBuilder_BuildMappingProperties_SetsAnalyzer(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	type person struct {
		CompanyName string `opensearch:"analyzer:keyword"`
	}

	builder := NewMappingPropertiesBuilder()
	mps, err := builder.BuildMappingProperties(person{})
	g.Expect(err).To(gomega.BeNil())
	g.Expect(mps).To(gomega.HaveLen(1))
	g.Expect(mps[0].Analyzer).ToNot(gomega.BeNil())
	g.Expect(*mps[0].Analyzer).To(gomega.Equal("keyword"))
}

func TestMappingPropertiesBuilder_BuildMappingProperties_SetsSearchAnalyzerAndCopyTo(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	type doc struct {
		Title string `opensearch:"type:text,search_analyzer:english,copy_to:all_text,analyzer:standard"`
		All   string `opensearch:"type:text"`
	}

	builder := NewMappingPropertiesBuilder()
	mps, err := builder.BuildMappingProperties(doc{})
	g.Expect(err).To(gomega.BeNil())
	g.Expect(mps).To(gomega.HaveLen(2))

	// Find Title and All
	var title MappingProperty
	var all MappingProperty
	for _, mp := range mps {
		if mp.FieldName == "title" {
			title = mp
		}
		if mp.FieldName == "all" {
			all = mp
		}
	}

	g.Expect(title.FieldType).To(gomega.Equal("text"))
	g.Expect(title.Analyzer).ToNot(gomega.BeNil())
	g.Expect(*title.Analyzer).To(gomega.Equal("standard"))
	g.Expect(title.SearchAnalyzer).ToNot(gomega.BeNil())
	g.Expect(*title.SearchAnalyzer).To(gomega.Equal("english"))
	g.Expect(title.CopyTo).To(gomega.Equal([]string{"all_text"}))

	g.Expect(all.FieldType).To(gomega.Equal("text"))
}

func TestMappingPropertiesBuilder_BuildMappingProperties_ErrorsByDefaultWithUnsupportedType(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	type person struct {
		addresses interface{} // no support
	}

	builder := NewMappingPropertiesBuilder()
	_, err := builder.BuildMappingProperties(person{})
	g.Expect(err).ToNot(gomega.BeNil())
	g.Expect(err.Error()).To(gomega.ContainSubstring("field not supported"))
}

func TestMappingPropertiesBuilder_BuildMappingProperties_DoesNotErrorWithUnsupportedTypeIfOptionProvided(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	type person struct {
		addresses interface{} // no support
	}

	builder := NewMappingPropertiesBuilder(OmitUnsupportedTypes())
	_, err := builder.BuildMappingProperties(person{})
	g.Expect(err).To(gomega.BeNil())
}

func TestMappingPropertiesBuilder_BuildMappingProperties_SupportsObjectSlices(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	type location struct {
		city string
	}
	type person struct {
		addresses []location
	}

	builder := NewMappingPropertiesBuilder()
	mps, err := builder.BuildMappingProperties(person{})
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(mps).To(gomega.HaveLen(1))
	g.Expect(mps[0]).To(gomega.Equal(MappingProperty{
		FieldName: "addresses",
		Children: []MappingProperty{
			{
				FieldName: "city",
				FieldType: "text",
			},
		},
	}))
}

func TestMappingPropertiesBuilder_BuildMappingProperties_SupportsPrimitiveSlices(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	type person struct {
		names []string
	}

	builder := NewMappingPropertiesBuilder()
	mps, err := builder.BuildMappingProperties(person{})
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(mps).To(gomega.HaveLen(1))
	g.Expect(mps[0]).To(gomega.Equal(MappingProperty{
		FieldName: "names",
		FieldType: "text",
	}))
}

func TestMappingPropertiesBuilder_BuildMappingProperties_SliceRecursiveMaxDepth(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	type person struct {
		siblings []person
	}

	const depth = 3

	builder := NewMappingPropertiesBuilder(WithMaxDepth(depth))

	expectedMappingProperties := []MappingProperty{ // Level 1
		{
			FieldName: "siblings",
			Children: []MappingProperty{ // Level 2
				{
					FieldName: "siblings",
					Children:  nil, // Level 3, nothing to map, recursion stopped
				},
			},
		},
	}

	mps, err := builder.BuildMappingProperties(person{})
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(mps).To(gomega.Equal(expectedMappingProperties))
	for _, mp := range mps {
		g.Expect(mp.GetDepth() <= depth).To(gomega.BeTrue())
	}
}
