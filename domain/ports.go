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

// DocExtractor extracts all the raw OmegaDoc contents from the provided
// io.Reader.
type DocExtractor interface {
	ExtractDocs(reader io.Reader) ([]string, error)
}

// Parses the contents of the OmegaDoc to create the OmegaDoc struct.
type DocParser interface {
	ParseDoc(srcpath, contents string) (OmegaDoc, error)
}

type DocPlacer interface {
	PlaceDoc(outpath string, odoc OmegaDoc) error
}
