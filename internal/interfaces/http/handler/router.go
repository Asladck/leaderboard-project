// Package handler HTTP handlers and server.
// @title Example API
// @version 1.0
// @description This is an example API server.
// @termsOfService http://example.com/terms/
// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
package handler

import (
	"context"
	"net/http"
	"time"
)

// Server wraps the underlying http.Server.
type Server struct {
	httpServer *http.Server
}

// Run starts the HTTP server on the given port with the provided handler.
func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server.
func (s *Server) Shutdown() error {
	return s.httpServer.Shutdown(context.Background())
}
