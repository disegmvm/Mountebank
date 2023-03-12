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
	testingID = "17"
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
	log.Printf("Payload sent to local service:\n%s", payload1)

	// Preparing a request to send it to 'transform' service
	request, err := http.NewRequest("POST", "http://localhost:8080/transform", bytes.NewBuffer([]byte(payload1)))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	request.Header.Add("Testing-Id", testingID)

	// Sending prepared request
	httpClient := &http.Client{}
	_, err = httpClient.Do(request)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	// Retrieving all requests present in Mountebank
	reqs, err := mbank.GetRequestsFromMountebank(ctx, 8181)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	// Looping through Mountebank requests to find the one with knows Testing ID
	for _, req := range reqs {
		if req.Headers.Get("Testing-Id") == testingID {
			log.Printf("Request from Mountebank:\n%s", req.Body.(string))
		}
	}
}
