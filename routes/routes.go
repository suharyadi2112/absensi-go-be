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
	e.GET("/get_absen_top", handler.GetAbsenTop)
	e.POST("/post_absen", handler.PostAbsen)

}
