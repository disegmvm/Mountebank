package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go/helpers"
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

func main() {
	router := gin.Default()
	router.GET("/cars", getCars)
	router.GET("/cars/:id", getCarByID)
	router.POST("/cars", post2)

	router.Run("localhost:8080")
}

func postCars(c *gin.Context) {
	var newCar car
	if err := c.BindJSON(&newCar); err != nil {
		return
	}
	cars = append(cars, newCar)
	jsonchik := helpers.ToJSON(newCar)

	baseUrl := "api-url.com"
	path := "/v1"
	aHttpHelper := helpers.NewHttpHelper(baseUrl)

	_, _ = aHttpHelper.Send(
		aHttpHelper.POST(path).
			WithContentTypeHeaderJSON().
			WithBody(strings.NewReader(jsonchik)))
	//Expect(httpErr).NotTo(HaveOccurred())
	//Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	mes := fmt.Sprintf("POST is made to %s%s", baseUrl, path)

	c.IndentedJSON(http.StatusCreated, gin.H{"message": mes})
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
