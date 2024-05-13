package infrastructure

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gookit/goutil/dump"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func GinDebugPrintRouteFunc(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	log.Info().
		Str("Path", absolutePath).
		Str("Method", httpMethod).
		Str("Handler", handlerName).
		Msgf("(%d handlers)", nuHandlers)
}

func InitLogger() {
	dev := os.Getenv("DEV")
	dump.P("DEV:", dev)
	noColor := false
	if dev == "" {
		noColor = true
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: noColor})

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func Logger(c *gin.Context) {
	c.Next()

	if c.Writer.Status() >= 400 {
		log.Warn().Strs("errors", c.Errors.Errors()).Msgf("%d - %s", c.Writer.Status(), c.Request.URL)
		if c.Writer.Status() == 500 {
			log.Error().Strs("errors", c.Errors.Errors()).Msgf("%d - %s", c.Writer.Status(), c.Request.URL)
		}

	} else {
		log.Debug().Msgf("%d - %s", c.Writer.Status(), c.Request.URL)
	}
}
