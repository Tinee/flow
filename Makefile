.PHONY: build


build:
	GOOS=linux GOARCH=amd64 go build -o ./bin/create-workspace ./cmd/workspace/create 