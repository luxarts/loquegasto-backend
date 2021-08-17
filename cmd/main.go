package main

import (
	"github.com/joho/godotenv"
	"log"
	"loquegasto-backend/internal/router"
)

func main() {
	_ = godotenv.Load()

	r := router.New()

	if err := r.Run(); err != nil {
		log.Fatalln(err)
	}
}
