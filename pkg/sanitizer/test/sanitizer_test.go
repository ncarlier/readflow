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
	{`<img src="test.png" class="test"  alt="test" />`, `<img src="test.png" alt="test">`},
	{`<img src="test.png" srcset="data:image/gif;base64,R0lGODlhAQABAIAAAAUEBAAAACwAAAAAAQABAAACAkQBADs= 2x" class="test"  alt="test" />`, `<img src="test.png" alt="test">`},
}

func TestSanitizer(t *testing.T) {
	bl, err := sanitizer.NewBlockList("file://block-list.txt", sanitizer.DefaultBlockList)
	assert.Nil(t, err)
	san := sanitizer.NewSanitizer(bl)

	for _, tt := range tests {
		cleaned := san.Sanitize(tt.content)
		assert.Equal(t, tt.expectation, cleaned)
	}
}

func TestSanitizerWithoutBlockList(t *testing.T) {
	san := sanitizer.NewSanitizer(nil)

	for idx, tt := range tests {
		cleaned := san.Sanitize(tt.content)
		if idx == 0 || idx == 3 {
			assert.NotEqual(t, tt.expectation, cleaned)
		} else {
			assert.Equal(t, tt.expectation, cleaned)
		}
	}
}
