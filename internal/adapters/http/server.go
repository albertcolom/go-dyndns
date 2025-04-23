package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-dyndns/internal/adapters/http/middleware"
	"net/http"
)

type Server struct {
	HttpServer *http.Server
}

func NewHTTPServer(handler *Handler, addr, token string) *Server {
	router := gin.New()
	router.Use(middleware.RequestIdMiddleware())
	router.Use(gin.Recovery())
	v1 := router.Group("/v1")
	{
		v1.GET("/health", handler.Health)
		protected := v1.Group("").Use(middleware.AuthMiddleware(token))
		{
			protected.GET("/update", handler.UpdateIp)
			protected.GET("/get", handler.GetIp)
		}
	}

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
