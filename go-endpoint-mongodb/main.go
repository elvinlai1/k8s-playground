package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.POST("/insertData", insertData) // load sample data
	router.POST("/movies/addByTitle", addMovieByTitle)
	router.POST("/movies/add", addMovie) // standard way of adding a document with a generated ObjectID

	router.GET("/movies/all", getAllMovies)
	router.GET("/movies/getByTitle/:id", getMovieByTitle)
	router.GET("/movie/:id", getMovie) // standard way of getting documents by ObjectID

	router.PUT("/movies/:id", updateMovie)

	router.DELETE("/movies/:id", deleteMovie)
	router.DELETE("/movies/all", removeAllMovies)

	err := router.Run(":8080")
	if err != nil {
		fmt.Println(err)
	}

}
