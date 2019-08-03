package main

import (
	"flag"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/epipho/tfh_scoreboard/api"
	"github.com/epipho/tfh_scoreboard/api/admin"
	"github.com/epipho/tfh_scoreboard/scorer"
)

func main() {
	key := flag.String("k", "", "API Key for admin endpoints")
	flag.Parse()

	if len(*key) == 0 {
		log.Fatalf("api key (-k) must be set")
	}

	sc := scorer.NewInMemoryScorer(nil, nil)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	ag := e.Group("/admin", middleware.KeyAuth(func(k string, c echo.Context) (bool, error) {
		return k == *key, nil
	}))

	ag.POST("/score", admin.CreateScore(sc))       // start creating a score (and switch display)
	ag.POST("/score/:id", admin.UpdateScore(sc))   // update or complete a pending score
	ag.DELETE("/score/:id", admin.DeleteScore(sc)) // cancel a score update

	e.GET("/", api.Index())                // main page
	e.GET("/scores/:cls", api.GetScores()) // scores for a specific class
	e.GET("/live", api.Live())             // live updates for switching pages and updating

	e.Logger.Fatal(e.Start(":8080"))
}
