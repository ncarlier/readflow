package test

import (
	"testing"

	"github.com/ncarlier/readflow/pkg/utils"
	"github.com/stretchr/testify/require"
)

var testCases = []struct {
	input    string
	expected string
}{
	{
		input:    "a.b.html",
		expected: "a-b-html",
	},
	{
		input:    "whatēverwëirduserînput",
		expected: "whateverweirduserinput",
	},
}

func TestSanitizeFilename(t *testing.T) {
	for _, testCase := range testCases {
		require.Equal(t, testCase.expected, utils.SanitizeFilename(testCase.input))
	}
}
