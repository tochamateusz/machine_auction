package http

import (
	"fmt"
	"os"
	"time"

	_ "image/jpeg"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/tochamateusz/machine_auction/infrastructure/events"
	"github.com/xuri/excelize/v2"
)

type HttpExporterApi struct {
	eventBus events.IEventBus
}

func checkErr(err error) {
	if err != nil {
		log.Fatal().Err(err).Msgf("")
	}
}

func Init(r *gin.Engine) {
	data, err := os.ReadFile("/home/toch/develop/golang/machine_auction/backend/Book1.xlsx")
	checkErr(err)
	// Write data to dst
	err = os.WriteFile("/home/toch/develop/golang/machine_auction/backend/Book1-"+time.Now().Format(time.DateOnly)+".xlsx", data, 0644)
	checkErr(err)
	f, err := excelize.OpenFile("/home/toch/develop/golang/machine_auction/backend/Book1.xlsx", excelize.Options{Password: ""})
	if err != nil {
		log.Err(err).Msgf("cann't init exporter")
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Insert a picture to worksheet with scaling.
	if err := f.AddPicture("Sheet1", "D2", "/home/toch/develop/golang/machine_auction/backend/scrapping-result/10000/0.jpg",
		&excelize.GraphicOptions{ScaleX: 0.5, ScaleY: 0.5}); err != nil {
		log.Err(err).Msgf("cann't init exporter")
	}
	// Save the spreadsheet with the origin path.
	if err = f.Save(); err != nil {
		log.Err(err).Msgf("cann't init exporter")
	}
}
