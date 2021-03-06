package docparser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/lelandbatey/omegadoc/domain"

	log "github.com/sirupsen/logrus"
)

// RequiredParseError represents an error which cannot be skipped and which is
// NOT safe to ignore.
type RequiredParseError struct {
	inner error
}

func (e *RequiredParseError) Error() string {
	return fmt.Sprintf("a non-ignorable error occured: %v", e.inner)
}

func (e *RequiredParseError) Unwrap() error {
	return e.inner
}

type docfinder struct {
	urlfinder gitURLFinder
}

func NewDocParser() domain.DocParser {
	return docfinder{
		urlfinder: newGitURLFinder(),
	}
}

type oatt struct {
	Key   string
	Value string
}

// parseOdoc tracks all the data necessary for us to parse the document, and
// when complete may be turned into a "real" OmegaDoc
type parseOdoc struct {
	SourceFilePath  string
	DestFilePath    []rune
	Contents        []rune
	Attrs           []oatt
	StartLineNumber int
	//HTTPUrl string
}

func (po *parseOdoc) AppCont(r ...rune) {
	po.Contents = append(po.Contents, r...)
}

func (po *parseOdoc) AppDestFP(r ...rune) {
	po.DestFilePath = append(po.DestFilePath, r...)
}

func (po *parseOdoc) AppAttr(key, val string) {
	po.Attrs = append(po.Attrs, oatt{Key: key, Value: val})
}

func (po *parseOdoc) MakeOmegaDoc() domain.OmegaDoc {
	attrs := []domain.OmegaAttribute{}
	for _, att := range po.Attrs {
		attrs = append(attrs, domain.OmegaAttribute(att))
	}
	return domain.OmegaDoc{
		SourceFilePath:  po.SourceFilePath,
		DestFilePath:    string(po.DestFilePath),
		Contents:        string(po.Contents),
		Attributes:      attrs,
		StartLineNumber: po.StartLineNumber,
	}
}

func (df docfinder) ParseDoc(srcpath string, data io.Reader) ([]domain.OmegaDoc, error) {
	l := log.WithField("srcpath", srcpath)
	odocs, err := df.parseDoc(srcpath, data)
	if err != nil {
		return nil, err
	}
	newodocs := []domain.OmegaDoc{}
	for _, od := range odocs {
		url, err := df.urlfinder.GetURL(od.SourceFilePath, od.StartLineNumber)
		if err != nil {
			l.Warnf("cannot find URL for document %q: %v", od.SourceFilePath, err)
		} else {
			l.WithField("url", url).Infof("URL for %q found", od.SourceFilePath)
		}
		od.HTTPUrl = url
		newodocs = append(newodocs, od)
	}
	return newodocs, nil
}

const COMMON_PREFIX string = "#!/usr/bin/env omegadoc "

var BEGINDOC_MAGICRUNES []rune = []rune(strings.ReplaceAll(domain.START_OMEGADOC, COMMON_PREFIX, ""))
var IGNORDOC_MAGICRUNES []rune = []rune(strings.ReplaceAll(domain.IGNORE_OMEGADOC, COMMON_PREFIX, ""))

type odScanner struct {
	Rdr    *bufio.Reader
	LineNo int
	Prior  rune
}

func (ods *odScanner) ReadRune() (rune, int, error) {
	r, s, err := ods.Rdr.ReadRune()
	ods.Prior = r
	if r == '\n' {
		ods.LineNo++
	}
	return r, s, err
}

func (ods *odScanner) UnreadRune() error {
	if ods.Prior == '\n' {
		ods.LineNo -= 1
	}
	ods.Prior = rune(0)
	return ods.Rdr.UnreadRune()
}

func (ods *odScanner) LineNumber() int {
	return ods.LineNo
}

// Returns the next group of runes in a logical group. There are three possible
// groups: a group of non-whitespace characters, a single newline, and a group
// of non-newline whitespace characters.
func (ods *odScanner) ReadRuneGroup() ([]rune, error) {
	buf := []rune{}

	ch, _, err := ods.ReadRune()
	if err != nil {
		return nil, err
	}
	buf = append(buf, ch)
	if ch == '\n' {
		// A single newline
		return buf, nil
	} else if unicode.IsSpace(ch) {
		// A group of non-newline whitespace characters
		for {
			ch, _, err = ods.ReadRune()
			if err != nil {
				return buf, err
			}
			if !unicode.IsSpace(ch) || ch == '\n' {
				err = ods.UnreadRune()
				if err != nil {
					return buf, err
				}
				return buf, nil
			} else {
				buf = append(buf, ch)
			}
		}
	} else {
		// A group of non-whitespace characters
		for {
			ch, _, err = ods.ReadRune()
			if err != nil {
				return buf, err
			}
			if unicode.IsSpace(ch) {
				err = ods.UnreadRune()
				if err != nil {
					return buf, err
				}
				return buf, nil
			}
			buf = append(buf, ch)
		}
	}
}

