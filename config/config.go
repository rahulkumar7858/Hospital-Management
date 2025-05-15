package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

const (
	DBHost     = "localhost"
	DBPort     = 5432
	DBUser     = "postgres"
	DBPassword = "rai78580"
	DBName     = "hospital"
	JWTSecret  = "secretkey"
)

func InitDB() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DBHost, DBPort, DBUser, DBPassword, DBName)
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
}
