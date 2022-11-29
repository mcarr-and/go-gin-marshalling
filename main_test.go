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

func Test_getAllAlbums(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/albums", nil)

	router.ServeHTTP(w, req)
	var albums []Album

	if err := json.Unmarshal(w.Body.Bytes(), &albums); err != nil {
		assert.Failf(t, "json unmarshal fail", "fail to unmarshall Albums")
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, listAlbums(), albums)
}

func Test_getAlbumById(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/albums/2", nil)

	router.ServeHTTP(w, req)
	var album Album

	if err := json.Unmarshal(w.Body.Bytes(), &album); err != nil {
		assert.Failf(t, "json unmarshal fail", "fail to unmarshall Albums")
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, listAlbums()[1], album)
	assert.Equal(t, listAlbums()[1].Title, album.Title)
}

func Test_getAlbumById_NotFound(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/albums/5666", nil)

	router.ServeHTTP(w, req)
	var serverError ServerError
	if err := json.Unmarshal(w.Body.Bytes(), &serverError); err != nil {
		assert.Empty(t, err, "json marshaling failed")
	}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "album not found", serverError.Message)
}

func Test_postAlbum(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()

	albumBody := `{"ID": "10", "Title": "The Ozzman Cometh", "Artist": "Black Sabbath", "Price": 56.99}`
	expectedAlbum := Album{ID: "10", Title: "The Ozzman Cometh", Artist: "Black Sabbath", Price: 56.99}
	req, _ := http.NewRequest(http.MethodPost, "/albums", strings.NewReader(albumBody))

	router.ServeHTTP(w, req)

	var album Album
	if err := json.Unmarshal(w.Body.Bytes(), &album); err != nil {
		assert.Empty(t, err, "json marshaling failed")
	}

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, album, expectedAlbum)
}

func Test_PostAlbum_BadRequest_BadJSON(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()

	album := `{"XID": "10", "Titlexx": "Blue Train", "Artistx": "John Coltrane", "Price": 56.99, "X": "asdf"}`
	var serverError ServerError
	req, _ := http.NewRequest(http.MethodPost, "/albums", strings.NewReader(album))

	router.ServeHTTP(w, req)

	body := w.Body.Bytes()
	fmt.Println(string(body))

	if err := json.Unmarshal(body, &serverError); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			fmt.Println(ve)
		}
		assert.Failf(t, "did not get error message from server", ve.Error())
	}
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, 1, len(serverError.BindingErrors))
	assert.Equal(t, "Title", serverError.BindingErrors[0].Field)
	assert.Equal(t, "This field is Required", serverError.BindingErrors[0].Message)
}

func BenchmarkGetAllAlbums(b *testing.B) {
	router := setupRouter()
	var albums []Album
	req, _ := http.NewRequest(http.MethodGet, "/albums", nil)
	w := httptest.NewRecorder()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
		if err := json.Unmarshal(w.Body.Bytes(), &albums); err != nil {
			assert.Fail(b, "should not fail", "GetAlbum", err)
		}
		w.Body.Reset()
		assert.Equal(b, http.StatusOK, w.Code)
	}
}

func BenchmarkGetAlbumById(b *testing.B) {
	router := setupRouter()
	var album Album
	req, _ := http.NewRequest(http.MethodGet, "/albums/2", nil)
	w := httptest.NewRecorder()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
		if err := json.Unmarshal(w.Body.Bytes(), &album); err != nil {
			assert.Fail(b, "should not fail", "GetAlbum", err)
		}
		assert.Equal(b, http.StatusOK, w.Code)
		assert.Equal(b, listAlbums()[1].Title, album.Title)
		w.Body.Reset()
	}
}

func BenchmarkGetAlbumById_BadRequest(b *testing.B) {
	router := setupRouter()
	req, _ := http.NewRequest(http.MethodGet, "/albums/5666", nil)
	w := httptest.NewRecorder()
	var serverError ServerError
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
		if err := json.Unmarshal(w.Body.Bytes(), &serverError); err != nil {
			assert.Fail(b, "should not fail", "GetAlbumById", err)
		}
		assert.Equal(b, http.StatusBadRequest, w.Code)
		assert.Equal(b, "album not found", serverError.Message)
		w.Body.Reset()
	}
}

func BenchmarkPostAlbum(b *testing.B) {
	router := setupRouter()
	album := `{"ID": "10", "Title": "The Ozzman Cometh", "Artist": "Black Sabbath", "Price": 56.99}`
	expectedAlbum := Album{ID: "10", Title: "The Ozzman Cometh", Artist: "Black Sabbath", Price: 56.99}
	w := httptest.NewRecorder()

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest(http.MethodPost, "/albums", strings.NewReader(album))
		router.ServeHTTP(w, req)
		body := w.Body.Bytes()

		var album Album
		if err := json.Unmarshal(body, &album); err != nil {
			assert.Empty(b, err, "json marshaling failed")
		}
		assert.Equal(b, http.StatusCreated, w.Code)
		assert.Equal(b, album, expectedAlbum)
		w.Body.Reset()
		w.Flush()
	}
}

func BenchmarkPostAlbum_BadRequest_BadJson(b *testing.B) {
	router := setupRouter()
	album := `{"XID": "10", "Titlexx": "Blue Train", "Artistx": "John Coltrane", "Price": 56.99, "X": "asdf"}`
	w := httptest.NewRecorder()

	for i := 0; i < b.N; i++ {
		var serverError ServerError
		req, _ := http.NewRequest(http.MethodPost, "/albums", strings.NewReader(album))
		router.ServeHTTP(w, req)
		body := w.Body.Bytes()
		assert.Equal(b, http.StatusBadRequest, w.Code)
		if err := json.Unmarshal(body, &serverError); err != nil {
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				fmt.Println(ve)
			}
			assert.Failf(b, "did not receive validation errors from server", ve.Error())
		}
		w.Body.Reset()
		w.Flush()
	}
}
