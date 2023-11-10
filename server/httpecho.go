package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/msf/cachingproxy/handler"
	"github.com/msf/cachingproxy/model"
)

func EchoPing(c echo.Context) error {
	c.B
	type r struct {
		M string `json:"message"`
	}
	return c.JSON(http.StatusOK, &r{M: "pong"})
}

func EchoMessage(c echo.Context) error {
	m := model.Message{
		ID:      c.Param("id"),
		Content: c.Param("cnt"),
	}
	r, err := handler.EchoMessage(m)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, r)
}
