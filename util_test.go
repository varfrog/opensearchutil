package opensearchutil

import (
	"reflect"
	"testing"
	"time"

	"github.com/onsi/gomega"
)

func Test_MakePtr(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	str := "aaa"
	num := 123
	obj := struct {
		name string
	}{name: "tom"}

	g.Expect(*MakePtr(str)).To(gomega.Equal(str))
	g.Expect(*MakePtr(num)).To(gomega.Equal(num))
	g.Expect(*MakePtr(obj)).To(gomega.Equal(obj))
}

func Test_getTagOptionValue(t *testing.T) {
	const tagKey = "opensearch"

	g := gomega.NewGomegaWithT(t)

	type foo struct {
		a string    // No tag
		b string    `opensearch:"type:keyword"`
		c time.Time `opensearch:"type:date, format:basic_time"`
		d string    `opensearch:"index_prefixes:min_chars=2;max_chars=10"`
		e string    `opensearch:"type:text, index_prefixes:min_chars=2;max_chars=10"`
	}

	v := reflect.TypeOf(foo{})
	g.Expect(getTagOptionValue(v.Field(0), tagKey, "type")).To(gomega.Equal(""))

	g.Expect(getTagOptionValue(v.Field(1), tagKey, "type")).To(gomega.Equal("keyword"))
	g.Expect(getTagOptionValue(v.Field(1), tagKey, "format")).To(gomega.Equal(""))

	g.Expect(getTagOptionValue(v.Field(2), tagKey, "type")).To(gomega.Equal("date"))
	g.Expect(getTagOptionValue(v.Field(2), tagKey, "format")).To(gomega.Equal("basic_time"))

	g.Expect(getTagOptionValue(v.Field(3), tagKey, "index_prefixes")).To(gomega.Equal(`min_chars=2;max_chars=10`))

	g.Expect(getTagOptionValue(v.Field(4), tagKey, "type")).To(gomega.Equal("text"))
	g.Expect(getTagOptionValue(v.Field(4), tagKey, "index_prefixes")).To(gomega.Equal(`min_chars=2;max_chars=10`))
}

func Test_parseCustomPropertyValue(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	g.Expect(parseCustomPropertyValue("min_chars=2;foo=bar")).To(gomega.Equal(map[string]string{
		"min_chars": "2",
		"foo":       "bar",
	}))
}
