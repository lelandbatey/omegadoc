package postprocess

import (
	"bytes"
	"io"
	"path/filepath"
	"strings"

	"github.com/lelandbatey/omegadoc/domain"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"

	"github.com/Kunde21/markdownfmt/v2/markdown"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"

	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterPostprocessor(MarkdownLinkRewriter{rank: 40})
}

type MarkdownLinkRewriter struct {
	rank int
}

func (mlr MarkdownLinkRewriter) Rank() int {
	return mlr.rank
}

func (mlr MarkdownLinkRewriter) Name() string {
	return "MarkdownLinkRewriter"
}
func (mlr MarkdownLinkRewriter) Description() string {
	doc := `#!/usr/bin/env omegadoc <<DELIMIDENT omegadoc/postprocessors/rewrite_mdlinks.md
MarkdownLinkRewriter changes all internal links which point to
'.md' paths so that instead they go to '.html'. External links to outside HTTP
pages will not be modified. The following are some examples of how a markdown
link might be written in original Markdown and how it may be changed by this
postprocessor:

    | Original markdown link | Markdown link after modification |
    |------------------------|----------------------------------|
    |[link](thing.md)        | [link](thing.html)               |
    |[link](stuff/thing.md)  | [link](stuff/thing.html)         |
    |[word](other.html)      | [word](other.html)               | # No change because original not linking to markdown file
    |[word](path/other.html) | [word](path/other.html)          | # No change because original not linking to markdown file
DELIMIDENT`
	lines := strings.Split(doc, "\n")
	// Trim off the in-band beginning and end of this OmegaDoc.
	return strings.Join(lines[1:len(lines)-1], "\n")
}

func (mlr MarkdownLinkRewriter) Postprocess(odocs []domain.OmegaDoc) ([]domain.OmegaDoc, error) {
	newdocs := []domain.OmegaDoc{}
	for _, odoc := range odocs {
		if odoc.DestFilePath != "" && !strings.HasSuffix(odoc.DestFilePath, ".md") {
			newdocs = append(newdocs, odoc)
			continue
		}
		// TODO submit a PR fixing the "trailing two spaces on a line" behavior
		// in the Kunde21/markdownfmt renderer, as that's handling things
		// incorrectly. The workaround for now is that I'll have to make sure
		// that the HTML is always rendered with the "HardWraps" option while
		// this markdown renderer _avoids_ using the "SoftWraps" option.
		fr := fakerender{
			realr: markdown.NewRenderer(),
			odoc:  odoc,
		}
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
	odoc  domain.OmegaDoc
}

func (f fakerender) Render(w io.Writer, source []byte, node ast.Node) error {
	// kickoff replacing all links with Links that have their '.md' rewritten to be '.html'
	//node.Dump(source, 0)
	err := ast.Walk(node, f.replaceLinks)
	if err != nil {
		return err
	}
	return f.realr.Render(w, source, node)
}

func (f fakerender) AddOptions(opt ...renderer.Option) {
	f.realr.AddOptions(opt...)
}

func (f fakerender) replaceLinks(node ast.Node, entering bool) (ast.WalkStatus, error) {
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
		if !strings.HasSuffix(dest, ".md") {
			continue
		}
		// Because of how web-servers and browsers handle clickable HREFs, in
		// cases when you don't know where exactly the root of the tree will
		// be, the most durable way to link to another file in the tree is to
		// use relative links. Thus, we rewrite all internal markdown links as
		// relative.
		l := log.WithFields(log.Fields{
			"doc_dest": f.odoc.DestFilePath,
			"link_url": dest,
		})
		ndest, err := filepath.Rel(filepath.Dir("/"+f.odoc.DestFilePath), "/"+dest)
		if err != nil {
			ndest = dest
		}
		l.Infof("new dest: %s", ndest)
		ndest = strings.TrimSuffix(ndest, ".md")
		ndest = ndest + ".html"
		newchild := ast.Link(*n)
		newchild.Destination = []byte(ndest)
		node.ReplaceChild(node, c, &newchild)
	}

	return ast.WalkContinue, nil
}
