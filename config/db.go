package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/pusher/pusher-http-go/v5"
	"github.com/sirupsen/logrus"
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

	pusherAppId   string
	pusherKey     string
	pusherSecret  string
	pusherCluster string
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

	pusherAppId = os.Getenv("APP_ID")
	pusherKey = os.Getenv("APP_KEY")
	pusherSecret = os.Getenv("APP_SECRET")
	pusherCluster = os.Getenv("APP_CLUSTER")
}

// LOGRUS
func InitLogRus() *logrus.Logger {

	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{})

	newFileName := time.Now().Format("2006-01-02") + ".log" //rename file
	logFileName := "log/" + newFileName
	createLogFile(logFileName)

	// Dapatkan informasi file
	fileInfo, err := os.Stat(logFileName)
	if err != nil {
		fmt.Println("Gagal mendapatkan informasi file log:", err)
	}

	// Cek ukuran file
	fileSize := fileInfo.Size()

	if fileSize > 1024*1024 { //pisah 1mb perfile
		newFileName := time.Now().Format("2006-01-02") + ".log" //rename file
		createLogFile(newFileName)                              // buat file
		logFileName = newFileName                               //file baru
	}
	file, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_APPEND, 0666) //buka file log
	if err != nil {
		fmt.Println("Gagal membuka file log:", err)
	}

	logger.SetOutput(file)
	logger.SetLevel(logrus.InfoLevel)

	return logger
}

// log config
func createLogFile(fileName string) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666) // buat file
	if err != nil {
		fmt.Println("Gagal membuat file log:", err)
	}
	defer file.Close()
}

// logconfig
func InitlogError(logger *logrus.Logger, context, addInfo string, err error, errorType string) {

	logger.SetOutput(os.Stdout) // untuk tetap di cetak di console
	entry := logger.WithFields(logrus.Fields{
		"context": context,
		"info":    addInfo,
	})

	switch errorType {
	case "info":
		entry.Info("Informational message")
	case "warning":
		entry.Warn("Warning message")
	case "error":
		if err != nil {
			entry = entry.WithError(err)
		}
		entry.Error("An error occurred")
	default:
		entry.Warn("Unknown log type")
	}
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

// RABBIT MQ
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

// PUSHER
func InitPusher() pusher.Client {

	// Inisialisasi client Pusher
	pusherClient := pusher.Client{
		AppID:   pusherAppId,
		Key:     pusherKey,
		Secret:  pusherSecret,
		Cluster: pusherCluster,
		Secure:  true,
	}

	return pusherClient

}
