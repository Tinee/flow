.PHONY: build

build:
	GOOS=linux GOARCH=amd64 go build -o ./bin/decode_zip_xlsx_flow ./cmd/flow/zip_xlsx