package models

import (
	"database/sql"
)

type Absensi struct {
	ID               sql.NullInt64  `json:"id"`
	IDPengajar       Pengajar       `json:"id_pengajar"`
	IDSiswa          Siswa          `json:"id_siswa"`
	IDKelas          Kelas          `json:"id_kelas"`
	Absensi          sql.NullString `json:"absensi"`
	Tanggal          sql.NullString `json:"tgl"`
	Masuk            sql.NullString `json:"masuk"`
	Keluar           sql.NullString `json:"keluar"`
	StatusMasuk      sql.NullInt64  `json:"status_in"`
	StatusKeluar     sql.NullInt64  `json:"status_out"`
	NotifikasiMasuk  sql.NullInt64  `json:"notif_in"`
	NotifikasiKeluar sql.NullInt64  `json:"notif_out"`
	Updated          sql.NullString `json:"updated"`
	UpdateAbsensi    sql.NullString `json:"update_absensi"`
}

type Siswa struct {
	ID            sql.NullInt64  `json:"id_siswa"`
	NIS           sql.NullString `json:"nis"`
	NamaLengkap   sql.NullString `json:"nama_lengkap"`
	UsernameLogin sql.NullString `json:"username_login"`
	PasswordLogin sql.NullString `json:"password_login"`
	Pswd          sql.NullString `json:"pswd"`
	IDKelas       Kelas          `json:"id_kelas"`
	Alamat        sql.NullString `json:"alamat"`
	Email         sql.NullString `json:"email"`
	NoHP          sql.NullString `json:"no_hp"`
	TempatLahir   sql.NullString `json:"tempat_lahir"`
	TanggalLahir  sql.NullString `json:"tgl_lahir"`
	JenisKelamin  sql.NullString `json:"jenis_kelamin"`
	Agama         sql.NullString `json:"agama"`
	Foto          sql.NullString `json:"foto"`
	TahunMasuk    sql.NullString `json:"th_masuk"`
	Blokir        sql.NullString `json:"blokir"`
	Created       sql.NullString `json:"created"`
	Updated       sql.NullString `json:"updated"`
	Poin          sql.NullInt64  `json:"poin"`
	Saldo         sql.NullInt64  `json:"saldo"`
	PIN           sql.NullInt64  `json:"pin"`
}

type Kelas struct {
	ID         sql.NullInt64  `json:"id"`
	Kelas      sql.NullString `json:"kelas"`
	IDPengajar sql.NullInt64  `json:"id_pengajar"`
	IDSemester sql.NullInt64  `json:"id_semester"`
	Created    sql.NullString `json:"created"`
}

type Pengajar struct {
	ID            sql.NullInt64  `json:"id"`
	NIP           sql.NullString `json:"nip"`
	NamaLengkap   sql.NullString `json:"nama_lengkap"`
	UsernameLogin sql.NullString `json:"username_login"`
	PasswordLogin sql.NullString `json:"password_login"`
	Pswd          sql.NullString `json:"pswd"`
	Alamat        sql.NullString `json:"alamat"`
	TempatLahir   sql.NullString `json:"tempat_lahir"`
	TanggalLahir  sql.NullString `json:"tgl_lahir"`
	JenisKelamin  sql.NullString `json:"jenis_kelamin"`
	Agama         sql.NullString `json:"agama"`
	NoTelp        sql.NullString `json:"no_telp"`
	Email         sql.NullString `json:"email"`
	Foto          sql.NullString `json:"foto"`
	Blokir        sql.NullString `json:"blokir"`
	Created       sql.NullString `json:"created"`
	Updated       sql.NullString `json:"updated"`
}
type CountSiswa struct {
	DataCountSiswa int
	Err            error
}
type CountGuru struct {
	DataCountGuru int
	Err           error
}
