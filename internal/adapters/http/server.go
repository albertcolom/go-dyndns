package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	HttpServer *http.Server
}

func NewHTTPServer(handler *Handler, addr string) *Server {
	router := gin.New()
	router.GET("/health", handler.Health)
	router.GET("/update", handler.UpdateIp)
	router.GET("/get", handler.GetIp)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return &Server{
		HttpServer: httpServer,
	}
}

func (s *Server) Start() error {
	return s.HttpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.HttpServer.Shutdown(ctx)
}
