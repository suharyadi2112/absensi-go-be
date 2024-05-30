package usecase

import (
	cont "absensi/controllers"
	"fmt"
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

	result := controller.GetAbsenTopQuery(tanggalhariIni)

	if result.Err != nil {
		return nil, result.Err
	}

	//custom return yang diperlukan response
	var absentopResp []map[string]interface{}
	for _, s := range result.DataAbsen {
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

// // use case simpan tap scan
func PostAbsenTopUsecase(formCode string) ([]map[string]interface{}, error) {

	fmt.Println("Form code - post absen:", formCode)

	result := controller.PostAbsenTopQuery(formCode)

	if result.Err != nil {
		return nil, result.Err
	}

	//custom return yang diperlukan response
	var absentopResp []map[string]interface{}
	for _, s := range result.DataAbsen {
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
