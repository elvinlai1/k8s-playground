package main

import (
	"github.com/gin-gonic/gin"
)

type album struct {
	ID     string  `json:"id" bson:"_id"`
	Title  string  `json:"title" bson:"title"`
	Artist string  `json:"artist" bson:"artist"`
	Price  float64 `json:"price" bson:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {

	// Set up Gin router
	router := gin.Default()
	router.GET("/", getAlbums)
	router.Run(":80")
}

func getAlbums(c *gin.Context) {
	var result []album
	result = append(result, albums...)
	c.IndentedJSON(200, result)
}
