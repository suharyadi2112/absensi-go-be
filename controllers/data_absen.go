package controllers

import (
	db "absensi/config"
	"absensi/models"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Handler untuk endpoint /users
func GetAbsenTop(c echo.Context) error {

	dbG, _ := db.InitDB()

	// Tentukan tanggal yang akan diambil
	date := "2024-05-27"

	// Eksekusi kueri SQL
	rows, err := dbG.Query(`
		SELECT
			absensi.*,
			siswa.nama_lengkap,
			siswa.foto,
			kelas.kelas,
			pengajar.nama_lengkap AS nm_guru,
			pengajar.foto AS foto_guru
		FROM
			absensi
		LEFT JOIN siswa ON absensi.id_siswa = siswa.id_siswa
		LEFT JOIN kelas ON siswa.id_kelas = kelas.id_kelas
		LEFT JOIN pengajar ON absensi.id_pengajar = pengajar.id_pengajar
		WHERE
			absensi.tgl = ?
		ORDER BY
			absensi.id DESC
		LIMIT 5`, date)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate through the result set
	var absensi []models.Absensi
	for rows.Next() {
		var a models.Absensi
		if err := rows.Scan(
			&a.ID, &a.IDSiswa.ID, &a.IDKelas.ID, &a.IDPengajar.ID,
			&a.Absensi, &a.Tanggal, &a.Masuk, &a.Keluar,
			&a.StatusMasuk, &a.StatusKeluar, &a.NotifikasiMasuk, &a.NotifikasiKeluar,
			&a.Updated, &a.UpdateAbsensi,
			&a.IDSiswa.NamaLengkap, &a.IDSiswa.Foto, &a.IDKelas.Kelas, &a.IDPengajar.NamaLengkap, &a.IDPengajar.Foto,
		); err != nil {
			log.Fatal(err)
		}
		absensi = append(absensi, a)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	//custom return yang diperlukan response
	var absentopResp []map[string]interface{}
	for _, s := range absensi {
		absentopResp = append(absentopResp, map[string]interface{}{
			"IDAbsensi": s.ID.Int64,
			"FotoSiswa": s.IDSiswa.Foto.String,
			"FotoGuru":  s.IDPengajar.Foto.String,
			"NamaSiswa": s.IDSiswa.NamaLengkap.String,
			"NamaGuru":  s.IDPengajar.NamaLengkap.String,
			"Kelas":     s.IDKelas.Kelas.String,
		})
	}

	defer rows.Close()

	response := map[string]interface{}{
		"AStatus":  "success",
		"BMessage": "Get top absen retrieved",
		"CData":    absentopResp,
	}

	return c.JSON(http.StatusOK, response)
}
