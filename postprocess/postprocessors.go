package postprocess

import (
	"sort"

	"github.com/lelandbatey/omegadoc/domain"
)

var pprocessors []domain.Postprocessor

func RegisterPostprocessor(p domain.Postprocessor) {
	pprocessors = append(pprocessors, p)
	sort.SliceStable(pprocessors, func(i, j int) bool {
		return pprocessors[i].Rank() < pprocessors[i].Rank()
	})
}

// GetPostprocessors returns the set of globally registered Postprocessors,
// sorted in ascending order by Rank. Postprocessors with a lower Rank() value
// will be earlier in the slice.
func GetPostprocessors() []domain.Postprocessor {
	tmp := make([]domain.Postprocessor, len(pprocessors))
	copy(tmp, pprocessors)
	sort.SliceStable(tmp, func(i, j int) bool {
		return tmp[i].Rank() < tmp[i].Rank()
	})
	return tmp
}

/*
#!/usr/bin/env omegadoc <<DELIMIDENT omegadoc/postprocessors/index.md
# Postprocessors

Postprocessors act as a kind of "super-middleware", an interface for things
which need to accept OmegaDocs and be able to make arbitrary modifications to
those OmegaDocs. Currently there are several built-in postprocessors which are
turned on by default; they are:

- [MarkdownLinkRewriter](omegadoc/postprocessors/rewrite_mdlinks.md) changes links so they point to .html files instead of .md files.
- [SourceLinkAdder](omegadoc/postprocessors/add_sourcelinks.md) changes links so they point to .html files instead of .md files.
- [SectionsCompiler](omegadoc/postprocessors/compile_sections.md) coallesces OmegaDocs which define separate sections/parts of the same file into a single file


DELIMIDENT
*/
