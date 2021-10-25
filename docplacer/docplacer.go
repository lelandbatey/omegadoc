package docplacer

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/lelandbatey/omegadoc/domain"

	log "github.com/sirupsen/logrus"
)

func NewDocPlacer() domain.DocPlacer {
	return docPlacer{}
}

type docPlacer struct {
	// do-not-overwrite (default) Will not overwrite existing file, returns error
	// ignore           Will not overwrite existing file, will not return error
	// yes-overwrite    Existing files will be overwritten
	handle_existing string
}

func (dpl docPlacer) PlaceDoc(outpath string, odoc domain.OmegaDoc) error {
	if dpl.handle_existing == "" {
		dpl.handle_existing = "do-not-overwrite"
	}
	reltpath := odoc.DestFilePath
	if reltpath[0] == '/' {
		reltpath = strings.TrimPrefix(reltpath, "/")
	}
	finpath := path.Join(outpath, reltpath)
	finbase := path.Dir(finpath)
	log.WithFields(log.Fields{
		"finpath": finpath,
		"finbase": finbase,
	}).Info("writing omegadoc to output")
	err := os.MkdirAll(finbase, 0775)
	if err != nil {
		return fmt.Errorf("cannot create parent directories for '%q': %w", finbase, err)
	}

	_, err = os.Stat(finpath)
	if err == nil {
		if dpl.handle_existing == "do-not-overwrite" {
			return fmt.Errorf("cannot overwrite existing file at location %q; file %q already exists and this program is configured not to overwrite existing files", finpath, finpath)
		} else if dpl.handle_existing == "ignore" {
			return nil
		} else if dpl.handle_existing == "yes-overwrite" {
			// Do nothing and proceed
		} else {
			return fmt.Errorf("unknown handle_existing value of %q, don't know how to proceed; exiting", dpl.handle_existing)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	err = os.WriteFile(finpath, []byte(odoc.Contents), 0644)
	if err != nil {
		return err
	}
	return nil
}
