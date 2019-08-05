package api

import (
	"github.com/labstack/echo/v4"
)

type Getter interface {
	GetAllScores(class string) ([]interface {
		Score() float32
		Attempts() int
	}, error)
}

func GetScores(db Getter) func(c echo.Context) error {
	return func(c echo.Context) error {
		return nil
	}
}
