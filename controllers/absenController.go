package controllers

import (
	db "absensi/config"
	"absensi/models"
	"database/sql"
	"fmt"
	"log"
)

// Conn struct yang menampung instance database
type Conn struct {
	DB *sql.DB
}

// Fungsi untuk inisialisasi handler dengan instance database
func NewCon() (*Conn, error) {
	dbG, err := db.InitDBMySql()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &Conn{
		DB: dbG,
	}, nil
}

type AbsenTopResult struct {
	DataAbsen []models.Absensi
	Err       error
}

func (h *Conn) GetAbsenTopQuery(dateS string) AbsenTopResult {

	// Eksekusi kueri SQL
	rows, err := h.DB.Query(`
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
		LIMIT 5`, dateS)

	if err != nil {
		return AbsenTopResult{nil, err}
	}

	defer rows.Close()

	var absensi []models.Absensi
	for rows.Next() {
		var a models.Absensi
		if err := rows.Scan(
			&a.ID, &a.IDPengajar.ID, &a.IDSiswa.ID, &a.IDKelas.ID,
			&a.Absensi, &a.Tanggal, &a.Masuk, &a.Keluar,
			&a.StatusMasuk, &a.StatusKeluar, &a.NotifikasiMasuk, &a.NotifikasiKeluar,
			&a.Updated, &a.UpdateAbsensi,
			&a.IDSiswa.NamaLengkap, &a.IDSiswa.Foto, &a.IDKelas.Kelas, &a.IDPengajar.NamaLengkap, &a.IDPengajar.Foto,
		); err != nil {
			log.Fatal(err)
		}
		absensi = append(absensi, a)
	}

	return AbsenTopResult{absensi, err}
}
