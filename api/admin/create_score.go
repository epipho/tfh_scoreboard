package admin

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CreateScoreRequest struct {
	UpdateScoreRequest
	Class string  `json:"class"`
	Email *string `json:"email",omitempty`
}

type CreateScoreResponse struct {
	ID string `json:"id"`
}

func CreateScore(s Scorer) func(c echo.Context) error {
	return func(c echo.Context) error {
		r := new(CreateScoreRequest)
		if err := c.Bind(r); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		id, err := s.Create(r.Name, r.Email, r.Class)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Unable to create score: %v", err))
		}

		if r.Score != nil {
			if err := s.Update(id, *r.Score); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Unable to update score: %v", err))
			}
			if err := s.Finalize(id, r.Replace); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Unable to finalize score: %v", err))
			}
		}

		resp := &CreateScoreResponse{ID: id}
		return c.JSON(http.StatusCreated, resp)
	}
}
