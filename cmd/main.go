package main

import (
	"io"
	"log"
	"lumo/pkg/flow"
	"lumo/xlsx"
	"lumo/zip"
	"os"
)

func main() {
	f, err := os.Open("Errors.zip")
	if err != nil {
		log.Fatalln(err)
	}
	f2, _ := os.Create("testar.zip")

	zipper := zip.New()
	excelDecoder := xlsx.New()
	s := flow.NewService(zipper, excelDecoder)
	flow, err := s.DecodeZipFile(f, "biscuit")
	if err != nil {
		log.Fatalln(err)
	}

	r, err := s.EncodeFlowToZip(flow, "walla.xlsx", "Testar")
	if err != nil {
		log.Fatalln(err)
	}
	io.Copy(f2, r)
}
