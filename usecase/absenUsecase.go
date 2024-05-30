package usecase

import (
	cont "absensi/controllers"
	"fmt"
)

type AbsenTopUseCase interface {
	GetAbsenTopUsecase(date string) ([]map[string]interface{}, error)
}

// Handler untuk endpoint /users
func GetAbsenTopUsecase(tanggalhariIni string) ([]map[string]interface{}, error) {

	fmt.Println("Tanggal sekarang - get absen top:", tanggalhariIni)

	handler, err := cont.NewCon()
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	result := handler.GetAbsenTopQuery(tanggalhariIni)
	rowsData := result.DataAbsen
	err = result.Err

	if err != nil {
		return nil, err
	}

	//custom return yang diperlukan response
	var absentopResp []map[string]interface{}
	for _, s := range rowsData {
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
