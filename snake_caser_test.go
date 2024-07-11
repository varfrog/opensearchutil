package opensearchutil

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestSnakeCaser_TransformFieldName(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	snakeCaser := NewSnakeCaser()

	tests := []struct {
		give string
		want string
	}{
		{
			give: "Foo",
			want: "foo",
		},
		{
			give: "FOo",
			want: "foo",
		},
		{
			give: "FOO",
			want: "foo",
		},
		{
			give: "foo",
			want: "foo",
		},
		{
			give: "fooBar",
			want: "foo_bar",
		},
		{
			give: "FooBar",
			want: "foo_bar",
		},
		{
			give: "FooBarBaz123",
			want: "foo_bar_baz_123",
		},
		{
			give: "Foo123BarBaz123",
			want: "foo_123_bar_baz_123",
		},
		{
			give: "123",
			want: "123",
		},
		{
			give: "123foo",
			want: "123foo",
		},
		{
			give: "123Foo",
			want: "123_foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			res, err := snakeCaser.TransformFieldName(tt.want)
			g.Expect(err).To(gomega.BeNil())
			g.Expect(res).To(gomega.Equal(tt.want))
		})
	}
}
