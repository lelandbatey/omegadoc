package application

import (
	"github.com/lelandbatey/omegadoc/domain"
)

type OmegaDocController struct {
	finder domain.DocFinder
	parser domain.DocParser
	placer domain.DocPlacer
}

func NewController(
	finder domain.DocFinder,
	parser domain.DocParser,
	placer domain.DocPlacer) OmegaDocController {
	return OmegaDocController{
		finder: finder,
		parser: parser,
		placer: placer,
	}
}

func (odcc OmegaDocController) GenerateOmegaTree(inpath, outpath string) error {
	readers, err := odcc.finder.FindReaders(inpath)
	if err != nil {
		return err
	}

	odocs := []domain.OmegaDoc{}
	for srcpath, rdr := range readers {
		newodocs, err := odcc.parser.ParseDoc(srcpath, rdr)
		if err != nil {
			return err
		}
		odocs = append(odocs, newodocs...)
	}

	for _, odoc := range odocs {
		err := odcc.placer.PlaceDoc(odoc.DestFilePath, odoc)
		if err != nil {
			return err
		}
	}
	return nil
}
