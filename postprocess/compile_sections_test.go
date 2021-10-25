package postprocess

import (
	"testing"

	"github.com/lelandbatey/omegadoc/domain"
	"github.com/stretchr/testify/require"
)

func TestCompileSections(t *testing.T) {
	ppr := SectionsCompiler{}
	type test struct {
		odocs    []domain.OmegaDoc
		expected []domain.OmegaDoc
	}

	for _, tst := range []test{
		{
			odocs: []domain.OmegaDoc{
				{Attributes: []domain.OmegaAttribute{{Key: "section", Value: "001"}}, Contents: "\nFirst section\n"},
				{Attributes: []domain.OmegaAttribute{{Key: "section", Value: "002"}}, Contents: "\nSecond section\n"},
			},
			expected: []domain.OmegaDoc{
				{Attributes: []domain.OmegaAttribute{{Key: "section", Value: "001"}, {Key: "section", Value: "002"}},
					Contents: "\nFirst section\n\nSecond section\n"},
			},
		},
	} {
		newdocs, err := ppr.Postprocess(tst.odocs)
		require.NoError(t, err)
		require.Len(t, newdocs, len(tst.expected))
		for idx, nd := range newdocs {
			ed := tst.expected[idx]
			require.Equal(t, ed.Contents, nd.Contents)
			require.Equal(t, ed.Attributes, nd.Attributes)
		}
	}
}
