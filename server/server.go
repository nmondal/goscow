package server

import (
	"fmt"
	"goscow/scripting"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func setupRoutes(e *echo.Echo, cfg *GoSCowConfig) error {
	verbHandlerMap := make(map[VerbType]func(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route)
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
		for uri, scriptInfo := range routeInfo {
			script := scriptInfo
			rtype := ""
			rtarr := strings.Split(scriptInfo, "@")
			if len(rtarr) != 1 {
				script = rtarr[1]
				rtype = rtarr[0]
			}

			if _, err := os.Stat(script); os.IsNotExist(err) {
				return err
			}
			handlerFunction(uri, func(context echo.Context) error {
				res, err := scripting.Execute(cfg.Reload, script, context)
				if err != nil {
					return err
				}
				switch rtype {
				case JSON:
					return context.JSON(http.StatusOK, res)
				case JSON_PRETTY:
					return context.JSONPretty(http.StatusOK, res, "\t")
				// Should we use XML at all? Possibly not
				default:
					return context.String(http.StatusOK, fmt.Sprintf("%s", res))
				}
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
