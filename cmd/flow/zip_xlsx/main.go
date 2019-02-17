package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"lumo/pkg/flow"
	"lumo/xlsx"
	"lumo/zip"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type decodeZIPFunction struct {
	s *flow.Service
}

func (f *decodeZIPFunction) handler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	password := r.QueryStringParameters["password"]
	bs, err := base64.StdEncoding.DecodeString(r.Body)

	flow, err := f.s.DecodeZipFile(bytes.NewReader(bs), password)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			Headers:    map[string]string{"Content-Type": "plain/text"},
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	bs, _ = json.Marshal(&flow)
	return events.APIGatewayProxyResponse{
		Body:       string(bs),
		Headers:    map[string]string{"Content-Type": "application/json"},
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	zipper := zip.New()
	decoder := xlsx.New()
	s := flow.NewService(zipper, decoder)
	lambdaFunction := decodeZIPFunction{s}
	lambda.Start(lambdaFunction.handler)

	// flow, err := s.DecodeZipFile(f, "biscuit")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// r, err := s.EncodeFlowToZip(flow, "walla.xlsx", "Testar")
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// io.Copy(f2, r)
}
