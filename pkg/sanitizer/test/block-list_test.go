package test

import (
	"testing"

	"github.com/ncarlier/readflow/pkg/sanitizer"
	"github.com/stretchr/testify/assert"
)

func TestLocalBlockList(t *testing.T) {
	location := "file://block-list.txt"
	bl, err := sanitizer.NewBlockList(location)
	assert.Nil(t, err)
	assert.Equal(t, location, bl.Location())
	assert.Equal(t, uint32(95), bl.Size())
	assert.True(t, bl.Contains("002.city"))
	assert.False(t, bl.Contains("02.city"))
}

func TestRemoteBlockList(t *testing.T) {
	location := "https://raw.githubusercontent.com/anudeepND/blacklist/master/adservers.txt"
	bl, err := sanitizer.NewBlockList(location)
	assert.Nil(t, err)
	assert.Equal(t, location, bl.Location())
	assert.Positive(t, bl.Size())
	assert.True(t, bl.Contains("feedburner.google.com"))
	assert.False(t, bl.Contains("google.com"))
}
