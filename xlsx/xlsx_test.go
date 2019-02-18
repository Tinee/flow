package xlsx_test

import (
	"lumo"
	"lumo/xlsx"
	"testing"
)

var XLSXEncodeTests = []struct {
	name         string
	flows        lumo.Flows
	resultLength int
	fileName     string
	wantErr      bool
}{
	{
		"should give me 1 flow.",
		lumo.Flows{{Type: "SAR", Fields: []string{"SAR"}}},
		1,
		"test.xlsx",
		false,
	},
	{
		"should give me two flows.",
		lumo.Flows{{Type: "SAR", Fields: []string{"SAR", "TESt"}}, {Type: "SAR", Fields: []string{"SAR", "Test"}}},
		2,
		"test.xlsx",
		false,
	},
	{
		"should give me two flows.",
		lumo.Flows{{Type: "SAR", Fields: []string{"SAR", "TEST"}}},
		0,
		"",
		true,
	},
}

func TestXLSX_Encode(t *testing.T) {
	xlsx := &xlsx.XLSX{}
	for _, tt := range XLSXEncodeTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := xlsx.Encode(tt.flows, tt.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("XLSX.Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			flows, err := xlsx.Decode(got)
			if err != nil {
				t.Errorf("XLSX.Encode(reader) got an error: %v", err)
			}
			if len(flows) != tt.resultLength {
				t.Errorf("XLSX.Encode(reader) expected result to have an length of %v but got %v", tt.resultLength, len(flows))
			}
		})
	}
}

func BenchmarkXLSXEncode(b *testing.B) {
	xlsx := &xlsx.XLSX{}
	for i := 0; i < b.N; i++ {
		for _, test := range XLSXEncodeTests {
			xlsx.Encode(test.flows, test.fileName)
		}
	}
}
