package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		span.SetAttributes(attribute.Key("proxy-service.status_code").Int(http.StatusBadRequest))
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
		span.SetAttributes(attribute.Key("proxy-service.status_code").Int(http.StatusBadRequest))
		span.SetAttributes(attribute.Key("proxy-service.response.body").String(responseBodyString))
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
	span.SetAttributes(attribute.Key("proxy-service.status_code").Int(http.StatusOK))
	span.SetStatus(codes.Ok, "")
	c.JSON(http.StatusOK, albumStoreResponseBodyJson)
}

func getAlbumByID(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetName("/albums/:id GET")
	defer span.End()
	id := c.Param("id")
	span.SetAttributes(attribute.Key("proxy-service.request.parameters").String(fmt.Sprintf("%s=%s", "ID", id)))
	albumID, err := strconv.Atoi(id)
	// param ID is expected to be a number so fail if cannot covert to integer
	if err != nil {
		errorMessage := fmt.Sprintf("%s [%s] %s", "error invalid ID", id, "requested")
		span.SetStatus(codes.Error, errorMessage)
		span.AddEvent(errorMessage)
		span.SetAttributes(attribute.Key("proxy-service.status_code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": errorMessage})
		return
	}
	// proxy call to album-Store
	resp, err := Get(c.Request.Context(), fmt.Sprintf("%v/albums/%v", albumStoreURL, albumID))
	if err != nil {
		errorMessage := fmt.Sprintf("error contacting album-store getAlbumById %v", err)
		span.AddEvent(errorMessage)
		span.SetStatus(codes.Error, errorMessage)
		span.SetAttributes(attribute.Key("proxy-service.status_code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": errorMessage})
		return
	}
	var albumStoreResponseBodyJson interface{}
	byteArray, err := io.ReadAll(resp.Body)
	responseBodyString := string(byteArray[:])
	span.SetAttributes(attribute.Key("proxy-service.response.body").String(responseBodyString))
	err = json.NewDecoder(strings.NewReader(responseBodyString)).Decode(&albumStoreResponseBodyJson)
	if err != nil {
		errorMessage := "error from album-store getAlbumById malformed JSON"
		span.AddEvent(errorMessage)
		span.SetStatus(codes.Error, errorMessage)
		span.SetAttributes(attribute.Key("proxy-service.status_code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": errorMessage})
		return
	}
	err = resp.Body.Close()
	if err != nil {
		log.Println(err)
		errorMessage := fmt.Sprintf("error closing body from album-store getAlbumById %v", err)
		span.AddEvent(errorMessage)
		span.SetStatus(codes.Error, errorMessage)
		span.SetAttributes(attribute.Key("proxy-service.status_code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": errorMessage})
		return
	}
	span.SetAttributes(attribute.Key("proxy-service.response").Int(resp.StatusCode))
	span.SetAttributes(attribute.Key("proxy-service.status_code").Int(http.StatusOK))
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
	span.SetAttributes(attribute.Key("proxy-service.request.body").String(requestBodyString))
	//any problems with conversion of request body to JSON then return failure
	if err != nil {
		log.Println(err)
		errorMessage := fmt.Sprintf("invalid request json body %v", requestBodyString)
		span.AddEvent(errorMessage)
		span.SetStatus(codes.Error, errorMessage)
		span.SetAttributes(attribute.Key("proxy-service.request.body").String(requestBodyString))
		span.SetAttributes(attribute.Key("proxy-service.response.body").String(fmt.Sprintf("{\"message\":\"%v\"}", errorMessage)))
		span.SetAttributes(attribute.Key("proxy-service.status_code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": errorMessage})
		return
	}
	// proxy call to album-Store post album
	resp, err := Post(c.Request.Context(), albumStoreURL+"/albums/", "application/json", strings.NewReader(fmt.Sprintf("%v", requestBodyString)))
	if err != nil {
		errorMessage := fmt.Sprintf("error contacting album-store postAlbum %v", err)
		span.AddEvent(errorMessage)
		span.SetStatus(codes.Error, errorMessage)
		span.SetAttributes(attribute.Key("proxy-service.status_code").Int(http.StatusBadRequest))
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
		span.SetAttributes(attribute.Key("proxy-service.response.body").String(responseBodyString))
		span.SetAttributes(attribute.Key("proxy-service.status_code").Int(http.StatusBadRequest))
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
	span.SetAttributes(attribute.Key("proxy-service.response.body").String(responseBodyString))
	span.SetAttributes(attribute.Key("proxy-service.status_code").Int(http.StatusOK))
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
	log.Println(fmt.Sprintf("version: %v-%v", version, gitHash))
	shutdownTraceProvider, err := initOtelProvider(serviceName, version, gitHash)
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

// Set up the context for this Application in Open Telemetry
// application name, application version, k8s namespace , k8s instance name (horizontal scaling)
func setupOtelResource(serviceName string, version string, gitHash string, ctx context.Context, namespace *string, instanceName *string) (*resource.Resource, error) {

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(version+"-"+gitHash),
			semconv.ServiceNamespaceKey.String(*namespace),
			semconv.ServiceInstanceIDKey.String(*instanceName),
		),
	)
	return res, err
}

// InitOtelProvider - Initializes an OTLP exporter, and configures the corresponding trace and metric providers.
func initOtelProvider(serviceName string, version string, gitHash string) (func(context.Context) error, error) {
	ctx := context.Background()

	namespace := os.Getenv("NAMESPACE")
	instanceName := os.Getenv("INSTANCE_NAME")
	otelLocation := os.Getenv("OTEL_LOCATION")
	if instanceName == "" || otelLocation == "" || namespace == "" {
		log.Fatalf("Env variables not assigned NAMESPACE=%v, INSTANCE_NAME=%v, OTEL_LOCATION=%v", namespace, instanceName, otelLocation)
	}

	otelResource, err := setupOtelResource(serviceName, version, gitHash, ctx, &namespace, &instanceName)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	//setup for Protobuff - model.proto works with sending this to port 14250 in Jaeger
	otelTraceExporter, err := setupOtelProtoBuffGrpcTrace(ctx, &otelLocation)
	if err != nil {
		return nil, err
	}

	traceProvider := setupOtelTraceProvider(otelTraceExporter, otelResource)
	return traceProvider.Shutdown, nil //return shutdown signal so the application can trigger shutting itself down
}

func setupOtelTraceProvider(traceExporter *otlptrace.Exporter, otelResource *resource.Resource) *sdktrace.TracerProvider {
	// Register the trace exporter with a TracerProvider, using a batch span processor to aggregate spans before export.
	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(otelResource),
		sdktrace.WithSpanProcessor(batchSpanProcessor),
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{}) // set global propagator to tracecontext (the default is no-op).
	return tracerProvider
}

func setupOtelProtoBuffGrpcTrace(ctx context.Context, otelLocation *string) (*otlptrace.Exporter, error) {
	// insecure transport here. DO NOT USE IN PROD
	conn, err := grpc.DialContext(ctx, *otelLocation,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to opentelemetry-collector: %w", err)
	}
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create opentelemetry trace exporter: %w", err)
	}
	return traceExporter, nil
}
