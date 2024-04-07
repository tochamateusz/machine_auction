package main

import (
	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-gonic/gin"
	"github.com/mandrigin/gin-spa/spa"
	health_controller "github.com/tochamateusz/machine_auction/app"
	scrapper_http "github.com/tochamateusz/machine_auction/app/scrapper/adapters/http"
	"github.com/tochamateusz/machine_auction/infrastructure"
	"github.com/tochamateusz/machine_auction/infrastructure/server"
)

func main() {
	infrastructure.InitLogger()
	infrastructure.InitEnv()

	gin.DebugPrintRouteFunc = infrastructure.GinDebugPrintRouteFunc
	r := gin.New()
	r.Use(infrastructure.Logger)
	r.Use(ginzerolog.Logger("gin"))

	health := health_controller.NewHandler()

	r.GET("health", health.DbHealthCheck)
	scrapper_http.Init(r)

	r.GET("/", spa.Middleware("/", "./web/build"))

	server.InitServer(r)

}
