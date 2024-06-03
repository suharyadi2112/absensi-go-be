package usecase

import (
	rabbit "absensi/config"
	cont "absensi/controllers"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var controller *cont.Conn

type AbsenUsecase struct {
	RabMQ *amqp.Channel
}

func init() {
	var err error
	controller, err = cont.NewCon()
	if err != nil {
		panic(err) // Handle error appropriately
	}
}

// Fungsi untuk inisialisasi handler dengan instance database
func NewConUsecase() (*AbsenUsecase, error) {
	rabMQD, err := rabbit.InitRabbitMQ()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize rabbit: %w", err)
	}
	return &AbsenUsecase{
		RabMQ: rabMQD,
	}, nil
}

// use case absen top
func GetAbsenTopUsecase(tanggalhariIni string) ([]map[string]interface{}, error) {

	fmt.Println("Tanggal sekarang - get absen top:", tanggalhariIni)

	result, err := controller.GetAbsenTopController(tanggalhariIni)

	fmt.Println(result)

	if err != nil {
		return nil, err
	}

	//custom return yang diperlukan response
	var absentopResp []map[string]interface{}
	for _, s := range result {
		absentopResp = append(absentopResp, map[string]interface{}{
			"IDAbsensi":  s.ID.Int64,
			"FotoSiswa":  s.IDSiswa.Foto.String,
			"FotoGuru":   s.IDPengajar.Foto.String,
			"NamaSiswa":  s.IDSiswa.NamaLengkap.String,
			"NamaGuru":   s.IDPengajar.NamaLengkap.String,
			"Kelas":      s.IDKelas.Kelas.String,
			"IDPengajar": s.IDPengajar.ID.Int64,
		})
	}

	return absentopResp, nil
}

