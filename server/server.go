package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"goscow/scripting"
	"net/http"
)

func handle(scriptFile string, context echo.Context) (interface{}, error) {
	argMap := make(map[string]interface{})
	argMap["c"] = context
	return scripting.JS(scriptFile, argMap)
}

func setupRoutes(e *echo.Echo, cfg *GoSCowConfig) error {
	verbHandlerMap := make(map[string]func(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route)
	// setup - how neat is this ?
	verbHandlerMap[CONNECT] = e.CONNECT
	verbHandlerMap[DELETE] = e.DELETE
	verbHandlerMap[GET] = e.GET
	verbHandlerMap[HEAD] = e.HEAD
	verbHandlerMap[OPTIONS] = e.OPTIONS
	verbHandlerMap[PATCH] = e.PATCH
	verbHandlerMap[POST] = e.POST
	verbHandlerMap[PUT] = e.PUT
	verbHandlerMap[TRACE] = e.TRACE
	// now setup shop
	for verb, routeInfo := range cfg.Routes {
		handlerFunction := verbHandlerMap[verb]
		for uri, script := range routeInfo.Table {
			handlerFunction(uri, func(c echo.Context) error {
				res, err := handle(script, c)
				if err != nil {
					return err
				}
				return c.String(http.StatusOK, fmt.Sprintf("%s", res))
			})
		}
	}
	return nil
}

func RunServer(configFile string) {
	cfg, err := From(configFile)
	if err != nil {
		panic(err)
	}
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	err = setupRoutes(e, cfg)
	if err != nil {
		panic(err)
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.Port)))
}
