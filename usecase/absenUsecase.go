package usecase

import (
	conFig "absensi/config"
	cont "absensi/controllers"
	helper "absensi/helper"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/pusher/pusher-http-go/v5"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	controller   *cont.Conn
	pusherChanel string
	pusherEvent  string
	logger       *logrus.Logger
	randomString string
	sch          string
)

type AbsenUsecase struct {
	RabMQ  *amqp.Channel
	Pusher pusher.Client
}

func init() {
	ctx := "Usecase-InitUsecase"
	var err error

	logger = conFig.InitLogRus()
	randomString = helper.GenerateRandomString(6)

	controller, err = cont.NewCon()
	if err != nil {
		conFig.InitLog(logger, ctx, "error terhubung dengan controller", err, "error") // catat log
	}
	err = godotenv.Load("../.env")
	if err != nil {
		conFig.InitLog(logger, ctx, "Error loading .env file on absensi usecase", err, "error") // catat log
	}

	pusherChanel = os.Getenv("APP_CHANNEL")
	pusherEvent = os.Getenv("APP_EVENT")
	sch = "PLT"

}

// Fungsi untuk inisialisasi handler dengan instance database
func NewConUsecase() (*AbsenUsecase, error) {
	ctx := "Usecase-NewConUsecase"
	rabMQD, err := conFig.InitRabbitMQ()
	if err != nil {
		conFig.InitLog(logger, ctx, "failed to initialize rabbit", err, "error") // catat log
		return nil, fmt.Errorf("failed to initialize rabbit: %w", err)
	}
	pusHER := conFig.InitPusher()

	return &AbsenUsecase{
		RabMQ:  rabMQD,
		Pusher: pusHER,
	}, nil
}

// use case absen top
func GetAbsenTopUsecase(tanggalhariIni string) ([]map[string]interface{}, error) {
	ctx := "Usecase-GetAbsenTopUsecase"

	fmt.Println("Tanggal sekarang - get absen top:", tanggalhariIni)
	result, err := controller.GetAbsenTopController(tanggalhariIni)
	fmt.Println(result)

	if err != nil {
		conFig.InitLog(logger, ctx, "error get data absen get top", err, "error") // catat log
		return nil, err
	}

	//custom return yang diperlukan response
	var absentopResp []map[string]interface{}
	for _, s := range result {
		absentopResp = append(absentopResp, map[string]interface{}{
			"IDAbsensi":  s.ID.Int64,
			"FotoSiswa":  s.IDSiswa.Foto.String,
			"FotoGuru":   s.IDPengajar.Foto.String,
			"Nama":       s.IDSiswa.NamaLengkap.String,
			"NamaGuru":   s.IDPengajar.NamaLengkap.String,
			"Kelas":      s.IDKelas.Kelas.String,
			"IDPengajar": s.IDPengajar.ID.Int64,
		})
	}

	return absentopResp, nil
}

