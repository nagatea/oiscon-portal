package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	// gorm mysql driver
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/nagatea/oiscon-portal/model"
)

func main() {

	db, err := model.EstablishConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.LogMode(true)

	err = model.Migrate()
	if err != nil {
		panic(err)
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Route => handler
	e.GET("/api/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	// Start server
	e.Logger.Fatal(e.Start(":3001"))
}
