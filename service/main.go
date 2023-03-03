package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// car represents data about each car's record.
type payload struct {
	Key1 string `json:"Key1"`
	Key2 string `json:"Key2"`
	Key3 string `json:"Key3"`
}

var mountebankUrl = "http://localhost:8181/transform"

func main() {
	router := gin.Default()
	router.POST("/transform", processCar)

	err := router.Run("localhost:8080")
	if err != nil {
		log.Printf("Run failed: %s", err)
		return
	}
}

// processCar updates received car's title and sends it to Mountebank
func processCar(receivedRequest *gin.Context) {

	var receivedPayload payload                                        // Declaring a new Car variable.
	if err := receivedRequest.BindJSON(&receivedPayload); err != nil { // Bind the request body to newCar variable.
		receivedRequest.IndentedJSON(http.StatusBadRequest, // Return 400 Status Code,
			gin.H{"message": "Failed to create a car"}) // and error message if binding has failed.
		return
	}
	//Transforming the payload
	receivedPayload.Key1 = "transformed key1 value"
	receivedPayload.Key2 = "transformed key2 value"
	receivedPayload.Key3 = "transformed key3 value"

	receivedRequest.IndentedJSON(http.StatusCreated, receivedPayload)
	transformedPayload, _ := json.Marshal(receivedPayload)

	request, err := http.NewRequest("POST", mountebankUrl, bytes.NewBuffer(transformedPayload))
	if err != nil {
		log.Printf("Request creation failed: %s", err)
		return
	}
	request.Header.Add("Testing-Id", receivedRequest.GetHeader("Testing-Id"))

	//Sending request to Mountebank
	httpClient := &http.Client{}
	_, err = httpClient.Do(request)
	if err != nil {
		log.Printf("Failed to send request: %s", err)
		return
	}
}