// FFTillMagicCommon moves through the odScanner till the underlying reader is
// just after the string "#!/usr/bin/env omegadoc ", which is the common prefix
// to both "magic strings" of OmegaDoc: the 'ignore directive' and the 'opening
// statement'.
func (ods *odScanner) FFTillMagicCommon() error {
	_, err := readTillSentinel([]rune(COMMON_PREFIX), ods)
	return err
}

func runesEqual(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func readTillSentinel(sentinel []rune, ods *odScanner) ([]rune, error) {
	buf := []rune{}
	pos := 0
	for {
		r, _, err := ods.ReadRune()
		if err != nil {
			return buf, err
		}
		buf = append(buf, r)
		if r == sentinel[pos] {
			pos++
		} else {
			pos = 0
		}

		if pos == len(sentinel) {
			// For compatibility with previous behavior, slice off the sentinel
			// from the end of buf
			return buf[:len(buf)-len(sentinel)], nil
		}
	}
}

// ParseDoc for docfinder parses a text file and extracts all OmegaDocs present
// in the file. This is currently implemented as a simple direct parser,
// without being broken down into scanner/lexer since the language is so
// simple. In the future this implementation may need to be further broken down
// though, as features such as automatic indentation removal or line-prefix
// removal may require a full lexer/parser.
func (df docfinder) parseDoc(srcpath string, data io.Reader) ([]domain.OmegaDoc, error) {
	l := log.WithField("srcpath", srcpath)
	var odocs []domain.OmegaDoc = []domain.OmegaDoc{}
	brdr := bufio.NewReader(data)
	rdr := &odScanner{brdr, 0, 0}

	var curodoc parseOdoc = parseOdoc{
		SourceFilePath: srcpath,
	}

	deriveCorrectExit := func(err error) ([]domain.OmegaDoc, error) {
		// End of file isn't necessarily an error, more a signal that we're
		// done here.
		if errors.Is(err, io.EOF) {
			return odocs, nil
		}
		return nil, err
	}
	rrg := func() ([]rune, string, error) {
		rg, _err := rdr.ReadRuneGroup()
		if _err != nil {
			return nil, "", _err
		}
		return rg, string(rg), nil
	}
	for {
	RESET_CONTINUE:
		err := rdr.FFTillMagicCommon()
		if err != nil {
			return deriveCorrectExit(err)
		}
		rg, err := rdr.ReadRuneGroup()
		if err != nil {
			return deriveCorrectExit(err)
		}
		if runesEqual(rg, IGNORDOC_MAGICRUNES) {
			if len(odocs) == 0 {
				return odocs, nil
			} else {
				goto RESET_CONTINUE
			}
		} else if strings.HasPrefix(string(rg), string(BEGINDOC_MAGICRUNES)) {
			var delimiting_ident []rune = []rune(strings.TrimPrefix(string(rg), string(BEGINDOC_MAGICRUNES)))
			curodoc.StartLineNumber = rdr.LineNumber()
			l = l.WithFields(log.Fields{
				"startline":        rdr.LineNumber(),
				"delimiting_ident": string(delimiting_ident),
			})
			l.Debug("found beginning of document")
			for {
				rg, s, err := rrg()
				if err != nil {
					return deriveCorrectExit(err)
				}
				if rg[0] == '\n' {
					goto RESET_CONTINUE
				}
				if !unicode.IsSpace(rg[0]) {
					if !strings.Contains(s, ":") {
						// Isn't an attribute, must be destination file path.
						curodoc.AppDestFP(rg...)
						l = l.WithFields(log.Fields{
							"dest_file_path": string(curodoc.DestFilePath),
						})
						l.Debug("found destination file path")
						for {
							r, _, err := rdr.ReadRune()
							if err != nil {
								return deriveCorrectExit(err)
							}
							if r == '\n' {
								// Newline marks end of file path and start of contents of the
								// omegadoc. The end of the contents will be marked by the
								// 'delimiting identifier' (or EOF) so read until that's reached.
								contents, err := readTillSentinel(delimiting_ident, rdr)
								curodoc.AppCont(contents...)
								if errors.Is(err, io.EOF) {
									// Ending the file in the middle of an OmegaDoc is considered a
									// valid ending to the OmegaDoc.
									odocs = append(odocs, curodoc.MakeOmegaDoc())
									return odocs, nil
								}
								if err != nil {
									return deriveCorrectExit(err)
								}
								// Found end of this current OmegaDoc, wrap it all up and reset
								odocs = append(odocs, curodoc.MakeOmegaDoc())
								curodoc = parseOdoc{
									SourceFilePath: srcpath,
								}
								goto RESET_CONTINUE
							} else {
								curodoc.AppDestFP(r)
							}
						}
					} else {
						// This is an attribute, there could be many
						split := strings.SplitN(s, ":", 2)
						if len(split) != 2 {
							return nil, fmt.Errorf("when parsing attribute on line %d of file %q, attribute split on ':' did not have length 2, had length %d: %v", rdr.LineNumber(), srcpath, len(split), split)
						}
						curodoc.AppAttr(split[0], split[1])
					}
				}
			}
		}
	}
}
