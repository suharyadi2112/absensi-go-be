package main

import (
	"absensi/routes"

	conFig "absensi/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

func init() {
	logger = conFig.InitLogRus()
}

func main() {

	ctx := "RunningServer"
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://103.175.216.178:5173",
			"http://localhost:5173",
		},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	routes.InitRoutes(e)

	conFig.InitLog(logger, ctx, "server running at - localhost:1323", nil, "info") // catat log
	e.Logger.Fatal(e.Start("localhost:1323"))
	// e.Logger.Fatal(e.Start("192.168.0.244:1323"))

	logrus.Info("Server run at localhost:1323")

}
