package thumbhash

import (
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"

	_ "golang.org/x/image/webp"

	"github.com/galdor/go-thumbhash"
)

// GetThumbhash get thumbhash from image
func GetThumbhash(r io.Reader) (string, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return "", err
	}
	width := img.Bounds().Size().X

	binHash := thumbhash.EncodeImage(img)
	hash := base64.StdEncoding.EncodeToString(binHash)

	return fmt.Sprintf("%d|%s", width, hash), nil
}
