package main

import (
	db "absensi/config"
	"log"
)

func main() {

	conn, ch, err := db.InitRabbitmq()
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %s", err.Error())
	}
	defer conn.Close()
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"absensi", // Nama antrian
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err.Error())
	}

	msgs, err := ch.Consume(
		q.Name, // Nama antrian
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err.Error())
	}

	log.Println("Waiting for messages...")
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Menerima pesan: %s", d.Body)
		}
	}()

	log.Println("Tekan CTRL+C untuk keluar")
	<-forever

}
