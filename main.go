package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var version = ""
var builddate = ""

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", index)

	e.Logger.Fatal(e.Start(":8080"))
}

func index(c echo.Context) error {
	return c.JSONBlob(http.StatusOK, []byte(fmt.Sprintf(`{"application": "alps", "version": "%s", "buildDate": "%s"}`, version, builddate)))
}
