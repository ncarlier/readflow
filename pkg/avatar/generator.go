package avatar

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
)

type Avatar struct {
	Size  int    `json:"size"`
	Parts []Part `json:"parts"`
}

type Part struct {
	Name string `json:"name"`
	Nb   int    `json:"nb"`
}

type Avatars map[string]*Avatar

// Generator instance
type Generator struct {
	avatars    Avatars
	directory  string
	defaultSet string
}

// Generate avatar
func (g *Generator) Generate(seed, set string) (*bytes.Buffer, error) {
	if set == "" {
		set = g.defaultSet
	}
	// get avatar set
	avatar, ok := g.avatars[set]
	if !ok {
		// avatar set doesn't exists, using default
		set = g.defaultSet
		avatar = g.avatars[set]
	}

	// init random seed
	hSeed := md5.Sum([]byte(seed))
	nSeed, err := strconv.ParseInt(hex.EncodeToString(hSeed[:])[:6], 16, 64)
	if err != nil {
		return nil, fmt.Errorf("unable to compute seed: %v", err)
	}
	random := rand.New(rand.NewSource(nSeed))

	// build avatar spec
	specs := make([]Part, len(avatar.Parts))
	for idx, part := range avatar.Parts {
		specs[idx] = Part{
			Name: part.Name,
			Nb:   random.Intn(part.Nb) + 1,
		}
	}

	// init image
	upLeft := image.Point{0, 0}
	lowRight := image.Point{avatar.Size, avatar.Size}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// build avatar image
	for _, part := range specs {
		imgSrc := filepath.Join(g.directory, set, fmt.Sprintf("%s_%d.png", part.Name, part.Nb))
		if _, err := os.Stat(imgSrc); err == nil {
			img, _ = blendWithImageFile(img, imgSrc)
		}
	}
	buff := new(bytes.Buffer)
	err = png.Encode(buff, img)
	if err != nil {
		return nil, fmt.Errorf("unable to encode image: %v", err)
	}
	return buff, nil
}

// NewServer creates new server instance
func NewGenerator(dir, defaultSet string) (*Generator, error) {
	if _, err := os.Stat(dir); err != nil {
		return nil, err
	}
	avatars := make(Avatars)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	defaultSetExists := false
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		avatar, e := readAvatarDir(filepath.Join(dir, file.Name()))
		if e != nil {
			return nil, err
		}
		if defaultSet == "" {
			defaultSet = file.Name()
		}
		if file.Name() == defaultSet {
			defaultSetExists = true
		}
		avatars[file.Name()] = avatar
	}
	if !defaultSetExists {
		return nil, fmt.Errorf("avatar default set not found: %s", defaultSet)
	}

	return &Generator{
		avatars:    avatars,
		directory:  dir,
		defaultSet: defaultSet,
	}, nil
}

func readAvatarDir(dir string) (*Avatar, error) {
	content, err := os.ReadFile(filepath.Join(dir, "_avatar.json"))
	if err != nil {
		return nil, err
	}

	avatar := Avatar{}
	if err := json.Unmarshal([]byte(content), &avatar); err != nil {
		return nil, err
	}
	return &avatar, nil
}

func blendWithImageFile(img *image.RGBA, src string) (*image.RGBA, error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer srcFile.Close()

	srcImg, err := png.Decode(srcFile)
	if err != nil {
		return nil, err
	}
	b := img.Bounds()
	p := image.Point{0, 0}
	result := image.NewRGBA(b)
	draw.Draw(result, b, img, p, draw.Over)
	draw.Draw(result, b, srcImg, p, draw.Over)
	return result, nil
}
