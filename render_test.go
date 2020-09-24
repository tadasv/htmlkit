package htmlkit

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkBuildAndRender(b *testing.B) {
	buf := &bytes.Buffer{}
	for i := 0; i < b.N; i++ {
		node := NewElementBuilder("strong").WithChild(NewRawNode("<span>a</span>")).WithChild(NewElementBuilder("a").Build(nil)).Build(nil)
		Render(buf, node)
	}
}

func TestRender(t *testing.T) {
	testCase := []struct {
		node     *Node
		expected string
	}{
		{
			node:     NewRawNode("<b>testing</b>"),
			expected: "<b>testing</b>",
		},
		{
			node:     NewCommentNode("this is a comment"),
			expected: "<!-- this is a comment -->",
		},
		{
			node:     NewTextNode("<b>text</b>"),
			expected: "&lt;b&gt;text&lt;/b&gt;",
		},
		{
			node:     NewDoctypeNode(),
			expected: "<!DOCTYPE html>",
		},
		{
			node: &Node{
				Type: DocumentNode,
				Children: []*Node{
					&Node{
						Type: DoctypeNode,
					},
					&Node{
						Type:       ElementNode,
						Data:       "html",
						Attributes: []Attribute{},
						Children:   []*Node{},
					},
				},
			},
			expected: "<!DOCTYPE html><html />",
		},
		{
			node:     NewElementBuilder("strong").Build(nil),
			expected: "<strong />",
		},
		{
			node:     NewElementBuilder("strong").WithAttribute("class", "a").WithAttribute("id", "test").Build(nil),
			expected: "<strong class=\"a\" id=\"test\" />",
		},
		{
			node:     NewElementBuilder("strong").WithChild(NewRawNode("<span>a</span>")).WithChild(NewElementBuilder("a").Build(nil)).Build(nil),
			expected: "<strong><span>a</span><a /></strong>",
		},
	}

	for _, test := range testCase {
		buf := &bytes.Buffer{}
		assert.NoError(t, Render(buf, test.node))
		assert.Equal(t, test.expected, buf.String())
	}
}
