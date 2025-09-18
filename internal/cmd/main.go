package main

import (
	"learning-cards/internal/startup"
	"log"
)

func main() {
	if err := startup.Run(); err != nil {
		log.Fatal(err)
	}
}
