package main

import (
	"time"

	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://machine-auction-0-0-1-6dpaaunyfa-lm.a.run.app"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	health := health_controller.NewHandler()

	r.GET("/health", health.DbHealthCheck)
	r.Static("/scrapped/", "./scrapping-result/")
	r.Static("/backup/", "./db/")
	r.Static("/static/", "./web/dist/")
	r.Static("/assets/", "./web/dist/assets")
	scrapper_http.Init(r)

	server.InitServer(r)

}