// use case simpan tap scan absen
func (r *AbsenUsecase) PostAbsenTopUsecase(formCode, tanggalhariIni, dateTimehariini, timeonlyHariini string) (res map[string]interface{}, codeHttp int, err error) {

	ctx := "Usecase-PostAbsenTopUsecase"

	fmt.Println("Form code - post absen:", formCode)

	countSiswa, err := controller.CountSiswaController(formCode)
	if err != nil {
		conFig.InitLog(logger, ctx, "error CountSiswaController", err, "error") // catat log
		return nil, 500, err
	}
	countGuru, err := controller.CountGuruController(formCode)
	if err != nil {
		conFig.InitLog(logger, ctx, "error CountGuruController", err, "error") // catat log
		return nil, 500, err
	}

	if countSiswa > 0 { //section siswa

		resSiswa, err := controller.GetSiswaController(formCode)
		if err != nil {
			conFig.InitLog(logger, ctx, "error GetSiswaController", err, "error") // catat log
			return nil, 500, err
		}

		id_siswa := resSiswa.ID.Int64
		id_kelas := resSiswa.IDKelas.ID.Int64
		no_hp_ortu := resSiswa.IDOrtu.NoHP.String
		cAbsen, err := controller.GetOneAbsensiSiswaController(id_siswa, id_kelas, tanggalhariIni)

		if cAbsen != nil {
			var jamFixCon int
			jam_masuk := cAbsen.JMasuk
			keluar := cAbsen.Absensi.Keluar

			if jam_masuk.Valid { //tidak boleh null
				jamFix, err := calculateHoursDifference(dateTimehariini, jam_masuk.String) //cek jeda jam
				if err != nil {
					conFig.InitLog(logger, ctx, "error calculateHoursDifference", err, "error") // catat log
					return nil, 500, err
				}
				jamFixCon = jamFix //assign jam fix
				fmt.Println(jamFix, "perbedaan jam")
			}
			if jamFixCon > 0 { //jeda waktu absen balek
				if !keluar.Valid {

					err = controller.PostUpdateAbsensiSiswaController(id_siswa, id_kelas, tanggalhariIni, timeonlyHariini, randomString) //untuk update absen keluar
					if err != nil {
						conFig.InitLog(logger, ctx, "error PostUpdateAbsensiSiswaController Keluar", err, "error")
						return nil, 500, err
					}
					// Generate message for pulang sekolah
					keluarMessage := helper.GenerateKeluarMessage(resSiswa.NamaLengkap.String, resSiswa.IDKelas.Kelas.String, dateTimehariini)
					err = controller.PostInsertAbsensiWaController(randomString, sch, no_hp_ortu, keluarMessage, "0", dateTimehariini, "out")

					if err != nil {
						conFig.InitLog(logger, ctx, "error PostInsertAbsensiWaController keluar", err, "error")
					}

					fmt.Println(keluarMessage)

					r.PushPusher(formCode) //sendPusher

					// Create response structure
					responseItem := map[string]interface{}{
						"FormCode": resSiswa.NIS.String,
						"Nama":     resSiswa.NamaLengkap.String,
						"Kelas":    resSiswa.IDKelas.Kelas.String,
						"Alamat":   resSiswa.Alamat.String,
						"Foto":     url.PathEscape(resSiswa.Foto.String),
						"AbsenAt":  dateTimehariini,
						"Tipe":     "siswa",
					}
					return responseItem, 200, nil
				} else {
					responseItem := map[string]interface{}{
						"Message": "Anda sudah melakukan absensi",
					}
					conFig.InitLog(logger, ctx, "sudah melakukan absen #2ess3", nil, "info") // catat log
					return responseItem, 400, nil
				}
			} else { //kemungkinan terjadi saat data absen ada tapi jam masuk maupun pulang kosong/null
				responseItem := map[string]interface{}{
					"Message": "Anda sudah melakukan absensi",
				}
				conFig.InitLog(logger, ctx, fmt.Sprintf("Sudah melakukan absen different jam %s #k3k3", jamFixCon), nil, "info") // catat log
				return responseItem, 400, nil
			}

		} else { // absen maasuk

			parsedTime, err := time.Parse("15:04:05", timeonlyHariini)
			if err != nil {
				conFig.InitLog(logger, ctx, "Error parsing time", err, "error") // catat log
				return nil, 500, err
			}
			morningStart, _ := time.Parse("15:04:05", "05:00:00")
			noonEnd, _ := time.Parse("15:04:05", "10:30:00")
			afternunStart, _ := time.Parse("15:04:05", "10:30:00")
			niteEnd, _ := time.Parse("15:04:05", "21:00:00")

			isMorning := parsedTime.After(morningStart) && parsedTime.Before(noonEnd)
			isNite := parsedTime.After(afternunStart) && parsedTime.Before(niteEnd)

			// isMorning = true //untuk testing
			// isNite = false   //untuk testing

			fmt.Println(isMorning, isNite, "cek ketentuan jam")

			var tipeAbsen string
			if isMorning { //pagi

				tipeAbsen = "masuk"
				// err = r.DeclarePublishAbsen(timeonlyHariini, tanggalhariIni, tipeAbsen, "siswa", id_siswa, id_kelas, "dua") //RABbITMQ
				err = controller.PostInsertAbsensiSiswaController(id_siswa, id_kelas, "H", dateTimehariini, timeonlyHariini, 0, randomString, tipeAbsen) //insert table absensi
				if err != nil {
					conFig.InitLog(logger, ctx, "error PostInsertAbsensiSiswaController Masuk", err, "error")
					return nil, 500, err
				}
				//insert table wa
				masukMessage := helper.GenerateMasukMessage(resSiswa.NamaLengkap.String, resSiswa.IDKelas.Kelas.String, dateTimehariini)
				err = controller.PostInsertAbsensiWaController(randomString, sch, no_hp_ortu, masukMessage, "0", dateTimehariini, "in")

				if err != nil {
					conFig.InitLog(logger, ctx, "error PostInsertAbsensiWaController masuk", err, "error")
				}

				fmt.Println(masukMessage)

			} else if isNite { //diatas jam 12 siang kemungkinan absensi pulang

				tipeAbsen = "keluar"
				// err = r.DeclarePublishAbsen(timeonlyHariini, tanggalhariIni, tipeAbsen, "siswa", id_siswa, id_kelas, "tiga") //RABBITMQ
				err = controller.PostInsertAbsensiSiswaController(id_siswa, id_kelas, "H", dateTimehariini, timeonlyHariini, 0, randomString, tipeAbsen)
				if err != nil {
					conFig.InitLog(logger, ctx, "error PostInsertAbsensiSiswaController Keluar", err, "error")
					return nil, 500, err
				}

				// Generate message for pulang sekolah
				keluarMessage := helper.GenerateKeluarMessage(resSiswa.NamaLengkap.String, resSiswa.IDKelas.Kelas.String, dateTimehariini)
				err = controller.PostInsertAbsensiWaController(randomString, sch, no_hp_ortu, keluarMessage, "0", dateTimehariini, "out")

				if err != nil {
					conFig.InitLog(logger, ctx, "error PostInsertAbsensiWaController keluar", err, "error")
				}

				fmt.Println(keluarMessage)

			} else {
				responseItem := map[string]interface{}{
					"Message": "Anda sudah melakukan absensi",
				}
				conFig.InitLog(logger, ctx, "bukan waktu sekolah", nil, "info") // catat log
				return responseItem, 400, nil
			}

			r.PushPusher(formCode) //sendPusher

			// Create response structure
			responseItem := map[string]interface{}{
				"FormCode":  resSiswa.NIS.String,
				"Nama":      resSiswa.NamaLengkap.String,
				"Kelas":     resSiswa.IDKelas.Kelas.String,
				"Alamat":    resSiswa.Alamat.String,
				"Foto":      url.PathEscape(resSiswa.Foto.String),
				"AbsenAt":   dateTimehariini,
				"Tipe":      "siswa",
				"TipeAbsen": tipeAbsen,
			}
			return responseItem, 200, nil

		}

		/*
			SECTIONN GURUUUUUUUUUUUUUUUU
		*/
	} else if countGuru > 0 {

		resGuru, err := controller.GetGuruController(formCode)
		if err != nil {
			conFig.InitLog(logger, ctx, "error GetGuruController", err, "error") // catat log
			return nil, 500, err
		}

		id_pengajar := resGuru.ID.Int64
		nip := resGuru.NIP.String
		nama_guru := resGuru.NamaLengkap.String
		alamat_guru := resGuru.Alamat.String
		foto_guru := resGuru.Foto.String

		cAbsenGuru, err := controller.GetOneAbsensiGuruController(id_pengajar, tanggalhariIni)
		if err != nil {
			conFig.InitLog(logger, ctx, "error GetOneAbsensiGuruController", err, "error") // catat log
			return nil, 500, err
		}

		if cAbsenGuru != nil { //jika absen guru tersebut ada ditanggal tsb
			var jamFixConGur int
			jam_masuk := cAbsenGuru.JMasuk
			keluarGur := cAbsenGuru.Absensi.Keluar
			if jam_masuk.Valid { //tidak boleh null
				jamFixGur, err := calculateHoursDifference(dateTimehariini, jam_masuk.String) //cek veda jam
				if err != nil {
					conFig.InitLog(logger, ctx, "error calculateHoursDifference", err, "error") // catat log
					return nil, 500, err
				}
				jamFixConGur = jamFixGur //assign jam fix
			}

			logrus.Info(jamFixConGur, " jeda absen guru dari masuk")
			if jamFixConGur > 0 { //jeda absen balek

				if !keluarGur.Valid {

					err = controller.PostUpdateAbsensiGuruController(id_pengajar, timeonlyHariini, tanggalhariIni) //untuk update absen keluar
					if err != nil {
						conFig.InitLog(logger, ctx, "error PostUpdateAbsensiGuruController Keluar", err, "error")
						return nil, 500, err
					}

					r.PushPusher(formCode) //sendPusher

					// Create response structure
					responseItem := map[string]interface{}{
						"FormCode":  nip,
						"Nama":      nama_guru,
						"Kelas":     "-",
						"Alamat":    alamat_guru,
						"Foto":      url.PathEscape(foto_guru),
						"AbsenAt":   dateTimehariini,
						"Tipe":      "guru",
						"TipeAbsen": "-",
					}
					return responseItem, 200, nil

				} else {

					logrus.Info("sudah absen guru 223")
					responseItem := map[string]interface{}{
						"Message": "Anda sudah melakukan absensi",
					}
					return responseItem, 400, nil
				}

			} else {
				logrus.Info("sudah absen guru 123")
				responseItem := map[string]interface{}{
					"Message": "Anda sudah melakukan absensi",
				}
				return responseItem, 400, nil
			}

		} else { //Guru baru absen

			logrus.Info(cAbsenGuru)

			parsedTime, err := time.Parse("15:04:05", timeonlyHariini)
			if err != nil {
				conFig.InitLog(logger, ctx, "Error parsing time", err, "error") // catat log
				return nil, 500, err
			}
			//set waktu untuk guru
			morningStart, _ := time.Parse("15:04:05", "05:00:00")
			noonEnd, _ := time.Parse("15:04:05", "10:30:00")
			afternunStart, _ := time.Parse("15:04:05", "10:30:00")
			niteEnd, _ := time.Parse("15:04:05", "21:00:00")

			isMorning := parsedTime.After(morningStart) && parsedTime.Before(noonEnd)
			isNite := parsedTime.After(afternunStart) && parsedTime.Before(niteEnd)

			isMorning = true //untuk testing
			// isNite = false   //untuk testing

			fmt.Println(isMorning, isNite, "cek ketentuan jam")

			if isMorning { //pagi
				err = controller.PostInsertAbsensiGuruController(id_pengajar, "H", dateTimehariini, timeonlyHariini, "masuk") //insert table absensi
				if err != nil {
					conFig.InitLog(logger, ctx, "error PostInsertAbsensiGuruController Masuk", err, "error")
					return nil, 500, err
				}

			} else {
				responseItem := map[string]interface{}{
					"Message": "Anda sudah melakukan absensi",
				}
				conFig.InitLog(logger, ctx, "bukan waktu sekolah setting untuk guru", nil, "info") // catat log
				return responseItem, 400, nil
			}

			r.PushPusher(formCode) //sendPusher

			// Create response structure
			responseItem := map[string]interface{}{
				"FormCode":  nip,
				"Nama":      nama_guru,
				"Kelas":     "-",
				"Alamat":    alamat_guru,
				"Foto":      url.PathEscape(foto_guru),
				"AbsenAt":   dateTimehariini,
				"Tipe":      "guru",
				"TipeAbsen": "-",
			}
			return responseItem, 200, nil
		}

	} else {

		responseItem := map[string]interface{}{
			"Message": "Kartu anda tidak terdaftar",
		}
		conFig.InitLog(logger, ctx, "Kartu anda tidak terdaftar", err, "warning") // catat log
		return responseItem, 400, nil

	}

	// return nil, 200, nil
}

