package main

import (
	"context"
	"encoding/json"
	"errors"
	"example/go-gin-example/models"
	"fmt"
	"go.opentelemetry.io/otel/codes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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
	span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
}

func getAlbumByID(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetName("/albums/:id GET")
	defer span.End()
	id := c.Param("id")
	span.SetAttributes(attribute.Key("Id").String(id))
	span.AddEvent("An Event", trace.WithAttributes(attribute.String("id", id)))

	if albumId, err := strconv.Atoi(id); err != nil {
		serverError := models.ServerError{Message: fmt.Sprintf("%s [%s] %s", "Album ID", id, "is not a valid number")}
		span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest), attribute.Key("http.request.id").String(id))
		span.SetStatus(codes.Error, serverError.Message)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": serverError.Message})
		return
	} else {
		for _, album := range albums {
			if album.ID == albumId {
				c.JSON(http.StatusOK, album)
				span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusOK))
				return
			}
		}
	}
	serverError := models.ServerError{Message: fmt.Sprintf("%s [%s] %s", "Album", id, "not found")}
	span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
	span.SetStatus(codes.Error, serverError.Message)
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
			span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
			span.SetStatus(codes.Error, fmt.Sprintf("%v", string(jsonBytes)))
			return
		}
	}
	albums = append(albums, newAlbum)
	span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusOK))
	c.JSON(http.StatusCreated, newAlbum)
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is Required"
	}
	return "Unknown Error"
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
	service = "album-store"
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
	openTelemetryCollectorServiceLocation := getEnvironmentValue("OTEL_SERVICE_LOCATION", "localhost:4317")
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
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := initProvider()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	// Attributes represent additional key-value descriptors that can be bound
	// to a metric observer or recorder.
	commonAttrs := []attribute.KeyValue{
		attribute.String("attrA", "chocolate"),
		attribute.String("attrB", "raspberry"),
		attribute.String("attrC", "vanilla"),
	}

	tracer := otel.Tracer("test-tracer")
	ctx, span := tracer.Start(
		ctx,
		"CollectorExporter-Example",
		trace.WithAttributes(commonAttrs...))
	defer span.End()

	router := setupRouter()

	runAddress := getEnvironmentValue("ALBUM_START_URL", "localhost:9080")
	if errRun := router.Run(runAddress); errRun != nil {
		log.Fatal("Could not start server on ", runAddress)
		return
	}
}

func getEnvironmentValue(searchValue, defaultValue string) string {
	envValue, returnedValue := os.LookupEnv(searchValue)
	if !returnedValue {
		return defaultValue
	}
	return envValue
}
