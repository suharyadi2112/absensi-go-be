package main

import (
	db "absensi/config"
	"log"

	"github.com/streadway/amqp"
)

func main() {

	conn, ch, err := db.InitRabbitmq()
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

	log.Println("Waiting for messages...")
	go consumeMessages(ch, q)

	log.Println("Tekan CTRL+C untuk keluar")
	forever := make(chan bool)
	<-forever
}

func consumeMessages(ch *amqp.Channel, q amqp.Queue) {

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

	for d := range msgs {
		log.Printf("Menerima pesan: %s", d.Body)
		d.Ack(false) // Kirim ack setelah pesan diproses berhasil
	}
}
