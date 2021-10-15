// noops contains non-functional/NOn-OPerational implementations of the
// interfaces defined in the domain.
package noops

import (
	"github.com/lelandbatey/omegadoc/domain"
	"io"
)

type NoOpDocFinder struct{}

func (ndf NoOpDocFinder) FindReaders(path string) (map[string]io.Reader, error) {
	return nil, nil
}

var _ domain.DocFinder = NoOpDocFinder{}

type NoOpDocExtractor struct{}

func (nde NoOpDocExtractor) ExtractDocs(reader io.Reader) ([]string, error) {
	return nil, nil
}

var _ domain.DocExtractor = NoOpDocExtractor{}

type NoOpDocParser struct{}

func (nde NoOpDocParser) ParseDoc(srcpath, contents string) (domain.OmegaDoc, error) {
	return domain.OmegaDoc{}, nil
}

var _ domain.DocParser = NoOpDocParser{}

type NoOpDocPlacer struct{}

func (nde NoOpDocPlacer) PlaceDoc(outpath string, odoc domain.OmegaDoc) error {
	return nil
}

var _ domain.DocPlacer = NoOpDocPlacer{}
