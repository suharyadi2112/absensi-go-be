package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// Package-level variables
var (
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string
)

func init() {

	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err.Error())
	}

	// Initialize package-level variables
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_NAME")
}

// MYSQL
func InitDBMySql() (*sql.DB, error) {

	fmt.Println(dbUser, dbPassword, dbHost, dbPort, dbName)

	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println(err.Error(), "koneksi mysql 2")
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error(), "koneksi mysql 2")
	}

	log.Println("Connected to MySQL database")

	return db, nil
}
