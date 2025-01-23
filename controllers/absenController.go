package controllers

import (
	db "absensi/config"
	"absensi/models"
	"database/sql"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

// Conn struct yang menampung instance database
type Conn struct {
	DB    *sql.DB
	DBVPS *sql.DB
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

	dbVPS, err := db.InitDBMySqlVPS()
	if err != nil {
		db.InitLog(logger, ctx, "failed to initialize VPS database", err, "error") // catat log
		return nil, err
	}

	return &Conn{
		DB:    dbG,
		DBVPS: dbVPS,
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
			&a.Ket,
			&a.Ket1,
			&a.RefIn,
			&a.RefOut,
			&a.StatusUpdated,
			&a.IDSiswa.NamaLengkap, &a.IDSiswa.Foto, &a.IDKelas.Kelas, &a.IDPengajar.NamaLengkap, &a.IDPengajar.Foto,
		); err != nil {
			db.InitLog(logger, ctx, "Error executing SQL statement GetAbsenTopController", err, "error") // catat log
			return nil, err
		}
		absensi = append(absensi, a)
	}

	return absensi, nil
}

func (h *Conn) GetOneAbsensiSiswaController(idSiswa, idKelas int64, date string) (dataOneAbsen *AbsensiDetailJamMasuk, err error) {

	ctx := "Controller-GetOneAbsensiSiswaController"
	query := `SELECT absensi.id, absensi.keluar, CONCAT(tgl, ' ', masuk) as j_masuk FROM absensi WHERE id_siswa = ? AND tgl = ? AND id_kelas = ? ORDER BY
	tgl DESC`

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
						s.foto,
						o.id as id_ortu,
						o.no_hp
					FROM siswa s
					LEFT JOIN kelas k ON s.id_kelas = k.id_kelas
					LEFT JOIN orangtua o ON s.nis = o.nis
					WHERE s.nis = ?`, nis).Scan(&s.ID, &s.IDKelas.ID, &s.NIS, &s.NamaLengkap, &s.IDKelas.Kelas, &s.Alamat, &s.Foto, &s.IDOrtu.ID, &s.NoHP)

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

// insert absensi Siswa
func (h *Conn) PostInsertAbsensiSiswaController(id_siswa, id_kelas int64, tipeMasuk string, dateTimehariini, timeonlyHariini string, notif int, randomString, tipeAbsen string) (err error) {

	ctx := "Controller-PostInsertAbsensiSiswaController"
	var query string

	if tipeAbsen == "masuk" { //cek tipe absen
		query = `
		INSERT INTO absensi (id_siswa, id_kelas, absensi, tgl, masuk, notif_in, ref_in)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	} else if tipeAbsen == "keluar" {
		query = `
		INSERT INTO absensi (id_siswa, id_kelas, absensi, tgl, keluar, notif_out, ref_out)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	} else { //tidak dalam kondisi apapun
		db.InitLog(logger, ctx, "Tidak dalam kondisi tipe absen apapun PostInsertAbsensiSiswaController", err, "error")
		return fmt.Errorf("invalid tipeAbsen: %s", tipeAbsen)
	}

	stmt, err := h.DB.Prepare(query)
	if err != nil {
		db.InitLog(logger, ctx, "Error preparing SQL statement PostInsertAbsensiSiswaController", err, "error")
		return err
	}
	defer stmt.Close()
	//exe
	_, err = stmt.Exec(id_siswa, id_kelas, tipeMasuk, dateTimehariini, timeonlyHariini, notif, randomString)

	if err != nil {
		db.InitLog(logger, ctx, "Error executing SQL statement PostInsertAbsensiSiswaController", err, "error")
		return err
	}

	db.InitLog(logger, ctx, fmt.Sprintf("Absensi for siswa %d successfully inserted", id_siswa), nil, "info")
	return nil
}

func (h *Conn) PostUpdateAbsensiSiswaController(id_siswa, id_kelas int64, date, time, randomString string) (err error) {

	ctx := "Controller-PostUpdateAbsensiController"

	query := `
		UPDATE absensi
		SET keluar = ?, notif_out = ?, ref_out = ?
		WHERE id_siswa = ? AND tgl = ? AND id_kelas = ?
	`
	stmt, err := h.DB.Prepare(query)
	if err != nil {
		db.InitLog(logger, ctx, "Error preparing SQL statement PostUpdateAbsensiSiswaController", err, "error")
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(time, "0", randomString, id_siswa, date, id_kelas)
	if err != nil {
		db.InitLog(logger, ctx, "Error executing SQL statement PostUpdateAbsensiSiswaController", err, "error")
		return err
	}

	db.InitLog(logger, ctx, fmt.Sprintf("Absensi for siswa %d successfully updated", id_siswa), nil, "info")

	return nil
}

// insert absensi Guru
func (h *Conn) PostInsertAbsensiGuruController(id_pengajar int64, absensi, dateTimehariini, timeonlyHariini string, tipeAbsen string) (err error) {

	ctx := "Controller-PostInsertAbsensiGuruController"
	var query string

	if tipeAbsen == "masuk" { //cek tipe absen
		query = `
		INSERT INTO absensi (id_pengajar, absensi, tgl, masuk, status_in)
		VALUES (?, ?, ?, ?, ?)`
	} else if tipeAbsen == "keluar" {
		// query = `
		// INSERT INTO absensi (id_pengajar, absensi, tgl, masuk, status_in)
		// VALUES (?, ?, ?, ?, ?, ?, ?)`
	} else { //tidak dalam kondisi apapun
		db.InitLog(logger, ctx, "Tidak dalam kondisi tipe absen apapun PostInsertAbsensiGuruController", err, "error")
		return fmt.Errorf("invalid tipeAbsen: %s", tipeAbsen)
	}

	stmt, err := h.DB.Prepare(query)
	if err != nil {
		db.InitLog(logger, ctx, "Error preparing SQL statement PostInsertAbsensiGuruController", err, "error")
		return err
	}
	defer stmt.Close()
	//exe
	_, err = stmt.Exec(id_pengajar, absensi, dateTimehariini, timeonlyHariini, "0")

	if err != nil {
		db.InitLog(logger, ctx, "Error executing SQL statement PostInsertAbsensiGuruController", err, "error")
		return err
	}

	db.InitLog(logger, ctx, fmt.Sprintf("Absensi for guru %d successfully inserted", id_pengajar), nil, "info")
	return nil
}

