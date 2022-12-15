package main

import (
	"encoding/json"
	"errors"
	"example/go-gin-example/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// Gin & span testing upon https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/gin-gonic/gin/otelgin/test/gintrace_test.go

func TestMain(m *testing.M) {
	//Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Run the other tests
	os.Exit(m.Run())
}

func setupTestRouter() (*tracetest.SpanRecorder, *gin.Engine) {
	sr := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	router := setupRouter()
	router.Use(otelgin.Middleware("test-otel"))
	return sr, router
}

func makeKeyMap(attributes []attribute.KeyValue) map[attribute.Key]attribute.Value {
	var attributeMap = make(map[attribute.Key]attribute.Value)
	for _, keyValue := range attributes {
		attributeMap[keyValue.Key] = keyValue.Value
	}
	return attributeMap
}

func Test_getAllAlbums(t *testing.T) {
	sr, router := setupTestRouter()

	testRecorder := httptest.NewRecorder()
	var albums []models.Album
	req, _ := http.NewRequest(http.MethodGet, "/albums", nil)
	router.ServeHTTP(testRecorder, req)
	if err := json.Unmarshal(testRecorder.Body.Bytes(), &albums); err != nil {
		assert.Fail(t, "json unmarshal fail", "should be []Albums ", albums)
	}

	finishedSpans := sr.Ended()
	assert.Len(t, finishedSpans, 1)
	attributeMap := makeKeyMap(finishedSpans[0].Attributes())
	assert.Contains(t, attributeMap["http.status_code"].Emit(), "200")
	assert.Equal(t, codes.Ok, finishedSpans[0].Status().Code)
	assert.Equal(t, http.StatusOK, testRecorder.Code)
	assert.Equal(t, listAlbums(), albums)
}

func Test_getAlbumById(t *testing.T) {
	sr, router := setupTestRouter()
	testRecorder := httptest.NewRecorder()
	var album models.Album

	req, _ := http.NewRequest(http.MethodGet, "/albums/2", nil)
	router.ServeHTTP(testRecorder, req)

	body := testRecorder.Body.Bytes()
	if err := json.Unmarshal(body, &album); err != nil {
		assert.Fail(t, "json unmarshal fail", "Should be Album ", string(body))
	}

	finishedSpans := sr.Ended()
	assert.Len(t, finishedSpans, 1)
	attributeMap := makeKeyMap(finishedSpans[0].Attributes())
	assert.Contains(t, attributeMap["http.status_code"].Emit(), "200")
	assert.Equal(t, codes.Ok, finishedSpans[0].Status().Code)
	assert.Equal(t, http.StatusOK, testRecorder.Code)
	assert.Equal(t, listAlbums()[1], album)
	assert.Equal(t, listAlbums()[1].Title, album.Title)
}

func Test_getAlbumById_BadId(t *testing.T) {
	sr, router := setupTestRouter()
	testRecorder := httptest.NewRecorder()
	var serverError models.ServerError

	req, _ := http.NewRequest(http.MethodGet, "/albums/X", nil)
	router.ServeHTTP(testRecorder, req)

	body := testRecorder.Body.Bytes()
	if err := json.Unmarshal(body, &serverError); err != nil {
		assert.Fail(t, "json unmarshal fail", "Should be Album ", string(body))
	}

	finishedSpans := sr.Ended()
	assert.Len(t, finishedSpans, 1)
	attributeMap := makeKeyMap(finishedSpans[0].Attributes())
	assert.Equal(t, "400", attributeMap["http.status_code"].Emit())
	assert.Equal(t, "X", attributeMap["Id"].Emit())
	assert.Equal(t, 1, len(finishedSpans[0].Events()))
	assert.Equal(t, "Get /album invalid ID X", finishedSpans[0].Events()[0].Name)
	assert.Equal(t, codes.Error, finishedSpans[0].Status().Code)
	assert.Equal(t, http.StatusBadRequest, testRecorder.Code)
	assert.Equal(t, "Album ID [X] is not a valid number", serverError.Message)
}

func Test_getAlbumById_NotFound(t *testing.T) {
	sr, router := setupTestRouter()
	testRecorder := httptest.NewRecorder()
	var serverError models.ServerError
	albumID := 5666

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%v", "/albums/", albumID), nil)
	router.ServeHTTP(testRecorder, req)

	body := testRecorder.Body.Bytes()
	if err := json.Unmarshal(body, &serverError); err != nil {
		assert.Fail(t, "json unmarshalling fail", "Should be ServerError ", string(body))
	}
	finishedSpans := sr.Ended()
	assert.Len(t, finishedSpans, 1)
	attributeMap := makeKeyMap(finishedSpans[0].Attributes())
	assert.Equal(t, "400", attributeMap["http.status_code"].Emit())
	assert.Equal(t, codes.Error, finishedSpans[0].Status().Code)
	assert.Equal(t, "Get /album not found with ID 5666", finishedSpans[0].Events()[0].Name)
	assert.Equal(t, http.StatusBadRequest, testRecorder.Code)
	assert.Equal(t, fmt.Sprintf("%s [%v] %s", "Album", albumID, "not found"), serverError.Message)
}

func Test_postAlbum(t *testing.T) {
	resetAlbums()
	sr, router := setupTestRouter()
	testRecorder := httptest.NewRecorder()
	var album models.Album

	albumBody := `{"id": 10, "title": "The Ozzman Cometh", "artist": "Black Sabbath", "price": 66.60}`
	expectedAlbum := models.Album{ID: 10, Title: "The Ozzman Cometh", Artist: "Black Sabbath", Price: 66.60}
	req, _ := http.NewRequest(http.MethodPost, "/albums", strings.NewReader(albumBody))

	assert.Equal(t, len(listAlbums()), 3)
	router.ServeHTTP(testRecorder, req)
	body := testRecorder.Body.Bytes()

	if err := json.Unmarshal(body, &album); err != nil {
		assert.Fail(t, "json unmarshalling fail", "Should be an Album ", string(body))
	}
	finishedSpans := sr.Ended()
	assert.Len(t, finishedSpans, 1)
	attributeMap := makeKeyMap(finishedSpans[0].Attributes())
	assert.Equal(t, "201", attributeMap["http.status_code"].Emit())
	assert.Equal(t, codes.Ok, finishedSpans[0].Status().Code)
	assert.Equal(t, http.StatusCreated, testRecorder.Code)
	assert.Equal(t, album, expectedAlbum)
	assert.Equal(t, len(listAlbums()), 4)
}

func Test_postAlbum_BadRequest_BadJSON_MissingValues(t *testing.T) {
	resetAlbums()
	sr, router := setupTestRouter()
	testRecorder := httptest.NewRecorder()

	album := `{"xid": 10, "titlex": "Blue Train", "artistx": "Lead Belly", "pricex": 56.99, "X": "asdf"}`
	var serverError models.ServerError
	req, _ := http.NewRequest(http.MethodPost, "/albums", strings.NewReader(album))
	router.ServeHTTP(testRecorder, req)
	body := testRecorder.Body.Bytes()

	if err := json.Unmarshal(body, &serverError); err != nil {
		var ve validator.ValidationErrors
		errors.As(err, &ve)
		assert.Fail(t, "json unmarshalling fail", "should be ServerError ", ve.Error(), string(body))
	}
	assert.Equal(t, http.StatusBadRequest, testRecorder.Code)
	finishedSpans := sr.Ended()
	assert.Len(t, finishedSpans, 1)
	attributeMap := makeKeyMap(finishedSpans[0].Attributes())
	assert.Equal(t, "400", attributeMap["http.status_code"].Emit())
	assert.Equal(t, codes.Error, finishedSpans[0].Status().Code)
	assert.Equal(t, "[{\"field\":\"id\",\"message\":\"below minimum value\"},{\"field\":\"title\",\"message\":\"required field\"},{\"field\":\"artist\",\"message\":\"required field\"},{\"field\":\"price\",\"message\":\"required field\"}]", finishedSpans[0].Events()[0].Name)
	assert.Equal(t, 4, len(serverError.BindingErrors))
	assert.Equal(t, "title", serverError.BindingErrors[1].Field)
	assert.Equal(t, "required field", serverError.BindingErrors[1].Message)
	assert.Equal(t, "artist", serverError.BindingErrors[2].Field)
	assert.Equal(t, "required field", serverError.BindingErrors[2].Message)
	assert.Equal(t, "price", serverError.BindingErrors[3].Field)
	assert.Equal(t, "required field", serverError.BindingErrors[3].Message)
	assert.Equal(t, len(listAlbums()), 3)

}

func Test_postAlbum_BadRequest_BadJSON_MinValues(t *testing.T) {
	resetAlbums()
	sr, router := setupTestRouter()
	testRecorder := httptest.NewRecorder()
	album := `{"id": -1, "title": "a", "artist": "z", "price": -0.1}`
	var serverError models.ServerError
	req, _ := http.NewRequest(http.MethodPost, "/albums", strings.NewReader(album))
	router.ServeHTTP(testRecorder, req)
	body := testRecorder.Body.Bytes()

	if err := json.Unmarshal(body, &serverError); err != nil {
		var ve validator.ValidationErrors
		errors.As(err, &ve)
		assert.Fail(t, "json unmarshalling fail", "should be ServerError ", ve.Error(), string(body))
	}
	assert.Equal(t, http.StatusBadRequest, testRecorder.Code)
	finishedSpans := sr.Ended()
	assert.Len(t, finishedSpans, 1)
	attributeMap := makeKeyMap(finishedSpans[0].Attributes())
	assert.Equal(t, "400", attributeMap["http.status_code"].Emit())
	assert.Equal(t, codes.Error, finishedSpans[0].Status().Code)
	assert.Equal(t, "[{\"field\":\"id\",\"message\":\"below minimum value\"},{\"field\":\"title\",\"message\":\"below minimum value\"},{\"field\":\"artist\",\"message\":\"below minimum value\"},{\"field\":\"price\",\"message\":\"below minimum value\"}]", finishedSpans[0].Events()[0].Name)
	assert.Equal(t, 4, len(serverError.BindingErrors))
	assert.Equal(t, "id", serverError.BindingErrors[0].Field)
	assert.Equal(t, "below minimum value", serverError.BindingErrors[0].Message)
	assert.Equal(t, "title", serverError.BindingErrors[1].Field)
	assert.Equal(t, "below minimum value", serverError.BindingErrors[1].Message)
	assert.Equal(t, "artist", serverError.BindingErrors[2].Field)
	assert.Equal(t, "below minimum value", serverError.BindingErrors[2].Message)
	assert.Equal(t, "price", serverError.BindingErrors[3].Field)
	assert.Equal(t, "below minimum value", serverError.BindingErrors[3].Message)
	assert.Equal(t, len(listAlbums()), 3)
}

func Test_postAlbum_BadRequest_Malformed_JSON(t *testing.T) {
	resetAlbums()
	sr, router := setupTestRouter()
	testRecorder := httptest.NewRecorder()
	album := `{"id": -1,`
	var serverError models.ServerError
	req, _ := http.NewRequest(http.MethodPost, "/albums", strings.NewReader(album))
	router.ServeHTTP(testRecorder, req)
	body := testRecorder.Body.Bytes()

	if err := json.Unmarshal(body, &serverError); err == nil {
		assert.Fail(t, "", "should be ServerError ")
	}
	assert.Equal(t, http.StatusBadRequest, testRecorder.Code)
	finishedSpans := sr.Ended()
	assert.Len(t, finishedSpans, 1)
	attributeMap := makeKeyMap(finishedSpans[0].Attributes())
	assert.Contains(t, attributeMap["http.status_code"].Emit(), "400")
	assert.Equal(t, attributeMap["http.request.body"].Emit(), "{\"id\": -1,")
	assert.Equal(t, codes.Error, finishedSpans[0].Status().Code)
	assert.Equal(t, "unexpected EOF", finishedSpans[0].Events()[0].Name)
	assert.Equal(t, 0, len(serverError.BindingErrors))
	assert.Equal(t, len(listAlbums()), 3)
}

func Test_getSwagger(t *testing.T) {
	resetAlbums()
	router := setupRouter()
	testRecorder := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v3/api-docs/", nil)
	router.ServeHTTP(testRecorder, req)
	bodyString := string(testRecorder.Body.Bytes())
	assert.Contains(t, bodyString, "swagger-initializer.js")
	assert.Equal(t, http.StatusOK, testRecorder.Code)
}

func Benchmark_getAllAlbums(b *testing.B) {
	router := setupRouter()
	testRecorder := httptest.NewRecorder()
	var albums []models.Album
	req, _ := http.NewRequest(http.MethodGet, "/albums", nil)

	for i := 0; i < b.N; i++ {
		router.ServeHTTP(testRecorder, req)
		body := testRecorder.Body.Bytes()
		if err := json.Unmarshal(body, &albums); err != nil {
			assert.Fail(b, "json unmarshalling fail", "should be []Album ", string(body))
		}
		testRecorder.Body.Reset() //get requests need resets else the returned body is concatenated
	}
}

func Benchmark_getAlbumById(b *testing.B) {
	router := setupRouter()
	testRecorder := httptest.NewRecorder()
	var album models.Album
	req, _ := http.NewRequest(http.MethodGet, "/albums/2", nil)

	for i := 0; i < b.N; i++ {
		router.ServeHTTP(testRecorder, req)
		body := testRecorder.Body.Bytes()
		if err := json.Unmarshal(body, &album); err != nil {
			assert.Fail(b, "json unmarshalling fail", "should be Album ", string(body))
		}
		testRecorder.Body.Reset() //get requests need resets else the returned body is concatenated
	}
}

func Benchmark_getAlbumById_BadRequest(b *testing.B) {
	router := setupRouter()
	testRecorder := httptest.NewRecorder()
	var serverError models.ServerError
	req, _ := http.NewRequest(http.MethodGet, "/albums/5666", nil)

	for i := 0; i < b.N; i++ {
		router.ServeHTTP(testRecorder, req)
		body := testRecorder.Body.Bytes()
		if err := json.Unmarshal(body, &serverError); err != nil {
			assert.Fail(b, "json unmarshalling fail", "should be ServerError ", string(body))
		}
		testRecorder.Body.Reset() //get requests need resets else the returned body is concatenated
	}
}

func Benchmark_postAlbum(b *testing.B) {
	router := setupRouter()
	testRecorder := httptest.NewRecorder()
	albumJson := `{"id": "10", "title": "The Ozzman Cometh", "artist": "Black Sabbath", "price": 56.99}`
	req, _ := http.NewRequest(http.MethodPost, "/albums", strings.NewReader(albumJson))
	var albumReturned models.Album

	for i := 0; i < b.N; i++ {
		router.ServeHTTP(testRecorder, req)
		bodyReturned := testRecorder.Body.Bytes()
		if err := json.Unmarshal(bodyReturned, &albumReturned); err != nil {
			assert.Fail(b, "json unmarshalling fail", "should be Album ", string(bodyReturned))
		}
		testRecorder.Body.Reset()
	}
}

func Benchmark_postAlbum_BadRequest_BadJson(b *testing.B) {
	router := setupRouter()
	testRecorder := httptest.NewRecorder()
	var returnedError models.ServerError
	albumJson := `{"xid": "10", "titlex": "Blue Train", "artistx": "John Coltrane", "pricex": 56.99, "X": "asdf"}`
	req, _ := http.NewRequest(http.MethodPost, "/albums", strings.NewReader(albumJson))

	for i := 0; i < b.N; i++ {
		router.ServeHTTP(testRecorder, req)
		bodyReturned := testRecorder.Body.Bytes()
		if err := json.Unmarshal(bodyReturned, &returnedError); err != nil {
			assert.Fail(b, "json unmarshalling fail", "Should be ServerError ", string(bodyReturned))
		}
		testRecorder.Body.Reset()
	}
}
