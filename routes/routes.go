package routes

import (
	"absensi/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Inisialisasi rute aplikasi
func InitRoutes(e *echo.Echo) {

	config := middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST},
	}

	// Mengaktifkan CORS dengan konfigurasi kustom
	e.Use(middleware.CORSWithConfig(config))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	e.GET("/get_absen_top", handler.GetAbsenTopHandler)
	e.POST("/post_absen", handler.PostAbsen)

	e.GET("/qrcode", handler.QrCode)

}
