package main

import (
	"log"
	"loquegasto-backend/internal/metrics"
	"loquegasto-backend/internal/router"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	r := router.New()

	metrics.StartServer("8081")

	if err := r.Run(); err != nil {
		log.Fatalln(err)
	}
}
