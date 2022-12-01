package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title" binding:"required"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type BindingErrorMsg struct {
	Field   string `json:"field" validate:"required"`
	Message string `json:"message" validate:"required"`
}

type ServerError struct {
	BindingErrors []*BindingErrorMsg `json:"errors"`
	Message       string             `json:"message"`
}

var albums = []Album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func listAlbums() []Album {
	return albums
}

func resetAlbums() {
	albums = []Album{
		{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
		{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	}
}

func getAlbums(c *gin.Context) {
	c.JSON(http.StatusOK, albums)
}

func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	for _, album := range albums {
		if album.ID == id {
			c.JSON(http.StatusOK, album)
			return
		}
	}
	serverError := ServerError{Message: "album not found"}
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": serverError.Message})
}

func postAlbum(c *gin.Context) {
	var newAlbum Album

	if err := c.ShouldBindJSON(&newAlbum); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			bindingErrorMessages := make([]BindingErrorMsg, len(ve))
			for i, fe := range ve {
				bindingErrorMessages[i] = BindingErrorMsg{Field: fe.Field(), Message: getErrorMsg(fe)}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": bindingErrorMessages})
			return
		}
	}
	albums = append(albums, newAlbum)
	c.JSON(http.StatusCreated, newAlbum)
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is Required"
	}
	return "Unknown Error"
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbum)
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
