package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/lelandbatey/omegadoc/application"
	"github.com/lelandbatey/omegadoc/docfinder"
	"github.com/lelandbatey/omegadoc/docparser"
	"github.com/lelandbatey/omegadoc/docplacer"
	"github.com/lelandbatey/omegadoc/postprocess"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	defaultOmegadocOut = path.Join(os.TempDir(), "omegadoc")
	outputpath         = pflag.StringP("output-path", "o", "", "Path to the directory in which to collect all found OmegaDocs")
	scanpath           = pflag.StringP("input-search-path", "i", "", "Path to the file or directory to search for OmegaDocs")
	helpFlag           = pflag.BoolP("help", "h", false, "Print usage")
	binName            = filepath.Base(os.Args[0])
	longDesc           = `OmegaDoc provides one solution to the documentation problems even medium-size
organizations face. We'd like to keep documentation located nearby the things
they're documenting, but doing that means actually finding and reading that
documentation requires going to where it's located across potentially many
codebases. OmegaDoc is meant to solve this by bringing together documentation
from anywhere, text files of any type, into a single collected directory.

An "OmegaDoc" is a series of bytes which can be recognized by the OmegaDoc
program and extracted into a separate file. In concept it behaves like a
specialized "here document"; an OmegaDoc is meant to be defined in-band with
code and configuration, so that it's nearby the things it documents.

Specification
-------------

An OmegaDoc is composed of an opening statement, then an output path, then a
body, and then a delimiting identifier, in that order. The opening statement is
a "magic string" followed by a delimiting identifier (e.g. the word EOF or
END), then some amount of whitespace, then an output path. The output path is
followed, starting on the next line, by the text to be quoted, and then closed
by the same delimiting identifier on its own line. The "magic string" which
marks the beginning of an OmegaDoc is the string:

	#!/usr/bin/` + `env omegadoc <<

An example of an OmegaDoc then is like so:

	#!/usr/bin/` + `env omegadoc <<DELIMIDENT exampleoutput/readme.md
	Hello I am a markdown document which will be recorded to
	a file at the relative path exampleoutput/readme.md.
	DELIMIDENT

Additionally, if the file ends before the delimiting identifier is reached,
that is considered to be the end of the omegadoc.

If a file contains an "ignore" directive in its bytes before an OmegaDoc
opening statement, then that file will be considered to have NO OmegaDocs in
it, even if it otherwise contains one or more valid OmegaDocs. If a file
contains an "ignore" directive after one or more valid OmegaDoc directives,
then that "ignore" directive will itself be ignored.

An "ignore" directive is the following string, present anywhere in the bytes of
a file:

    #!/usr/bin/env omegadoc ignore-this-file
`
)

func init() {
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nUsage: %s --input-search-path SEARCHPATH --output-path OUTPUTPATH\n", binName)
		fmt.Fprintf(os.Stderr, "\nA documentation extraction and collection program.\n")
		fmt.Fprintf(os.Stderr, "\n%s", longDesc)
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		pflag.PrintDefaults()
	}
}

func main() {
	pflag.Parse()

	if *helpFlag {
		pflag.Usage()
		os.Exit(0)
	}

	if *scanpath == "" && *outputpath == "" {
		fmt.Fprintf(os.Stderr, "\nError: you must provide at least one of --input-search-path or --output-path\n")
		pflag.Usage()
		os.Exit(0)
	}

	if *scanpath == "" {
		*scanpath = "./"
		log.Infof("--input-search-path not provided, defaulting to %s", *scanpath)
	}
	if *outputpath == "" {
		*outputpath = defaultOmegadocOut
		log.Infof("--output-path not provided, defaulting to %s", *outputpath)
	}

	inppath, err := filepath.Abs(*scanpath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	outpath, err := filepath.Abs(*outputpath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	log.SetLevel(log.DebugLevel)

	docfndr := docfinder.NewDocFinder()
	docprsr := docparser.NewDocParser()
	docplcr := docplacer.NewDocPlacer()
	odcc := application.NewController(
		docfndr,
		docprsr,
		postprocess.GetPostprocessors(),
		docplcr,
	)

	err = odcc.GenerateOmegaTree(inppath, outpath)
	if err != nil {
		log.Errorf("Error encountered while attempting to generate OmegaDocs: %q", err.Error())
		os.Exit(1)
	}
}
