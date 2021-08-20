package epub

import (
	"archive/zip"
	"html"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/ncarlier/readflow/pkg/constant"
)

const contentDir = "content/"

type epubItem struct {
	ContentType string
	Filename    string
}

type Writer struct {
	id      string
	title   string
	archive *zip.Writer
	items   []epubItem
}

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
	f.Write([]byte(constant.ContentTypeEpub))
	return &Writer{
		id:      uid.String(),
		title:   title,
		archive: archive,
	}, nil
}

func (w *Writer) NewContainer() error {
	f, err := w.archive.Create("META-INF/container.xml")
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(epubContainerDescriptor))
	return err
}

func (w *Writer) NewItem(filename string, ContentType string) (io.Writer, error) {
	w.items = append(w.items, epubItem{
		Filename:    filename,
		ContentType: ContentType,
	})
	return w.archive.Create(contentDir + filename)
}

func (w *Writer) WriteOPF(filename string, spineRef string) error {
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

func (w *Writer) Close() error {
	return w.archive.Close()
}
