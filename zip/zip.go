package zip

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/yeka/zip"
)

// Zip is the structure that can work with decrypted files.
type Zip struct{}

// New creates a new Zip struct.
func New() *Zip {
	return &Zip{}
}

// WithPassword takes a source reader, name and password
func (z *Zip) WithPassword(r io.Reader, fileName, password string) (*bytes.Reader, error) {
	var (
		buf  = bytes.Buffer{}
		zipw = zip.NewWriter(&buf)
	)

	if filepath.Ext(fileName) == "" {
		return nil, errors.New("missing extension on file name")
	}

	wr, err := zipw.Encrypt(fileName, password, zip.StandardEncryption)
	if err != nil {
		return nil, err
	}
	io.Copy(wr, r)
	zipw.Close()
	return bytes.NewReader(buf.Bytes()), nil
}

// Unlock takes a src reader and a password, it then gives back a byte reader with all those files.
func (z *Zip) Unlock(r io.Reader, password string) (*bytes.Reader, error) {
	var (
		result = bytes.Buffer{}
		bs, _  = ioutil.ReadAll(r)
		buf    = bytes.NewReader(bs)
	)
	ziprd, err := zip.NewReader(buf, buf.Size())
	if err != nil {
		return nil, err
	}

	for _, file := range ziprd.File {
		if file.IsEncrypted() {
			file.SetPassword(password)
		}
		frd, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer frd.Close()

		_, err = io.Copy(&result, frd)
		if err != nil {
			return nil, err
		}
	}

	return bytes.NewReader(result.Bytes()), nil
}
