package htmlkit

import (
	"context"
)

type AttributeFunc func(context.Context) (*Attribute, error)
type NodeFunc func(context.Context) ([]*Node, error)

type ElementBuilder struct {
	name           string
	attributes     []Attribute
	attributeFuncs []AttributeFunc
	childrenFuncs  []NodeFunc
}

func NewElementBuilder(name string) *ElementBuilder {
	return &ElementBuilder{
		name:           name,
		attributes:     []Attribute{},
		attributeFuncs: []AttributeFunc{},
		childrenFuncs:  []NodeFunc{},
	}
}

func (b *ElementBuilder) Build(ctx context.Context) *Node {
	node := &Node{
		Type:       ElementNode,
		Data:       b.name,
		Attributes: b.attributes,
		Children:   []*Node{},
	}

	for _, fnc := range b.attributeFuncs {
		attr, err := fnc(ctx)
		if err != nil {
			continue
		}

		if attr != nil {
			node.Attributes = append(node.Attributes, *attr)
		}
	}

	for _, fnc := range b.childrenFuncs {
		child, err := fnc(ctx)
		if err != nil {
			continue
		}

		if child != nil && len(child) > 0 {
			node.Children = append(node.Children, child...)
		}
	}

	return node
}

func (b *ElementBuilder) WithAttribute(key, value string) *ElementBuilder {
	b.attributes = append(b.attributes, Attribute{
		Key:   key,
		Value: value,
	})
	return b
}

func (b *ElementBuilder) WithAttributeFunc(fnc AttributeFunc) *ElementBuilder {
	b.attributeFuncs = append(b.attributeFuncs, fnc)
	return b
}

func (b *ElementBuilder) WithChildren(children ...*Node) *ElementBuilder {
	for _, child := range children {
		b.WithChild(child)
	}
	return b
}

func (b *ElementBuilder) WithChild(child *Node) *ElementBuilder {
	if child == nil {
		return b
	}

	b.childrenFuncs = append(b.childrenFuncs, func(ctx context.Context) ([]*Node, error) {
		return []*Node{child}, nil
	})
	return b
}

func (b *ElementBuilder) WithChildFunc(fnc NodeFunc) *ElementBuilder {
	b.childrenFuncs = append(b.childrenFuncs, fnc)
	return b
}
