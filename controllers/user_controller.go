package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

// Handler untuk endpoint /users
func GetUsers(c echo.Context) error {
	// Logika bisnis untuk mendapatkan daftar pengguna dari database
	users := []User{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Doe", Email: "jane@example.com"},
	}
	return c.JSON(http.StatusOK, users)
}
