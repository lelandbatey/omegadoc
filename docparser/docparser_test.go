package docparser

import (
	//"fmt"
	"strings"
	"testing"

	"github.com/lelandbatey/omegadoc/domain"
	"github.com/stretchr/testify/require"
)

// Don't treat the omegadocs in this file as omegadocs if this were to be
// scanned by the CLI. Though the strings with OmegaDocs in them will be fed as
// valid OmegaDocs through tests.
//#!/usr/bin/env omegadoc ignore-this-file

func TestParseDocBasic(t *testing.T) {
	rdr := strings.NewReader(`
some stuff
other stuff
#!/usr/bin/env omegadoc <<EOOD omegadoc/testdoc.md
this is a testing document
EOOD
other stuff`)
	dp := NewDocParser()
	odocs, err := dp.ParseDoc("tmp/testfile.md", rdr)
	require.NoError(t, err)
	require.Len(t, odocs, 1)
	require.Equal(t, odocs[0].Contents, "this is a testing document\n")
}

func mkoat(vals ...string) []domain.OmegaAttribute {
	oatts := []domain.OmegaAttribute{}
	for i := 0; i < len(vals); i += 2 {
		oatts = append(oatts, domain.OmegaAttribute{Key: vals[i], Value: vals[i+1]})
	}
	return oatts
}

func TestParseAssorted(t *testing.T) {
	type exp struct {
		Contents string
		DestFP   string
		Attrs    []domain.OmegaAttribute
		Err      string
	}
	type tst struct {
		Def  string
		Exps []exp
	}

	var tests []tst = []tst{
		{Def: "#!/usr/bin/env omegadoc <<EXT r/a.md\nfoobarEXT", Exps: []exp{{Contents: "foobar", DestFP: "r/a.md"}}},
		// Extra whitespace after the output path but before the newline should
		// be interpreted as part of the output path.
		{Def: "#!/usr/bin/env omegadoc <<EXT r/a.md \nfoobarEXT", Exps: []exp{{Contents: "foobar", DestFP: "r/a.md "}}},
		// Missing the output path means not parsed as an OmegaDoc.
		{Def: "#!/usr/bin/env omegadoc <<EXT \nfoobarEXT", Exps: []exp{}},
		// Ending the file in the middle of an OmegaDoc is considered a valid way to end the OmegaDoc.
		{Def: "#!/usr/bin/env omegadoc <<NLEXT r/a.md\nfoobar", Exps: []exp{{Contents: "foobar", DestFP: "r/a.md"}}},
		// Including the ignore directive causes the file to be ignored
		{Def: "#!/usr/bin/env omegadoc ignore-this-file\n\n" +
			"#!/usr/bin/env omegadoc <<EXT r/a.md\nfoobarEXT", Exps: []exp{}},
		// Including an ignore directive _after_ a valid omegadoc definition
		// causes nothing to happen; the ignore directive is itself ignored.
		{Def: "#!/usr/bin/env omegadoc <<EXT r/a.md\nfoobarEXT]\n" +
			"#!/usr/bin/env omegadoc ignore-this-file\n\n", Exps: []exp{{Contents: "foobar", DestFP: "r/a.md"}}},
		// Ensure destination file path can contain spaces.
		{Def: "#!/usr/bin/env omegadoc <<EXT tmp/with spaces.md\nfoobarEXT",
			Exps: []exp{{Contents: "foobar", DestFP: "tmp/with spaces.md"}}},
		// Basic attributes are parsed
		{Def: "#!/usr/bin/env omegadoc <<EXT fizz:pow r/a.md\nfoobarEXT",
			Exps: []exp{{Contents: "foobar", DestFP: "r/a.md", Attrs: mkoat("fizz", "pow")}}},
		// Duplicate attribute keys are allowed
		{Def: "#!/usr/bin/env omegadoc <<EXT fizz:pow  fizz:wahoo r/a.md\nfoobarEXT",
			Exps: []exp{{Contents: "foobar", DestFP: "r/a.md", Attrs: mkoat("fizz", "pow", "fizz", "wahoo")}}},
	}

	dp := NewDocParser()
	for tidx, test := range tests {
		rdr := strings.NewReader(test.Def)
		odocs, err := dp.ParseDoc("/tmp/testfile.md", rdr)
		require.Lenf(t, odocs, len(test.Exps), "test #%d expects to have created %d OmegaDocs but instead created %d, err: %v", tidx, len(test.Exps), len(odocs), err)
		for idx, expect := range test.Exps {
			odoc := odocs[idx]
			if expect.Err == "" {
				require.Equal(t, expect.Contents, odoc.Contents, "extracted and expected contents of OmegaDoc differ")
				require.Equal(t, expect.DestFP, odoc.DestFilePath, "extracted and expected destination file paths differ")
			} else {
				require.Equal(t, err.Error(), expect.Err)
			}
		}
	}
}