func calculateHoursDifference(datetime, j_masuk string) (jMasukDiff int, err error) {

	ctx := "Usecase-calculateHoursDifference"

	fmt.Println(j_masuk, "cek j_masuk")
	layout := "2006-01-02 15:04:05" //layout format time yang dikonvert

	tSatu, err := time.Parse(layout, datetime)
	if err != nil {
		conFig.InitLog(logger, ctx, "Error parsing datetime", err, "error") // catat log
		return 0, err
	}
	tDua, err := time.Parse(layout, j_masuk)
	if err != nil {
		conFig.InitLog(logger, ctx, "Error parsing j_masuk", err, "error") // catat log
		return 0, err
	}

	diff := tSatu.Sub(tDua).Seconds()
	jam := int(diff / 3600)

	return jam, nil
}

func (r *AbsenUsecase) PushPusher(formCode string) {

	ctx := "Usecase-PushPusher"

	r.Pusher.Trigger(pusherChanel, pusherEvent, map[string]string{"formCode": formCode})
	conFig.InitLog(logger, ctx, "push to pusher", nil, "info") // catat log

}

func (r *AbsenUsecase) DeclarePublishAbsen(timeKeluarORmasuk, tanggalHariIni, tipeAbsen, tipeOrang string, idSiswaOrGuru, idKelas int64, jenisQue string) (err error) {

	ctx := "Usecase-DeclarePublishAbsen"

	q, err := r.RabMQ.QueueDeclare(
		"absensi", // queue name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		conFig.InitLog(logger, ctx, "Failed to declare a queue", err, "error") // catat log
		return err
	}

	//struktur data
	data := map[string]interface{}{
		"timeOnly":       timeKeluarORmasuk,
		"TanggalHariIni": tanggalHariIni,
		"IdSiswaOrGuru":  idSiswaOrGuru,
		"IdKelas":        idKelas,
		"TipeAbsen":      tipeAbsen,
		"TipeOrang":      tipeOrang,
		"JenisProses":    jenisQue,
	}

	// Meng-marshal data menjadi JSON
	body, err := json.Marshal(data)
	if err != nil {
		conFig.InitLog(logger, ctx, "Failed to marshal JSON", err, "error") // catat log
		return err
	}

	err = r.RabMQ.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		conFig.InitLog(logger, ctx, "Failed to publish a message", err, "error") // catat log
		return err
	}
	return nil
}
