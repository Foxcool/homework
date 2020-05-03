module github.com/foxcool/homework/7-prometheus/app

go 1.13

require (
	github.com/deepmap/oapi-codegen v1.3.7
	github.com/getkin/kin-openapi v0.3.1
	github.com/getsentry/sentry-go v0.6.0
	github.com/go-openapi/strfmt v0.19.5
	github.com/google/uuid v1.1.1
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/echo-contrib v0.9.0
	github.com/labstack/echo/v4 v4.1.16
	github.com/sirupsen/logrus v1.5.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.5.1 // indirect
	go.mongodb.org/mongo-driver v1.3.2
)

replace github.com/labstack/echo-contrib => github.com/foxcool/echo-contrib v0.9.1-0.20200501171235-b24a71759558
