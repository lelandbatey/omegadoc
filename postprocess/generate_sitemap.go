package postprocess

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/lelandbatey/omegadoc/domain"
)

func init() {
	RegisterPostprocessor(GenerateSiteMap{rank: 35})
}

type GenerateSiteMap struct {
	rank int
}

func (gsm GenerateSiteMap) Rank() int {
	return gsm.rank
}

func (gsm GenerateSiteMap) Name() string {
	return "GenerateSiteMap"
}
func (gsm GenerateSiteMap) Description() string {
	doc := `#!/usr/bin/env omegadoc <<DELIMIDENT omegadoc/postprocessors/generate_sitemap.md
GenerateSiteMap generates a page which links to all OmegaDocs. If there's no
toplevel 'index.md' document defined then that 'index.md' will be created blank
by this postprocessor. Then 'index.md' will have the sitemap appended to it.
DELIMIDENT`
	lines := strings.Split(doc, "\n")
	// Trim off the in-band beginning and end of this OmegaDoc.
	return strings.Join(lines[1:len(lines)-1], "\n")
}

func (gsm GenerateSiteMap) Postprocess(odocs []domain.OmegaDoc) ([]domain.OmegaDoc, error) {
	root := &node{}
	for _, d := range odocs {
		root.Children = AddToTree(root.Children, strings.Split(d.DestFilePath, "/"))
	}
	sort.SliceStable(root.Children, func(i, j int) bool {
		return root.Children[i].Name < root.Children[j].Name
	})
	buf := bytes.NewBuffer(nil)
	for _, c := range root.Children {
		writeMDSiteMap(buf, c, nil)
	}
	var index *domain.OmegaDoc
	for idx := range odocs {
		x := &odocs[idx]
		if x.DestFilePath == "index.md" {
			index = x
			break
		}
	}
	if index == nil {
		index = &domain.OmegaDoc{
			DestFilePath: "index.md",
			Contents:     "",
		}
		odocs = append(odocs, *index)
		index = &odocs[len(odocs)-1]
	}
	index.Contents = fmt.Sprintf("%s\n# Sitemap\n\n%s\n", index.Contents, string(buf.Bytes()))
	return odocs, nil
}

type node struct {
	Name     string
	Children []*node
}

func AddToTree(root []*node, names []string) []*node {
	if len(names) > 0 {
		var i int
		for i = 0; i < len(root); i++ {
			if root[i].Name == names[0] { //already in tree
				break
			}
		}
		if i == len(root) {
			root = append(root, &node{Name: names[0]})
		}
		root[i].Children = AddToTree(root[i].Children, names[1:])
	}
	return root
}

func writeMDSiteMap(w io.Writer, n *node, pieces []string) error {
	pieces = append(pieces, n.Name)
	depth := len(pieces) - 1
	if len(n.Children) == 0 {
		for x := 0; x < depth; x++ {
			fmt.Fprintf(w, "\t")
		}
		fmt.Fprintf(w, "- [%s](%s)\n", n.Name, strings.Join(pieces, "/"))
		return nil
	}
	btn := func(x bool) int {
		if x {
			return 1
		}
		return 0
	}
	sort.SliceStable(n.Children, func(i, j int) bool {
		a := n.Children[i]
		b := n.Children[j]
		acmp := fmt.Sprintf("%d%d%s", btn(len(a.Children) > 0), 1-btn(strings.HasPrefix(a.Name, "index")), a.Name)
		bcmp := fmt.Sprintf("%d%d%s", btn(len(b.Children) > 0), 1-btn(strings.HasPrefix(b.Name, "index")), b.Name)
		return acmp < bcmp
	})
	for x := 0; x < depth; x++ {
		fmt.Fprintf(w, "\t")
	}
	fmt.Fprintf(w, "- [%s/](%s/)\n", n.Name, strings.Join(pieces, "/"))
	for _, c := range n.Children {
		tmp := make([]string, len(pieces))
		copy(tmp, pieces)
		writeMDSiteMap(w, c, tmp)
	}
	return nil
}
