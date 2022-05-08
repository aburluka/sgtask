package server

import (
	"bytes"
	"encoding/json"
	"io"
	"runtime/debug"

	"github.com/aburluka/sgtask/internal/logger"
	"github.com/aburluka/sgtask/internal/storage"
	"github.com/aburluka/sgtask/pkg/model/event"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Server struct {
	logger  *logger.Logger
	storage *storage.Client
}

func New(logger *logger.Logger, storage *storage.Client) *Server {
	return &Server{
		logger:  logger,
		storage: storage,
	}
}

func (s *Server) event(reqCtx *fasthttp.RequestCtx) {
	var events []*event.Event
	dec := json.NewDecoder(bytes.NewReader(reqCtx.Request.Body()))
	for {
		var event event.Event
		err := dec.Decode(&event)
		if err == io.EOF {
			break
		}
		if err != nil {
			s.logger.With(zap.Error(err)).Error("malformed event")
			reqCtx.Error("malformed event", fasthttp.StatusBadRequest)
			return
		}
		event.Enrich()
		events = append(events, &event)
	}

	s.storage.PutEvents(events)
	reqCtx.SetStatusCode(fasthttp.StatusOK)
}

func panicHandler(logger *logger.Logger) func(*fasthttp.RequestCtx, interface{}) {
	return func(requestCtx *fasthttp.RequestCtx, v interface{}) {
		requestCtx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		if err, ok := v.(error); ok {
			logger.WithOptions(zap.AddStacktrace(zapcore.ErrorLevel)).Error("panic recovered in request handler",
				zap.Any("panic", err),
				zap.ByteString("url", requestCtx.URI().FullURI()),
				zap.ByteString("stacktrace", debug.Stack()),
			)
		} else {
			logger.WithOptions(zap.AddStacktrace(zapcore.ErrorLevel)).Error("panic recovered in request handler",
				zap.Any("panic", err),
				zap.ByteString("url", requestCtx.URI().FullURI()),
				zap.ByteString("stacktrace", debug.Stack()),
			)
		}
	}
}

func (s *Server) Run() {
	router := router.New()
	router.PanicHandler = panicHandler(s.logger)

	router.POST("/event", s.event)

	httpServer := &fasthttp.Server{
		Handler: router.Handler,
	}
	addr := "0.0.0.0:12345"
	s.logger.With(zap.String("addr", addr)).Info(("http server is starting..."))
	err := httpServer.ListenAndServe("0.0.0.0:12345")
	if err != nil {
		s.logger.Fatal("error while starting http server", zap.String("addr", addr), zap.Error(err))
	}
}
