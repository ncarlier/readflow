package test

import (
	"testing"

	"github.com/ncarlier/readflow/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestHashtagsExtraction(t *testing.T) {
	require.Contains(t, utils.ExtractHashtags("hello #world"), "#world")
}

func TestReplaceHashtagsPrefix(t *testing.T) {
	require.Equal(t, "hello ~world", utils.ReplaceHashtagsPrefix("hello #world", "~"))
}
