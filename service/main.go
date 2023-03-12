package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 'payload' represents data accepted by transform function
type payload struct {
	Key1 string `json:"Key1"`
	Key2 string `json:"Key2"`
}

var mountebankUrl = "http://localhost:8181/transform"

func main() {
	router := gin.Default()
	router.POST("/transform", transform)
	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

// 'transform' function updates received payload and sends it further
func transform(receivedRequest *gin.Context) {
	var payload payload
	err := receivedRequest.BindJSON(&payload)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Transforming the payload
	payload.Key1 = "transformed value 1"
	payload.Key2 = "transformed value 2"

	// Sending the response with transformed payload
	receivedRequest.IndentedJSON(http.StatusCreated, payload)

	// Converting transformed payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Preparing a request to send it further to external service
	request, err := http.NewRequest("POST", mountebankUrl, bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	request.Header.Add("Testing-Id", receivedRequest.GetHeader("Testing-Id"))

	// Sending prepared request
	httpClient := &http.Client{}
	_, err = httpClient.Do(request)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
