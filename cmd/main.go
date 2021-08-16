package main

import (
	"log"
	"loquegasto-backend/internal/router"
)

func main(){
	r := router.New()

	if err := r.Run(); err != nil {
		log.Fatalln(err)
	}
}
