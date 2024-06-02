package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

// Package-level variables
var (
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string

	rabbitHost     string
	rabbitPort     string
	rabbitUser     string
	rabbitPassword string
)

func init() {

	// Load environment variables from .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file", err.Error())
	}

	// Initialize package-level variables
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_NAME")

	rabbitHost = os.Getenv("RABBIT_HOST")
	rabbitPort = os.Getenv("RABBIT_PORT")
	rabbitUser = os.Getenv("RABBIT_USER")
	rabbitPassword = os.Getenv("RABBIT_PASSWORD")
}

// MYSQL
func InitDBMySql() (*sql.DB, error) {

	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println(err.Error(), "koneksi mysql 2")
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error(), "koneksi mysql 2")
	}

	return db, nil
}

func InitRabbitMQ() (*amqp.Channel, error) {

	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitUser, rabbitPassword, rabbitHost, rabbitPort)
	conn, err := amqp.Dial(connStr)
	if err != nil {
		fmt.Println(err.Error(), "koneksi rabbitMQ 1")
	}

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err.Error(), "koneksi rabbitMQ 2")
	}

	return ch, nil
}
