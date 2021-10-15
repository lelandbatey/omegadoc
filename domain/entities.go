package domain

type OmegaDoc struct {
	// SourceFilePath is the path on-disk of the file which originally defined this OmegaDoc
	SourceFilePath string
	// DestFilePath is the path to where the OmegDoc should be stored in the
	// output FileTree. Even if this path is absolute (starts with a '/') it
	// will be treated as relative to the output directory when output.
	DestFilePath string
	Contents     string
}
