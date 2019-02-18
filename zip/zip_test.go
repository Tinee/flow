package zip_test

import (
	"bytes"
	"io/ioutil"
	"lumo/zip"
	"os"
	"testing"
)

var unLockTestCases = []struct {
	name     string
	password string
	wantErr  bool
}{
	{
		"should be able to unlock the file.",
		"password",
		false,
	},
	{
		"should not be able to unlock the file with the wrong password.",
		"wrongPassword",
		true,
	},
}

func TestZip_Unlock(t *testing.T) {
	zipper := zip.New()

	for _, tt := range unLockTestCases {
		t.Run(tt.name, func(t *testing.T) {
			f := getTestZipFile(t)
			defer f.Close()
			out, err := zipper.Unlock(f, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Zip.Unlock(reader, %q) error = %v, wantErr %v", tt.password, err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			bs, _ := ioutil.ReadAll(out)
			if !bytes.Contains(bs, []byte("Some random values.")) {
				t.Errorf("Zip.Unlock(reader, %q) doesn't give back a document that contains those values.", tt.password)
				return
			}
		})
	}
}

func BenchmarkZipUnlock(b *testing.B) {
	zipper := zip.New()
	f, _ := os.Open("./testdata/test.zip")
	for i := 0; i < b.N; i++ {
		for _, test := range unLockTestCases {
			zipper.Unlock(f, test.password)
		}
	}
}

var zipWithPasswordTestCases = []struct {
	name     string
	password string
	fileName string
	want     string
	wantErr  bool
}{
	{
		"should be able to unlock the file.",
		"password",
		"text.txt",
		"Some random values.",
		false,
	},
	{
		"should not have the correct content.",
		"password",
		"",
		"Some random values.",
		true,
	},
	{
		"should not have the correct content.",
		"",
		"text.txt",
		"Some random values.",
		false,
	},
}

func TestZip_WithPassword(t *testing.T) {
	zipper := &zip.Zip{}
	for _, tt := range zipWithPasswordTestCases {
		t.Run(tt.name, func(t *testing.T) {
			f := getTextFile(t)
			defer f.Close()
			out, err := zipper.WithPassword(f, tt.fileName, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Zip.WithPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			file, _ := zipper.Unlock(out, tt.password)
			bs, _ := ioutil.ReadAll(file)
			if !bytes.Contains(bs, []byte(tt.want)) {
				t.Errorf("Zip.Unlock(reader, %q, %q) doesn't give back a document that contains those values.", tt.fileName, tt.password)
			}
		})
	}
}

func BenchmarkZipWithPassword(b *testing.B) {
	zipper := zip.New()
	f, _ := os.Open("./testdata/test.txt")
	for i := 0; i < b.N; i++ {
		defer f.Close()
		for _, test := range zipWithPasswordTestCases {
			zipper.WithPassword(f, test.fileName, test.password)
		}
	}
}

func getTestZipFile(t *testing.T) *os.File {
	path := "./testdata/test.zip"
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("couldn't open file at %s. Error: %s", path, err)
	}
	return f
}

func getTextFile(t *testing.T) *os.File {
	path := "./testdata/test.txt"
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("couldn't open file at %s. Error: %s", path, err)
	}
	return f
}
