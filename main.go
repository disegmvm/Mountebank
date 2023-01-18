package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go/helpers"
	"net/http"
	"strings"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

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
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/cars", postCars)

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

// postAlbumsX responds with the list of all albums as JSON.
func postAlbumsX(c *gin.Context) {

	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}
	albums = append(albums, newAlbum)
	jsonchik := helpers.ToJSON(newAlbum)

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
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// getCars responds with the list of all albums as JSON.
func getCars(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, cars)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop through the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