// guru update
func (h *Conn) PostUpdateAbsensiGuruController(id_pengajar int64, timeOnly, date string) (err error) {

	ctx := "Controller-PostUpdateAbsensiController"

	query := `
		UPDATE absensi
		SET keluar = ?, status_out = ?
		WHERE id_pengajar = ? AND tgl = ?
	`
	stmt, err := h.DB.Prepare(query)
	if err != nil {
		db.InitLog(logger, ctx, "Error preparing SQL statement PostUpdateAbsensiGuruController", err, "error")
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(timeOnly, "0", id_pengajar, date)
	if err != nil {
		db.InitLog(logger, ctx, "Error executing SQL statement PostUpdateAbsensiGuruController", err, "error")
		return err
	}

	db.InitLog(logger, ctx, fmt.Sprintf("Absensi for guru %d successfully updated", id_pengajar), nil, "info")

	return nil
}

// Insert absensi text WA
func (h *Conn) PostInsertAbsensiWaController(randomString, sch, no_hp_ortu, masukMessage, status, dateTimehariini, tipe string) (err error) {

	ctx := "Controller-InsertAbsensiWaController"

	query := `
		INSERT INTO absensi_test (ref, school, destination, message, status, created, type)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	stmt, err := h.DBVPS.Prepare(query)
	if err != nil {
		db.InitLog(logger, ctx, "Error preparing SQL statement InsertAbsensiWAController", err, "error")
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(randomString, sch, no_hp_ortu, masukMessage, status, dateTimehariini, tipe)
	if err != nil {
		db.InitLog(logger, ctx, "Error executing SQL statement InsertAbsensiWAController", err, "error")
		return err
	}

	db.InitLog(logger, ctx, fmt.Sprintf("Absensi for siswa WA %s successfully inserted", randomString), nil, "info")

	return nil

}
