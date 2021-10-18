package docfinder

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lelandbatey/omegadoc/domain"
)

type docfinder struct {
	ignorepaths []string
	searchfunc  func(string, ...string) ([]string, error)
}

var _ domain.DocFinder = docfinder{}

func NewDocFinder(ignorepaths ...string) domain.DocFinder {
	// TODO use exec.LookPath to look up 'rg', 'ag', and 'grep' to choose the
	// underlying search program.
	return docfinder{
		ignorepaths: ignorepaths,
		searchfunc:  grepFind,
	}
}

func (df docfinder) FindReaders(path string) (map[string]io.Reader, error) {
	filepaths, err := df.searchfunc(path, df.ignorepaths...)
	if err != nil {
		return nil, err
	}
	var readers map[string]io.Reader = map[string]io.Reader{}
	{
		for _, fp := range filepaths {
			f, err := os.OpenFile(fp, os.O_RDONLY, 0644)
			if err != nil {
				return nil, err
			}
			readers[fp] = f
		}
	}
	return readers, nil
}

// grepFind finds all files recursively in srcpath which contain an OmegaDoc.
// srcpath and ignorepaths are absolute paths. The returned slice of strings is
// absolute paths to files which contain OmegaDoc(s).
func grepFind(srcpath string, ignorepaths ...string) ([]string, error) {
	var matches []string = []string{}
	{
		// TODO ignore the paths passed in ignorepaths. Right now no files are ignored.
		cmd := exec.Command("grep",
			// If 'type' passed to `--binary-files=type` is 'without-match', when grep
			// discovers null input binary data it assumes that the rest of the file
			// does not match; this is equivalent to the -I option.
			// https://www.gnu.org/software/grep/manual/grep.html#index-_002d_002dbinary_002dfiles
			"--binary-files=without-match",
			// Suppress normal output; instead print the name of each input file from
			// which output would normally have been printed. Scanning each input file
			// stops upon first match.
			// https://www.gnu.org/software/grep/manual/grep.html#index-_002dl
			"--files-with-matches",
			"-r", srcpath, domain.START_OMEGADOC)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return nil, err
		}
		err = cmd.Run()
		if err != nil {
			return nil, err
		}
		all, err := io.ReadAll(stdout)
		if err != nil {
			return nil, err
		}
		for _, line := range strings.Split(string(all), "\n") {
			foundp, err := filepath.Abs(line)
			if err != nil {
				return nil, err
			}
			matches = append(matches, foundp)
		}
		return matches, nil
	}
}
