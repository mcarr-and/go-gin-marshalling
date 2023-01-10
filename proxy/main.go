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
		span.SetStatus(codes.Error, fmt.Sprintf("error contacting album-store getAlbums %v", err))
		span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error calling album-store getAlbums"})
		return
	}
	var albumStoreResponseBodyJson interface{}
	err = json.NewDecoder(resp.Body).Decode(&albumStoreResponseBodyJson)
	if err != nil {
		span.SetStatus(codes.Error, "error from album-store getAlbums malformed JSON")
		span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error from album-store getAlbums malformed JSON"})
		return
	}
	err = resp.Body.Close()
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("error album-store closing response %v", err))
		span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error album-store closing response"})
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
	// proxy call to album-Store
	resp, err := Get(c.Request.Context(), albumStoreURL+"/albums/"+id)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("error contacting album-store getAlbumById %v", err))
		span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error calling album-store"})
		return
	}
	var albumStoreResponseBodyJson interface{}
	err = json.NewDecoder(resp.Body).Decode(&albumStoreResponseBodyJson)
	if err != nil {
		span.SetStatus(codes.Error, "error from album-store getAlbumById malformed JSON")
		span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error getting body from album-store getAlbumById malformed JSON"})
		return
	}
	err = resp.Body.Close()
	if err != nil {
		log.Println(err)
		span.SetStatus(codes.Error, fmt.Sprintf("error closing body from album-store getAlbumById %v", err))
		span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error getting body from album-store getAlbumById"})
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

	byteArray := make([]byte, 500)
	_, err := c.Request.Body.Read(byteArray)
	if err != nil {
		log.Println(err)
		span.SetStatus(codes.Error, fmt.Sprintf("could not get Post body %v", err))
	}
	str1 := string(byteArray[:])
	// proxy call to album-Store
	resp, err := Post(c.Request.Context(), albumStoreURL+"/albums/", "application/json", strings.NewReader(str1))
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("error contacting album-store postAlbum %v", err))
		span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error contacting album-store postAlbum"})
		return
	}
	var albumStoreResponseBodyJson interface{}
	err = json.NewDecoder(resp.Body).Decode(&albumStoreResponseBodyJson)
	if err != nil {
		span.SetStatus(codes.Error, "error from album-store postAlbum malformed JSON")
		span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error getting body from album-store postAlbum malformed JSON"})
		return
	}
	err = resp.Body.Close()
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("error closing body from album-store postAlbum %v", err))
		span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error closing body from album-store postAlbum"})
		return
	}
	span.SetAttributes(attribute.Key("proxy-service.response").Int(resp.StatusCode))
	span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusOK))
	span.SetStatus(codes.Ok, "")
	c.JSON(http.StatusOK, albumStoreResponseBodyJson)
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
