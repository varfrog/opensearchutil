package opensearchutil

import (
	"github.com/onsi/gomega"
	"reflect"
	"testing"
	"time"
)

func Test_makePtr(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	str := "aaa"
	num := 123
	obj := struct {
		name string
	}{name: "tom"}

	g.Expect(*makePtr(str)).To(gomega.Equal(str))
	g.Expect(*makePtr(num)).To(gomega.Equal(num))
	g.Expect(*makePtr(obj)).To(gomega.Equal(obj))
}

func Test_getTagOptionValue(t *testing.T) {
	const tagKey = "opensearch"

	g := gomega.NewGomegaWithT(t)

	type foo struct {
		a string    // No tag
		b string    `opensearch:"type:keyword"`
		c time.Time `opensearch:"type:date,format:basic_time"`
	}

	v := reflect.TypeOf(foo{})
	g.Expect(getTagOptionValue(v.Field(0), tagKey, "type")).To(gomega.Equal(""))

	g.Expect(getTagOptionValue(v.Field(1), tagKey, "type")).To(gomega.Equal("keyword"))
	g.Expect(getTagOptionValue(v.Field(1), tagKey, "format")).To(gomega.Equal(""))

	g.Expect(getTagOptionValue(v.Field(2), tagKey, "type")).To(gomega.Equal("date"))
	g.Expect(getTagOptionValue(v.Field(2), tagKey, "format")).To(gomega.Equal("basic_time"))
}
