package helper

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"
)

var wordDecoder = new(mime.WordDecoder)

// ExtractMailHeader extract some mail header
func ExtractMailHeader(header mail.Header) (from string, subject string) {
	from, _ = wordDecoder.DecodeHeader(header.Get("From"))
	subject, _ = wordDecoder.DecodeHeader(header.Get("Subject"))
	return
}

// ParseMultipart parse multipart mail body
func ParseMultipart(body io.Reader, boundary string) (html string, text string, err error) {
	reader := multipart.NewReader(body, boundary)
	if reader == nil {
		err = errors.New("unable to create mutipart reader")
		return
	}
	var part *multipart.Part
	for {
		part, err = reader.NextPart()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			break
		}

		mediaType, _, _ := mime.ParseMediaType(part.Header.Get("Content-Type"))
		isText := strings.HasPrefix(mediaType, "text/plain")
		isHTML := strings.HasPrefix(mediaType, "text/html")
		if !(isHTML || isText) {
			break
		}
		var b []byte
		b, err = io.ReadAll(part)
		if err != nil {
			return
		}
		if strings.HasPrefix(mediaType, "text/plain") {
			text = string(b)
		}
		if strings.HasPrefix(mediaType, "text/html") {
			html = string(b)
		}
	}
	return
}

// ExtractMailContent extract HTML and TEXT part of a mail body
func ExtractMailContent(msg *mail.Message) (html string, text string, err error) {
	// extract Media Type
	contentType := msg.Header.Get("Content-Type")
	var mediaType string
	var params map[string]string
	mediaType, params, err = mime.ParseMediaType(contentType)
	if err != nil {
		return
	}
	isText := strings.HasPrefix(mediaType, "text/plain")
	isHTML := strings.HasPrefix(mediaType, "text/html")
	isMultiPart := strings.HasPrefix(mediaType, "multipart")

	if !(isText || isHTML || isMultiPart) {
		err = fmt.Errorf("unsuported Media Type: %s", mediaType)
		return
	}
	// read body
	// TODO limit body size
	if isMultiPart {
		return ParseMultipart(msg.Body, params["boundary"])
	}
	// decode TEXT and HTML content
	var b []byte
	b, err = io.ReadAll(msg.Body)
	if err != nil {
		err = fmt.Errorf("unable to read body: %w", err)
		return
	}
	if isHTML {
		html = string(b)
	} else {
		text = string(b)
	}
	return
}
