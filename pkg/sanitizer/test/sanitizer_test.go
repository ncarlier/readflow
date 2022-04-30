package test

import (
	"testing"

	"github.com/ncarlier/readflow/pkg/sanitizer"
	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	content     string
	expectation string
}{
	{`<a href="https://001print.com/foo.html" alt="test">foo</a>`, `foo`},
	{`<a href="https://print.com/foo.html" alt="test">foo</a>`, `<a href="https://print.com/foo.html" rel="nofollow noopener" target="_blank">foo</a>`},
	{`<img src="test.png" class="test"  alt="test" />`, `<img src="test.png" alt="test"/>`},
}

func TestSanitizer(t *testing.T) {
	bl, err := sanitizer.NewBlockList("file://block-list.txt")
	assert.Nil(t, err)
	sanitizer := sanitizer.NewSanitizer(bl)

	for _, tt := range tests {
		cleaned := sanitizer.Sanitize(tt.content)
		assert.Equal(t, tt.expectation, cleaned)
	}
}

func TestSanitizerWithoutBlockList(t *testing.T) {
	sanitizer := sanitizer.NewSanitizer(nil)

	for idx, tt := range tests {
		cleaned := sanitizer.Sanitize(tt.content)
		if idx == 0 {
			assert.NotEqual(t, tt.expectation, cleaned)
		} else {
			assert.Equal(t, tt.expectation, cleaned)
		}
	}
}
