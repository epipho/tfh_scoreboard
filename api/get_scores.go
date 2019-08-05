package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Score struct {
	Name     string  `json:"name"`
	Score    float32 `json:"score"`
	Attempts int     `json:"attempts"`
}

type GetScoresResponse struct {
	Scores []Score `json:"scores"`
}

type Getter interface {
	GetAllScores(class string) ([]interface {
		Name() string
		Score() float32
		Attempts() int
	}, error)
}

func GetScores(db Getter) func(c echo.Context) error {
	return func(c echo.Context) error {
		class := c.Param("cls")
		s, err := db.GetAllScores(class)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Unable to get scores: %v", err))
		}

		ret := &GetScoresResponse{
			Scores: make([]Score, len(s)),
		}
		for i, v := range s {
			ret.Scores[i].Name = v.Name()
			ret.Scores[i].Score = v.Score()
			ret.Scores[i].Attempts = v.Attempts()
		}

		return c.JSON(http.StatusOK, ret)
	}
}
