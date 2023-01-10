package main

import (
	"context"
	"encoding/json"
	"example.com/album-store/otelGinSetup"
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type OtelHttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	DefaultClient OtelHttpClient
)

func init() {
	DefaultClient = otelhttp.DefaultClient
}

func getAlbums(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetName("/albums GET")
	defer span.End()

	// proxy call to album-Store
	resp, err := Get(c.Request.Context(), albumStoreURL+"/albums")
	if err != nil {
		errorMessage := fmt.Sprintf("error contacting album-store getAlbums %v", err)
		span.AddEvent(errorMessage)
		span.SetStatus(codes.Error, errorMessage)
		span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": errorMessage})
		return
	}
	var albumStoreResponseBodyJson interface{}
	byteArray, err := io.ReadAll(resp.Body)
	responseBodyString := string(byteArray[:])
	err = json.NewDecoder(strings.NewReader(responseBodyString)).Decode(&albumStoreResponseBodyJson)
	if err != nil {
		errorMessage := "error from album-store getAlbums malformed JSON"
		span.AddEvent(errorMessage)
		span.SetStatus(codes.Error, errorMessage)
		span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
		span.SetAttributes(attribute.Key("http.response.body").String(responseBodyString))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": errorMessage})
		return
	}
	err = resp.Body.Close()
	if err != nil {
		errorMessage := fmt.Sprintf("error album-store closing response %v", err)
		span.AddEvent(errorMessage)
		span.SetStatus(codes.Error, errorMessage)
		span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": errorMessage})
		return
	}
	span.SetAttributes(attribute.Key("proxy-service.response").Int(resp.StatusCode))
	span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusOK))
	span.SetStatus(codes.Ok, "")
	c.JSON(http.StatusOK, albumStoreResponseBodyJson)
}

func getAlbumByID(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetName("/albums/:id GET")
	defer span.End()
	id := c.Param("id")
	span.SetAttributes(attribute.Key("http.request.parameters").String(fmt.Sprintf("%s=%s", "ID", id)))
	albumID, err := strconv.Atoi(id)
	// param ID is expected to be a number so fail if cannot covert to integer
	if err != nil {
		errorMessage := fmt.Sprintf("%s [%s] %s", "error invalid ID", id, "requested")
		span.SetStatus(codes.Error, errorMessage)
		span.AddEvent(errorMessage)
		span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": errorMessage})
		return
	}
	// proxy call to album-Store
	resp, err := Get(c.Request.Context(), fmt.Sprintf("%v/albums/%v", albumStoreURL, albumID))
	if err != nil {
		errorMessage := fmt.Sprintf("error contacting album-store getAlbumById %v", err)
		span.AddEvent(errorMessage)
		span.SetStatus(codes.Error, errorMessage)
		span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": errorMessage})
		return
	}
	var albumStoreResponseBodyJson interface{}
	byteArray, err := io.ReadAll(resp.Body)
	responseBodyString := string(byteArray[:])
	err = json.NewDecoder(strings.NewReader(responseBodyString)).Decode(&albumStoreResponseBodyJson)
	if err != nil {
		errorMessage := "error from album-store getAlbumById malformed JSON"
		span.AddEvent(errorMessage)
		span.SetStatus(codes.Error, errorMessage)
		span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
		span.SetAttributes(attribute.Key("http.response.body").String(responseBodyString))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": errorMessage})
		return
	}
	err = resp.Body.Close()
	if err != nil {
		log.Println(err)
		errorMessage := fmt.Sprintf("error closing body from album-store getAlbumById %v", err)
		span.AddEvent(errorMessage)
		span.SetStatus(codes.Error, errorMessage)
		span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": errorMessage})
		return
	}
	span.SetAttributes(attribute.Key("proxy-service.response").Int(resp.StatusCode))
	span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusOK))
	span.SetStatus(codes.Ok, "")
	c.JSON(http.StatusOK, albumStoreResponseBodyJson)
}

