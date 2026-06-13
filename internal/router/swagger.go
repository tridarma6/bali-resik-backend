package router

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/indim/bali-resik-backend/docs"
)

func SetupSwagger(e *echo.Echo) {
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler())
}
