package xlsx

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"flow"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/tealeg/xlsx"
)

// XLSX is the type that can work with .xlsx files.
type XLSX struct{}

// New creates a XLSX structure.
func New() *XLSX {
	return &XLSX{}
}

// Decode takes a src reader and attempts to make it into a domain Flow.
func (x *XLSX) Decode(r io.Reader) (flow.Flows, error) {
	var (
		result = flow.Flows{}
		bs, _  = ioutil.ReadAll(r)
		buf    = bytes.NewReader(bs)
	)
	f, err := xlsx.OpenReaderAt(buf, buf.Size())
	if err != nil {
		return result, err
	}

	for _, sheet := range f.Sheets {
		for _, row := range sheet.Rows {
			if len(row.Cells) == 0 {
				continue
			}
			currentFlow := flow.Flow{Type: row.Cells[0].String()}

			for _, cell := range row.Cells {
				currentFlow.Fields = append(currentFlow.Fields, cell.String())
			}
			result = append(result, currentFlow)
		}
	}

	return result, nil
}

// Encode takes flows, and tries to write that to an xlsx file row by row.
// EX. [
// { Type:"SAR",Fields:["SAR","Something"]},
// { Type:"SAR",Fields:["SAR","Something"]}
//] -> xlsx file with a sheet and two rows, one for each element in flows.
func (x *XLSX) Encode(flows flow.Flows, name string) (io.Reader, error) {
	var (
		f   = xlsx.NewFile()
		buf = bytes.Buffer{}
		// tempName = generateTempFileName()
	)
	s, err := f.AddSheet(name)
	if err != nil {
		return nil, err
	}
	for _, flow := range flows {
		row := s.AddRow()
		for _, field := range flow.Fields {
			row.AddCell().SetValue(field)
		}
	}

	// err = f.Save("asd.xlsx")
	// if err != nil {
	// 	return nil, err
	// }

	// tempFile, err := os.OpenFile(tempName, os.O_RDWR, os.ModePerm)
	// defer os.Remove(tempFile.Name())
	// if err != nil {
	// 	return nil, err
	// }
	f.Write(&buf)
	return &buf, nil
}

func generateTempFileName() string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), hex.EncodeToString(randBytes)) + ".xlsx"
}
