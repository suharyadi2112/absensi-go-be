package main

import (
	"absensi/routes"

	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()

	routes.InitRoutes(e)

	e.Logger.Fatal(e.Start("localhost:1323"))
}
