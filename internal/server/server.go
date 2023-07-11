package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {

		response := Response{
			Message: "Hello, World!",
		}

		return c.JSON(http.StatusOK, response)
	})

	e.POST("/hello", func(c echo.Context) error {
		json_map := make(map[string]interface{})
		err := json.NewDecoder(c.Request().Body).Decode(&json_map)
		if err != nil {
			return c.String(http.StatusUnprocessableEntity, err.Error())
		}

		name, ok := json_map["name"]
		if !ok {
			name = "noname"
		}
		response := Response{
			Message: fmt.Sprintf("Hello, %s!", name),
		}

		return c.JSON(http.StatusOK, response)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
