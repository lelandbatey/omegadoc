package domain

const (
	// START_OMEGADOC is the magic string which indicates the first part of an
	// opening statement for an OmegaDoc. This magic string is one of two parts
	// of an "opening statement", with the second part being a "delimiting
	// identifier" (e.g. the word "EOF" or "END"). When the magic string is
	// followed by a "delimiting identifier", together those form an "opening
	// statement" of an OmegaDoc.
	START_OMEGADOC = "#!/usr/bin/env" + " omegadoc <<"
	// IGNORE_OMEGADOC is a magic string which acts as an "ignore directive"
	// for an OmegaDoc. If a file contains an "ignore directive" in its bytes
	// before an OmegaDoc opening statement, then that file will be considered
	// to have NO OmegaDocs in it, even if it otherwise contains one or more
	// valid OmegaDocs. If a file contains an "ignore directive" after one or
	// more valid OmegaDoc directives, then that "ignore directive" will itself
	// be ignored.
	IGNORE_OMEGADOC = "#!/usr/bin/env" + " omegadoc ignore-this-file"
)
