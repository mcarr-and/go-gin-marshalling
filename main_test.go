package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"example.com/album-store/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// inspired by this for setting up gin & otel to test spans
// https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/gin-gonic/gin/otelgin/test/gintrace_test.go

func TestMain(m *testing.M) {
	//Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Run the other tests
	os.Exit(m.Run())
}

func setupTestRouter() (*httptest.ResponseRecorder, *tracetest.SpanRecorder, *gin.Engine) {
	spanRecorder := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder)))
	router := setupRouter()
	testRecorder := httptest.NewRecorder()
	router.Use(otelgin.Middleware("test-otel"))
	return testRecorder, spanRecorder, router
}

func makeKeyMap(attributes []attribute.KeyValue) map[attribute.Key]attribute.Value {
	var attributeMap = make(map[attribute.Key]attribute.Value)
	for _, keyValue := range attributes {
		attributeMap[keyValue.Key] = keyValue.Value
	}
	return attributeMap
}

func Test_getAllAlbums(t *testing.T) {
	testRecorder, spanRecorder, router := setupTestRouter()

	var albums []models.Album

	req := newTestHttpRequest(http.MethodGet, "/albums", nil)
	router.ServeHTTP(testRecorder, req)
	if err := json.Unmarshal(testRecorder.Body.Bytes(), &albums); err != nil {
		assert.Fail(t, "json unmarshal fail", "should be []Albums ", albums)
	}

	assert.Equal(t, http.StatusOK, testRecorder.Code)

	finishedSpans := spanRecorder.Ended()
	assert.Len(t, finishedSpans, 1)

	assert.Equal(t, codes.Ok, finishedSpans[0].Status().Code)
	assert.Equal(t, "", finishedSpans[0].Status().Description)

	assert.Equal(t, 0, len(finishedSpans[0].Events()))

	attributeMap := makeKeyMap(finishedSpans[0].Attributes())
	assert.Equal(t, "200", attributeMap["album-store.response.code"].Emit())
	assert.Equal(t, "/albums", attributeMap["http.target"].Emit())

	assert.Equal(t, listAlbums(), albums)
}

func Test_getAlbumById(t *testing.T) {
	testRecorder, spanRecorder, router := setupTestRouter()

	var album models.Album

	req := newTestHttpRequest(http.MethodGet, "/albums/2", nil)
	router.ServeHTTP(testRecorder, req)
	if err := json.Unmarshal(testRecorder.Body.Bytes(), &album); err != nil {
		assert.Fail(t, "json unmarshal fail", "Should be Album ", testRecorder.Body.String())
	}

	assert.Equal(t, http.StatusOK, testRecorder.Code)

	finishedSpans := spanRecorder.Ended()
	assert.Len(t, finishedSpans, 1)

	assert.Equal(t, codes.Ok, finishedSpans[0].Status().Code)
	assert.Equal(t, "", finishedSpans[0].Status().Description)

	assert.Equal(t, 0, len(finishedSpans[0].Events()))

	attributeMap := makeKeyMap(finishedSpans[0].Attributes())
	assert.Equal(t, "200", attributeMap["album-store.response.code"].Emit())
	assert.Equal(t, `{"id":2,"title":"Jeru","artist":"Gerry Mulligan","price":17.99}`, attributeMap["album-store.response.body"].Emit())
	assert.Equal(t, "/albums/2", attributeMap["http.target"].Emit())

	assert.Equal(t, listAlbums()[1], album)
	assert.Equal(t, listAlbums()[1].Title, album.Title)
}

