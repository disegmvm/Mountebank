package main

import (
	"bytes"
	"context"
	"go/mountebank"
	"go/mountebank/stubs"
	"log"
	"net/http"
	"testing"
)

var (
	ctx       = context.Background()
	testingID = "1001"
)

func TestTransform(t *testing.T) {

	// Initializing Mountebank's components
	mbank, err := mountebank.NewMountebank()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	err = mbank.CreateImposter(ctx, mountebank.Imposter)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	err = mbank.AddStub(ctx, mountebank.Imposter.Port, stubs.TransformStub)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	// Creating a testing data
	payload1 := "{\"Key1\": \"initial value 1\", \"Key2\": \"initial value 2\"}"

	// Sending the request to 'transform' service
	request, _ := http.NewRequest("POST", "http://localhost:8080/transform", bytes.NewBuffer([]byte(payload1)))
	request.Header.Add("Testing-Id", testingID)
	httpClient.Do(request)
	log.Printf("Payload sent to local service:\n%s", payload1)

	// Retrieving all requests present in Mountebank
	reqs, _ := mbank.GetRequestsFromMountebank(ctx, 8181)

	// Looping through Mountebank requests to find the one with knows Testing ID
	for _, req := range reqs {
		if req.Headers.Get("Testing-Id") == testingID {
			log.Printf("Request from Mountebank:\n%s", req.Body.(string))
		}
	}
}
