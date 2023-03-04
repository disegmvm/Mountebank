package main

import (
	"bytes"
	"encoding/json"
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
	router.Run("localhost:8080")
}

// 'transform' function updates received payload and sends it further
func transform(receivedRequest *gin.Context) {
	var payload payload
	receivedRequest.BindJSON(&payload)

	// Transforming the payload
	payload.Key1 = "transformed value 1"
	payload.Key2 = "transformed value 2"

	// Sending the response with transformed payload
	receivedRequest.IndentedJSON(http.StatusCreated, payload)

	// Converting transformed payload to JSON
	jsonPayload, _ := json.Marshal(payload)

	// Preparing a request to send it further to external service
	request, _ := http.NewRequest("POST", mountebankUrl, bytes.NewBuffer(jsonPayload))
	request.Header.Add("Testing-Id", receivedRequest.GetHeader("Testing-Id"))

	// Sending prepared request
	httpClient := &http.Client{}
	httpClient.Do(request)
}
