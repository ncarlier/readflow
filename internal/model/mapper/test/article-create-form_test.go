package test

import (
	"encoding/json"
	"testing"

	"github.com/ncarlier/readflow/internal/model/mapper"
	"github.com/stretchr/testify/require"
)

const conf = `
_=a
title=b.c

# this is a comment
url=b.d
`

func TestArticleCreateFormMapper(t *testing.T) {
	var object map[string]interface{}
	json.Unmarshal([]byte(`{
		"a": [
			"b": {
				"c": "foo",
				"d": "https://aaa"
			}
		]
 	}`), &object)

	m, err := mapper.NewArticleCreateFormMapper(conf)
	require.Nil(t, err)

	result, err := m.Build(object)
	require.Nil(t, err)
	require.Len(t, result, 1)
	require.Equal(t, "foo", result[0].Title)
	require.NotNil(t, result[0].URL)
	require.Equal(t, "https://aaa", *result[0].URL)
}
