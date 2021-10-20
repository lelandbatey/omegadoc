package docparser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/lelandbatey/omegadoc/domain"
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
}

func NewDocParser() domain.DocParser {
	return docfinder{}
}

// parseOdoc tracks all the data necessary for us to parse the document, and
// when complete may be turned into a "real" OmegaDoc
type parseOdoc struct {
	SourceFilePath string
	DestFilePath   []rune
	Contents       []rune
	//HTTPUrl string
}

func (po *parseOdoc) AppCont(r ...rune) {
	//for idx, x := range r {
	//fmt.Printf("Appending rune to Contents: %s %5v %d\n", string(x), x, idx)
	//}
	po.Contents = append(po.Contents, r...)
}

func (po *parseOdoc) AppDestFP(r ...rune) {
	po.DestFilePath = append(po.DestFilePath, r...)
}

func (po *parseOdoc) MakeOmegaDoc() domain.OmegaDoc {
	return domain.OmegaDoc{
		SourceFilePath: po.SourceFilePath,
		DestFilePath:   string(po.DestFilePath),
		Contents:       string(po.Contents),
	}
}

// ParseDoc for docfinder parses a text file and extracts all OmegaDocs present
// in the file. This is currently implemented as a simple direct parser,
// without being broken down into scanner/lexer since the language is so
// simple. In the future this implementation may need to be further broken down
// though, as features such as automatic indentation removal or line-prefix
// removal may require a full lexer/parser.
func (df docfinder) ParseDoc(srcpath string, data io.Reader) ([]domain.OmegaDoc, error) {
	var odocs []domain.OmegaDoc = []domain.OmegaDoc{}
	rdr := bufio.NewReader(data)

	// Outside of odoc
	// Inside OmegaDoc opening statement
	//     Inside magic string
	//         Inside a delimiting identifier
	//         OR Inside an ignore directive
	// Inside an output path
	// Inside the OmegaDoc body
	// Inside a closing delimiting identifier

	const common_prefix string = "#!/usr/bin/env omegadoc "
	var begindoc_magicrunes []rune = []rune(strings.ReplaceAll(domain.START_OMEGADOC, common_prefix, ""))
	var ignordoc_magicrunes []rune = []rune(strings.ReplaceAll(domain.IGNORE_OMEGADOC, common_prefix, ""))
	common_magicrunes := []rune(common_prefix)

	var curodoc parseOdoc = parseOdoc{
		SourceFilePath: srcpath,
	}
	var delimiting_ident []rune = []rune{}

	deriveCorrectExit := func(err error) ([]domain.OmegaDoc, error) {
		// End of file isn't necessarily an error, more a signal that we're
		// done here.
		if errors.Is(err, io.EOF) {
			return odocs, nil
		}
		return nil, err
	}
	for {
	RESET_CONTINUE:
		r, _, err := rdr.ReadRune()
		if err != nil {
			return deriveCorrectExit(err)
		}
		// We're in the common_prefix which could be an ignore statement or
		// "beginning statement"
		if r == common_magicrunes[0] {
			var commonpos int = 0
			//fmt.Printf("commonpos: %v\n", commonpos)
			for {
				commonpos += 1
				// We reached the end of the common_prefix, now figure out if
				// it's an ignoredoc or a beginning statement
				if commonpos == len(common_magicrunes) {
					//fmt.Printf("Reached end of common prefix\n")
					r, _, err = rdr.ReadRune()
					if err != nil {
						return deriveCorrectExit(err)
					}
					// It could be a "beginning statement"
					if r == begindoc_magicrunes[0] {
						beginpos := 0
						//fmt.Printf("first beginpos: %v\n", beginpos)
						for {
							beginpos += 1
							//fmt.Printf("beginpos: %v\n", beginpos)
							// Yes, this is a beginning statement. Now gather
							// the delimiting identifier
							if beginpos == len(begindoc_magicrunes) {
								//fmt.Printf("Reached end of beginning statement\n")
								for {
									r, _, err = rdr.ReadRune()
									if err != nil {
										return deriveCorrectExit(err)
									}
									if unicode.IsSpace(r) && len(delimiting_ident) != 0 {
										// We parsed the delimiting identifier so now we parse the output path
										pathLoc := "before"
										for {
											r, _, err = rdr.ReadRune()
											if err != nil {
												return deriveCorrectExit(err)
											}
											if pathLoc == "before" {
												if r == '\n' {
													goto RESET_CONTINUE
												} else if !unicode.IsSpace(r) {
													curodoc.AppDestFP(r)
													pathLoc = "within"
												}
											} else if pathLoc == "within" {
												if r == '\n' {
													// We parse the output path, now we're on to copying the entire document into
													// the body.
													tmp_pdi := []rune{}
													dipos := 0
													for {
														r, _, err = rdr.ReadRune()
														if errors.Is(err, io.EOF) {
															// Ending the file in the middle of an OmegaDoc is considered a
															// valid ending to the OmegaDoc.
															odocs = append(odocs, curodoc.MakeOmegaDoc())
															return odocs, nil
														}
														if err != nil {
															return deriveCorrectExit(err)
														}
														if r == delimiting_ident[dipos] {
															tmp_pdi = append(tmp_pdi, r)
															dipos += 1
															if dipos == len(delimiting_ident) {
																// Found the end of this current OmegaDoc, wrap it all up and reset.
																odocs = append(odocs, curodoc.MakeOmegaDoc())
																curodoc = parseOdoc{
																	SourceFilePath: srcpath,
																}
																goto RESET_CONTINUE
															}
														} else {
															if len(tmp_pdi) > 0 {
																curodoc.AppCont(tmp_pdi...)
																tmp_pdi = []rune{}
															}
															dipos = 0
															curodoc.AppCont(r)
														}
													}
												} else {
													curodoc.AppDestFP(r)
												}
											}
										}
									} else if unicode.IsSpace(r) && len(delimiting_ident) == 0 {
										goto RESET_CONTINUE
									}
									delimiting_ident = append(delimiting_ident, r)
								}
							}
							r, _, err = rdr.ReadRune()
							if err != nil {
								return deriveCorrectExit(err)
							}
							if r != begindoc_magicrunes[beginpos] {
								goto RESET_CONTINUE
							}
						}
					} else if r == ignordoc_magicrunes[0] {
						ignorpos := 0
						for {
							ignorpos += 1
							// If the ignore directive comes before any OmegaDocs have been
							// defined, then the whole file is ignored. Otherwise, the ignore
							// directive is itself ignored.
							if ignorpos == len(ignordoc_magicrunes) {
								if len(odocs) == 0 {
									return odocs, nil
								} else {
									goto RESET_CONTINUE
								}
							}
							r, _, err = rdr.ReadRune()
							if err != nil {
								return deriveCorrectExit(err)
							}
							if ignordoc_magicrunes[ignorpos] != r {
								goto RESET_CONTINUE
							}
						}
					}
				}
				r, _, err = rdr.ReadRune()
				if err != nil {
					return deriveCorrectExit(err)
				}
				if r != common_magicrunes[commonpos] {
					goto RESET_CONTINUE
				}
			}
		}
	}
}
