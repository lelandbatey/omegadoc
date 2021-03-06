package postprocess

import (
	"fmt"
	"strings"

	"github.com/lelandbatey/omegadoc/domain"
)

func init() {
	RegisterPostprocessor(SourceLinkAdder{rank: 50})
}

type SourceLinkAdder struct {
	rank int
}

func (sla SourceLinkAdder) Rank() int {
	return sla.rank
}

func (sla SourceLinkAdder) Name() string {
	return "SourceLinkAdder"
}
func (sla SourceLinkAdder) Description() string {
	doc := `#!/usr/bin/env omegadoc <<DELIMIDENT omegadoc/postprocessors/add_sourcelinks.md
SourceLinkAdder assumes that each OmegaDoc is written in markdown
format. To each OmegaDoc, it adds an HTTP link to the original source file
which contains the Omegadoc. That way, if when reading the Omegadoc you wonder
"where is this written? I want to add something to it..." then you may visit
the original file directly in your web-browser with just a click.

For example, let's say that in the 'main.go' file at the root of this
repository, I were to declare an omegadoc like the following:

	#!/usr/bin/env omegadoc <<END exampleoutput.md
	Hello, I should have a URL back to this source document on the line below
	this one, just down here ↓↓↓↓↓↓↓↓↓
	END

If I were to then invoke Omegadoc and instruct it to use this postprocessor
step, then when the 'exampleoutput.md' file was placed in the final location,
it'd have the following modification:

	Hello, I should have a URL back to this source document on the line below
	this one, just down here ↓↓↓↓↓↓↓↓↓

	[Link to this original document: https://github.com/lelandbatey/omegadoc/tree/32bb1a36ee9bd0a5437eb952d6e4cab09125ca47/main.go#L30](https://github.com/lelandbatey/omegadoc/tree/32bb1a36ee9bd0a5437eb952d6e4cab09125ca47/main.go#L30)

This is a very useful addition when your documentation is spread widely across
a large file structure crossing many repositories, which is the exact case
where Omegadoc is meant to be used. By having a link to the original location,
you can _much_ more easily figure out where you need to go in order to make
necessary changes to documentation.
DELIMIDENT
`
	lines := strings.Split(doc, "\n")
	// Trim off the in-band beginning and end of this OmegaDoc.
	return strings.Join(lines[1:len(lines)-1], "\n")
}

func (sla SourceLinkAdder) Postprocess(odocs []domain.OmegaDoc) ([]domain.OmegaDoc, error) {
	newdocs := []domain.OmegaDoc{}
	for _, odoc := range odocs {
		nodoc := domain.OmegaDoc(odoc)
		if nodoc.HTTPUrl != "" {
			nodoc.Contents = odoc.Contents + fmt.Sprintf("\n\n[Link to the source of this document: %s](%s)\n", nodoc.HTTPUrl, nodoc.HTTPUrl)
		}
		newdocs = append(newdocs, nodoc)
	}
	return newdocs, nil
}
