module github.com/foxcool/homework/7-prometheus/app

go 1.13

require (
	github.com/deepmap/oapi-codegen v1.3.8
	github.com/getkin/kin-openapi v0.9.0
	github.com/getsentry/sentry-go v0.6.0
	github.com/go-openapi/strfmt v0.19.5
	github.com/google/uuid v1.1.1
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/echo-contrib v0.9.0
	github.com/labstack/echo/v4 v4.1.16
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.5.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.5.1 // indirect
	go.mongodb.org/mongo-driver v1.3.2
	golang.org/x/crypto v0.0.0-20200604202706-70a84ac30bf9 // indirect
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9 // indirect
	golang.org/x/sys v0.0.0-20200602225109-6fdc65e7d980 // indirect
	gopkg.in/hlandau/easymetric.v1 v1.0.0 // indirect
	gopkg.in/hlandau/measurable.v1 v1.0.1 // indirect
	gopkg.in/hlandau/passlib.v1 v1.0.10
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

replace github.com/labstack/echo-contrib => github.com/foxcool/echo-contrib v0.9.1-0.20200501171235-b24a71759558
