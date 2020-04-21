package server

import (
	"github.com/foxcool/homework/5-k8s-crud/logger"
	"github.com/foxcool/homework/5-k8s-crud/storage"
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
