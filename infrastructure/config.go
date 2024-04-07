package infrastructure

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func InitEnv() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Printf("Can't load env's %+v\n", err)
	}

	envLogin := os.Getenv("LOGIN")
	fmt.Printf("LOGIN: %+v\n", envLogin)
	envPassword := os.Getenv("PASSWORD")
	fmt.Printf("PASSWORD: %+v\n", envPassword)
}
