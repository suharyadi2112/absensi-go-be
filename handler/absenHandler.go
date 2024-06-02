package handler

import (
	usec "absensi/usecase"
	"encoding/base64"

	// usec "absensi/usecase"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/skip2/go-qrcode"
)

var (
	currentTime     time.Time
	tanggalHariIni  string
	datetimeHariini string
	timeonlyHariini string
	usecase         *usec.AbsenUsecase
)

// Struct untuk menangkap data JSON
type AbsenForm struct {
	FormCode string `json:"form_code" validate:"required"`
}

func init() {
	var err error
	usecase, err = usec.NewConUsecase()
	if err != nil {
		panic(err.Error())
	}

	currentTime = time.Now()
	tanggalHariIni = currentTime.Format("2006-01-02")
	datetimeHariini = currentTime.Format("2006-01-02 15:04:05")
	timeonlyHariini = currentTime.Format("15:04:05")
}

// get absen top
func GetAbsenTopHandler(c echo.Context) error {

	absenTopData, err := usec.GetAbsenTopUsecase(tanggalHariIni)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	responseUsecase := map[string]interface{}{
		"AStatus":  "Success",
		"BMessage": "Get top absen retrieved",
		"CData":    absenTopData,
	}

	return c.JSON(http.StatusOK, responseUsecase)

}

// post absen
func PostAbsenHandler(c echo.Context) error {

	u := &AbsenForm{}
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request payload",
		})
	}

	validator := validator.New()
	err := validator.Struct(u)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	formCode := u.FormCode
	dataAbsenPost, status, err := usecase.PostAbsenTopUsecase(formCode, tanggalHariIni, datetimeHariini, timeonlyHariini)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if status == 400 {
		responseUsecase := map[string]interface{}{
			"AStatus":  "Failed",
			"BMessage": dataAbsenPost,
			"CData":    nil,
		}
		return c.JSON(http.StatusBadRequest, responseUsecase)
	}

	responseUsecase := map[string]interface{}{
		"AStatus":  "Success",
		"BMessage": "Succees process absen",
		"CData":    dataAbsenPost,
	}

	return c.JSON(http.StatusOK, responseUsecase)

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
