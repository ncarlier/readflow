package test

import (
	"bytes"
	"testing"

	"github.com/ncarlier/readflow/pkg/template"
	_ "github.com/ncarlier/readflow/pkg/template/all"
	"github.com/stretchr/testify/assert"
)

const goTemplateTestCase = `
Article title: {{ .title | urlquery }}
Article link: {{ .url }}
`

func TestGoTemplateEngine(t *testing.T) {
	provider, err := template.NewTemplateEngine("gotemplate", goTemplateTestCase)
	assert.Nil(t, err)
	assert.NotNil(t, provider)
	var buf bytes.Buffer
	err = provider.Execute(&buf, data)
	assert.Nil(t, err)
	result := buf.String()
	assert.Contains(t, result, "Article title: Foo+%26+Bar")
	assert.Contains(t, result, "link: http://foo.bar")
}
