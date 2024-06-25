package controllers

import (
	db "absensi/config"
	"absensi/models"
	"database/sql"
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

// Conn struct yang menampung instance database
type Conn struct {
	DB *sql.DB
}
type AbsensiDetailJamMasuk struct {
	Absensi *models.Absensi
	JMasuk  sql.NullString
}

func init() {
	logger = db.InitLogRus()
}

// Fungsi untuk inisialisasi handler dengan instance database
func NewCon() (*Conn, error) {
	ctx := "Controller-NewCon"
	dbG, err := db.InitDBMySql()
	if err != nil {
		db.InitLog(logger, ctx, "failed to initialize database", err, "error") // catat log
		return nil, err
	}

	return &Conn{
		DB: dbG,
	}, nil
}

func (h *Conn) GetAbsenTopController(dateS string) (DataAbsen []*models.Absensi, err error) {

	ctx := "Controller-GetAbsenTopController"
	conn, err := NewCon()
	if err != nil {
		db.InitLog(logger, ctx, "error koneksi database", err, "error") // catat log
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
		db.InitLog(logger, ctx, "error query database GetAbsenTopController", err, "error") // catat log
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
			db.InitLog(logger, ctx, "Error executing SQL statement GetAbsenTopController", err, "error") // catat log
			return nil, err
		}
		absensi = append(absensi, a)
	}

	return absensi, nil
}

// Update absensi
func (h *Conn) UpdateAbsenController(Keluar, tanggalHariIni string, idSiswa, idKelas int64) (err error) {

	ctx := "Controller-UpdateAbsenController"

	// Prepare the UPDATE query
	stmt, err := h.DB.Prepare("UPDATE absensi SET keluar = ?, notif_out = 0 WHERE id_siswa = ? AND tgl = ? AND id_kelas = ?")
	if err != nil {
		db.InitLog(logger, ctx, "error query database UpdateAbsenController", err, "error") // catat log
		return err
	}
	defer stmt.Close()
	// Execute the UPDATE query
	_, err = stmt.Exec(Keluar, idSiswa, tanggalHariIni, idKelas)
	if err != nil {
		db.InitLog(logger, ctx, "Error executing SQL statement UpdateAbsenController", err, "error") // catat log
		return err
	}
	return nil
}

