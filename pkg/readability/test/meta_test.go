package test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/readability"
)

var testCase = `<head>
<title>Test case</title>
<meta charset="iso-8859-1" />
<meta property="og:title" content="test case" />
<meta name="description" content="general description">
<meta property="twitter:description" content="twitter description" />
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<meta name="author" content="John Doe">
</head>`

func TestExtract(t *testing.T) {
	metas, err := readability.ExtractMetas(strings.NewReader(testCase))
	assert.Nil(t, err)
	assert.Equal(t, 6, len(metas))
	assert.Equal(t, "", metas["og:title"].Name)
	assert.Equal(t, "og:title", metas["og:title"].Property)
	assert.Equal(t, "test case", metas["og:title"].Content)
	assert.Equal(t, "twitter description", *metas.GetContent("twitter:description", "description"))
	assert.Equal(t, "iso-8859-1", metas["charset"].Content)
}
