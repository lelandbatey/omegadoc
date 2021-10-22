package postprocess

import (
	"bytes"
	"io"
	"strings"

	"github.com/lelandbatey/omegadoc/domain"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"

	"github.com/Kunde21/markdownfmt/v2/markdown"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
)

func init() {
	RegisterPostprocessor(MarkdownLinkRewriter{})
}

type MarkdownLinkRewriter struct {
}

func (mlr MarkdownLinkRewriter) Name() string {
	return "MarkdownLinkRewriter"
}
func (mlr MarkdownLinkRewriter) Description() string {
	return `MarkdownLinkRewriter changes all internal links which point to
'.md' paths so that instead they go to '.html'. External links to outside HTTP
pages will not be modified.
`
}

func (mlr MarkdownLinkRewriter) Postprocess(odocs []domain.OmegaDoc) ([]domain.OmegaDoc, error) {
	newdocs := []domain.OmegaDoc{}
	for _, odoc := range odocs {
		fr := fakerender{markdown.NewRenderer()}
		g := goldmark.New(
			goldmark.WithExtensions(extension.GFM),
			goldmark.WithRenderer(
				fr,
			),
		)
		nodoc := domain.OmegaDoc(odoc)
		var buf bytes.Buffer
		err := g.Convert([]byte(odoc.Contents), &buf)
		if err != nil {
			return nil, err
		}
		nodoc.Contents = string(buf.Bytes())
		newdocs = append(newdocs, nodoc)
	}
	return newdocs, nil
}

type fakerender struct {
	realr *markdown.Renderer
}

func (f fakerender) Render(w io.Writer, source []byte, node ast.Node) error {
	// kickoff replacing all links with Links that have their '.md' rewritten to be '.html'
	err := ast.Walk(node, replaceLinks)
	if err != nil {
		return err
	}
	return f.realr.Render(w, source, node)
}

func (f fakerender) AddOptions(opt ...renderer.Option) {
	f.realr.AddOptions(opt...)
}

func replaceLinks(node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		return ast.WalkContinue, nil
	}
	if !node.HasChildren() {
		return ast.WalkContinue, nil
	}
	linkChildren := []ast.Node{}
	child := node.FirstChild()
	if child == nil {
		return ast.WalkContinue, nil
	}
	for i := 0; i < node.ChildCount(); i++ {
		if child.Kind() == ast.KindLink {
			linkChildren = append(linkChildren, child)
		}
		child = child.NextSibling()
	}

	for _, c := range linkChildren {
		n := c.(*ast.Link)
		dest := string(n.Destination)
		if strings.HasSuffix(dest, ".md") {
			ndest := strings.TrimSuffix(dest, ".md")
			ndest = ndest + ".html"
			newchild := ast.Link(*n)
			newchild.Destination = []byte(ndest)
			node.ReplaceChild(node, c, &newchild)
		}
	}

	return ast.WalkContinue, nil
}
