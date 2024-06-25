package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	dbRab "absensi/config"
	cont "absensi/controllers"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	controller *cont.Conn
	logger     *logrus.Logger
	waHost     string
	waPort     string
	waUsername string
	waPassword string
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
	waHost = os.Getenv("WA_HOST")
	waPort = os.Getenv("WA_PORT")
	waUsername = os.Getenv("WA_USERNAME")
	waPassword = os.Getenv("WA_PASSWORD")

	logger = dbRab.InitLogRus()

	var err error
	controller, err = cont.NewCon()
	if err != nil {
		panic(err) // Handle error appropriately
	}
}

// Fungsi untuk inisialisasi handler dengan instance database
func NewCon() (*ConnAmqpAbsen, error) {
	ctx := "Amqp-NewConAbsensi"
	rabG, err := dbRab.InitRabbitMQ()
	if err != nil {
		dbRab.InitLog(logger, ctx, "failed to initialize rabbit MQ", err, "error") // catat log
		return nil, err
	}

	return &ConnAmqpAbsen{
		RabMQ: rabG,
	}, nil
}

func main() {

	ctx := "Amqp-MainAbsensi"
	myInstance, err := NewCon() //koneksi
	if err != nil {
		dbRab.InitLog(logger, ctx, "Gagal menginisialisasi koneksi RabbitMQ", err, "error") // catat log
	}
	defer myInstance.RabMQ.Close()

	msgChannel := make(chan amqp.Delivery, 100) //buat chanel
	var wg sync.WaitGroup                       //WaitGroup untuk menunggu semua worker selesai

	numWorkers := 5 // jumlah worker
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(msgChannel, &wg)
	}

	go myInstance.ProcessAmqpAbsen(msgChannel)                      // proses pesan
	dbRab.InitLog(logger, ctx, "Menunggu pesan queue", nil, "info") // catat log
	wg.Wait()                                                       // Menunggu semua worker selesai

}

func worker(msgChannel <-chan amqp.Delivery, wg *sync.WaitGroup) {
	ctx := "Amqp-WorkerAbsensi"
	defer wg.Done()

	for d := range msgChannel {
		log.Printf("Menerima pesan: %s", d.Body)
		dbRab.InitLog(logger, ctx, fmt.Sprintf("Menerima pesan disini `%s`", d.Body), nil, "info") // catat log

		var payload PayloadRabbit
		if err := json.Unmarshal(d.Body, &payload); err != nil {
			dbRab.InitLog(logger, ctx, "Error decoding JSON", err, "error") // catat log
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
				dbRab.InitLog(logger, ctx, "Error update absen controller amqp satu", err, "error") // catat log
				continue
			}

		} else if jenisProses == "dua" || jenisProses == "tiga" {

			// untuk absen masuk
			err := controller.PostAbsenSiswaController(TimeOnly, TanggalHariIni, TipeAbsen, IdSiswaOrGuru, IdKelas)
			if err != nil {
				dbRab.InitLog(logger, ctx, "Error post absen controller amqp dua | tiga:", err, "error") // catat log
				continue
			}

		} else if jenisProses == "empat" { //zona guru

			// untuk absen masuk
			err := controller.UpdateAbsenGuruController(TimeOnly, TanggalHariIni, IdSiswaOrGuru)
			if err != nil {
				dbRab.InitLog(logger, ctx, "Error post absen controller amqp empat", err, "error") // catat log
				continue
			}

		} else if jenisProses == "lima" {

			// untuk absen masuk
			err := controller.PostAbsenGuruController(TimeOnly, TanggalHariIni, IdSiswaOrGuru)
			if err != nil {
				dbRab.InitLog(logger, ctx, "Error post absen controller amqp lime", err, "error") // catat log
				continue
			}

		} else {
			dbRab.InitLog(logger, ctx, "tanpa jenis proses", nil, "warning") // catat log
		}

		logrus.Info("succes")

		// Mengakui pesan yang diterima
		d.Ack(false)
	}
}

func (amqp *ConnAmqpAbsen) ProcessAmqpAbsen(msgChannel chan<- amqp.Delivery) {
	ctx := "Amqp-ProcessAmqpAbsen"
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
		dbRab.InitLog(logger, ctx, "Gagal mendeklarasikan antrian", err, "error") // catat log
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
		dbRab.InitLog(logger, ctx, "Gagal mendaftarkan konsumen", err, "error") // catat log
	}

	// Mengirim pesan ke channel
	for d := range msgs {
		msgChannel <- d
	}

	close(msgChannel)
}

// Handler untuk melakukan permintaan HTTP POST ke URL eksternal
func GetTokenWA(c echo.Context) {

	ctx := "Amqp-GetTokenWA"
	url := fmt.Sprintf("http://%s:%s/one/login", waHost, waPort)
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{
		"username" : "%s",
		"password" : "%s"
	}`, waUsername, waPassword))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		dbRab.InitLog(logger, ctx, "Request creation wa token failed", err, "error")
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		dbRab.InitLog(logger, ctx, "Client do request failed", err, "error")
	}
	defer res.Body.Close()

	_, err = io.ReadAll(res.Body)
	if err != nil {
		dbRab.InitLog(logger, ctx, "Read body failed", err, "error")
	} else {
		dbRab.InitLog(logger, ctx, "Success send wa to orang tua", nil, "info")
	}

}
