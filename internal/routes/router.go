package routes

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/klaxn/klaxn-api/internal/config"
	"github.com/klaxn/klaxn-api/internal/data"
	"github.com/klaxn/klaxn-api/pkg/outbound"
)

type Router struct {
	db     *data.Manager
	logger logrus.FieldLogger
	tracer trace.Tracer
	ob     []outbound.Sender
	cfg    *config.Config
}

func New(db *data.Manager, ob []outbound.Sender, logger logrus.FieldLogger, cfg *config.Config) *Router {
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(cfg.App.Name),
	)

	ctx := context.Background()

	initTracerProvider(ctx, res)

	//responseTimeInstrument, err := meter.SyncInt64().Histogram("http.server.duration")
	//if err != nil {
	//	log.Panicf("failed to initialize instrument: %v", err)
	//}

	return &Router{
		db:     db,
		logger: logger,
		tracer: otel.GetTracerProvider().Tracer("my-tracer"),
		ob:     ob,
		cfg:    cfg,
	}
}

func initTracerProvider(ctx context.Context, res *resource.Resource) {
	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		log.Fatalf("%s: %v", "failed to create metric exporter", err)
	}

	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(res),
	)
	otel.SetTracerProvider(tracerProvider)
}

// Deprecated: Please use `SendErr` instead
func (r *Router) SendErrorMessage(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"message": message,
	})
}

func (r *Router) SendErr(c *gin.Context, status int, err error, span trace.Span) {
	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err)
	if status != http.StatusBadRequest {
		r.logger.Errorf("Sending error to client: %s", err.Error())
	}
	r.SendJsonResponse(c, status, gin.H{
		"message": err.Error(),
	}, span)
}

func (r *Router) SendJsonResponse(c *gin.Context, status int, h interface{}, span trace.Span) {
	attributes := []attribute.KeyValue{
		semconv.HTTPMethodKey.String(c.Request.Method),
		semconv.HTTPStatusCodeKey.Int(status),
		semconv.HTTPUserAgentKey.String(c.Request.UserAgent()),
		semconv.HTTPURLKey.String(c.Request.URL.String()),
	}
	span.SetAttributes(attributes...)
	c.JSON(status, h)
}

func (r *Router) SendNullResponse(c *gin.Context, status int, span trace.Span) {
	attributes := []attribute.KeyValue{
		semconv.HTTPMethodKey.String(c.Request.Method),
		semconv.HTTPStatusCodeKey.Int(status),
		semconv.HTTPUserAgentKey.String(c.Request.UserAgent()),
		semconv.HTTPURLKey.String(c.Request.URL.String()),
	}
	span.SetAttributes(attributes...)
	c.Status(status)
}

// GetConfig godoc
//
//	@Summary	Get currently running config
//	@Schemes
//	@Description	Get currently running config
//	@Tags			debug
//	@Produce		json
//	@Success		200	{object}	config.Config
//	@Router			/debug/config [get]
func (r *Router) GetConfig(c *gin.Context) {
	_, span := r.tracer.Start(c, "GetConfig")
	defer span.End()

	r.SendJsonResponse(c, http.StatusOK, r.cfg, span)
}
