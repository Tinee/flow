package main

import (
	"encoding/json"
	"io/ioutil"
	"lumo"
	"lumo/pkg/flow"
	"lumo/xlsx"
	"lumo/zip"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type encodeZIPFunction struct {
	s *flow.Service
}

type requestBody struct {
	Password string     `json:"password"`
	Name     string     `json:"name"`
	Flows    lumo.Flows `json:"flows"`
}

func (r requestBody) validate() bool {
	if r.Name == "" {
		return false
	}
	if len(r.Flows) == 0 {
		return false
	}

	return true
}

func (f *encodeZIPFunction) handler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var request requestBody
	json.Unmarshal([]byte(r.Body), &request)
	if ok := request.validate(); !ok {
		return events.APIGatewayProxyResponse{
			Headers:    map[string]string{"Content-Type": "plain/text"},
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	rd, err := f.s.EncodeFlowToZip(request.Flows, request.Name+".xlsx", request.Password)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			Headers:    map[string]string{"Content-Type": "plain/text"},
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	bs, _ := ioutil.ReadAll(rd)

	return events.APIGatewayProxyResponse{
		Body:       string(bs),
		Headers:    map[string]string{"Content-Type": "application/octet-stream"},
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	zipper := zip.New()
	decoder := xlsx.New()
	s := flow.NewService(zipper, decoder)
	lambdaFunction := encodeZIPFunction{s}
	lambda.Start(lambdaFunction.handler)
}
