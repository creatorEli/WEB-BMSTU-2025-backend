package dsn

import (
	"fmt"
	"os"
)

// функция для создания строки подключения к БД.
func FromEnv() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		return ""
	}
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	db_name := os.Getenv("DB_NAME")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, db_name)
}
