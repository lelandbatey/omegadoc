package postprocess

import (
	"testing"

	"github.com/lelandbatey/omegadoc/domain"
	"github.com/stretchr/testify/require"
)

func TestMarkdownLinkRewriter(t *testing.T) {
	ppr := MarkdownLinkRewriter{}

	docs := []domain.OmegaDoc{
		{
			Contents: "hello [foo](zap/bar.md)\n",
		},
	}

	newdocs, err := ppr.Postprocess(docs)
	require.NoError(t, err)
	require.Equal(t, "hello [foo](zap/bar.html)\n", newdocs[0].Contents)
}
