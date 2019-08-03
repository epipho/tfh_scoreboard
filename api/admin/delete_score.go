package admin

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func DeleteScore(s Scorer) func(c echo.Context) error {
	return func(c echo.Context) error {
		r := new(UpdateScoreRequest)
		if err := c.Bind(r); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		id := c.Param("id")
		s.Cancel(id)

		return c.NoContent(http.StatusAccepted)
	}
}
