package test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ncarlier/readflow/pkg/html"
)

func TestHTML2Text(t *testing.T) {
	text, err := html.HTML2Text("<div></div>")
	require.Nil(t, err)
	require.Equal(t, "", text)

	text, err = html.HTML2Text("<div>simple text</div>")
	require.Nil(t, err)
	require.Equal(t, "simple text", text)

	text, err = html.HTML2Text("<ul><li>1</li><li>2</li><li>3</li></ul>")
	require.Nil(t, err)
	require.Equal(t, "1\n2\n3\n", text)

	text, err = html.HTML2Text("click <a href=\"test\">here</a>")
	require.Nil(t, err)
	require.Equal(t, "click here", text)

	text, err = html.HTML2Text("<html><head><title>今日は</title></head><body>hétérogénéité<script type=\"javascript\">console.log('hello')</script></body>")
	require.Nil(t, err)
	require.Equal(t, "今日は\nhétérogénéité", text)

	text, err = html.HTML2Text("&quot;I'm sorry, Dave. I'm afraid I can't do that.&quot; – HAL, 2001: A Space Odyssey")
	require.Nil(t, err)
	require.Equal(t, "\"I'm sorry, Dave. I'm afraid I can't do that.\" – HAL, 2001: A Space Odyssey", text)
}
