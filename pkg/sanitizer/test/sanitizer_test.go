package test

import (
	"testing"

	"github.com/ncarlier/readflow/pkg/sanitizer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSanitizer_WithBlockList(t *testing.T) {
	t.Parallel()

	blockList, err := sanitizer.NewBlockList(
		"file://block-list.txt",
		sanitizer.DefaultBlockList,
	)
	require.NoError(t, err)

	san := sanitizer.NewSanitizer(blockList)

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "removes link to blocked host",
			content:  `<a href="https://001print.com/foo.html" alt="test">foo</a>`,
			expected: `foo`,
		},
		{
			name:     "removes blocked host using ambiguous four slash URL",
			content:  `<a href="https:////001print.com/foo.html" alt="test">foo</a>`,
			expected: `foo`,
		},
		{
			name:     "removes blocked host using malformed three slash URL",
			content:  `<a href="https:///001print.com/foo.html" alt="test">foo</a>`,
			expected: `foo`,
		},
		{
			name:     "removes blocked host using malformed one slash URL",
			content:  `<a href="https:/001print.com/foo.html" alt="test">foo</a>`,
			expected: `foo`,
		},
		{
			name:     "removes blocked host using protocol relative URL",
			content:  `<img src="//001print.com/image.png" alt="test">`,
			expected: `<img alt="test">`,
		},
		{
			name:     "removes blocked host with uppercase hostname",
			content:  `<img src="https://001PRINT.COM/image.png" alt="test">`,
			expected: `<img alt="test">`,
		},
		{
			name:     "removes blocked host with trailing dot",
			content:  `<img src="https://001print.com./image.png" alt="test">`,
			expected: `<img alt="test">`,
		},
		{
			name:     "keeps non blocked absolute link",
			content:  `<a href="https://print.com/foo.html" alt="test">foo</a>`,
			expected: `<a href="https://print.com/foo.html" rel="nofollow noopener" target="_blank">foo</a>`,
		},
		{
			name:     "keeps relative image URL",
			content:  `<img src="test.png" class="test" alt="test" />`,
			expected: `<img src="test.png" alt="test">`,
		},
		{
			name:     "removes invalid HTTP URL without authority",
			content:  `<img src="https:///example.com/image.png" alt="test">`,
			expected: `<img alt="test">`,
		},
		{
			name:     "removes URL containing credentials",
			content:  `<img src="https://user:password@print.com/image.png" alt="test">`,
			expected: `<img alt="test">`,
		},
		{
			name:     "removes unsupported javascript scheme",
			content:  `<a href="javascript:alert(1)">test</a>`,
			expected: `test`,
		},
		{
			name:     "removes entire srcset containing data URL",
			content:  `<img src="test.png" srcset="data:image/gif;base64,R0lGODlhAQABAIAAAAUEBAAAACwAAAAAAQABAAACAkQBADs= 2x" class="test" alt="test" />`,
			expected: `<img src="test.png" alt="test">`,
		},
		{
			name:     "removes srcset when first candidate is blocked",
			content:  `<img src="test.png" srcset="https://001print.com/blocked.png 1x, https://print.com/allowed.png 2x" alt="test">`,
			expected: `<img src="test.png" alt="test">`,
		},
		{
			name:     "removes srcset when later candidate is blocked",
			content:  `<img src="test.png" srcset="https://print.com/allowed.png 1x, https://001print.com/blocked.png 2x" alt="test">`,
			expected: `<img src="test.png" alt="test">`,
		},
		{
			name:     "keeps srcset when every candidate is allowed",
			content:  `<img src="test.png" srcset="https://print.com/image-1x.png 1x, https://print.com/image-2x.png 2x" alt="test">`,
			expected: `<img src="test.png" srcset="https://print.com/image-1x.png 1x, https://print.com/image-2x.png 2x" alt="test">`,
		},
		{
			name:     "removes blocked source srcset inside picture",
			content:  `<picture><source srcset="https://001print.com/image.webp 1x"><img src="test.png" alt="test"></picture>`,
			expected: `<picture><img src="test.png" alt="test"></picture>`,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := san.Sanitize(tt.content)

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestSanitizer_WithoutBlockList(t *testing.T) {
	t.Parallel()

	san := sanitizer.NewSanitizer(nil)

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "still applies bluemonday policy to javascript URL",
			content:  `<a href="javascript:alert(1)">test</a>`,
			expected: `test`,
		},
		{
			name:     "keeps allowed HTTP link",
			content:  `<a href="https://print.com/foo.html">foo</a>`,
			expected: `<a href="https://print.com/foo.html" rel="nofollow noopener" target="_blank">foo</a>`,
		},
		{
			name:     "keeps relative image URL",
			content:  `<img src="test.png" class="test" alt="test">`,
			expected: `<img src="test.png" alt="test">`,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := san.Sanitize(tt.content)

			assert.Equal(t, tt.expected, actual)
		})
	}
}
