package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// car represents data about each car's record.
type car struct {
	Year  string
	Title string
	Color string
}

var mountebankUrl = "http://localhost:8181/test/cars"

func main() {
	router := gin.Default()
	router.POST("/cars", processCar)

	err := router.Run("localhost:8080")
	if err != nil {
		log.Printf("Run failed: %s", err)
		return
	}
}

// processCar updates received car's title and sends it to Mountebank
func processCar(receivedRequest *gin.Context) {
	var receivedCar car
	if err := receivedRequest.BindJSON(&receivedCar); err != nil {
		return
	}
	//Updating the Title
	receivedCar.Title = "Brand new title"
	newCarJson, _ := json.Marshal(receivedCar)
	log.Printf("XXXX new car:")
	log.Printf(string(newCarJson))

	//Creating HTTP POST request
	request, err := http.NewRequest("POST", mountebankUrl, bytes.NewBuffer(newCarJson))
	if err != nil {
		log.Printf("Request creation failed: %s", err)
		return
	}
	request.Header.Add("Testing-Id", receivedRequest.GetHeader("Testing-Id"))

	log.Print("OOOOOOOOOOOOOOtpravluyau B MBANK")
	//Sending request to Mountebank
	httpClient := &http.Client{}
	_, err = httpClient.Do(request)
	if err != nil {
		log.Printf("Failed to send request: %s", err)
		return
	}

	log.Printf("otpravil, bb, staraya tachka:")
	receivedRequest.IndentedJSON(http.StatusCreated, receivedCar)
}
