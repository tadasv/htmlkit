# htmlkit

htmlkit is a library for building HTML documents in Go. It is meant to serve as
an alternative to `html/template` and other template libraries.

## Why this project?

There's something missing in `html/template` package. It is really useful, but
it lacks compile time safety. I find template reusability a bit problematic also.

Ok, can we do better? I don't know. But here is one attempt. It's not
necessarily better, just different.

htmlkit is a very simple library for constructing HTML documents. It let's you
write HTML DOM in more or less declarative way in Go. This is actually really
nice since you get static type checks at build time and you can also easilly
compose larger documents from reusable HTML nodes (think React components).

## Examples

See `examples` folder on how to use it.

The following code:

```go
package main

import (
	"context"
	"fmt"
	"github.com/tadasv/htmlkit"
	"os"
)

// Define a builder for a reusable comment list component
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

// An example illustrating full HTML page template or what it might look like.
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
```

Produces the following HTML document (pretty printed for readability):

```html
<!DOCTYPE html>
<html>
  <head>
    <title>htmlkit demo</title>
  </head>
  <body>
    <div class="post">This is an example of htmlkit</div>
    <ul class="comments">
      <li>it demonstrates declarative html document construction</li>
      <li>htmlkit let&#39;s you easilly mix in raw HTML when needed</li>
      <li>it offers dynamic html document rendering with type safety</li>
    </ul>
  </body>
</html>
```
