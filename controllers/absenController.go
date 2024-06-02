package controllers

import (
	db "absensi/config"
	"absensi/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
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

type AbsensiDetailJamMasuk struct {
	Absensi *models.Absensi
	JMasuk  sql.NullString
}

func (h *Conn) GetAbsenTopController(dateS string) (DataAbsen []*models.Absensi, err error) {

	conn, err := NewCon()
	if err != nil {
		log.Fatal(err)
	}
	// Eksekusi kueri SQL
	rows, err := conn.DB.Query(`
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
		return nil, nil
	}

	defer rows.Close()

	var absensi []*models.Absensi
	for rows.Next() {
		var a = new(models.Absensi)
		if err := rows.Scan(
			&a.ID, &a.IDPengajar.ID, &a.IDSiswa.ID, &a.IDKelas.ID,
			&a.Absensi, &a.Tanggal, &a.Masuk, &a.Keluar,
			&a.StatusMasuk, &a.StatusKeluar, &a.NotifikasiMasuk, &a.NotifikasiKeluar,
			&a.Updated, &a.UpdateAbsensi,
			&a.IDSiswa.NamaLengkap, &a.IDSiswa.Foto, &a.IDKelas.Kelas, &a.IDPengajar.NamaLengkap, &a.IDPengajar.Foto,
		); err != nil {
			return nil, err
		}
		absensi = append(absensi, a)
	}

	return absensi, nil
}

// Update absensi
func (h *Conn) UpdateAbsenController(Keluar, tanggalHariIni string, idSiswa, idKelas int64) (err error) {

	fmt.Println(Keluar, "tytyty")

	// Prepare the UPDATE query
	stmt, err := h.DB.Prepare("UPDATE absensi SET keluar = ?, notif_out = 0 WHERE id_siswa = ? AND tgl = ? AND id_kelas = ?")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()
	// Execute the UPDATE query
	_, err = stmt.Exec(Keluar, idSiswa, tanggalHariIni, idKelas)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("success update absen", idSiswa, idKelas)

	return nil
}

func (h *Conn) PostAbsenController(timeOnly, tanggalHariIni, tipeAbsen string, idSiswa, idKelas int64) (err error) {

	var query string

	// Prepare the SQL statement
	if tipeAbsen == "masuk" {
		query = "INSERT INTO absensi (id_siswa, id_kelas, absensi, tgl, masuk, notif_in) VALUES (?, ?, ?, ?, ?, ?)"
	} else if tipeAbsen == "keluar" {
		query = "INSERT INTO absensi (id_siswa, id_kelas, absensi, tgl, keluar, notif_out) VALUES (?, ?, ?, ?, ?, ?)"
	} else {
		return errors.New("invalid tipeAbsen")
	}
	stmt, err := h.DB.Prepare(query)
	if err != nil {
		log.Fatal("Error preparing SQL statement:", err)
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(idSiswa, idKelas, "H", tanggalHariIni, timeOnly, "0")
	if err != nil {
		log.Fatal("Error executing SQL statement:", err)
		return err
	}

	return nil
}

func (h *Conn) GetOneAbsensiController(idSiswa, idKelas int64, date string) (dataOneAbsen *AbsensiDetailJamMasuk, err error) {

	query := `SELECT absensi.id, absensi.keluar, CONCAT(tgl, ' ', masuk) as j_masuk FROM absensi WHERE id_siswa = ? AND tgl = ? AND id_kelas = ?`

	row := h.DB.QueryRow(query, idSiswa, date, idKelas)

	var absensi models.Absensi
	var jMasuk sql.NullString

	err = row.Scan(&absensi.ID, &absensi.Keluar, &jMasuk)

	if err == sql.ErrNoRows {
		// Data absen tidak ditemukan
		return nil, nil // Bukan error 500
	} else if err != nil {
		return nil, err
	}

	absensiDetail := &AbsensiDetailJamMasuk{
		Absensi: &absensi,
		JMasuk:  jMasuk, // Convert 'jMasuk' to string
	}

	return absensiDetail, nil
}

func (h *Conn) GetSiswaController(nis string) (DataSiswa *models.Siswa, err error) {

	var s models.Siswa
	err = h.DB.QueryRow(`SELECT
						s.id_siswa,
						s.id_kelas,
						s.nis,
						s.nama_lengkap,
						k.kelas,
						s.alamat,
						s.foto
					FROM siswa s
					LEFT JOIN kelas k ON s.id_kelas = k.id_kelas
					LEFT JOIN orangtua o ON s.nis = o.nis
					WHERE s.nis = ?`, nis).Scan(&s.ID, &s.IDKelas.ID, &s.NIS, &s.NamaLengkap, &s.IDKelas.Kelas, &s.Alamat, &s.Foto)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (h *Conn) CountSiswaController(formCode string) (DataCountSiswa int, err error) {

	sanitizedCode := strings.ReplaceAll(formCode, ".", "")
	sanitizedCode = strings.ReplaceAll(sanitizedCode, " ", "")

	var qSiswa int
	err = h.DB.QueryRow(`SELECT COUNT(*) FROM siswa WHERE nis = ?`, sanitizedCode).Scan(&qSiswa)
	if err == sql.ErrNoRows {
		// Data absen tidak ditemukan
		return 0, nil // Bukan error 500
	} else if err != nil {
		return 0, err
	}

	return qSiswa, nil

}
func (h *Conn) CountGuruController(formCode string) (DataCountGuru int, err error) {

	sanitizedCode := strings.ReplaceAll(formCode, ".", "")
	sanitizedCode = strings.ReplaceAll(sanitizedCode, " ", "")

	var qGuru int
	err = h.DB.QueryRow("SELECT COUNT(*) FROM pengajar WHERE nip = ?", sanitizedCode).Scan(&qGuru)
	if err == sql.ErrNoRows {
		// Data absen tidak ditemukan
		return 0, nil // Bukan error 500
	} else if err != nil {
		return 0, err
	}

	return qGuru, nil

}
