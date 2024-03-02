package epub

import (
	"archive/zip"
	"html"
	"io"
	"time"

	"github.com/google/uuid"

	"github.com/ncarlier/readflow/pkg/mediatype"
)

const contentDir = "content/"

type epubItem struct {
	ContentType string
	Filename    string
}

// Writer of a ePub format
type Writer struct {
	id      string
	title   string
	archive *zip.Writer
	items   []epubItem
}

// NewWriter create new ePub writer
func NewWriter(w io.Writer, title string) (*Writer, error) {
	// create ZIP archive
	archive := zip.NewWriter(w)
	// generate ID
	uid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	// Store mimetype
	mimetypeHeader := &zip.FileHeader{
		Name:   "mimetype",
		Method: zip.Store,
	}
	f, err := archive.CreateHeader(mimetypeHeader)
	if err != nil {
		return nil, err
	}
	f.Write([]byte(mediatype.Epub))
	return &Writer{
		id:      uid.String(),
		title:   title,
		archive: archive,
	}, nil
}

// NewContainer create new ePub container
func (w *Writer) NewContainer() error {
	f, err := w.archive.Create("META-INF/container.xml")
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(epubContainerDescriptor))
	return err
}

// NewItem create new ePub container item
func (w *Writer) NewItem(filename, contentType string) (io.Writer, error) {
	w.items = append(w.items, epubItem{
		Filename:    filename,
		ContentType: contentType,
	})
	return w.archive.Create(contentDir + filename)
}

// WriteOPF write ePub file
func (w *Writer) WriteOPF(filename, spineRef string) error {
	f, err := w.archive.Create(contentDir + filename)
	if err != nil {
		return err
	}
	return epubOpfTmpl.Execute(f, epubOpfTmplArgs{
		ID:       html.EscapeString(w.id),
		Title:    html.EscapeString(w.title),
		SpineRef: spineRef,
		Time:     time.Now().Format(time.RFC3339),
		Items:    w.items,
	})
}

// Close the writer
func (w *Writer) Close() error {
	return w.archive.Close()
}
