package main

import (
	"log"
	"loquegasto-backend/internal/router"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	r := router.New()

	if err := r.Run(); err != nil {
		log.Fatalln(err)
	}
}
