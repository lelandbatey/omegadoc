package application

import (
	"github.com/lelandbatey/omegadoc/domain"
)

type OmegaDocController struct {
	finder domain.DocFinder
	xtrctr domain.DocExtractor
	parser domain.DocParser
	placer domain.DocPlacer
}

func NewController(
	finder domain.DocFinder,
	xtrctr domain.DocExtractor,
	parser domain.DocParser,
	placer domain.DocPlacer) OmegaDocController {
	return OmegaDocController{
		finder: finder,
		xtrctr: xtrctr,
		parser: parser,
		placer: placer,
	}
}

func (odcc OmegaDocController) GenerateOmegaTree(inpath, outpath string) error {
	readers, err := odcc.finder.FindReaders(inpath)
	if err != nil {
		return err
	}

	sodoc := make(map[string][]string)
	for srcpath, rdr := range readers {
		extracts, err := odcc.xtrctr.ExtractDocs(rdr)
		if err != nil {
			return err
		}
		sodoc[srcpath] = extracts
	}

	odocs := []domain.OmegaDoc{}
	for srcpath, rawodocs := range sodoc {
		for _, rawodoc := range rawodocs {
			odoc, err := odcc.parser.ParseDoc(srcpath, rawodoc)
			if err != nil {
				return err
			}
			odocs = append(odocs, odoc)
		}
	}

	for _, odoc := range odocs {
		err := odcc.placer.PlaceDoc(odoc.DestFilePath, odoc)
		if err != nil {
			return err
		}
	}
	return nil
}
