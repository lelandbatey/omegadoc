package domain

import (
	"io"
)

// DocFinder finds the on-disk files which contain Omegadoc documents and
// returns them as a slice of io.Readers. While the goal is to find actual
// files on the disk, implementations of DocFinder could return whatever
// io.Readers they want, which is useful for internal testing.
type DocFinder interface {
	// FindReaders returns a map of absolute on-disk paths to io.Readers of
	// each file containing at least one Omegadoc in the filesystem at
	// 'path'. The filesystem at 'path' is searched recursively for files
	// containing Omegadocs.
	FindReaders(path string) (map[string]io.Reader, error)
}

// Parses the contents of the file to extract all the OmegaDocs in that file.
type DocParser interface {
	ParseDoc(srcpath string, data io.Reader) ([]OmegaDoc, error)
}

type DocPlacer interface {
	PlaceDoc(outpath string, odoc OmegaDoc) error
}

// Postprocessors act as a kind of "super-middleware", an interface for things
// which need to accept OmegaDocs and be able to make arbitrary modifications
// to those OmegaDocs.
type Postprocessor interface {
	Postprocess([]OmegaDoc) ([]OmegaDoc, error)
	Name() string
	Description() string
	// Rank is a way for a Postprocessor to indicate it's own "relative
	// importance" compared to other Postprocessors. There is no inherant
	// meaning to any number returned by Rank(); it is meant only as a way to
	// sort a collection of Postprocessors so that you know the other in which
	// to run them.
	Rank() int
}