func Test_getAlbumById_InvalidID_Character(t *testing.T) {
	testRecorder, spanRecorder, router := setupTestRouter()

	var serverError models.ServerError

	req := newTestHttpRequest(http.MethodGet, "/albums/X", nil)
	router.ServeHTTP(testRecorder, req)
	if err := json.Unmarshal(testRecorder.Body.Bytes(), &serverError); err != nil {
		assert.Fail(t, "json unmarshal fail", "Should be ServerError ", testRecorder.Body.String())
	}

	assert.Equal(t, http.StatusBadRequest, testRecorder.Code)

	finishedSpans := spanRecorder.Ended()
	assert.Len(t, finishedSpans, 1)

	expectedErrorMessage := "Album [X] not found, invalid request"

	assert.Equal(t, codes.Error, finishedSpans[0].Status().Code)
	assert.Equal(t, expectedErrorMessage, finishedSpans[0].Status().Description)

	assert.Equal(t, 1, len(finishedSpans[0].Events()))
	assert.Equal(t, expectedErrorMessage, finishedSpans[0].Events()[0].Name)

	attributeMap := makeKeyMap(finishedSpans[0].Attributes())
	assert.Equal(t, "400", attributeMap["album-store.response.code"].Emit())
	assert.Equal(t, "/albums/X", attributeMap["http.target"].Emit())

	assert.Equal(t, expectedErrorMessage, serverError.Message)
}

func Test_getAlbumById_NotFound(t *testing.T) {
	testRecorder, spanRecorder, router := setupTestRouter()

	var serverError models.ServerError
	invalidAlbumID := -1666

	req := newTestHttpRequest(http.MethodGet, fmt.Sprintf("%s%v", "/albums/", invalidAlbumID), nil)
	router.ServeHTTP(testRecorder, req)
	if err := json.Unmarshal(testRecorder.Body.Bytes(), &serverError); err != nil {
		assert.Fail(t, "json unmarshalling fail", "Should be ServerError ", testRecorder.Body.String())
	}

	assert.Equal(t, http.StatusBadRequest, testRecorder.Code)

	finishedSpans := spanRecorder.Ended()
	assert.Len(t, finishedSpans, 1)

	expectedErrorMessage := fmt.Sprintf("Album [%v] not found", invalidAlbumID)

	assert.Equal(t, codes.Error, finishedSpans[0].Status().Code)
	assert.Equal(t, expectedErrorMessage, finishedSpans[0].Status().Description)

	assert.Equal(t, 1, len(finishedSpans[0].Events()))
	assert.Equal(t, expectedErrorMessage, finishedSpans[0].Events()[0].Name)

	attributeMap := makeKeyMap(finishedSpans[0].Attributes())
	assert.Equal(t, "400", attributeMap["album-store.response.code"].Emit())
	assert.Equal(t, fmt.Sprintf("/albums/%v", invalidAlbumID), attributeMap["http.target"].Emit())

	assert.Equal(t, expectedErrorMessage, serverError.Message)
}

func Test_postAlbum(t *testing.T) {
	resetAlbums()

	testRecorder, spanRecorder, router := setupTestRouter()
	var album models.Album

	expectedAlbum := models.Album{ID: 10, Title: "The Ozzman Cometh", Artist: "Black Sabbath", Price: 66.60}
	albumBody := `{"id": 10, "title": "The Ozzman Cometh", "artist": "Black Sabbath", "price": 66.60}`

	req := newTestHttpRequest(http.MethodPost, "/albums", strings.NewReader(albumBody))
	router.ServeHTTP(testRecorder, req)
	if err := json.Unmarshal(testRecorder.Body.Bytes(), &album); err != nil {
		assert.Fail(t, "json unmarshalling fail", "Should be a valid Album ", testRecorder.Body.String())
	}

	assert.Equal(t, http.StatusCreated, testRecorder.Code)

	finishedSpans := spanRecorder.Ended()
	assert.Len(t, finishedSpans, 1)

	assert.Equal(t, codes.Ok, finishedSpans[0].Status().Code)
	assert.Equal(t, "", finishedSpans[0].Status().Description)

	assert.Equal(t, 0, len(finishedSpans[0].Events()))

	attributeMap := makeKeyMap(finishedSpans[0].Attributes())
	assert.Equal(t, `{"id": 10, "title": "The Ozzman Cometh", "artist": "Black Sabbath", "price": 66.60}`, attributeMap["album-store.request.body"].Emit())
	assert.Equal(t, `{"id":10,"title":"The Ozzman Cometh","artist":"Black Sabbath","price":66.6}`, attributeMap["album-store.response.body"].Emit())
	assert.Equal(t, "201", attributeMap["album-store.response.code"].Emit())
	assert.Equal(t, "/albums", attributeMap["http.target"].Emit())

	assert.Equal(t, album, expectedAlbum)
	assert.Equal(t, len(listAlbums()), 4)
}

