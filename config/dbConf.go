package config

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
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
	// hostPort := os.Getenv("PORT")

	// pgConnStrings := fmt.Sprintf("port=%s host=%s user=%s "+"password=%s dbname=%s sslmode=disable", hostPort, hostName, userName, password, dbName)
	url := fmt.Sprintf("%v:%v@tcp(%v)/%v?parseTime=true",
		userName,
		password,
		hostName,
		// hostPort,
		dbName)
	// log.Println("TESS ", url)
	d, err := sqlx.Connect("mysql", url)
	if err != nil {
		log.Println("DATA ", err)
		panic(err)
	}
	d.SetMaxIdleConns(10)
	d.SetMaxOpenConns(10)

	dbConn.SQL = d
	return dbConn, err
}
