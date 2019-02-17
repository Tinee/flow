package flow

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/tealeg/xlsx"

	"github.com/yeka/zip"
)

type FlowRow []string

// Errors within the flow package.
var (
	ErrInvalidPassword = errors.New("invalid password to the zip file")
	ErrInvalidFileType = errors.New("not valid zip file in this program")
)

// ZipXLSXWithPassword writes a zip file which is encrypted with the passed in password.
func ZipXLSXWithPassword(w io.Writer, row FlowRow, name, password string) error {

	zipw := zip.NewWriter(w)
	defer zipw.Close()
	wr, err := zipw.Encrypt(name+".zip", password, zip.AES256Encryption)
	io.Copy(wr, &buf)
	return zipw.Flush()
}

// UnZipXLSXWithPassword reads from r and decrypts the zip file with the password.
// if everything goes well it will parse all the rows inside the underlaying excel file.
func UnZipXLSXWithPassword(r io.Reader, password string) (FlowRow, error) {
	buf, _ := ioutil.ReadAll(r)
	rd, err := decodeZip(bytes.NewReader(buf), password)
	if err != nil {
		return nil, err
	}
	return parseFlowRow(rd)
}

func decodeZip(r *bytes.Reader, password string) (*bytes.Reader, error) {
	ziprd, err := zip.NewReader(r, r.Size())
	if err != nil {
		return nil, fmt.Errorf("not a valid zip file. Error: %s", err)
	}

	if len(ziprd.File) == 0 {
		return nil, fmt.Errorf("zip file doesn't contain any files. Error: %s", err)
	}
	file := ziprd.File[0]
	if file.IsEncrypted() {
		file.SetPassword(password)
	}

	frd, err := file.Open()
	if err != nil {
		if err == zip.ErrAuthentication {
			return nil, ErrInvalidPassword
		}
		return nil, err
	}

	defer frd.Close()

	buf, err := ioutil.ReadAll(frd)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(buf), nil
}

func parseFlowRow(r *bytes.Reader) (FlowRow, error) {
	var result []string
	f, err := xlsx.OpenReaderAt(r, r.Size())
	if err != nil {
		return result, err
	}
	for _, sheet := range f.Sheets {
		for _, row := range sheet.Rows {
			for _, cell := range row.Cells {
				result = append(result, cell.Value)
			}
		}
	}
	return result, nil
}

func buildXLSXSheet(row FlowRow, name string) (*bytes.Buffer, error) {
	var (
		buf  = bytes.Buffer{}
		file = xlsx.NewFile()
	)

	s, err := file.AddSheet(name)
	if err != nil {
		return nil, err
	}
	r := s.AddRow()
	r.WriteSlice(&row, len(row))

	file.Write(&buf)
}
