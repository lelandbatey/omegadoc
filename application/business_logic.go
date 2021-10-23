package application

import (
	"fmt"
	"strings"

	"github.com/lelandbatey/omegadoc/domain"

	log "github.com/sirupsen/logrus"
)

type OmegaDocController struct {
	finder domain.DocFinder
	parser domain.DocParser
	pprocs []domain.Postprocessor
	placer domain.DocPlacer
}

func NewController(
	finder domain.DocFinder,
	parser domain.DocParser,
	pprocs []domain.Postprocessor,
	placer domain.DocPlacer) OmegaDocController {
	return OmegaDocController{
		finder: finder,
		parser: parser,
		pprocs: pprocs,
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
		log.Infof("Found reader: %s", srcpath)
	}

	odocs := []domain.OmegaDoc{}
	for srcpath, rdr := range readers {
		newodocs, err := odcc.parser.ParseDoc(srcpath, rdr)
		if err != nil {
			return err
		}
		odocs = append(odocs, newodocs...)
	}

	for _, pproc := range odcc.pprocs {
		odocs, err = pproc.Postprocess(odocs)
		if err != nil {
			return err
		}
	}

	if len(readers) > len(odocs) {
		skipped := len(readers) - len(odocs)
		log.Infof("Some files with potential OmegaDocs in them were ignored, count of ignored: %d, count of files with potential OmegaDocs: %d", skipped, len(readers))
	}
	odocs = append(odocs, MakeIndex(odocs))

	for _, odoc := range odocs {
		err := odcc.placer.PlaceDoc(outpath, odoc)
		if err != nil {
			return err
		}
	}
	return nil
}

func MakeIndex(odocs []domain.OmegaDoc) domain.OmegaDoc {
	odoc := domain.OmegaDoc{
		DestFilePath: "index.md",
	}
	s := ""
	for _, x := range odocs {
		s += fmt.Sprintf("[%s](./%s)  \n", x.DestFilePath, strings.ReplaceAll(x.DestFilePath, ".md", ".html"))
	}
	odoc.Contents = s
	return odoc
}
