package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go/helpers"
	"io"
	"log"
	"net/http"
	"strings"
)

// album represents data about a record album.
type car struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Color string `json:"color"`
}

var cars = []car{
	{ID: "1", Title: "BMW", Color: "Black"},
	{ID: "2", Title: "Tesla", Color: "Red"},
}

type Post struct {
	Userid string `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func main() {
	router := gin.Default()
	router.GET("/cars", getCars)
	router.GET("/cars/:id", getCarByID)
	router.POST("/cars", post3)

	router.Run("localhost:8080")
}

// postAlbums adds an album from JSON received in the request body.
func post2(c *gin.Context) {
	var newCar car
	if err := c.BindJSON(&newCar); err != nil {
		return
	}
	cars = append(cars, newCar)
	jsonchik := helpers.ToJSON(newCar)
	resp, err := http.Post("http://localhost:4545/test", "application/json", strings.NewReader(jsonchik))

	if err != nil {
		log.Printf("Request Failed: %s", err)
		return
	}

	log.Printf("SENT SUCCESSFULLY!!!!!!!: %s", resp.StatusCode)
}

func post3(c *gin.Context) {
	var newCar car
	if err := c.BindJSON(&newCar); err != nil {
		return
	}
	cars = append(cars, newCar)
	body, _ := json.Marshal(newCar)

	// Create a HTTP post request
	r, err := http.NewRequest("POST", "http://localhost:4545/test", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Testing-Id", "777")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	log.Printf("SENT SUCCESSFULLY!!!!!!!: %s", res.StatusCode)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)
}

// getCars responds with the list of all albums as JSON.
func getCars(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, cars)
}

// getCarByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getCarByID(c *gin.Context) {
	id := c.Param("id")

	// Loop through the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range cars {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
