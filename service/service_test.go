package main

import (
	"context"
	"go/mountebank"
	"go/mountebank/stubs"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/senseyeio/mbgo"
)

//var ctx context.Context

func Test_ServiceTest(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{}

	mbank, err := mountebank.NewMountebank()
	if err != nil {
		t.Fatalf("Error on creating mountebank: %v", err)
	}
	err = mbank.CreateImposters(ctx, []mbgo.Imposter{mountebank.Imposter})
	if err != nil {
		panic("create imposter error")
	}
	err = mbank.AddStub(ctx, mountebank.Imposter.Port, stubs.CarsStub)
	if err != nil {
		panic("add stub error")
	}

	/*p := payload{key1: "asd",
		key2: "dsa",
		key3: "KEY Z"}
	pp, _ := json.Marshal(p)*/

	payload1 := "{\"key1\": \"initial value 1\", \"key2\": \"initial value 2\", \"key3\": \"initial value 3\"}"
	testingID := "16662"

	//log.Printf("JSON, that will be sent to local service:\n%s", pp)

	request, err := http.NewRequest("POST", "http://localhost:8080/transform", strings.NewReader(payload1))
	if err != nil {
		log.Printf("Request creation failed: %s", err)
		return
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Testing-Id", testingID)

	log.Printf("SENDING......")
	_, err = client.Do(request)
	if err != nil {
		log.Printf("Failed to send request: %s", err)
		return
	}
	log.Print("Sent")

	var requestFromMBank string
	reqs, _ := mbank.GetRecordedRequestsByPath(ctx, 8181, "all-paths")
	for _, req := range reqs {
		if req.Headers.Get("Testing-Id") == "16662" {
			requestFromMBank = req.Body.(string)
		}
	}

	log.Print("Request from Mountebank:")
	log.Print(requestFromMBank)
}
