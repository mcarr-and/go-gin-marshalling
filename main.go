package main

import (
	"context"
	"encoding/json"
	"errors"
	"example/go-gin-example/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	setStatusOnSpan(span, http.StatusOK, codes.Ok, okMessage)
	c.JSON(http.StatusOK, albums)
}

func getAlbumByID(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetName("/albums/:id GET")
	defer span.End()
	id := c.Param("id")
	span.SetAttributes(attribute.Key("Id").String(id))

	albumId, err := strconv.Atoi(id)
	if err != nil {
		serverError := models.ServerError{Message: fmt.Sprintf("%s [%s] %s", "Album ID", id, "is not a valid number")}
		setStatusOnSpan(span, http.StatusBadRequest, codes.Error, serverError.Message)
		addSpanEventAndLog(span, fmt.Sprintf("Get /album invalid ID %s", id))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": serverError.Message})
		return
	}

	for _, album := range albums {
		if album.ID == albumId {
			c.JSON(http.StatusOK, album)
			setStatusOnSpan(span, http.StatusOK, codes.Ok, okMessage)
			return
		}
	}

	serverError := models.ServerError{Message: fmt.Sprintf("%s [%s] %s", "Album", id, "not found")}
	setStatusOnSpan(span, http.StatusBadRequest, codes.Error, serverError.Message)
	addSpanEventAndLog(span, fmt.Sprintf("Get /album not found with ID %s", id))
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": serverError.Message})
}

func postAlbum(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetName("/albums POST")
	defer span.End()
	var newAlbum models.Album

	if err := c.ShouldBindJSON(&newAlbum); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			bindingErrorMessages := make([]models.BindingErrorMsg, len(ve))
			for i, fe := range ve {
				bindingErrorMessages[i] = models.BindingErrorMsg{Field: fe.Field(), Message: getErrorMsg(fe)}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": bindingErrorMessages})
			jsonBytes, _ := json.Marshal(bindingErrorMessages)
			addSpanEventAndLog(span, string(jsonBytes))
		} else {
			addSpanEventAndLog(span, fmt.Sprintf("%s", err))
		}
		setStatusOnSpan(span, http.StatusBadRequest, codes.Error, "could not bind JSON posted to method")
		return
	}
	albums = append(albums, newAlbum)
	setStatusOnSpan(span, http.StatusOK, codes.Ok, okMessage)
	c.JSON(http.StatusCreated, newAlbum)
}

func setStatusOnSpan(span trace.Span, httpStatusCode int, spanStatusCode codes.Code, message string) {
	span.SetAttributes(attribute.Key("http.status_code").Int(httpStatusCode))
	span.SetStatus(spanStatusCode, message)
}

func addSpanEventAndLog(span trace.Span, errorMsg string) {
	span.AddEvent(errorMsg)
	log.Println(errorMsg)
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
	router.Use(otelgin.Middleware(service)) // weave in OpenTelemetry
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbum)
	return router
}

const (
	service   = "album-store"
	okMessage = "OK"
)

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func initProvider() (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(service),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	openTelemetryCollectorServiceLocation := getEnvironmentValue("OTEL_SERVICE_LOCATION", "localhost:4327")
	conn, err := grpc.DialContext(ctx, openTelemetryCollectorServiceLocation,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to opentelemetry-collector: %w", err)
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}

func main() {
	shutdownTraceProvider, err := initProvider()
	if err != nil {
		log.Fatal(err)
	}

	router := setupRouter()
	srv := &http.Server{
		Addr:    getEnvironmentValue("ALBUM_SERVICE_URL", "localhost:9080"),
		Handler: router,
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

	log.Println("OpenTelemetry TraceProvider shutting down")
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

func getEnvironmentValue(searchValue, defaultValue string) string {
	envValue, returnedValue := os.LookupEnv(searchValue)
	if !returnedValue {
		return defaultValue
	}
	return envValue
}
