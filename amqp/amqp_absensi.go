package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	dbRab "absensi/config"
	cont "absensi/controllers"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	controller *cont.Conn
)

type ConnAmqpAbsen struct {
	RabMQ *amqp.Channel
}
type PayloadRabbit struct {
	IdKelas        int64  `json:"IdKelas"`
	IdSiswaOrGuru  int64  `json:"IdSiswaOrGuru"`
	JenisProses    string `json:"JenisProses"`
	TanggalHariIni string `json:"TanggalHariIni"`
	TipeAbsen      string `json:"TipeAbsen"`
	TimeOnly       string `json:"timeOnly"`
}

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	var err error
	controller, err = cont.NewCon()
	if err != nil {
		panic(err) // Handle error appropriately
	}
}

// Fungsi untuk inisialisasi handler dengan instance database
func NewCon() (*ConnAmqpAbsen, error) {
	rabG, err := dbRab.InitRabbitMQ()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize rabbit MQ: %w", err)
	}

	return &ConnAmqpAbsen{
		RabMQ: rabG,
	}, nil
}

func main() {

	myInstance, err := NewCon() //koneksi
	if err != nil {
		log.Fatalf("Gagal menginisialisasi koneksi RabbitMQ: %v", err)
	}
	defer myInstance.RabMQ.Close()

	msgChannel := make(chan amqp.Delivery, 100) //buat chanel
	var wg sync.WaitGroup                       //WaitGroup untuk menunggu semua worker selesai

	numWorkers := 5 // jumlah worker
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(msgChannel, &wg)
	}

	go myInstance.ProcessAmqpAbsen(msgChannel) // proses pesan\
	logrus.Info("menunggu pesan")
	wg.Wait() // Menunggu semua worker selesai

}

func worker(msgChannel <-chan amqp.Delivery, wg *sync.WaitGroup) {
	defer wg.Done()
	// Konfigurasi Logrus
	logrus.SetLevel(logrus.InfoLevel)

	for d := range msgChannel {
		log.Printf("Menerima pesan: %s", d.Body)

		var payload PayloadRabbit
		if err := json.Unmarshal(d.Body, &payload); err != nil {
			log.Printf("Error decoding JSON: %s", err.Error())
			continue
		}

		jenisProses := payload.JenisProses
		IdKelas := payload.IdKelas
		IdSiswaOrGuru := payload.IdSiswaOrGuru
		TanggalHariIni := payload.TanggalHariIni
		TipeAbsen := payload.TipeAbsen
		TimeOnly := payload.TimeOnly

		if jenisProses == "satu" {

			err := controller.UpdateAbsenController(TimeOnly, TanggalHariIni, IdSiswaOrGuru, IdKelas) //tanpa tipeabsen
			if err != nil {
				logrus.Errorf("Error update absen controller amqp satu: %s", err.Error())
				continue
			}

		} else if jenisProses == "dua" || jenisProses == "tiga" {

			// untuk absen masuk
			err := controller.PostAbsenSiswaController(TimeOnly, TanggalHariIni, TipeAbsen, IdSiswaOrGuru, IdKelas)
			if err != nil {
				logrus.Errorf("Error post absen controller amqp dua | tiga: %s", err.Error())
				continue
			}

		} else if jenisProses == "empat" { //zona guru

		} else if jenisProses == "lima" {

			// untuk absen masuk
			err := controller.PostAbsenGuruController(TimeOnly, TanggalHariIni, IdSiswaOrGuru)
			if err != nil {
				logrus.Errorf("Error post absen controller amqp dua | tiga: %s", err.Error())
				continue
			}

		} else {
			logrus.Info("tanpa jenis proses")
		}

		logrus.Info("succes")

		// Mengakui pesan yang diterima
		d.Ack(false)
	}
}

func (amqp *ConnAmqpAbsen) ProcessAmqpAbsen(msgChannel chan<- amqp.Delivery) {
	// Mendeklarasikan antrian
	q, err := amqp.RabMQ.QueueDeclare(
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

	// Menunggu dan menerima pesan
	msgs, err := amqp.RabMQ.Consume(
		q.Name, // antrian
		"",     // konsumen
		false,  // auto-ack (false untuk manual ack)
		false,  // eksklusif
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Gagal mendaftarkan konsumen: %v", err)
	}

	// Mengirim pesan ke channel
	for d := range msgs {
		msgChannel <- d
	}

	close(msgChannel)
}