func Test_postAlbum_BadRequest_BadJSON_MissingValues(t *testing.T) {
	resetAlbums()
	testRecorder, spanRecorder, router := setupTestRouter()

	var serverError models.ServerError
	album := `{"xid": 10, "titlex": "Blue Train", "artistx": "Lead Belly", "pricex": 56.99, "X": "asdf"}`
	bindingErrorMessage := `[{"field":"id","message":"below minimum value"},{"field":"title","message":"required field"},{"field":"artist","message":"required field"},{"field":"price","message":"required field"}]`

	req := newTestHttpRequest(http.MethodPost, "/albums", strings.NewReader(album))
	router.ServeHTTP(testRecorder, req)
	if err := json.Unmarshal(testRecorder.Body.Bytes(), &serverError); err != nil {
		var ve validator.ValidationErrors
		errors.As(err, &ve)
		assert.Fail(t, "json unmarshalling fail", "should be ServerError ", ve.Error(), testRecorder.Body.String())
	}

	assert.Equal(t, http.StatusBadRequest, testRecorder.Code)

	finishedSpans := spanRecorder.Ended()
	assert.Len(t, finishedSpans, 1)

	assert.Equal(t, codes.Error, finishedSpans[0].Status().Code)
	assert.Equal(t, "Album JSON field validation failed", finishedSpans[0].Status().Description)

	assert.Equal(t, 1, len(finishedSpans[0].Events()))
	assert.Equal(t, bindingErrorMessage, finishedSpans[0].Events()[0].Name)

	attributeMap := makeKeyMap(finishedSpans[0].Attributes())
	assert.Equal(t, "400", attributeMap["album-store.response.code"].Emit())
	assert.Equal(t, album, attributeMap["album-store.request.body"].Emit())
	assert.Equal(t, fmt.Sprintf("{\"errors\":%v}", bindingErrorMessage), attributeMap["album-store.response.body"].Emit())
	assert.Equal(t, `{"xid": 10, "titlex": "Blue Train", "artistx": "Lead Belly", "pricex": 56.99, "X": "asdf"}`, attributeMap["album-store.request.body"].Emit())
	assert.Equal(t, "/albums", attributeMap["http.target"].Emit())

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
	testRecorder, spanRecorder, router := setupTestRouter()

	album := `{"id": -1, "title": "a", "artist": "z", "price": -0.1}`
	bindingErrorMessage := `[{"field":"id","message":"below minimum value"},{"field":"title","message":"below minimum value"},{"field":"artist","message":"below minimum value"},{"field":"price","message":"below minimum value"}]`
	var serverError models.ServerError

	req := newTestHttpRequest(http.MethodPost, "/albums", strings.NewReader(album))
	router.ServeHTTP(testRecorder, req)
	if err := json.Unmarshal(testRecorder.Body.Bytes(), &serverError); err != nil {
		var ve validator.ValidationErrors
		errors.As(err, &ve)
		assert.Fail(t, "json unmarshalling fail", "should be ServerError ", ve.Error(), testRecorder.Body.String())
	}

	assert.Equal(t, http.StatusBadRequest, testRecorder.Code)

	finishedSpans := spanRecorder.Ended()
	assert.Len(t, finishedSpans, 1)

	assert.Equal(t, codes.Error, finishedSpans[0].Status().Code)
	assert.Equal(t, "Album JSON field validation failed", finishedSpans[0].Status().Description)

	assert.Equal(t, 1, len(finishedSpans[0].Events()))
	assert.Equal(t, bindingErrorMessage, finishedSpans[0].Events()[0].Name)

	attributeMap := makeKeyMap(finishedSpans[0].Attributes())
	assert.Equal(t, "400", attributeMap["album-store.response.code"].Emit())
	assert.Equal(t, `{"id": -1, "title": "a", "artist": "z", "price": -0.1}`, attributeMap["album-store.request.body"].Emit())
	assert.Equal(t, fmt.Sprintf("{\"errors\":%v}", bindingErrorMessage), attributeMap["album-store.response.body"].Emit())
	assert.Equal(t, "/albums", attributeMap["http.target"].Emit())

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
	testRecorder, spanRecorder, router := setupTestRouter()

	var serverError models.ServerError
	requestBody := `{"id": -1,`

	req := newTestHttpRequest(http.MethodPost, "/albums", strings.NewReader(requestBody))
	router.ServeHTTP(testRecorder, req)
	if err := json.Unmarshal(testRecorder.Body.Bytes(), &serverError); err != nil {
		assert.Fail(t, "", "should be ServerError ")
	}

	assert.Equal(t, http.StatusBadRequest, testRecorder.Code)

	finishedSpans := spanRecorder.Ended()
	assert.Len(t, finishedSpans, 1)

	assert.Equal(t, codes.Error, finishedSpans[0].Status().Code)
	assert.Equal(t, "Malformed JSON. Not valid for Album", finishedSpans[0].Status().Description)

	assert.Equal(t, 1, len(finishedSpans[0].Events()))
	assert.Equal(t, "Malformed JSON. unexpected EOF", finishedSpans[0].Events()[0].Name)

	attributeMap := makeKeyMap(finishedSpans[0].Attributes())
	assert.Equal(t, "400", attributeMap["album-store.response.code"].Emit())
	assert.Equal(t, requestBody, attributeMap["album-store.request.body"].Emit())
	assert.Equal(t, `{"message":"Malformed JSON. Not valid for Album"}`, attributeMap["album-store.response.body"].Emit())
	assert.Equal(t, "/albums", attributeMap["http.target"].Emit())

	assert.Equal(t, "Malformed JSON. Not valid for Album", serverError.Message)
	assert.Equal(t, 0, len(serverError.BindingErrors))

	assert.Equal(t, len(listAlbums()), 3)
}

