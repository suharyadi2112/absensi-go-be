package handler

import (
	cont "absensi/controllers"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// Handler untuk endpoint /users
func GetAbsenTop(c echo.Context) error {

	currentTime := time.Now()
	date := currentTime.Format("2006-01-02")

	fmt.Println("Tanggal sekarang - get absen top:", date)

	handler, err := cont.NewCon()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	result := handler.GetAbsenTopQuery(date)
	rowsData := result.DataAbsen
	err = result.Err

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	//custom return yang diperlukan response
	var absentopResp []map[string]interface{}
	for _, s := range rowsData {
		absentopResp = append(absentopResp, map[string]interface{}{
			"IDAbsensi": s.ID.Int64,
			"FotoSiswa": s.IDSiswa.Foto.String,
			"FotoGuru":  s.IDPengajar.Foto.String,
			"NamaSiswa": s.IDSiswa.NamaLengkap.String,
			"NamaGuru":  s.IDPengajar.NamaLengkap.String,
			"Kelas":     s.IDKelas.Kelas.String,
		})
	}

	response := map[string]interface{}{
		"AStatus":  "success",
		"BMessage": "Get top absen retrieved",
		"CData":    absentopResp,
	}

	return c.JSON(http.StatusOK, response)
}

// Handler untuk endpoint /users
func PostAbsen(c echo.Context) error {

	// Tentukan tanggal yang akan diambil
	// Mendapatkan tanggal sekarang
	currentTime := time.Now()

	// Mengonversi tanggal ke dalam format yang diinginkan (YYYY-MM-DD)
	date := currentTime.Format("2006-01-02")

	fmt.Println("Tanggal sekarang:", date)

	return nil
}
