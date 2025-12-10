package main

import (
	"fmt"
	"time_of_armies/internal/app/config"
	"time_of_armies/internal/app/dsn"
	handler "time_of_armies/internal/app/hadler"
	"time_of_armies/internal/app/redis"
	"time_of_armies/internal/app/repository"
	"time_of_armies/internal/pkg"

	_ "time_of_armies/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @title Travel Times API
// @version 1.0
// @description Система расчёта дневного перехода армии
// @contact.name Мефодьев Илья
// @contact.url https://t.me/Thir5tyF0r1ife
// @contact.email creatoreli8@gmail.com
// @license.name AS IS (NO WARRANTY)
// @host  localhost:8084
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите: Bearer {ваш_JWT_токен}
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

	redisClient, err := redis.New(&gin.Context{}, conf.Redis)
	if err != nil {
		logrus.Fatalf("error connecting redis!: %v", errRep)
	}

	hand := handler.NewHandler(rep, conf, redisClient)

	defer redisClient.Close() // если есть такой метод

	application := pkg.NewApp(conf, router, hand, redisClient)
	application.RunApp()
}
