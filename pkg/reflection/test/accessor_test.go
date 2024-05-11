package test

import (
	"encoding/json"
	"testing"

	"github.com/ncarlier/readflow/pkg/reflection"
	"github.com/stretchr/testify/require"
)

func TestGetField(t *testing.T) {
	var object map[string]interface{}
	json.Unmarshal([]byte(`{
		"a": [{
			"b": {
				"c": "foo",
				"d": 3,
				"e": true
			}
		}]
 	}`), &object)

	p := reflection.GetField(object, "a[0].b.c")
	v1, ok := p.String()
	require.True(t, ok)
	require.Equal(t, "foo", v1)

	p = reflection.GetField(object, "a[0].b.d")
	v2, ok := p.Int()
	require.True(t, ok)
	require.Equal(t, int64(3), v2)

	p = reflection.GetField(object, "a[0].b.e")
	v3, ok := p.Bool()
	require.True(t, ok)
	require.Equal(t, true, v3)
}
