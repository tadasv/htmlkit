package htmlkit

import (
	"fmt"
	"html"
	"io"
)

func Render(writer io.Writer, nodes ...*Node) error {
	return NewRenderer(writer).Render(nodes...)
}

type Renderer struct {
	writer io.Writer
}

func NewRenderer(w io.Writer) *Renderer {
	return &Renderer{
		writer: w,
	}
}

func (r *Renderer) Render(nodes ...*Node) error {
	for _, node := range nodes {
		switch node.Type {
		case DocumentNode:
			if err := r.Render(node.Children...); err != nil {
				return err
			}
		case ElementNode:
			if err := r.renderElement(node); err != nil {
				return err
			}
		case TextNode:
			if _, err := r.writer.Write([]byte(html.EscapeString(node.Data))); err != nil {
				return err
			}
		case CommentNode:
			if _, err := fmt.Fprintf(r.writer, "<!-- %s -->", node.Data); err != nil {
				return err
			}
		case DoctypeNode:
			if _, err := fmt.Fprintf(r.writer, "<!DOCTYPE html>"); err != nil {
				return err
			}
		case RawNode:
			if _, err := r.writer.Write([]byte(node.Data)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Renderer) renderElement(node *Node) error {
	fmt.Fprintf(r.writer, "<%s", node.Data)

	attribLen := len(node.Attributes)
	for i, attrib := range node.Attributes {
		if i == 0 {
			fmt.Fprintf(r.writer, " ")
		}

		if i == attribLen-1 {
			fmt.Fprintf(r.writer, "%s=\"%s\"", attrib.Key, attrib.Value)
		} else {
			fmt.Fprintf(r.writer, "%s=\"%s\" ", attrib.Key, attrib.Value)
		}
	}

	if len(node.Children) > 0 {
		if _, err := fmt.Fprintf(r.writer, ">"); err != nil {
			return err
		}

		if err := r.Render(node.Children...); err != nil {
			return err
		}

		fmt.Fprintf(r.writer, "</%s>", node.Data)
	} else {
		fmt.Fprintf(r.writer, " />")
	}

	return nil
}