func Test_getSwagger(t *testing.T) {
	resetAlbums()
	testRecorder, _, router := setupTestRouter()

	req := newTestHttpRequest(http.MethodGet, "/v3/api-docs/", nil)
	router.ServeHTTP(testRecorder, req)
	bodyString := testRecorder.Body.String()

	assert.Contains(t, bodyString, "swagger-initializer.js")
	assert.Equal(t, http.StatusOK, testRecorder.Code)
}

func Test_getStatus(t *testing.T) {
	testRecorder, spanRecorder, router := setupTestRouter()

	req := newTestHttpRequest(http.MethodGet, "/status", nil)
	router.ServeHTTP(testRecorder, req)

	responseBodyString := testRecorder.Body.String()

	assert.Equal(t, http.StatusOK, testRecorder.Code)
	assert.Equal(t, `{"status":"OK"}`, responseBodyString)

	finishedSpans := spanRecorder.Ended()
	assert.Len(t, finishedSpans, 1)

	assert.Equal(t, codes.Ok, finishedSpans[0].Status().Code)
	assert.Equal(t, "", finishedSpans[0].Status().Description)

	assert.Equal(t, 0, len(finishedSpans[0].Events()))

	attributeMap := makeKeyMap(finishedSpans[0].Attributes())
	assert.Equal(t, "/status", attributeMap["http.target"].Emit())
}

func Test_getMetrics(t *testing.T) {
	testRecorder, spanRecorder, router := setupTestRouter()

	req := newTestHttpRequest(http.MethodGet, "/metrics", nil)
	router.ServeHTTP(testRecorder, req)

	responseBodyString := testRecorder.Body.String()

	assert.Equal(t, http.StatusOK, testRecorder.Code)
	assert.Contains(t, responseBodyString, `go_gc_duration_seconds`)

	finishedSpans := spanRecorder.Ended()
	assert.Len(t, finishedSpans, 1)

	assert.Equal(t, codes.Ok, finishedSpans[0].Status().Code)
	assert.Equal(t, "", finishedSpans[0].Status().Description)

	assert.Equal(t, 0, len(finishedSpans[0].Events()))

	attributeMap := makeKeyMap(finishedSpans[0].Attributes())
	assert.Equal(t, "/metrics", attributeMap["http.target"].Emit())
}

