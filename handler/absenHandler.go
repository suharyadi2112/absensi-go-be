package handler

import (
	"absensi/usecase"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/skip2/go-qrcode"
)

// Handler untuk endpoint /users
func GetAbsenTopHandler(c echo.Context) error {

	currentTime := time.Now()
	tanggalhariIni := currentTime.Format("2006-01-02")

	absenTopData, err := usecase.GetAbsenTopUsecase(tanggalhariIni)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	responseUsecase := map[string]interface{}{
		"AStatus":  "success",
		"BMessage": "Get top absen retrieved",
		"CData":    absenTopData,
	}

	return c.JSON(http.StatusOK, responseUsecase)

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

func QrCode(c echo.Context) error {

	// Data to be encoded into QR code
	data := "rss"

	// Generate QR code as []byte
	qrCode, err := qrcode.Encode(data, qrcode.Medium, 256)
	if err != nil {
		log.Fatal("Error generating QR code: ", err)
	}

	// Convert []byte to base64 string
	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCode)

	// Print or return the base64 string
	log.Println("QR code base64:", qrCodeBase64)

	response := map[string]interface{}{
		"AStatus":  "success",
		"BMessage": "Success generate",
		"CData":    qrCodeBase64,
	}

	return c.JSON(http.StatusOK, response)
}
