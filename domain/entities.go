package domain

type OmegaDoc struct {
	// SourceFilePath is the path on-disk of the file which originally defined this OmegaDoc
	SourceFilePath string
	// DestFilePath is the path to where the OmegDoc should be stored in the
	// output FileTree. Even if this path is absolute (starts with a '/') it
	// will be treated as relative to the output directory when output.
	DestFilePath string
	// The contents of the OmegaDoc, found between the opening statement (which
	// defines the delimiting identifier) and the delimiting identifier.
	Contents string
	// HTTPURL contains a single full HTTP URL where you can read the source of
	// this OmegaDoc in your web-browser. This URL is not present in the
	// original document and if present will have been derived from the git
	// repository which the document belongs within. If the file at
	// "SourceFilePath" is not within a git repository or there is not
	// configuration sufficient to derive the HTTP url for that file inside
	// that repository, then this HTTPUrl will be blank.
	HTTPUrl string
}