func postAlbum(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetName("/albums POST")
	defer span.End()
	var requestBody interface{}
	byteArray, err := io.ReadAll(c.Request.Body)
	requestBodyString := string(byteArray[:])
	err = json.NewDecoder(strings.NewReader(requestBodyString)).Decode(&requestBody)
	//any problems with conversion of request body to JSON then return failure
	if err != nil {
		log.Println(err)
		errorMessage := fmt.Sprintf("invalid request json body %v", requestBodyString)
		span.AddEvent(errorMessage)
		span.SetStatus(codes.Error, errorMessage)
		span.SetAttributes(attribute.Key("http.request.body").String(requestBodyString))
		span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": errorMessage})
		return
	}
	// proxy call to album-Store post album
	resp, err := Post(c.Request.Context(), albumStoreURL+"/albums/", "application/json", strings.NewReader(fmt.Sprintf("%v", requestBody)))
	if err != nil {
		errorMessage := fmt.Sprintf("error contacting album-store postAlbum %v", err)
		span.AddEvent(errorMessage)
		span.SetStatus(codes.Error, errorMessage)
		span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": errorMessage})
		return
	}
	var responseBody interface{}
	byteArray, err = io.ReadAll(resp.Body)
	responseBodyString := string(byteArray[:])
	err = json.NewDecoder(strings.NewReader(responseBodyString)).Decode(&responseBody)
	if err != nil {
		errorMessage := "error from album-store postAlbum malformed JSON"
		span.AddEvent(errorMessage)
		span.SetStatus(codes.Error, errorMessage)
		span.SetAttributes(attribute.Key("http.response.body").String(responseBodyString))
		span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": errorMessage})
		return
	}
	err = resp.Body.Close()
	if err != nil {
		errorMessage := fmt.Sprintf("error closing body from album-store postAlbum %v", err)
		span.AddEvent(errorMessage)
		span.SetStatus(codes.Error, errorMessage)
		span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": errorMessage})
		return
	}
	span.SetAttributes(attribute.Key("proxy-service.response").Int(resp.StatusCode))
	span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusOK))
	span.SetStatus(codes.Ok, "")
	c.JSON(http.StatusOK, responseBody)
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(otelgin.Middleware(serviceName)) // add OpenTelemetry to Gin
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbum)
	return router
}

const (
	serviceName  = "proxy-service"
	startAddress = "0.0.0.0:9070"
)

var version = "No-Version"
var gitHash = "No-Hash"
var albumStoreURL = "http://localhost:9080"

func main() {
	shutdownTraceProvider, err := otelGinSetup.InitOtelProvider(serviceName, version, gitHash)
	if err != nil {
		log.Fatal(err)
	}

	albumStoreUrlEnv := os.Getenv("ALBUM_STORE_URL")
	if albumStoreURL != "" {
		albumStoreURL = albumStoreUrlEnv
	}

	router := setupRouter()
	//serve requests until termination signal is sent.

	srv := &http.Server{
		Addr:    startAddress,
		Handler: h2c.NewHandler(router, &http2.Server{}),
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	<-quit

	log.Println("Server shutdown with 500ms timeout...")
	ctxServer, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	log.Println("OpenTelemetry TraceProvider flushing & shutting down")
	if err := shutdownTraceProvider(ctxServer); err != nil {
		log.Fatal("OpenTelemetry TracerProvider shutdown failure: %w", err)
	}
	log.Println("OpenTelemetry TraceProvider exited")

	if err := srv.Shutdown(ctxServer); err != nil {
		log.Fatal("Server shutdown failure:", err)
	}
	<-ctxServer.Done()

	log.Println("Server exiting")
}

// Extracted methods from https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/net/http/otelhttp/client.go v0.37.0
// this is to allow use of interface for httpClient and be able to mock out responses

func Get(ctx context.Context, targetURL string) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return nil, err
	}
	return DefaultClient.Do(req)
}

// Post is a convenient replacement for http.Post that adds a span around the request.
func Post(ctx context.Context, targetURL, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, "POST", targetURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return DefaultClient.Do(req)
}
