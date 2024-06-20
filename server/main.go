package main

import (
	"absensi/routes"

	// conFig "absensi/config"

	"github.com/labstack/echo/v4"
	// "github.com/sirupsen/logrus"
)

// var (
// 	logger *logrus.Logger
// )

// func init() {
// 	logger = conFig.InitLogRus()
// }

func main() {

	// ctx := "RunningServer"
	e := echo.New()

	routes.InitRoutes(e)

	// conFig.InitlogError(logger, ctx, "server running at - localhost:1323", nil, "info") // catat log
	e.Logger.Fatal(e.Start("localhost:1323"))
	// e.Logger.Fatal(e.Start("192.168.0.244:1323"))
}
