package htmlkit

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestElementBuilder(t *testing.T) {
	testCases := []struct {
		builder  *ElementBuilder
		context  context.Context
		expected *Node
	}{
		{
			builder: NewElementBuilder("b"),
			expected: &Node{
				Data:       "b",
				Type:       ElementNode,
				Children:   []*Node{},
				Attributes: []Attribute{},
			},
			context: nil,
		},
		{
			builder: NewElementBuilder("parent").WithChild(
				NewElementBuilder("child1").Build(nil),
			).WithChild(
				nil,
			).WithChildFunc(func(ctx context.Context) ([]*Node, error) {
				return nil, nil
			}).WithChildFunc(func(ctx context.Context) ([]*Node, error) {
				return nil, errors.New("will be ignored")
			}).WithChildFunc(func(ctx context.Context) ([]*Node, error) {
				value := ctx.Value("key").(string)
				return []*Node{NewElementBuilder(value).Build(nil)}, nil
			}),
			context: context.WithValue(context.Background(), "key", "child2"),
			expected: &Node{
				Data: "parent",
				Type: ElementNode,
				Children: []*Node{
					&Node{
						Data:       "child1",
						Type:       ElementNode,
						Attributes: []Attribute{},
						Children:   []*Node{},
					},
					&Node{
						Data:       "child2",
						Type:       ElementNode,
						Attributes: []Attribute{},
						Children:   []*Node{},
					},
				},
				Attributes: []Attribute{},
			},
		},
		{
			builder: NewElementBuilder("b").WithAttribute(
				"a", "b",
			).WithAttributeFunc(func(ctx context.Context) (*Attribute, error) {
				val := ctx.Value("key").(string)
				return &Attribute{
					Key:   "key",
					Value: val,
				}, nil
			}).WithAttributeFunc(func(ctx context.Context) (*Attribute, error) {
				// nulls will be skipped
				return nil, nil
			}).WithAttributeFunc(func(ctx context.Context) (*Attribute, error) {
				return nil, errors.New("errors will be skipped")
			}),
			expected: &Node{
				Data:     "b",
				Type:     ElementNode,
				Children: []*Node{},
				Attributes: []Attribute{
					{
						Key:   "a",
						Value: "b",
					},
					{
						Key:   "key",
						Value: "value",
					},
				},
			},
			context: context.WithValue(context.Background(), "key", "value"),
		},
	}

	for _, testCase := range testCases {
		node := testCase.builder.Build(testCase.context)
		assert.Equal(t, testCase.expected, node)
	}
}
