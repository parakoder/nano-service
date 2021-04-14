package config

import (
	"fmt"
	"log"
	"os"
	"time"

	// _ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DB struct {
	SQL *sqlx.DB
}

var dbConn = &DB{}

func ConnectSQL() (*DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbName := os.Getenv("DB_NAME")
	password := os.Getenv("PASS")
	hostName := os.Getenv("HOST")
	userName := os.Getenv("USER_DB")
	hostPort := os.Getenv("DB_PORT")

	// pgConnStrings := fmt.Sprintf("port=%s host=%s user=%s "+"password=%s dbname=%s sslmode=disable", hostPort, hostName, userName, password, dbName)
	url := fmt.Sprintf("postgres://%v:%v@%v:%v/%v",
		userName,
		password,
		hostName,
		hostPort,
		dbName)

	d, err := sqlx.Open("postgres", url)
	// defer d.Close()
	if err != nil {
		panic(err)
	}
	d.SetMaxOpenConns(25)
	d.SetMaxIdleConns(25)
	d.SetConnMaxLifetime(5*time.Minute)

	dbConn.SQL = d
	
	return dbConn, err
}
