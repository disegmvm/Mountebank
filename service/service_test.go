package main

import (
	"context"
	"github.com/senseyeio/mbgo"
	"go/mountebank"
	"go/mountebank/stubs"
	"log"
	"net/http"
	"strings"
	"testing"
)

var ctx context.Context

func Test_ServiceTest(t *testing.T) {
	ctx = context.Background()

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
	//log.Printf("всё, я спать нахуй.... (2 секунды)")
	//time.Sleep(time.Second * time.Duration(2))

	testingID := "1172"
	newCar := "{\"Year\": \"1998\", \"Title\": \"GM\", \"Color\": \"PINK\"}"

	log.Printf("JSON, that will be sent to local service:\n%s", newCar)

	log.Printf("готовлю к отправке её в http://localhost:8080/cars")
	request, err := http.NewRequest("POST", "http://localhost:8080/cars", strings.NewReader(newCar))
	if err != nil {
		log.Printf("Request creation failed: %s", err)
		return
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Testing-Id", testingID)

	log.Printf("SENDING......")
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		log.Printf("Failed to send request: %s", err)
		return
	}
	log.Print(res)
	//log.Printf("Status code: %d", res)

	log.Printf("ща буду лупаться") //------------------------------------FIX FROM HERE-----------------------------------
	for _, request := range mountebank.Imposter.Requests {
		log.Printf("зашел в Requests...")
		if request.(*mbgo.HTTPRequest).Headers.Get("Testing-Id") == testingID {
			log.Printf("Payload received by Mountebank using TestingID %s:\n%s", testingID, request.(*mbgo.HTTPRequest).Body)
		} else {
			log.Printf("нихуя не нашел")
		}
	}
}
