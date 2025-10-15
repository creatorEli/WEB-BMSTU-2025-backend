package main

import (
	"time_of_armies/internal/app/ds"
	"time_of_armies/internal/app/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database!!")
	}
	err = db.AutoMigrate(
		&ds.Army{},
		&ds.Historian{},
		&ds.TravelTime{},
		&ds.ConnectArmyTT{},
	)
	if err != nil {
		panic("cannot migrate db!")
	}
}
