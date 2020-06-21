package server

import (
	"github.com/foxcool/homework/7-prometheus/app/logger"
	"github.com/foxcool/homework/7-prometheus/app/storage"
	"github.com/labstack/echo"
	"gopkg.in/hlandau/passlib.v1"
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

// VerifyPassword validates password and regenerate hash
func VerifyPassword(password, hash string) (string, error) {
	passlib.UseDefaults(passlib.DefaultsLatest)

	return passlib.Verify(password, hash)
}

// HashPassword returns password hash
func HashPassword(password string) (string, error) {
	passlib.UseDefaults(passlib.DefaultsLatest)

	return passlib.Hash(password)
}

