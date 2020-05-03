package server

import (
	"github.com/foxcool/homework/7-prometheus/app/logger"
	"github.com/foxcool/homework/7-prometheus/app/storage"
	"github.com/labstack/echo"
)

var (
	SuccessMessage = "SUCCESS"
)

type Server struct {
	Version string

	storage *storage.Storage

	authLog logger.AuthEntrier

	e *echo.Echo
}

func New(version string, storage *storage.Storage, l logger.Logger) *Server {
	return &Server{
		Version: version,

		storage: storage,

		authLog: l.Auth(),
	}
}
