package sanitizer

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/ncarlier/readflow/pkg/utils"
)

var DefaultBlockList = []string{
	"data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7",
	"data:image/gif;base64,R0lGODlhAQABAIAAAAUEBAAAACwAAAAAAQABAAACAkQBADs=",
	"feedburner.google.com",
}

// BlockList is a block-list structure
type BlockList struct {
	location string
	filter   *bloom.BloomFilter
}

// NewBlockList create new block-list form a text file
func NewBlockList(location string, init []string) (*BlockList, error) {
	if location == "" {
		return nil, nil
	}
	// open block-list file
	input, err := utils.OpenResource(location)
	if err != nil {
		return nil, err
	}
	defer input.Close()

	// count number of lines
	var buf bytes.Buffer
	tee := io.TeeReader(input, &buf)
	size, err := utils.CountLines(tee)
	if err != nil {
		return nil, err
	}

	// initialize bloom filter
	filter := bloom.NewWithEstimates(size+uint(len(init)), 0.01)
	for _, v := range init {
		filter.AddString(v)
	}

	// read block-list file
	scanner := bufio.NewScanner(&buf)
	for scanner.Scan() {
		line := scanner.Text()
		if line[0] == '#' {
			continue
		}
		filter.AddString(strings.TrimPrefix(line, "0.0.0.0 "))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &BlockList{
		location: location,
		filter:   filter,
	}, nil
}

// Contains check if the provided value is into the block-list
func (bl *BlockList) Contains(value string) bool {
	return bl.filter.TestString(value)
}

// Location of the block-list
func (bl *BlockList) Location() string {
	return bl.location
}

// Size of the block-list
func (bl *BlockList) Size() uint32 {
	return bl.filter.ApproximatedSize()
}
