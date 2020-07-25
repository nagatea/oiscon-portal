package router

import (
	"net/http"
	"github.com/labstack/echo"

	"github.com/nagatea/oiscon-portal/model"
)

// GetUsersMe GET /users/me
func GetUsersMe(c echo.Context) error {
	user := c.Get("user").(model.User)
	res, err := model.GetUserByName(user.Name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if res.Name == "" {
		res, err = model.CreateUser(user)
	}

	return c.JSON(http.StatusOK, res)
} 
