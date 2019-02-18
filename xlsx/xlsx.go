package xlsx

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"lumo"

	"github.com/tealeg/xlsx"
)

// XLSX is the type that can work with .xlsx files.
type XLSX struct{}

// New creates a XLSX structure.
func New() *XLSX {
	return &XLSX{}
}

// Decode takes a src reader and attempts to make it into a domain Flow.
func (x *XLSX) Decode(r io.Reader) (lumo.Flows, error) {
	var (
		result = lumo.Flows{}
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
			currentFlow := lumo.Flow{Type: row.Cells[0].String()}

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
func (x *XLSX) Encode(flows lumo.Flows, name string) (io.Reader, error) {
	var (
		f   = xlsx.NewFile()
		buf = bytes.Buffer{}
	)
	if name == "" {
		return nil, errors.New("invalid argument")
	}

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
	f.Write(&buf)
	return &buf, nil
}
