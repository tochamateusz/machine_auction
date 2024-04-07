package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func cleanUp() chan struct{} {
	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl)

	done := make(chan struct{}, 1)

	go func() {
		for {
			sig := <-sigchnl
			switch sig {
			case os.Interrupt, os.Kill, syscall.SIGTERM:
				{
					done <- struct{}{}
				}
			default:
			}

		}
	}()

	return done
}

func InitServer(r *gin.Engine) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	WaitingForServer := func(done chan struct{}) {
		address := "localhost"
		address = "http://" + address
		address = address + "" + server.Addr

		completed := false
		for !completed {
			select {
			case <-done:
				completed = true
				break
			default:
				{
					log.Info().Msgf("Requesting: %s", address+"/health")
					req, _ := http.NewRequest("GET", address+"/health", nil)
					client := &http.Client{}
					resp, err := client.Do(req)
					if err != nil {
						log.Info().Msg("Waiting for server...")
						time.Sleep(time.Second)
						continue
					}

					if resp.StatusCode == http.StatusOK {
						log.Info().Msg("Server is running...")
						completed = true
						break
					}

				}
			}

		}
	}

	log.Info().Msg("Starting server")
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.WithLevel(zerolog.FatalLevel).Err(err).Msg("server crashed")
		}

		if err == http.ErrServerClosed {
			log.Err(err)
			log.Warn().Msg("Closing server")
			return
		}
		if err != nil {
			log.Err(err)
			return
		}
	}()

	done := cleanUp()
	WaitingForServer(done)
	<-done

	if err := server.Shutdown(context.Background()); err != nil {
		log.Info().Msgf("Server Shutdown: %s", err.Error())
	}

	log.Info().Msgf("Server exiting")

}
