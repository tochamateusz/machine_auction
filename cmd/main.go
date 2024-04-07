package main

import (
	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-gonic/gin"
	"github.com/tochamateusz/machine_auction/infrastructure"
)

func main() {
	infrastructure.InitLogger()
	infrastructure.InitEnv()

	gin.DebugPrintRouteFunc = infrastructure.GinDebugPrintRouteFunc
	r := gin.New()
	r.Use(infrastructure.Logger)
	r.Use(ginzerolog.Logger("gin"))

	r.GET("test", func(ctx *gin.Context) {})

}
