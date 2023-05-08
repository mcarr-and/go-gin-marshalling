package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"

	_ "github.com/mcarr-and/go-gin-otelcollector/album-store/api"
	"github.com/mcarr-and/go-gin-otelcollector/album-store/model"

	"github.com/gin-gonic/gin/binding"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.17.0"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var albums = []model.Album{
	{ID: 1, Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: 2, Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: 3, Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func listAlbums() []model.Album {
	return albums
}

func resetAlbums() {
	albums = []model.Album{
		{ID: 1, Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		{ID: 2, Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
		{ID: 3, Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	}
}

// @title           Album Store API
// @version         1.0
// @description     Simple golang album store CRUD application
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:9080
// @BasePath /

// GetAlbums godoc
// @Summary Get all Albums
// @Schemes
// @Description get all the albums in the store
// @Tags albums
// @Produce json
// @Success 200 {array} model.Album
// @Router /albums [get]
func getAlbums(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetName("/albums GET")
	defer span.End()
	span.SetStatus(codes.Ok, "")
	span.SetAttributes(attribute.Key("album-store.response.code").Int(http.StatusOK))
	c.JSON(http.StatusOK, albums)
}

// GetAlbumById godoc
// @Summary Get Album by id
// @Schemes
// @Description get as single album by id
// @Tags albums
// @Param  id query int true  "int valid" minimum(1)
// @Produce json
// @Success 200 {object} model.Album
// @Failure 400 {object} model.ServerError
// @Router /albums/{id} [get]
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

// PostAlbum godoc
// @Summary Create album
// @Schemes
// @Description add a new album to the store
// @Tags albums
// @Param request body model.Album true "album"
// @Accept json
// @Produce json
// @Success 201 {object} model.Album
// @Failure 400 {object} model.ServerError
// @Router /albums [post]
func postAlbum(log zerolog.Logger) gin.HandlerFunc {
	fn := func(context *gin.Context) {
		span := trace.SpanFromContext(context.Request.Context())
		span.SetName("/albums POST")
		defer span.End()
		//c.ShouldBindBodyWith() // the old way to get the JSON body and did get body and bind
		requestBodyString, errBody := getRequestBody(context, span)
		if errBody {
			return
		}
		hasError, albumValue := bindJsonBody(context, span, requestBodyString, log)
		if hasError {
			return
		}
		albums = append(albums, albumValue)

		buildSuccessResponse(context, span, requestBodyString, albumValue)
	}
	return fn
}

// Status godoc
// @Summary Status of service
// @Schemes
// @Description get the status of the service
// @Tags albums
// @Produce json
// @Success 200 {string} status
// @Router /status [get]
func status(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetName("/status")
	span.SetStatus(codes.Ok, "")
	defer span.End()
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

// Metrics godoc
// @Summary Prometheus metrics
// @Schemes
// @Description get Prometheus metrics for the service
// @Tags albums
// @Produce plain
// @Success 200 {string} metrics
// @Router /status [get]
func metrics(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetName("/metrics")
	span.SetStatus(codes.Ok, "")
	defer span.End()
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
	serverError := model.ServerError{Message: errorMessage}
	span.SetStatus(codes.Error, serverError.Message)
	span.AddEvent(errorMessage)
	span.SetAttributes(attribute.Key("album-store.response.code").Int(http.StatusBadRequest))
	c.AbortWithStatusJSON(http.StatusBadRequest, serverError)
}

func bindJsonToModelFails(c *gin.Context, err error, id string, span trace.Span) bool {
	if err != nil {
		errorMessage := fmt.Sprintf("Album [%s] not found, invalid request", id)
		serverError := model.ServerError{Message: errorMessage}
		span.SetStatus(codes.Error, serverError.Message)
		span.AddEvent(errorMessage)
		// span.RecordError(err, )// todo - figure out when to use this instead of event
		span.SetAttributes(attribute.Key("album-store.response.code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, serverError)
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

func buildSuccessResponse(c *gin.Context, span trace.Span, requestBodyString string, responseAlbum model.Album) {
	span.SetStatus(codes.Ok, "")
	span.SetAttributes(attribute.Key("album-store.request.body").String(requestBodyString))
	span.SetAttributes(attribute.Key("album-store.response.code").Int(http.StatusCreated))
	jsonByteArr, _ := json.Marshal(responseAlbum)
	span.SetAttributes(attribute.Key("album-store.response.body").String(string(jsonByteArr)))
	c.JSON(http.StatusCreated, responseAlbum)
}

func bindJsonBody(c *gin.Context, span trace.Span, requestBodyString string, log zerolog.Logger) (bool, model.Album) {
	var album model.Album
	if err := binding.JSON.BindBody([]byte(requestBodyString), &album); err != nil {
		if processValidationBindingError(c, err, span, requestBodyString, log) {
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
	c.AbortWithStatusJSON(http.StatusBadRequest, model.ServerError{Message: "Malformed JSON. Not valid for Album"})
	return true
}

func processValidationBindingError(c *gin.Context, err error, span trace.Span, requestBodyJSON string, log zerolog.Logger) bool {
	var newAlbum model.Album
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		bindingErrorMessages := make([]*model.BindingErrorMsg, len(validationErrors))
		for index, fieldError := range validationErrors {
			field, _ := reflect.TypeOf(&newAlbum).Elem().FieldByName(fieldError.Field())
			fieldJSONName, okay := field.Tag.Lookup("json")
			if !okay {
				log.Fatal().Msg(fmt.Sprintf("No json type on Struct model.Album %s Expecting : `json:\"title\" ...`", fieldError.Field()))
			}
			bindingErrorMessages[index] = &model.BindingErrorMsg{Field: fieldJSONName, Message: getErrorMsg(fieldError)}
		}
		bindingErrorMessage, _ := json.Marshal(bindingErrorMessages)
		span.SetStatus(codes.Error, "Album JSON field validation failed")
		span.AddEvent(string(bindingErrorMessage))
		span.SetAttributes(attribute.Key("album-store.request.body").String(requestBodyJSON))
		span.SetAttributes(attribute.Key("album-store.response.body").String(fmt.Sprintf(`{"errors":%s}`, bindingErrorMessage)))
		span.SetAttributes(attribute.Key("album-store.response.code").Int(http.StatusBadRequest))
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ServerError{BindingErrors: bindingErrorMessages})
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
	case "max":
		return "above maximum value"
	default:
		return fmt.Sprintf("Unknown Error %s", fe.Tag())
	}
}

func setupRouter(log zerolog.Logger) *gin.Engine {
	router := gin.Default()
	router.Use(otelgin.Middleware(serviceName)) // add OpenTelemetry to Gin
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbum(log))
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
	logError := zerolog.New(os.Stderr).With().Timestamp().Logger()
	logInfo := zerolog.New(os.Stdout).With().Timestamp().Logger()

	logInfo.Info().Msg(fmt.Sprintf("version: %v-%v", version, gitHash))
	shutdownTraceProvider, err := initOtelProvider(serviceName, version, gitHash, logInfo)
	if err != nil {
		logError.Fatal().Err(err)
	}

	router := setupRouter(logInfo)
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
			logError.Fatal().Err(err)
		}
	}()
	<-quit

	logInfo.Info().Msg("Server shutdown with 500ms timeout...")
	ctxServer, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	logInfo.Info().Msg("OpenTelemetry TraceProvider flushing & shutting down")
	if err := shutdownTraceProvider(ctxServer); err != nil {
		logError.Fatal().Err(err)
	}
	logInfo.Info().Msg("OpenTelemetry TraceProvider exited")

	if err := srv.Shutdown(ctxServer); err != nil {
		logError.Fatal().Err(err)
	}
	<-ctxServer.Done()

	logInfo.Info().Msg("Server exiting")
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
func initOtelProvider(serviceName string, version string, gitHash string, log zerolog.Logger) (func(context.Context) error, error) {
	ctx := context.Background()

	namespace := os.Getenv("NAMESPACE")
	instanceName := os.Getenv("INSTANCE_NAME")
	otelLocation := os.Getenv("OTEL_LOCATION")
	if instanceName == "" || otelLocation == "" || namespace == "" {
		log.Fatal().Msg(fmt.Sprintf("Env variables not assigned NAMESPACE=%v, INSTANCE_NAME=%v, OTEL_LOCATION=%v", namespace, instanceName, otelLocation))
	}

	otelResource, err := setupOtelResource(serviceName, version, gitHash, ctx, &namespace, &instanceName)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	otelTraceExporter, err := setupOtelHttpTrace(ctx, &otelLocation)
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

func setupOtelHttpTrace(ctx context.Context, otelLocation *string) (*otlptrace.Exporter, error) {
	// insecure transport here DO NOT USE IN PROD
	client := otlptracehttp.NewClient(
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(*otelLocation),
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
	)
	err := client.Start(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start http client: %w", err)
	}
	traceExporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}
	return traceExporter, nil
}