// setup test requests to have the requestURI. This helps enable the http.target to be set when reading the request in Open-Telemetry
func newTestHttpRequest(method string, url string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, url, body)
	req.RequestURI = req.URL.RequestURI()
	return req
}

func Benchmark_getAllAlbums(b *testing.B) {
	testRecorder, _, router := setupTestRouter()

	var albums []models.Album
	req, _ := http.NewRequest(http.MethodGet, "/albums", nil)

	for i := 0; i < b.N; i++ {
		router.ServeHTTP(testRecorder, req)
		if err := json.Unmarshal(testRecorder.Body.Bytes(), &albums); err != nil {
			assert.Fail(b, "json unmarshalling fail", "should be []Album ", testRecorder.Body.String())
		}
		testRecorder.Body.Reset() //get requests need resets else the returned body is concatenated
	}
}

func Benchmark_getAlbumById(b *testing.B) {
	testRecorder, _, router := setupTestRouter()

	var album models.Album
	req, _ := http.NewRequest(http.MethodGet, "/albums/2", nil)

	for i := 0; i < b.N; i++ {
		router.ServeHTTP(testRecorder, req)
		if err := json.Unmarshal(testRecorder.Body.Bytes(), &album); err != nil {
			assert.Fail(b, "json unmarshalling fail", "should be Album ", testRecorder.Body.String())
		}
		testRecorder.Body.Reset() //get requests need resets else the returned body is concatenated
	}
}

func Benchmark_getAlbumById_BadRequest(b *testing.B) {
	testRecorder, _, router := setupTestRouter()

	var serverError models.ServerError
	req, _ := http.NewRequest(http.MethodGet, "/albums/5666", nil)

	for i := 0; i < b.N; i++ {
		router.ServeHTTP(testRecorder, req)
		if err := json.Unmarshal(testRecorder.Body.Bytes(), &serverError); err != nil {
			assert.Fail(b, "json unmarshalling fail", "should be ServerError ", testRecorder.Body.String())
		}
		testRecorder.Body.Reset() //get requests need resets else the returned body is concatenated
	}
}

func Benchmark_postAlbum(b *testing.B) {
	testRecorder, _, router := setupTestRouter()

	var albumReturned models.Album
	albumJson := `{"id": "10", "title": "The Ozzman Cometh", "artist": "Black Sabbath", "price": 56.99}`
	req, _ := http.NewRequest(http.MethodPost, "/albums", strings.NewReader(albumJson))

	for i := 0; i < b.N; i++ {
		router.ServeHTTP(testRecorder, req)
		if err := json.Unmarshal(testRecorder.Body.Bytes(), &albumReturned); err != nil {
			assert.Fail(b, "json unmarshalling fail", "should be Album ", testRecorder.Body.String())
		}
		testRecorder.Body.Reset()
	}
}

func Benchmark_postAlbum_BadRequest_BadJson(b *testing.B) {
	testRecorder, _, router := setupTestRouter()

	var returnedError models.ServerError
	albumJson := `{"xid": "10", "titlex": "Blue Train", "artistx": "John Coltrane", "pricex": 56.99, "X": "asdf"}`
	req, _ := http.NewRequest(http.MethodPost, "/albums", strings.NewReader(albumJson))

	for i := 0; i < b.N; i++ {
		router.ServeHTTP(testRecorder, req)
		if err := json.Unmarshal(testRecorder.Body.Bytes(), &returnedError); err != nil {
			assert.Fail(b, "json unmarshalling fail", "Should be ServerError ", testRecorder.Body.String())
		}
		testRecorder.Body.Reset()
	}
}
