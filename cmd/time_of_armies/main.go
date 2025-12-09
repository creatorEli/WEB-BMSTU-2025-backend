package main

import (
	"fmt"
	"time_of_armies/internal/app/config"
	"time_of_armies/internal/app/dsn"
	handler "time_of_armies/internal/app/hadler"
	"time_of_armies/internal/app/repository"
	"time_of_armies/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	router := gin.Default()
	conf, err := config.NewConfig()

	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	fmt.Println("conn to db: %s", postgresString)

	rep, errRep := repository.New(postgresString)

	if errRep != nil {
		logrus.Fatalf("error initializing repository!: %v", errRep)
	}

	hand := handler.NewHandler(rep)

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}
