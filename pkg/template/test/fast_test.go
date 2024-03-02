package test

import (
	"bytes"
	"testing"

	"github.com/ncarlier/readflow/pkg/template"
	_ "github.com/ncarlier/readflow/pkg/template/all"
	"github.com/stretchr/testify/assert"
)

const fastTestCase = `
Article title: {{ title | urlquery }}
Article link: {{ url }}
Article text: {{ text | json }}
`

var data = map[string]interface{}{
	"id":    uint(1),
	"title": "Foo & Bar",
	"text":  `let's call this a "test"`,
	"url":   "http://foo.bar",
}

func TestFastTemplateEngine(t *testing.T) {
	provider, err := template.NewTemplateEngine("fast", fastTestCase)
	assert.Nil(t, err)
	assert.NotNil(t, provider)
	var buf bytes.Buffer
	err = provider.Execute(&buf, data)
	assert.Nil(t, err)
	result := buf.String()
	assert.Contains(t, result, "Article title: Foo+%26+Bar")
	assert.Contains(t, result, "link: http://foo.bar")
	assert.Contains(t, result, "text: \"let's call this a \\\"test\\\"\"")
}
