package admin

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UpdateScoreRequest struct {
	Name string `json:"name"`
	// if score is set, the score will be immediatly updated and finalized
	Score       *float32 `json:"score",omitempty`
	Incremental bool     `json:"incremental"` //incrementall add to score
	Replace     bool     `json:"replace"`     // replace current score (even if not higher)
	Finalize    bool     `json:"finalize"`    // last score update, write to db
}

func UpdateScore(s Scorer) func(c echo.Context) error {
	return func(c echo.Context) error {
		r := new(UpdateScoreRequest)
		if err := c.Bind(r); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		id := c.Param("id")
		if r.Score != nil {
			if err := s.Update(id, *r.Score, r.Incremental); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Unable to update score: %v", err))
			}
		}

		if r.Finalize {
			if err := s.Finalize(id, r.Replace); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Unable to finalize score: %v", err))
			}
		}

		return c.NoContent(http.StatusAccepted)
	}
}
