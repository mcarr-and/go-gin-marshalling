package main

import (
	"context"
	"encoding/json"
	"errors"
	"example.com/album-store/models"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var albums = []models.Album{
	{ID: 1, Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: 2, Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: 3, Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func listAlbums() []models.Album {
	return albums
}

func resetAlbums() {
	albums = []models.Album{
		{ID: 1, Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		{ID: 2, Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
		{ID: 3, Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	}
}

func getAlbums(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetName("/albums GET")
	defer span.End()
	span.SetStatus(codes.Ok, "")
	span.SetAttributes(attribute.Key("album-store.response.code").Int(http.StatusOK))
	c.JSON(http.StatusOK, albums)
}

func getAlbumByID(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetName("/albums/:id GET")
	defer span.End()
	id := c.Param("id")
	span.SetAttributes(attribute.Key("album-store.request.parameters").String(fmt.Sprintf("%s=%s", "ID", id)))

	albumId, err := strconv.Atoi(id)
	if bindJsonToModelFails(c, err, id, span) {
		return
	}
	findAlbum(c, albumId, span)
}

func postAlbum(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetName("/albums POST")
	defer span.End()
	//c.ShouldBindBodyWith() // the old way to get the JSON body and did get body and bind
	requestBodyString, err := getRequestBody(c, span)
	if err {
		return
	}
	hasError, albumValue := bindJsonBody(c, span, requestBodyString)
	if hasError {
		return
	}
	albums = append(albums, albumValue)
	buildSuccessResponse(c, span, requestBodyString, albumValue)
}

func status(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetName("/status")
	span.SetStatus(codes.Ok, "")
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func metrics(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetName("/metrics")
	span.SetStatus(codes.Ok, "")
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

func findAlbum(c *gin.Context, albumId int, span trace.Span) {
	for _, album := range albums {
		if album.ID == albumId {
			span.SetStatus(codes.Ok, "")
			span.SetAttributes(attribute.Key("album-store.response.code").Int(http.StatusOK))
			jsonVal, _ := json.Marshal(album)
			span.SetAttributes(attribute.Key("album-store.response.body").String(string(jsonVal)))
			c.JSON(http.StatusOK, album)
			return
		}
	}
	errorMessage := fmt.Sprintf("Album [%v] not found", albumId)
	serverError := models.ServerError{Message: errorMessage}
	span.SetStatus(codes.Error, serverError.Message)
	span.AddEvent(errorMessage)
	span.SetAttributes(attribute.Key("album-store.response.code").Int(http.StatusBadRequest))
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": serverError.Message})
}

func bindJsonToModelFails(c *gin.Context, err error, id string, span trace.Span) bool {
	if err != nil {
		errorMessage := fmt.Sprintf("Album [%s] not found, invalid request", id)
		serverError := models.ServerError{Message: errorMessage}
		span.SetStatus(codes.Error, serverError.Message)
		span.AddEvent(errorMessage)
		span.SetAttributes(attribute.Key("album-store.response.code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": serverError.Message})
		return true
	}
	return false
}

func getRequestBody(c *gin.Context, span trace.Span) (string, bool) {
	var requestBody interface{}
	byteArray, err := io.ReadAll(c.Request.Body)
	requestBodyString := string(byteArray[:])
	if err = json.NewDecoder(strings.NewReader(requestBodyString)).Decode(&requestBody); err != nil {
		buildMalformedJsonErrorResponse(c, span, err, requestBodyString)
		return "", true
	}
	return requestBodyString, false
}

func buildSuccessResponse(c *gin.Context, span trace.Span, requestBodyString string, responseAlbum models.Album) {
	span.SetStatus(codes.Ok, "")
	span.SetAttributes(attribute.Key("album-store.request.body").String(requestBodyString))
	span.SetAttributes(attribute.Key("album-store.response.code").Int(http.StatusCreated))
	jsonByteArr, _ := json.Marshal(responseAlbum)
	span.SetAttributes(attribute.Key("album-store.response.body").String(string(jsonByteArr)))
	c.JSON(http.StatusCreated, responseAlbum)
}

func bindJsonBody(c *gin.Context, span trace.Span, requestBodyString string) (bool, models.Album) {
	var album models.Album
	if err := binding.JSON.BindBody([]byte(requestBodyString), &album); err != nil {
		if processValidationBindingError(c, err, span, requestBodyString) {
			return true, album
		}
	}
	return false, album
}

func buildMalformedJsonErrorResponse(c *gin.Context, span trace.Span, err error, requestBodyJSON string) bool {
	span.SetStatus(codes.Error, "Malformed JSON. Not valid for Album")
	span.AddEvent(fmt.Sprintf("Malformed JSON. %s", err))
	span.SetAttributes(attribute.Key("album-store.request.body").String(requestBodyJSON))
	span.SetAttributes(attribute.Key("album-store.response.body").String(`{"message":"Malformed JSON. Not valid for Album"}`))
	span.SetAttributes(attribute.Key("album-store.response.code").Int(http.StatusBadRequest))
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Malformed JSON. Not valid for Album"})
	return true
}

func processValidationBindingError(c *gin.Context, err error, span trace.Span, requestBodyJSON string) bool {
	var newAlbum models.Album
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		bindingErrorMessages := make([]models.BindingErrorMsg, len(ve))
		for i, fe := range ve {
			field, _ := reflect.TypeOf(&newAlbum).Elem().FieldByName(fe.Field())
			fieldJSONName, okay := field.Tag.Lookup("json")
			if !okay {
				log.Fatal(fmt.Sprintf("No json type on Struct model.Album %s Expecting : `json:\"title\" ...`", fe.Field()))
			}
			bindingErrorMessages[i] = models.BindingErrorMsg{Field: fieldJSONName, Message: getErrorMsg(fe)}
		}
		bindingErrorMessage, _ := json.Marshal(bindingErrorMessages)
		span.SetStatus(codes.Error, "Album JSON field validation failed")
		span.AddEvent(string(bindingErrorMessage))
		span.SetAttributes(attribute.Key("album-store.request.body").String(requestBodyJSON))
		span.SetAttributes(attribute.Key("album-store.response.body").String(fmt.Sprintf(`{"errors":%s}`, bindingErrorMessage)))
		span.SetAttributes(attribute.Key("album-store.response.code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": bindingErrorMessages})
		return true
	}
	return false
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "required field"
	case "min":
		return "below minimum value"
	case "gte":
		return "below minimum value"
	case "max":
		return "above maximum value"
	case "lte":
		return "above maximum value"
	default:
		return fmt.Sprintf("Unknown Error %s", fe.Tag())
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(otelgin.Middleware(serviceName)) // add OpenTelemetry to Gin
	router.Static("/v3/api-docs/", "cmd/api/swaggerui")
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbum)
	router.GET("/status", status)
	router.GET("/metrics", metrics)
	return router
}

const (
	serviceName  = "album-store"
	startAddress = "0.0.0.0:9080"
)

var version = "No-Version"
var gitHash = "No-Hash"

func main() {
	log.Println(fmt.Sprintf("version: %v-%v", version, gitHash))
	shutdownTraceProvider, err := initOtelProvider(serviceName, version, gitHash)
	if err != nil {
		log.Fatal(err)
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
