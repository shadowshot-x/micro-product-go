//go:build integration
// +build integration

package tests

import (
	"github.com/jarcoal/httpmock"
	"github.com/shadowshot-x/micro-product-go/ordertransformerservice"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// we would again need http mocks to prevent external calls
// but here we will test our code end to end.
// We just want to test that the call runs from starting controller to the end http call
func TestHealthCheckHandler(t *testing.T) {
	log, _ := zap.NewProduction()
	defer log.Sync()

	transc := ordertransformerservice.NewTransformerController(log)
	transc.Store_json_dir = "../ordertransformerservice/json_store/"
	transc.Region_rules_dir = "../ordertransformerservice/region_rules/"

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://httpbin.org/post", func(req *http.Request) (*http.Response, error) {
		resp, err := httpmock.NewJsonResponse(200, "200 OK")
		if err != nil {
			return httpmock.NewStringResponse(500, "Error Generating Response"), nil
		}
		return resp, nil
	})

	req, err := http.NewRequest("GET", "/transformer/transform", nil)
	if err != nil {
		t.Fatalf("Got error in creating Http Test Request, %v", err)
	}

	outputCatcher := httptest.NewRecorder()
	h := http.HandlerFunc(transc.TransformerHandler)
	h.ServeHTTP(outputCatcher, req)

	outputBody := outputCatcher.Body

	want := "Files Parsed and Validated"
	if !strings.Contains(outputBody.String(), want) {
		t.Errorf("Got Unexpected Output in the HTTP Response. want: %v, got: %v\n", outputBody.String(), want)
	}
}
