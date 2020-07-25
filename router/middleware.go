package router

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/dgrijalva/jwt-go"

	"github.com/nagatea/oiscon-portal/model"
)

func middlewareAuthUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := c.Cookie("oiscon_session")
		if err != nil {
			return c.String(http.StatusUnauthorized, "認証に失敗しました")
		}
		token, err := VerifyToken(sess.Value)
		if err != nil || !token.Valid {
			return c.String(http.StatusUnauthorized, "認証に失敗しました")
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		userName := claims["user"].(string)
		user, err := model.GetUserByName(userName)
		if err != nil {
			return c.String(http.StatusUnauthorized, "認証に失敗しました")
		}
		c.Set("user", user)
		return next(c)
	}
}

// VerifyToken jwtの検証をする
func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return token, err
	}
	return token, nil
}
