package main

import (
	"flag"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	key := flag.String("k", "", "API Key for admin endpoints")
	flag.Parse()

	if len(*key) == 0 {
		log.Fatalf("api key (-k) must be set")
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	admin := e.Group("/admin", middleware.KeyAuth(func(k string, c echo.Context) (bool, error) {
		return k == *key, nil
	}))

	admin.POST("/score", nil)       // start creating a score (and switch display)
	admin.POST("/score/:id", nil)   // update or complete a pending score
	admin.DELETE("/score/:id", nil) // cancel a score update

	e.GET("/", nil)       // main page
	e.GET("/scores/:cls") // scores for a specific class
	e.GET("/live", nil)   // live updates for switching pages and updating

	e.Logger.Fatal(e.Start(":8080"))
}
