package handler

import (
	usec "absensi/usecase"
	"encoding/base64"

	conFig "absensi/config"

	// usec "absensi/usecase"

	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
)

var (
	currentTime     time.Time
	tanggalHariIni  string
	datetimeHariini string
	timeonlyHariini string
	usecase         *usec.AbsenUsecase
	logger          *logrus.Logger
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

	logger = conFig.InitLogRus()

	currentTime = time.Now()
	tanggalHariIni = currentTime.Format("2006-01-02")
	datetimeHariini = currentTime.Format("2006-01-02 15:04:05")
	timeonlyHariini = currentTime.Format("15:04:05")
}

// get absen top
func GetAbsenTopHandler(c echo.Context) error {
	ctx := "Handler-GetAbsenTopHandler"
	absenTopData, err := usec.GetAbsenTopUsecase(tanggalHariIni)
	// err = errors.New("math: square root of negative number")

	if err != nil {
		conFig.InitLog(logger, ctx, "error get absen top", err, "error") // catat log
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	//response
	responseUsecase := map[string]interface{}{
		"AStatus":  "Success",
		"BMessage": "Get top absen retrieved",
		"CData":    absenTopData,
	}
	conFig.InitLog(logger, ctx, "success get absen top", nil, "info") // catat log

	return c.JSON(http.StatusOK, responseUsecase)

}

// post absen
func PostAbsenHandler(c echo.Context) error {
	ctx := "Handler-PostAbsenHandler"
	u := &AbsenForm{}
	if err := c.Bind(u); err != nil {
		conFig.InitLog(logger, ctx, "error bind payload", err, "error") // catat log
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request payload",
		})
	}

	validator := validator.New()
	err := validator.Struct(u)
	if err != nil {
		conFig.InitLog(logger, ctx, "error validator payload", err, "error") // catat log
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	formCode := u.FormCode
	dataAbsenPost, status, err := usecase.PostAbsenTopUsecase(formCode, tanggalHariIni, datetimeHariini, timeonlyHariini)

	if err != nil {
		conFig.InitLog(logger, ctx, "error post data to PostAbsenTopUsecase", err, "error") // catat log
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if status == 400 {
		responseUsecase := map[string]interface{}{
			"AStatus":  "Failed",
			"BMessage": dataAbsenPost,
			"CData":    nil,
		}
		conFig.InitLog(logger, ctx, "success send response 400", nil, "info") // catat log
		return c.JSON(http.StatusBadRequest, responseUsecase)
	}

	responseUsecase := map[string]interface{}{
		"AStatus":  "Success",
		"BMessage": "Succees process absen",
		"CData":    dataAbsenPost,
	}
	conFig.InitLog(logger, ctx, "success send response 200", nil, "info") // catat log

	return c.JSON(http.StatusOK, responseUsecase)

}

func QrCode(c echo.Context) error {
	ctx := "Handler-QrCode"
	// Data to be encoded into QR code
	data := "rss"

	// Generate QR code as []byte
	qrCode, err := qrcode.Encode(data, qrcode.Medium, 256)
	if err != nil {
		conFig.InitLog(logger, ctx, "Error generating QR code: ", err, "error") // catat log
	}

	// Convert []byte to base64 string
	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCode)

	// Print or return the base64 string
	conFig.InitLog(logger, ctx, "success generate qrcode", nil, "info") // catat log

	response := map[string]interface{}{
		"AStatus":  "success",
		"BMessage": "Success generate",
		"CData":    qrCodeBase64,
	}

	conFig.InitLog(logger, ctx, "success send response 200", nil, "info") // catat log

	return c.JSON(http.StatusOK, response)
}
