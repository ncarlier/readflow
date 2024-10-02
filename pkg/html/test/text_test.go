package test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ncarlier/readflow/pkg/html"
)

var testCases = []struct {
	input    string
	expected string
}{
	{input: "<div></div>", expected: ""},
	{input: "<div>simple text</div>", expected: "simple text"},
	{input: "<ul><li>1</li><li>2</li><li>3</li></ul>", expected: "1\n2\n3\n"},
	{input: "click <a href=\"test\">here</a>", expected: "click here"},
	{input: "<html><head><title>今日は</title></head><body>hétérogénéité<script type=\"javascript\">console.log('hello')</script></body>", expected: "今日は\nhétérogénéité"},
	{input: "&quot;I'm sorry, Dave. I'm afraid I can't do that.&quot; – HAL, 2001: A Space Odyssey", expected: "\"I'm sorry, Dave. I'm afraid I can't do that.\" – HAL, 2001: A Space Odyssey"},
}

func TestHTML2Text(t *testing.T) {
	for _, testCase := range testCases {
		text, err := html.HTML2Text(testCase.input)
		require.Nil(t, err)
		require.Equal(t, testCase.expected, text)
	}
}
