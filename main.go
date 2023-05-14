package main

import (
	"net/http"
	
	"github.com/labstack/echo/v4"

	"fmt"

	//"reflect"
)

func main() {
	e := echo.New()

	type User struct {
		Nama  string `json:"name"`
		Email string `json:"email"`
	}

	e.POST("/chat", func(c echo.Context) error {
		u := new(User)

		if err := c.Bind(u); err != nil {
			return err
		}

		//fmt.Println(u)
		fmt.Println(u.Nama)

		return c.JSON(http.StatusCreated, u)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
