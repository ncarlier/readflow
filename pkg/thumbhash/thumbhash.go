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
	size := img.Bounds().Size()

	binHash := thumbhash.EncodeImage(img)
	hash := base64.StdEncoding.EncodeToString(binHash)

	return fmt.Sprintf("%dx%d|%s", size.X, size.Y, hash), nil
}