// use case simpan tap scan absen
func (r *AbsenUsecase) PostAbsenTopUsecase(formCode, tanggalhariIni, dateTimehariini, timeonlyHariini string) (res map[string]interface{}, codeHttp int, err error) {

	fmt.Println("Form code - post absen:", formCode)

	countSiswa, err := controller.CountSiswaController(formCode)
	if err != nil {
		return nil, 500, err
	}
	countGuru, err := controller.CountGuruController(formCode)
	if err != nil {
		return nil, 500, err
	}

	if countSiswa > 0 {

		resSiswa, err := controller.GetSiswaController(formCode)
		if err != nil {
			return nil, 500, err
		}

		id_siswa := resSiswa.ID.Int64
		id_kelas := resSiswa.IDKelas.ID.Int64
		cAbsen, err := controller.GetOneAbsensiSiswaController(id_siswa, id_kelas, tanggalhariIni)

		if err != nil {
			return nil, 500, err
		}

		if cAbsen != nil {
			var jamFixCon int
			jam_masuk := cAbsen.JMasuk
			keluar := cAbsen.Absensi.Keluar

			if jam_masuk.Valid { //tidak boleh null
				jamFix, err := calculateHoursDifference(dateTimehariini, jam_masuk.String) //cek veda jam
				if err != nil {
					return nil, 500, err
				}
				jamFixCon = jamFix //assign jam fix
				fmt.Println(jamFix, "masokkk")
			}
			if jamFixCon > 0 {
				if !keluar.Valid {

					err = r.DeclarePublishAbsen(timeonlyHariini, tanggalhariIni, "absenready", "siswa", id_siswa, id_kelas, "satu")
					if err != nil {
						return nil, 500, err
					}

					// Create response structure
					responseItem := map[string]interface{}{
						"FormCode":  resSiswa.NIS.String,
						"NamaSiswa": resSiswa.NamaLengkap.String,
						"Kelas":     resSiswa.IDKelas.Kelas.String,
						"Alamat":    resSiswa.Alamat.String,
						"Foto":      url.PathEscape(resSiswa.Foto.String),
						"AbsenAt":   dateTimehariini,
						"Tipe":      "siswa",
					}
					return responseItem, 200, nil
				} else {
					responseItem := map[string]interface{}{
						"Message": "Anda sudah melakukan absensi",
					}
					return responseItem, 400, nil
				}
			} else { //kemungkinan terjadi saat data absen ada tapi jam masuk maupun pulang kosong/null
				responseItem := map[string]interface{}{
					"Message": "Terjadi kesalahan hubungi admin #sks88",
				}
				return responseItem, 400, nil
			}

		} else { // absen maasuk

			parsedTime, err := time.Parse("15:04:05", timeonlyHariini)
			if err != nil {
				fmt.Println("Error parsing time:", err)
				return nil, 500, err
			}
			morningStart, _ := time.Parse("15:04:05", "05:00:00")
			noonEnd, _ := time.Parse("15:04:05", "12:00:00")

			afternunStart, _ := time.Parse("15:04:05", "12:00:00")
			niteEnd, _ := time.Parse("15:04:05", "21:00:00")

			//JAM TESTING
			// morningStart, _ := time.Parse("15:04:05", "21:00:00")
			// noonEnd, _ := time.Parse("15:04:05", "23:00:00")
			// afternunStart, _ := time.Parse("15:04:05", "12:00:00")
			// niteEnd, _ := time.Parse("15:04:05", "21:00:00")

			isMorning := parsedTime.After(morningStart) && parsedTime.Before(noonEnd)
			isNite := parsedTime.After(afternunStart) && parsedTime.Before(niteEnd)

			fmt.Println(isMorning, isNite)

			var tipeAbsen string
			if isMorning {

				tipeAbsen = "masuk"
				err = r.DeclarePublishAbsen(timeonlyHariini, tanggalhariIni, tipeAbsen, "siswa", id_siswa, id_kelas, "dua")
				if err != nil {
					return nil, 500, err
				}

			} else if isNite { //diatas jam 12 siang kemungkinan absensi pulang

				tipeAbsen = "keluar"
				err = r.DeclarePublishAbsen(timeonlyHariini, tanggalhariIni, tipeAbsen, "siswa", id_siswa, id_kelas, "tiga")
				if err != nil {
					return nil, 500, err
				}

			} else {
				responseItem := map[string]interface{}{
					"Message": "Terjadi kesalahan hubungi admin #kn3k2",
				}
				return responseItem, 400, nil
			}
			// Create response structure
			responseItem := map[string]interface{}{
				"FormCode":  resSiswa.NIS.String,
				"NamaSiswa": resSiswa.NamaLengkap.String,
				"Kelas":     resSiswa.IDKelas.Kelas.String,
				"Alamat":    resSiswa.Alamat.String,
				"Foto":      url.PathEscape(resSiswa.Foto.String),
				"AbsenAt":   dateTimehariini,
				"Tipe":      "siswa",
				"TipeAbsen": tipeAbsen,
			}
			return responseItem, 200, nil

		}

	} else if countGuru > 0 { //section guru

		resGuru, err := controller.GetGuruController(formCode)
		if err != nil {
			return nil, 500, err
		}

		id_pengajar := resGuru.ID.Int64
		nip := resGuru.NIP.String
		nama_guru := resGuru.NamaLengkap.String
		alamat_guru := resGuru.Alamat.String
		foto_guru := resGuru.Foto.String

		logrus.Info(id_pengajar, nip, nama_guru, alamat_guru, foto_guru) //pakai & untuk cepat lok pointer

		cAbsenGuru, err := controller.GetOneAbsensiGuruController(id_pengajar, tanggalhariIni)
		if err != nil {
			return nil, 500, err
		}

		if cAbsenGuru != nil {
			logrus.Info(cAbsenGuru, "isi")
			var jamFixConGur int
			jam_masuk := cAbsenGuru.JMasuk
			keluarGur := cAbsenGuru.Absensi.Keluar
			if jam_masuk.Valid { //tidak boleh null
				jamFixGur, err := calculateHoursDifference(dateTimehariini, jam_masuk.String) //cek veda jam
				if err != nil {
					return nil, 500, err
				}
				jamFixConGur = jamFixGur //assign jam fix
				fmt.Println(jamFixConGur, "masok guru")
			}

			if jamFixConGur > 0 {

				if !keluarGur.Valid {
					err = r.DeclarePublishAbsen(timeonlyHariini, tanggalhariIni, "", "guru", id_pengajar, 0, "empat")
					if err != nil {
						return nil, 500, err
					}

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
					"Message": "Terjadi kesalahan hubungi admin #sks88",
				}
				return responseItem, 400, nil
			}

		} else {

			logrus.Info(cAbsenGuru)
			err = r.DeclarePublishAbsen(timeonlyHariini, tanggalhariIni, "", "guru", id_pengajar, 0, "lima")
			if err != nil {
				return nil, 500, err
			}
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
		logrus.Info(formCode, "inputan no kartu")
		return responseItem, 400, nil

	}

	return nil, 200, nil
}

func calculateHoursDifference(datetime, j_masuk string) (jMasukDiff int, err error) {

	fmt.Println(j_masuk, "cek j_masuk")
	layout := "2006-01-02 15:04:05" //layout format time yang dikonvert

	tSatu, err := time.Parse(layout, datetime)
	if err != nil {
		fmt.Println("error parsing datetime:", err.Error())
		return 0, err
	}
	tDua, err := time.Parse(layout, j_masuk)
	if err != nil {

		fmt.Println("error parsing j_masuk:", err.Error())
		return 0, err
	}

	diff := tSatu.Sub(tDua).Seconds()
	jam := int(diff / 3600)

	return jam, nil
}

func (r *AbsenUsecase) DeclarePublishAbsen(timeKeluarORmasuk, tanggalHariIni, tipeAbsen, tipeOrang string, idSiswaOrGuru, idKelas int64, jenisQue string) (err error) {

	q, err := r.RabMQ.QueueDeclare(
		"absensi", // queue name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
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
		log.Fatalf("Failed to marshal JSON: %v", err)
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
		log.Fatalf("Failed to publish a message: %v", err)
		return err
	}
	return nil
}
