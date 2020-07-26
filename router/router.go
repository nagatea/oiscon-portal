package router

import (
	"net/http"

	"github.com/labstack/echo"
)

// SetupRouting APIのルーティングを行います
func SetupRouting(e *echo.Echo) {
	e.GET("/api/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	api := e.Group("/api")
	{
		apiAuth := api.Group("/auth")
		{
			apiAuth.GET("/github", AuthGitHub)
			apiAuth.GET("/github/callback", AuthGitHubCallback)
			apiAuth.GET("/logout", RevokeSession)
		}
	}

	apiWithAuth := e.Group("/api", middlewareAuthUser)
	{
		apiUsers := apiWithAuth.Group("/users")
		{
			apiUsers.GET("/me", GetUsersMe)
		}
	}
} 
