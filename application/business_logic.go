package application

import (
	"github.com/lelandbatey/omegadoc/domain"

	log "github.com/sirupsen/logrus"
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
	log.Debug("Beginnning operation")
	readers, err := odcc.finder.FindReaders(inpath)
	if err != nil {
		return err
	}
	for srcpath := range readers {
		log.Info("Found reader", srcpath)
	}

	odocs := []domain.OmegaDoc{}
	for srcpath, rdr := range readers {
		newodocs, err := odcc.parser.ParseDoc(srcpath, rdr)
		if err != nil {
			return err
		}
		odocs = append(odocs, newodocs...)
	}

	if len(readers) > len(odocs) {
		skipped := len(readers) - len(odocs)
		log.Infof("Some files with potential OmegaDocs in them were ignored, count of ignored: %d, count of files with potential OmegaDocs: %d", skipped, len(readers))
	}

	for _, odoc := range odocs {
		err := odcc.placer.PlaceDoc(odoc.DestFilePath, odoc)
		if err != nil {
			return err
		}
	}
	return nil
}
