package main

import (
	"context"
	"encoding/json"
	"errors"
	"example/go-gin-example/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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
	span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
	defer span.End()
	c.JSON(http.StatusOK, albums)
}

func getAlbumByID(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	defer span.End()
	id := c.Param("id")
	span.SetAttributes(attribute.Key("Id").String(id))
	span.AddEvent("An Event", trace.WithAttributes(attribute.String("id", id)))

	if albumId, err := strconv.Atoi(id); err != nil {
		serverError := models.ServerError{Message: fmt.Sprintf("%s [%s] %s", "Album ID", id, "is not a valid number")}
		span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
		span.SetAttributes(attribute.Key("http.request.id").String(id))
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
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": serverError.Message})
}

func postAlbum(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
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
			span.SetAttributes(attribute.Key("http.status_code").Int(http.StatusBadRequest))
			jsonBytes, _ := json.Marshal(bindingErrorMessages)
			span.SetAttributes(attribute.Key("postAlbum.error.binding.message").String(fmt.Sprintf("%v", string(jsonBytes))))
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
	service     = "album-store"
	environment = "development"
	id          = 1
)

// tracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	)
	return tp, nil
}

func main() {
	tp, err := tracerProvider(getEnvironmentValue("JAEGER_TRACES_URL", "http://localhost:14268/api/traces"))
	if err != nil {
		log.Fatal(err)
	}

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cleanly shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	router := setupRouter()

	runAddress := getEnvironmentValue("ALBUM_START_URL", "localhost:9080")
	if errRun := router.Run(runAddress); errRun != nil {
		log.Fatal("Could not start server on ", runAddress)
		return
	}
}

func getEnvironmentValue(searchValue, defaultValue string) string {
	envValue := os.Getenv(searchValue)
	if envValue == "" {
		return defaultValue
	}
	return envValue
}
