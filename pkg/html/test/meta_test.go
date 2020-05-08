package test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/html"
)

var testCase = `<head>
<title>Test Case</title>
<meta charset="iso-8859-1" />
<meta property="og:title" content="test case" />
<meta name="description" content="general description">
<meta property="twitter:description" content="twitter description" />
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<meta name="author" content="John Doe">
</head>`

func TestExtract(t *testing.T) {
	meta, err := html.ExtractMeta(strings.NewReader(testCase))
	assert.Nil(t, err)
	assert.Equal(t, 7, len(meta))
	assert.Equal(t, "title", meta["title"].Name)
	assert.Equal(t, "Test Case", meta["title"].Content)
	assert.Equal(t, "", meta["og:title"].Name)
	assert.Equal(t, "og:title", meta["og:title"].Property)
	assert.Equal(t, "test case", meta["og:title"].Content)
	assert.Equal(t, "twitter description", meta.GetContent("twitter:description", "description"))
	assert.Equal(t, "iso-8859-1", meta["charset"].Content)
}
