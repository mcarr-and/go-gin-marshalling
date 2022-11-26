package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title" binding:"required"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type ErrorMsg struct {
	Field   string `json:"field" validate:"required"`
	Message string `json:"message" validate:"required"`
}

type Error struct {
	Errors []*ErrorMsg `json:"errors" validate:"required"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func listAlbums() []album {
	return albums
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is Required"
	}
	return "Unknown Error"

}

func postAlbums(c *gin.Context) {
	var newAlbum album

	if err := c.ShouldBindJSON(&newAlbum); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = ErrorMsg{Field: fe.Field(), Message: getErrorMsg(fe)}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
		}
		return
	}
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	return router
}

func main() {
	router := setupRouter()
	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatal("Could not start server on port 8080")
		return
	}
}
