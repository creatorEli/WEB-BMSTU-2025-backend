package main

import (
	"log"
	"time_of_armies/internal/api"
)

func main() {
	log.Println("Launching the app!")
	api.StartServer()
	log.Println("App is terminated!")
}
