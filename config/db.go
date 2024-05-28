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

var DB *sql.DB

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
	rabbitQueue    string
)

func init() {
	// wd, err := os.Getwd()
	// if err != nil {
	// 	log.Fatalf("Error getting current working directory: %v", err)
	// }
	// log.Printf("Current working directory: %s", wd)

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
	rabbitQueue = os.Getenv("RABBIT_QUEUE")

}

// MYSQL
func InitDBMySql() (*sql.DB, error) {

	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	DB = db

	log.Println("Connected to MySQL database")

	return db, nil
}

// RABBIT MQ
func InitRabbitmq() (*amqp.Connection, *amqp.Channel, error) {

	fmt.Println(rabbitHost, rabbitPort, rabbitUser, rabbitPassword)

	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitUser, rabbitPassword, rabbitHost, rabbitPort)
	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Fatal(err, "koneksi")
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err, "koneksi 2")
	}

	return conn, ch, err

}
