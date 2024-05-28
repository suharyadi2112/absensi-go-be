package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

var (
	rabbitHost     string
	rabbitPort     string
	rabbitUser     string
	rabbitPassword string
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file", err.Error())
	}

	rabbitHost = os.Getenv("RABBIT_HOST")
	rabbitPort = os.Getenv("RABBIT_PORT")
	rabbitUser = os.Getenv("RABBIT_USER")
	rabbitPassword = os.Getenv("RABBIT_PASSWORD")
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

func main() {

	conn, ch, err := InitRabbitmq()
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %s", err.Error())
	}
	defer conn.Close()
	defer ch.Close()

	log.Println("Mendeklarasikan antrian...")
	q, err := ch.QueueDeclare(
		"absensi", // Nama antrian
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Gagal mendeklarasikan antrian: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // antrian
		"",     // konsumen
		false,  // auto-ack (false untuk manual ack)
		false,  // eksklusif
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Printf("Gagal mendaftarkan konsumen: %v", err)
	}

	log.Println("Menunggu pesan...")
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Menerima pesan: %s", d.Body)
			d.Ack(false)
		}
	}()

	log.Println("Tekan CTRL+C untuk keluar")
	<-forever
}
