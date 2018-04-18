package courier

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"golang.org/x/net/html/charset"
)

// Attachment describes a simple interface for files attached to an email.
type Attachment interface {
	// Name returns the name of the attachment
	Name() string

	// Content returns the string content of the attachment.
	Content() ([]byte, error)

	// Base64Content returns the content of the attachment as an encoded base64 string.
	Base64Content() (string, error)

	// ContentType returns the MIME content type of the attached file. For example: image/jpeg.
	ContentType() string

	// Charset returns the character set of the attached file. For example: utf-8.
	Charset() string

	// ContentID returns the Content ID header of the attached file. This is used when embedding attached files into
	// emails.
	ContentID() string
}

// AttachmentHeaders defines the necessary HTTP headers for an attachment
type AttachmentHeaders struct {
	// name defines an explicit name to use on the attachment
	Name string

	// contentType defines an explicit content type to set on the attachment.
	ContentType string

	// charset defines an explicit character set encoding to set on the attachment.
	Charset string

	// contentID defines the Content ID header to set on the attachment.
	ContentID string
}

// FileAttachment is a type of attachment referenced by a file path.
type FileAttachment struct {
	// The path to the file this attachment references.
	File *os.File

	Headers AttachmentHeaders

	// content is the cached string value of the data in the file.
	content []byte
}

// Name returns the name of the attachment
func (f FileAttachment) Name() string {
	if f.Headers.Name != "" {
		return f.Headers.Name
	}

	return path.Base(f.File.Name())
}

// Content returns the string content of the attachment.
func (f FileAttachment) Content() ([]byte, error) {
	if f.content != nil {
		return f.content, nil
	}

	return ioutil.ReadAll(f.File)
}

// Base64Content returns the content of the attachment as an encoded base64 string.
func (f FileAttachment) Base64Content() (string, error) {
	c, err := f.Content()

	if err != nil {
		return "", err
	}

	return string(c), err
}

// ContentType returns the MIME content type of the attached file. For example: image/jpeg.
func (f FileAttachment) ContentType() string {
	if f.Headers.ContentType != "" {
		return f.Headers.ContentType
	}

	// @TODO handle errors here
	c, _ := f.Content()

	return http.DetectContentType(c)
}

// Charset returns the character set of the attached file. For example: utf-8.
func (f FileAttachment) Charset() string {
	if f.Headers.Charset != "" {
		return f.Headers.Charset
	}

	// @TODO handle errors here
	c, _ := f.Content()

	_, n, _ := charset.DetermineEncoding(c, f.ContentType())

	return n
}

// ContentID returns the Content ID header of the attached file. This is used when embedding attached files into
// emails.
func (f FileAttachment) ContentID() string {
	return f.Headers.ContentID
}
