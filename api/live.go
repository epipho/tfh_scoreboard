package api

import (
	"github.com/labstack/echo/v4"
)

func Live() func(c echo.Context) error {
	return func(c echo.Context) error {
		return nil
	}
}
