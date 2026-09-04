package main

import (
	"log"

	"github.com/NooraPanahi/Weblog-Application.git/internal/database"
	"github.com/labstack/echo/v4"
)

func main () {
	connStr := "user=web password=web dbname=web host=localhost port=5433 sslmode=disable"
	db,err := database.NewPostgres(connStr)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Weblog Application")
	})

	e.Logger.Fatal(e.Start(":8080"))
}