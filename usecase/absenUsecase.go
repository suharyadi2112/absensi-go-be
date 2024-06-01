package usecase

import (
	cont "absensi/controllers"
	"fmt"
	"net/url"
	"time"
)

var controller *cont.Conn

func init() {
	var err error
	controller, err = cont.NewCon()
	if err != nil {
		panic(err) // Handle error appropriately
	}
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

// // use case simpan tap scan absen
func PostAbsenTopUsecase(formCode, tanggalhariIni, dateTimehariini, timeonlyHariini string) (res map[string]interface{}, codeHttp int, err error) {

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
		cAbsen, err := controller.GetOneAbsensiController(id_siswa, id_kelas, tanggalhariIni)

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
					err := controller.UpdateAbsenController(timeonlyHariini, tanggalhariIni, id_siswa, id_kelas)
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
					"Message": "Terjadi kesalahan hubungi admin",
				}
				return responseItem, 400, nil
			}
		} else {
			fmt.Println("sini")
		}
	}

	if countGuru > 0 {
		fmt.Println("ada guru")
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