func (h *Conn) PostAbsenSiswaController(timeOnly, tanggalHariIni, tipeAbsen string, idSiswa, idKelas int64) (err error) {

	ctx := "Controller-PostAbsenSiswaController"
	var query string

	// Prepare the SQL statement
	if tipeAbsen == "masuk" {
		query = "INSERT INTO absensi (id_siswa, id_kelas, absensi, tgl, masuk, notif_in) VALUES (?, ?, ?, ?, ?, ?)"
	} else if tipeAbsen == "keluar" {
		query = "INSERT INTO absensi (id_siswa, id_kelas, absensi, tgl, keluar, notif_out) VALUES (?, ?, ?, ?, ?, ?)"
	} else {
		db.InitLog(logger, ctx, "warning invalid tipeAbsen", nil, "warning") // catat log
		return errors.New("invalid tipeAbsen")
	}
	stmt, err := h.DB.Prepare(query)
	if err != nil {
		db.InitLog(logger, ctx, "Error preparing SQL statement PostAbsenSiswaController", err, "error") // catat log
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(idSiswa, idKelas, "H", tanggalHariIni, timeOnly, "0")
	if err != nil {
		db.InitLog(logger, ctx, "Error executing SQL statement PostAbsenSiswaController", err, "error") // catat log
		return err
	}

	return nil
}

func (h *Conn) PostAbsenGuruController(timeOnly, tanggalHariIni string, idPengajar int64) (err error) {

	ctx := "Controller-PostAbsenGuruController"
	query := "INSERT INTO absensi (id_pengajar, absensi, tgl, masuk, status_in) VALUES (?, ?, ?, ?, ?)"

	stmt, err := h.DB.Prepare(query)
	if err != nil {
		db.InitLog(logger, ctx, "Error preparing SQL statement PostAbsenGuruController", err, "error") // catat log
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(idPengajar, "H", tanggalHariIni, timeOnly, "0")
	if err != nil {
		db.InitLog(logger, ctx, "Error Execute SQL statement PostAbsenGuruController", err, "error") // catat log
		return err
	}

	return nil
}

func (h *Conn) GetOneAbsensiSiswaController(idSiswa, idKelas int64, date string) (dataOneAbsen *AbsensiDetailJamMasuk, err error) {

	ctx := "Controller-GetOneAbsensiSiswaController"
	query := `SELECT absensi.id, absensi.keluar, CONCAT(tgl, ' ', masuk) as j_masuk FROM absensi WHERE id_siswa = ? AND tgl = ? AND id_kelas = ?`

	row := h.DB.QueryRow(query, idSiswa, date, idKelas)

	var absensi models.Absensi
	var jMasuk sql.NullString

	err = row.Scan(&absensi.ID, &absensi.Keluar, &jMasuk)

	if err == sql.ErrNoRows {
		// Data absen tidak ditemukan
		db.InitLog(logger, ctx, "check data siswa not found GetOneAbsensiSiswaController", nil, "info") // catat log
		return nil, nil                                                                                 // Bukan error 500
	} else if err != nil {
		db.InitLog(logger, ctx, "Error preparing SQL statement GetOneAbsensiSiswaController", err, "error") // catat log
		return nil, err
	}

	absensiDetail := &AbsensiDetailJamMasuk{
		Absensi: &absensi,
		JMasuk:  jMasuk, // Convert 'jMasuk' to string
	}

	return absensiDetail, nil
}

func (h *Conn) GetSiswaController(nis string) (DataSiswa *models.Siswa, err error) {

	ctx := "Controller-GetSiswaController"
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
		db.InitLog(logger, ctx, "Error preparing SQL statement GetSiswaController", err, "error") // catat log
		return nil, err
	}

	return &s, nil
}

func (h *Conn) GetGuruController(nip string) (DataGuru *models.Pengajar, err error) {

	ctx := "Controller-GetGuruController"
	query := `SELECT id_pengajar, nip, nama_lengkap, username_login, password_login, pswd, alamat, tempat_lahir, tgl_lahir, jenis_kelamin, agama, no_telp, email, foto, blokir FROM pengajar WHERE nip = ?`

	row := h.DB.QueryRow(query, nip)

	var pengajar models.Pengajar
	err = row.Scan(
		&pengajar.ID,
		&pengajar.NIP,
		&pengajar.NamaLengkap,
		&pengajar.UsernameLogin,
		&pengajar.PasswordLogin,
		&pengajar.Pswd,
		&pengajar.Alamat,
		&pengajar.TempatLahir,
		&pengajar.TanggalLahir,
		&pengajar.JenisKelamin,
		&pengajar.Agama,
		&pengajar.NoTelp,
		&pengajar.Email,
		&pengajar.Foto,
		&pengajar.Blokir,
	)
	if err != nil {
		db.InitLog(logger, ctx, "Error Execute SQL statement GetGuruController", err, "error") // catat log
		return nil, err
	}

	return &pengajar, nil
}

func (h *Conn) CountSiswaController(formCode string) (DataCountSiswa int, err error) {

	ctx := "Controller-GetGuruController"
	sanitizedCode := strings.ReplaceAll(formCode, ".", "")
	sanitizedCode = strings.ReplaceAll(sanitizedCode, " ", "")

	var qSiswa int
	err = h.DB.QueryRow(`SELECT COUNT(*) FROM siswa WHERE nis = ?`, sanitizedCode).Scan(&qSiswa)
	if err == sql.ErrNoRows {
		// Data absen tidak ditemukan
		db.InitLog(logger, ctx, "check data siswa not found CountSiswaController", nil, "info") // catat log
		return 0, nil                                                                           // Bukan error 500
	} else if err != nil {
		db.InitLog(logger, ctx, "Error preparing SQL statement CountSiswaController", err, "error") // catat log
		return 0, err
	}

	return qSiswa, nil

}
func (h *Conn) CountGuruController(formCode string) (DataCountGuru int, err error) {

	ctx := "Controller-CountGuruController"
	sanitizedCode := strings.ReplaceAll(formCode, ".", "")
	sanitizedCode = strings.ReplaceAll(sanitizedCode, " ", "")

	var qGuru int
	err = h.DB.QueryRow("SELECT COUNT(*) FROM pengajar WHERE nip = ?", sanitizedCode).Scan(&qGuru)
	if err == sql.ErrNoRows {
		// Data absen tidak ditemukan
		db.InitLog(logger, ctx, "check data guru not found CountGuruController", nil, "info") // catat log
		return 0, nil                                                                         // Bukan error 500
	} else if err != nil {
		db.InitLog(logger, ctx, "Error preparing SQL statement CountGuruController", err, "error") // catat log
		return 0, err
	}

	return qGuru, nil

}

func (h *Conn) GetOneAbsensiGuruController(idPengajar int64, date string) (dataOneAbsen *AbsensiDetailJamMasuk, err error) {

	ctx := "Controller-GetOneAbsensiGuruController"
	query := `SELECT absensi.id, absensi.keluar, CONCAT(tgl, ' ', masuk) as j_masuk FROM absensi WHERE id_pengajar = ? AND tgl = ?`

	row := h.DB.QueryRow(query, idPengajar, date)

	var absensi models.Absensi
	var jMasuk sql.NullString

	err = row.Scan(&absensi.ID, &absensi.Keluar, &jMasuk)

	if err == sql.ErrNoRows {
		// Data absen tidak ditemukan
		db.InitLog(logger, ctx, "check data guru not found GetOneAbsensiGuruController", nil, "check") // catat log
		return nil, nil                                                                                // Bukan error 500
	} else if err != nil {
		db.InitLog(logger, ctx, "Error preparing SQL statement GetOneAbsensiGuruController", err, "error") // catat log
		return nil, err
	}

	absensiDetail := &AbsensiDetailJamMasuk{
		Absensi: &absensi,
		JMasuk:  jMasuk, // Convert 'jMasuk' to string
	}

	return absensiDetail, nil
}

// Update absensi guru
func (h *Conn) UpdateAbsenGuruController(Keluar, tanggalHariIni string, idPengajar int64) (err error) {

	ctx := "Controller-UpdateAbsenGuruController"
	// Prepare the UPDATE query
	stmt, err := h.DB.Prepare("UPDATE absensi SET keluar = ?, notif_out = 0 WHERE id_pengajar = ? AND tgl = ?")
	if err != nil {
		db.InitLog(logger, ctx, "Error preparing SQL statement UpdateAbsenGuruController", err, "error") // catat log
		return err
	}
	defer stmt.Close()
	// Execute the UPDATE query
	_, err = stmt.Exec(Keluar, idPengajar, tanggalHariIni)
	if err != nil {
		db.InitLog(logger, ctx, "Error Execute SQL statement UpdateAbsenGuruController", err, "error") // catat log
		return err
	}

	return nil
}
