package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	//Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Run the other tests
	os.Exit(m.Run())
}

func Test_PostBadJson(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()

	albums := `{"XID": "10", "Titlexx": "Blue Train", "Artistx": "John Coltrane", "Price": 56.99, "X": "asdf"}`
	var serverError Error
	req, _ := http.NewRequest(http.MethodPost, "/albums", strings.NewReader(albums))

	router.ServeHTTP(w, req)

	body := w.Body.Bytes()
	if err := json.Unmarshal(body, &serverError); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			fmt.Println(ve)
		}
		assert.Failf(t, "unmarshalling error message from server problem", ve.Error())
	}

	assert.Equal(t, 1, len(serverError.Errors))
	assert.Equal(t, "Title", serverError.Errors[0].Field)
	assert.Equal(t, "This field is Required", serverError.Errors[0].Message)

}

func Test_postAlbums(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()

	albums := `{"ID": "10", "Title": "War Pigs", "Artist": "Black Sabbath", "Price": 56.99}`
	expectedAlbum := album{ID: "10", Title: "War Pigs", Artist: "Black Sabbath", Price: 56.99}
	req, _ := http.NewRequest(http.MethodPost, "/albums", strings.NewReader(albums))

	router.ServeHTTP(w, req)

	var album album
	err := json.Unmarshal(w.Body.Bytes(), &album)
	if err != nil {
		assert.Empty(t, err, "json marshaling failed")
	}

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, album, expectedAlbum)
}

func Test_getAlbums(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/albums", nil)

	router.ServeHTTP(w, req)
	var albums []album

	err := json.Unmarshal(w.Body.Bytes(), &albums)
	if err != nil {
		assert.Failf(t, "json unmarshal fail", "fail to unmarshall Albums")
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, listAlbums(), albums)

}
