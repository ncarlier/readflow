package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/html"
)

func TestHTML2Text(t *testing.T) {
	text, err := html.HTML2Text("<div></div>")
	assert.Nil(t, err)
	assert.Equal(t, "", text)

	text, err = html.HTML2Text("<div>simple text</div>")
	assert.Nil(t, err)
	assert.Equal(t, "simple text", text)

	text, err = html.HTML2Text("<ul><li>1</li><li>2</li><li>3</li></ul>")
	assert.Nil(t, err)
	assert.Equal(t, "1\n2\n3\n", text)

	text, err = html.HTML2Text("click <a href=\"test\">here</a>")
	assert.Nil(t, err)
	assert.Equal(t, "click here", text)

	text, err = html.HTML2Text("<html><head><title>A title</title></head><body>foo<script type=\"javascript\">console.log('hello')</script></body>")
	assert.Nil(t, err)
	assert.Equal(t, "A title\nfoo", text)

	text, err = html.HTML2Text("&quot;I'm sorry, Dave. I'm afraid I can't do that.&quot; – HAL, 2001: A Space Odyssey")
	assert.Nil(t, err)
	assert.Equal(t, "\"I'm sorry, Dave. I'm afraid I can't do that.\" – HAL, 2001: A Space Odyssey", text)
}
