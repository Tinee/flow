package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Tinee/flow/flow"
)

type app struct{}
type response struct {
	FlowRow []string `json:"flowRow"`
}

func (a app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	password := r.URL.Query()["password"][0]
	defer r.Body.Close()
	rows, err := flow.UnZipXLSXWithPassword(r.Body, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	bs, _ := json.Marshal(&rows)
	w.Header().Add("Content-Type", "application-json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func main() {
	bs := bytes.Buffer{}
	err := flow.ZipXLSXWithPassword(&bs, []string{"Hej", "Da"}, "test", "wax")
	if err != nil {
		log.Fatalln(err)
	}
	s, err := flow.UnZipXLSXWithPassword(&bs, "wax")
	if err != nil {
		log.Fatalln(err)
	}
	for _, v := range s {
		fmt.Println(v)
	}
	// a := app{}
	// http.ListenAndServe(":3000", a)

}
