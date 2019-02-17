.PHONY: build

build:
	GOOS=linux GOARCH=amd64 go build -o ./bin/decode_zip_xlsx_flow ./cmd/flow/zip_xlsx
	GOOS=linux GOARCH=amd64 go build -o ./bin/encode_flow_zip_xlsx ./cmd/flow/encode_flow_xlsx