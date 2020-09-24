package htmlkit

type Attribute struct {
	Key   string
	Value string
}

type NodeType uint32

const (
	ErrorNode NodeType = iota
	TextNode
	DocumentNode
	ElementNode
	CommentNode
	DoctypeNode
	RawNode
)

type Node struct {
	Type       NodeType
	Data       string
	Children   []*Node
	Attributes []Attribute
}

func (n *Node) AppendChild(children ...*Node) {
	n.Children = append(n.Children, children...)
}

func NewElementNode(name string, attributes []Attribute, children []*Node) *Node {
	node := &Node{
		Type:       ElementNode,
		Data:       name,
		Children:   []*Node{},
		Attributes: []Attribute{},
	}

	if attributes != nil {
		node.Attributes = append(node.Attributes, attributes...)
	}

	if children != nil {
		node.Children = append(node.Children, children...)
	}

	return node
}

func NewDoctypeNode() *Node {
	return &Node{
		Type: DoctypeNode,
	}
}

func NewDocumentNode() *Node {
	return &Node{
		Type:     DocumentNode,
		Children: []*Node{},
	}
}

func NewTextNode(text string) *Node {
	return &Node{
		Type: TextNode,
		Data: text,
	}
}

func NewRawNode(value string) *Node {
	return &Node{
		Type: RawNode,
		Data: value,
	}
}

func NewCommentNode(comment string) *Node {
	return &Node{
		Type: CommentNode,
		Data: comment,
	}
}
