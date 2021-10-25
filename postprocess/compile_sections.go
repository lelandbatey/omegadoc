package postprocess

import (
	"sort"
	"strings"

	"github.com/lelandbatey/omegadoc/domain"
)

func init() {
	RegisterPostprocessor(SectionsCompiler{rank: 60})
}

// SectionsCompiler compiles OmegaDocs of the same destination but with
// different "section" attributes into a single OmegaDoc with the contents
// concatenated in the order of the "section" values.
type SectionsCompiler struct {
	rank int
}

func (sc SectionsCompiler) Rank() int {
	return sc.rank
}

func (sc SectionsCompiler) Name() string {
	return "SectionsCompiler"
}
func (sc SectionsCompiler) Description() string {
	doc := `#!/usr/bin/env omegadoc <<DELIMIDENT section:part01-overview omegadoc/postprocessors/compile_sections.md
SectionsCompiler compiles OmegaDocs of the same destination but with different
"section" attributes into a single OmegaDoc with the contents concatenated in
the order of the "section" values. This is mostly useful if we want to define
small sections of one logical output document in several different places. For
example:

	#!/usr/bin/env omegadoc <<EOF section:part01-overview tmp/examples/large_document.md
	# Introduction
	Let's assume that this is the first section (part 01) of a larger document;
	let's imagine it's the into to a single document covering a detailed piece
	of software.  Since this is the really "general" stuff, let's say we put
	this in the "main" file, the entrypoint of our program. And since we want
	it to be the first section of the document, we give it the name "part01" so
	that it'll come before other parts.
	EOF


	#!/usr/bin/env omegadoc <<EOF section:part04-detailed-inner-workings tmp/examples/large_document.md
	# Detailed inner workings
	This document has a section of "part04", so we want this to come after the
	intro, maybe to describe in detail some smaller set of inner workings of
	the software. Maybe we put this OmegaDoc inside the code for some module
	within the program.
	EOF

	#!/usr/bin/env omegadoc <<EOF section:part09-cleanup tmp/examples/large_document.md
	# Cleanup and Conclusion
	This OmegaDoc has a section of "part09", so we want it to come at the end.
	Let's imagine that this is discussing the final pieces of operation of a
	piece of software and the cleanup which that software does.
	EOF

These three separate OmegaDocs could be defined anywhere; each "section" could
be located as close to the relevant implementation as necessary. Then this
SectionsCompiler postprocessor will collect and concatenate these separate
OmegaDocs into a single output OmegaDoc which would look like the following:

	#!/usr/bin/env omegadoc <<EOF tmp/examples/large_document.md
	# Introduction
	Let's assume that this is the first section (part 01) of a larger document;
	let's imagine it's the into to a single document covering a detailed piece
	of software.  Since this is the really "general" stuff, let's say we put
	this in the "main" file, the entrypoint of our program. And since we want
	it to be the first section of the document, we give it the name "part01" so
	that it'll come before other parts.

	# Detailed inner workings
	This document has a section of "part04", so we want this to come after the
	intro, maybe to describe in detail some smaller set of inner workings of
	the software. Maybe we put this OmegaDoc inside the code for some module
	within the program.

	# Cleanup and Conclusion
	This OmegaDoc has a section of "part09", so we want it to come at the end.
	Let's imagine that this is discussing the final pieces of operation of a
	piece of software and the cleanup which that software does.
	EOF
DELIMIDENT`
	lines := strings.Split(doc, "\n")
	// Trim off the in-band beginning and end of this OmegaDoc.
	return strings.Join(lines[1:len(lines)-1], "\n")
}

func getattr(odoc domain.OmegaDoc, key, defval string) string {
	for _, attr := range odoc.Attributes {
		if strings.ToLower(attr.Key) == key {
			return attr.Value
		}
	}
	return defval
}

func (sc SectionsCompiler) Postprocess(odocs []domain.OmegaDoc) ([]domain.OmegaDoc, error) {
	newdocs := []domain.OmegaDoc{}
	existing := map[string][]domain.OmegaDoc{}
	for _, odoc := range odocs {
		if _, ok := existing[odoc.DestFilePath]; !ok {
			existing[odoc.DestFilePath] = []domain.OmegaDoc{}
		}
		existing[odoc.DestFilePath] = append(existing[odoc.DestFilePath], odoc)
	}

	for _, v := range existing {
		if len(v) < 2 {
			newdocs = append(newdocs, v...)
			continue
		}
		withsections := []domain.OmegaDoc{}
		lacksections := []domain.OmegaDoc{}
		for _, d := range v {
			if getattr(d, "section", "") != "" {
				withsections = append(withsections, d)
				continue
			}
			lacksections = append(lacksections, d)
		}
		newdocs = append(newdocs, lacksections...)

		sort.SliceStable(withsections, func(i, j int) bool {
			return getattr(withsections[i], "section", "9999") < getattr(withsections[j], "section", "9999")
		})
		nd := domain.OmegaDoc{}
		for _, d := range withsections {
			nd.DestFilePath = d.DestFilePath
			nd.Attributes = append(nd.Attributes, d.Attributes...)
			nd.Contents += d.Contents

		}
		newdocs = append(newdocs, nd)
	}
	return newdocs, nil
}
