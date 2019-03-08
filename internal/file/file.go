package file

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

// Download downloads the file by its link
func Download(link string) (string, error) {
	resp, err := http.DefaultClient.Get(link)
	if err != nil {
		return "", errors.Wrap(err, "unable to get file")
	}
	defer resp.Body.Close()

	tmpFile, err := NewTemp(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "unable to create temp file")
	}
	defer tmpFile.Close()

	return tmpFile.Name(), nil
}

// NewTemp creates a new temporary file
func NewTemp(content io.Reader) (*os.File, error) {
	file, err := ioutil.TempFile("", "*")
	if err != nil {
		return nil, errors.Wrap(err, "unable to create file")
	}

	if content != nil {
		if _, err = io.Copy(file, content); err != nil {
			file.Close()
			return nil, errors.Wrap(err, "unable to write file content")
		}
	}
	return file, nil
}
