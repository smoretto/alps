package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	log "github.com/sirupsen/logrus"
)

var version = ""
var builddate = ""

func main() {
	log.SetFormatter(&log.JSONFormatter{TimestampFormat: "2006-01-02T15:04:05.000Z0700"})
	log.Infof("started alps version: %s build-date: %s", version, builddate)

	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/health"
		},
		Format: `{"time":"${time_custom}","log":"access","id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency":${latency},` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
		CustomTimeFormat: "2006-01-02T15:04:05.000Z0700",
	}))

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	e.GET("/", index)
	e.GET("/health", health)

	e.Logger.Fatal(e.Start(":8080"))
}

func index(c echo.Context) error {
	return c.JSONBlob(http.StatusOK, []byte(fmt.Sprintf(`{"application": "alps", "version": "%s", "buildDate": "%s"}`, version, builddate)))
}

func health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"alive": true,
	})
}
