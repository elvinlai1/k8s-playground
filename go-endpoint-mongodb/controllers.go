package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Movie struct {
	ID       string `json:"_id,omitempty" bson:"_id,omitempty"`
	Title    string `json:"title,omitempty" bson:"title,omitempty"`
	Year     int32  `json:"year,omitempty" bson:"year,omitempty"`
	Rated    string `json:"rated,omitempty" bson:"rated,omitempty"`
	Released string `json:"released,omitempty" bson:"released,omitempty"`
}

var coll = mongoDB().Database("sample_mflix").Collection("movies")

func insertData(ctx *gin.Context) {
	// Insert Titles as IDs as well to make it easier to query and read by title
	docs := []interface{}{
		Movie{ID: "Back to the Future", Title: "Back to the Future", Year: 1985, Rated: "PG", Released: "03 Jul 1985"},
		Movie{ID: "Back to the Future Part II", Title: "Back to the Future Part II", Year: 1989, Rated: "PG", Released: "22 Nov 1989"},
		Movie{ID: "Back to the Future Part III", Title: "Back to the Future Part III", Year: 1990, Rated: "PG", Released: "25 May 1990"},
		Movie{ID: "The Terminator", Title: "The Terminator", Year: 1984, Rated: "R", Released: "26 Oct 1984"},
		Movie{ID: "Terminator 2: Judgement Day", Title: "Terminator 2: Judgment Day", Year: 1991, Rated: "R", Released: "03 Jul 1991"},
		Movie{ID: "Terminator 3: Rise of the Machines", Title: "Terminator 3: Rise of the Machines", Year: 2003, Rated: "R", Released: "02 Jul 2003"},
		Movie{ID: "Terminator Salvation", Title: "Terminator Salvation", Year: 2009, Rated: "PG-13", Released: "21 May 2009"},
	}

	result, err := coll.InsertMany(context.TODO(), docs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert movies"})
	}

	for _, id := range result.InsertedIDs {
		ctx.JSON(http.StatusOK, gin.H{"message": "Movies added successfully", "insertedID": id})
	}

}

func addMovieByTitle(ctx *gin.Context) {
	var movie Movie
	if err := ctx.BindJSON(&movie); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use the movie's title as its ID
	movie.ID = movie.Title

	insertResult, err := coll.InsertOne(context.Background(), movie)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert movie"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Movie added successfully", "insertedID": insertResult.InsertedID})
}

func addMovie(ctx *gin.Context) {
	var movie Movie
	if err := ctx.BindJSON(&movie); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	insertResult, err := coll.InsertOne(context.Background(), movie)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert movie"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Movie added successfully", "insertedID": insertResult.InsertedID})

}

func getMovie(ctx *gin.Context) {
	id := ctx.Param("id")
	idprim, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"\nerror": "Document not found\n"})
		return
	}

	var result Movie
	err = coll.FindOne(context.TODO(), bson.M{"_id": idprim}).Decode(&result)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"\nerror": "Document not found\n"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Found a single document", "data": result})
}

func getMovieByTitle(ctx *gin.Context) {
	id := ctx.Param("id")

	var result Movie
	err := coll.FindOne(context.TODO(), Movie{ID: id}).Decode(&result)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"\nerror": "Document not found\n"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Found a single document", "data": result})

}

func getAllMovies(ctx *gin.Context) {
	filter := bson.D{{}} // An empty filter matches all documents in the collection

	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"\nerror": err})
	}
	defer cursor.Close(context.Background())

	var movies []Movie
	if err = cursor.All(context.Background(), &movies); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"\nerror": err})
	}

	ctx.JSON(http.StatusOK, movies)

}

func updateMovie(ctx *gin.Context) {
	id := ctx.Param("id")

	var movie Movie
	if err := ctx.BindJSON(&movie); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.D{}
	// checks if fields are empty
	// no loops as dataset is tiny and we know the fields
	if movie.Title != "" {
		update = append(update, bson.E{Key: "title", Value: movie.Title})
	}
	if movie.Year != 0 {
		update = append(update, bson.E{Key: "year", Value: movie.Year})
	}
	if movie.Rated != "" {
		update = append(update, bson.E{Key: "rated", Value: movie.Rated})
	}
	if movie.Released != "" {
		update = append(update, bson.E{Key: "released", Value: movie.Released})
	}

	if len(update) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	filter := bson.D{{Key: "_id", Value: id}}
	result, err := coll.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Movie updated", "data": result})
}

func deleteMovie(ctx *gin.Context) {
	id := ctx.Param("id")

	filter := bson.D{{Key: "_id", Value: id}}

	// checks if exists
	var movie Movie
	err := coll.FindOne(context.TODO(), filter).Decode(&movie)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Movie not found"})
		return
	}

	result, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	// check if DeletedCount is 0 as deleteOne will return a success even if no document was deleted
	if result.DeletedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No movie found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Movie deleted", "data": result})

}

func removeAllMovies(ctx *gin.Context) {
	filter := bson.D{{}} // An empty filter matches all documents in the collection
	_, err := coll.DeleteMany(context.Background(), filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movies"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "All movies deleted successfully"})
}
