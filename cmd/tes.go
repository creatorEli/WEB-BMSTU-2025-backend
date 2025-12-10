package main

import (
	"log"
	handler "time_of_armies/internal/app/hadler"

	"github.com/sirupsen/logrus"
)

func main() {
	err := handler.InitMinIO()
	if err != nil {
		logrus.Info("errsdf = " + err.Error())
	}
	logrus.Info("successd")

	if err := handler.InitMinIO(); err != nil {
		log.Fatalf("MinIO initialization failed: %v", err)
	}

}
