package main

import (
	"encoding/json"
	"errors"
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
	var albums []Album
	req, _ := http.NewRequest(http.MethodGet, "/albums", nil)

	router.ServeHTTP(w, req)

	if err := json.Unmarshal(w.Body.Bytes(), &albums); err != nil {
		assert.Failf(t, "json unmarshal fail", "fail to unmarshall Albums")
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, listAlbums(), albums)
}

func Test_getAlbumById(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	var album Album

	req, _ := http.NewRequest(http.MethodGet, "/albums/2", nil)
	router.ServeHTTP(w, req)

	body := w.Body.Bytes()
	if err := json.Unmarshal(body, &album); err != nil {
		assert.Failf(t, "json unmarshal fail", "fail to unmarshall Albums %v", string(body))
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, listAlbums()[1], album)
	assert.Equal(t, listAlbums()[1].Title, album.Title)
}

func Test_getAlbumById_NotFound(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	var serverError ServerError

	req, _ := http.NewRequest(http.MethodGet, "/albums/5666", nil)
	router.ServeHTTP(w, req)

	if err := json.Unmarshal(w.Body.Bytes(), &serverError); err != nil {
		assert.Empty(t, err, "json marshaling failed")
	}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "album not found", serverError.Message)
}

func Test_postAlbum(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	var album Album

	albumBody := `{"ID": "10", "Title": "The Ozzman Cometh", "Artist": "Black Sabbath", "Price": 56.99}`
	expectedAlbum := Album{ID: "10", Title: "The Ozzman Cometh", Artist: "Black Sabbath", Price: 56.99}
	req, _ := http.NewRequest(http.MethodPost, "/albums", strings.NewReader(albumBody))

	router.ServeHTTP(w, req)
	body := w.Body.Bytes()

	if err := json.Unmarshal(body, &album); err != nil {
		assert.Empty(t, err, "json marshaling failed", string(body))
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

	if err := json.Unmarshal(body, &serverError); err != nil {
		var ve validator.ValidationErrors
		errors.As(err, &ve)
		assert.Failf(t, "did not get error message from server", ve.Error())
	}
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, 1, len(serverError.BindingErrors))
	assert.Equal(t, "Title", serverError.BindingErrors[0].Field)
	assert.Equal(t, "This field is Required", serverError.BindingErrors[0].Message)
}

func BenchmarkGetAllAlbums(b *testing.B) {
	router := setupRouter()
	w := httptest.NewRecorder()
	var albums []Album
	req, _ := http.NewRequest(http.MethodGet, "/albums", nil)

	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
		if err := json.Unmarshal(w.Body.Bytes(), &albums); err != nil {
			assert.Fail(b, "should bind JSON to []Albums", "GetAlbum", err)
		}
		assert.Equal(b, http.StatusOK, w.Code)
		w.Body.Reset() //get requests need resets else the returned body is concatenated
	}
}

func BenchmarkGetAlbumById(b *testing.B) {
	router := setupRouter()
	w := httptest.NewRecorder()
	var album Album
	req, _ := http.NewRequest(http.MethodGet, "/albums/2", nil)

	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
		if err := json.Unmarshal(w.Body.Bytes(), &album); err != nil {
			assert.Fail(b, "should bind to JSON to Album", "GetAlbum", err)
		}
		assert.Equal(b, http.StatusOK, w.Code)
		assert.Equal(b, listAlbums()[1].Title, album.Title)
		w.Body.Reset() //get requests need resets else the returned body is concatenated
	}
}

func BenchmarkGetAlbumById_BadRequest(b *testing.B) {
	router := setupRouter()
	w := httptest.NewRecorder()
	var serverError ServerError
	req, _ := http.NewRequest(http.MethodGet, "/albums/5666", nil)

	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
		if err := json.Unmarshal(w.Body.Bytes(), &serverError); err != nil {
			assert.Fail(b, "should not fail binding", "GetAlbumById", err)
		}
		assert.Equal(b, http.StatusBadRequest, w.Code)
		assert.Equal(b, "album not found", serverError.Message)
		w.Body.Reset() //get requests need resets else the returned body is concatenated
	}
}

func BenchmarkPostAlbum(b *testing.B) {
	router := setupRouter()
	w := httptest.NewRecorder()
	albumJson := `{"ID": "10", "Title": "The Ozzman Cometh", "Artist": "Black Sabbath", "Price": 56.99}`
	expectedAlbum := Album{ID: "10", Title: "The Ozzman Cometh", Artist: "Black Sabbath", Price: 56.99}
	req, _ := http.NewRequest(http.MethodPost, "/albums", strings.NewReader(albumJson))

	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
		assert.Equal(b, http.StatusCreated, w.Code)
		bodyReturned := w.Body.Bytes()
		var albumReturned Album
		if err := json.Unmarshal(bodyReturned, &albumReturned); err != nil {
			assert.Fail(b, "binding JSON returned failure", albumJson, err.Error(), string(bodyReturned), albumReturned)
		}
		assert.Equal(b, albumReturned, expectedAlbum)
	}
}

func BenchmarkPostAlbum_BadRequest_BadJson(b *testing.B) {
	router := setupRouter()
	w := httptest.NewRecorder()
	var serverError ServerError
	albumJson := `{"XID": "10", "Titlexx": "Blue Train", "Artistx": "John Coltrane", "Price": 56.99, "X": "asdf"}`
	req, _ := http.NewRequest(http.MethodPost, "/albums", strings.NewReader(albumJson))

	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
		bodyReturned := w.Body.Bytes()
		if err := json.Unmarshal(bodyReturned, &serverError); err != nil {
			assert.Fail(b, "binding JSON returned failure", err.Error(), string(bodyReturned), serverError.Message, serverError.BindingErrors[0])
		}
		assert.Equal(b, http.StatusBadRequest, w.Code)
	}
}
