package mail

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"strings"
)

var wordDecoder = new(mime.WordDecoder)

// extractMailHeader extract some mail header
func extractMailHeader(header mail.Header) (from, subject string) {
	from, _ = wordDecoder.DecodeHeader(header.Get("From"))
	subject, _ = wordDecoder.DecodeHeader(header.Get("Subject"))
	return
}

func decodeEmailBody(body io.Reader, encoding string) (string, error) {
	var reader io.Reader
	switch strings.ToLower(encoding) {
	case "quoted-printable":
		reader = quotedprintable.NewReader(body)
	case "base64":
		reader = base64.NewDecoder(base64.StdEncoding, body)
	case "", "8bit", "7bit":
		reader = body
	default:
		return "", fmt.Errorf("unsuported encoding: %s", encoding)
	}
	b, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("unable to read body: %w", err)
	}
	return string(b), nil
}

// parseMultipart parse multipart mail body
func parseMultipart(body io.Reader, boundary string) (html, text string, err error) {
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

// extractMailContent extract HTML and TEXT part of a mail body
func extractMailContent(msg *mail.Message) (html, text string, err error) {
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

	encoding := msg.Header.Get("Content-Transfer-Encoding")

	// read body
	// TODO limit body size
	if isMultiPart {
		return parseMultipart(msg.Body, params["boundary"])
	}
	// decode TEXT and HTML content
	body, err := decodeEmailBody(msg.Body, encoding)
	if err != nil {
		err = fmt.Errorf("unable to read body: %w", err)
		return
	}
	if isHTML {
		html = body
	} else {
		text = body
	}
	return
}
