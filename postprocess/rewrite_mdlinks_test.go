package postprocess

import (
	"testing"

	"github.com/lelandbatey/omegadoc/domain"
	"github.com/stretchr/testify/require"
)

func TestMarkdownLinkRewriter(t *testing.T) {
	ppr := MarkdownLinkRewriter{}

	type test struct {
		odocs    []domain.OmegaDoc
		expected string
	}

	for idx, tst := range []test{
		{
			odocs:    []domain.OmegaDoc{{Contents: "hello [foo](zap/bar.md)\n"}},
			expected: "hello [foo](zap/bar.html)\n",
		},
		// The "two spaces before newline" case is broken upstream in Kunde21.
		// {
		// 	odocs:    []domain.OmegaDoc{{Contents: "hello  \nwhat\n[foo](zap/bar.md)\n"}},
		// 	expected: "hello  \nwhat\n[foo](zap/bar.html)\n",
		// },
		{
			odocs:    []domain.OmegaDoc{{Contents: "hello friends\n\n[foo](zap/bar.md)\n"}},
			expected: "hello friends\n\n[foo](zap/bar.html)\n",
		},
		{
			odocs:    []domain.OmegaDoc{{Contents: "[foo](zap/bar.html)\n"}},
			expected: "[foo](zap/bar.html)\n",
		},
		{
			odocs:    []domain.OmegaDoc{{Contents: "[foo](zap/bar.md.html)\n"}},
			expected: "[foo](zap/bar.md.html)\n",
		},
	} {
		newdocs, err := ppr.Postprocess(tst.odocs)
		require.NoError(t, err, "test #%d", idx)
		require.Equal(t, tst.expected, newdocs[0].Contents, "test #%d", idx)
	}
}
