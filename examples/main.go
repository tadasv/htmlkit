package main

import (
	"context"
	"fmt"
	"github.com/tadasv/htmlkit"
	"os"
)

type Item struct {
	Name string
	URL  string
}

var commentListBuilder = htmlkit.NewElementBuilder(
	"ul",
).WithAttribute("class", "comments").WithChildFunc(func(ctx context.Context) ([]*htmlkit.Node, error) {
	comments := ctx.Value("comments").([]string)
	children := []*htmlkit.Node{}

	for _, comment := range comments {
		children = append(children, htmlkit.NewElementBuilder("li").WithChild(htmlkit.NewTextNode(comment)).Build(nil))
	}

	return children, nil
})

func buildComments(comments []string) *htmlkit.Node {
	return nil
}

func blogPostPageTemplate(title string, body string, comments []string) *htmlkit.Node {
	templateContext := context.WithValue(context.Background(), "comments", comments)

	document := htmlkit.NewDocumentNode()
	document.AppendChild(
		htmlkit.NewDoctypeNode(),
		htmlkit.NewElementBuilder(
			"html",
		).WithChildren(
			htmlkit.NewRawNode(fmt.Sprintf("<head><title>%s</title></head>", title)),
			htmlkit.NewRawNode(`<body>`),
			htmlkit.NewElementBuilder("div").WithAttribute("class", "post").WithChild(htmlkit.NewRawNode(body)).Build(templateContext),
			commentListBuilder.Build(templateContext),
			htmlkit.NewRawNode(`</body>`),
		).Build(templateContext),
	)
	return document
}

func main() {
	document := blogPostPageTemplate("htmlkit demo", "This is an example of htmlkit", []string{
		"it demonstrates declarative html document construction",
		"htmlkit let's you easilly mix in raw HTML when needed",
		"it offers dynamic html document rendering with type safety",
	})

	htmlkit.Render(os.Stdout, document)
}
