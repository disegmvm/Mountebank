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
	testingID = "12345"
)

func TestTransform(t *testing.T) {

	// Initializing Mountebank's components
	mbank, _ := mountebank.NewMountebank()
	_ = mbank.CreateImposter(ctx, mountebank.Imposter)
	_ = mbank.AddStub(ctx, mountebank.Imposter.Port, stubs.TransformStub)

	// Creating a testing data
	payload1 := "{\"Key1\": \"initial value 1\", \"Key2\": \"initial value 2\"}"
	log.Print("Payload sent to local service:")
	log.Print(payload1)

	// Preparing a request to send it to 'transform' service
	request, _ := http.NewRequest("POST", "http://localhost:8080/transform", bytes.NewBuffer([]byte(payload1)))
	request.Header.Add("Testing-Id", testingID)

	// Sending prepared request
	httpClient := &http.Client{}
	httpClient.Do(request)

	// Retrieving all requests present in Mountebank
	reqs, _ := mbank.GetRequestsFromMountebank(ctx, 8181)

	// Looping through Mountebank requests to find the one with knows Testing ID
	for _, req := range reqs {
		if req.Headers.Get("Testing-Id") == testingID {
			log.Print("Request from Mountebank:")
			log.Print(req.Body.(string))
		}
	}
}
