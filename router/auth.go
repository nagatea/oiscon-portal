package router

import (
	"os"
	"fmt"
	"time"
	"context"
	"net/http"
	"crypto/rand"
	"encoding/base64"

	"github.com/labstack/echo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"github.com/dgrijalva/jwt-go"
	json "github.com/json-iterator/go"

	"github.com/nagatea/oiscon-portal/model"
)

const (
	cookieName       = "oiscon_portal_auth_cookie"
	githubProfileURL = "https://api.github.com/user"
)

var (
	clientID     = os.Getenv("GITHUB_OAUTH2_CLIENT_ID")
	clientSecret = os.Getenv("GITHUB_OAUTH2_CLIENT_SECRET")
	secretKey    = os.Getenv("JWT_SECRET_KEY")
	config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     github.Endpoint,
		RedirectURL:  "http://localhost:3000/api/auth/github/callback",
		Scopes:       []string{"read:user"},
	}
)

// AuthGitHub GitHub認証
func AuthGitHub(c echo.Context) error {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	c.SetCookie(&http.Cookie{
		Name:     cookieName,
		Value:    state,
		Path:     "/",
		HttpOnly: true,
	})

	authURL:= config.AuthCodeURL(state)
	return c.Redirect(http.StatusFound, authURL)
}

// AuthGitHubCallback GitHub認証のcallback
func AuthGitHubCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")
	if len(code) == 0 || len(state) == 0 {
		return c.String(http.StatusBadRequest, "missing code or state")
	}

	cookie, err := c.Cookie(cookieName)
	if err != nil {
		return c.String(http.StatusBadRequest, "missing cookie")
	}
	if cookie.Value != state {
		return c.String(http.StatusBadRequest, "invalid state")
	}

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return c.String(http.StatusBadRequest, "token exchange failed")
	}

	u, err := FetchUserInfo(token)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// userが登録されてなかったら認証失敗
	user, err := model.GetUserByName(u.Name)
	if err != nil {
		return c.String(http.StatusUnauthorized, "ログインに失敗しました(登録されていないユーザーです)")
	}
	// profileImageURL, displayNameをセット
	if (user.DisplayName == "" || user.ProfileImageURL == "") {
		if (user.DisplayName == "") {
			user.DisplayName = u.DisplayName
		}
		if (user.ProfileImageURL == "") {
			user.ProfileImageURL = u.ProfileImageURL
		}
		_, err = model.UpdateUser(user)
		if err != nil {
			c.String(http.StatusInternalServerError, "ユーザー情報の更新に失敗しました")
		}
	}

	jwtToken, err := CreateJWTToken(user.Name)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	c.SetCookie(&http.Cookie{
		Name:     "oiscon_session",
		Value:    jwtToken,
		Path:     "/",
		HttpOnly: true,
	})

	return c.Redirect(http.StatusFound, "/")
}

// FetchUserInfo GitHubからユーザー情報を取得する
func FetchUserInfo(token *oauth2.Token) (model.User, error) {
	var userInfo struct {
		Name    string `json:"name"`
		Login   string `json:"login"`
		Picture string `json:"avatar_url"`
	}

	c := config.Client(context.Background(), token)

	resp, err := c.Get(githubProfileURL)
	if err != nil {
		return model.User{}, err
	}
	defer resp.Body.Close()
	if err := json.ConfigFastest.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return model.User{}, err
	}

	var ui model.User
	ui.Name = userInfo.Login
	ui.DisplayName = userInfo.Name
	ui.ProfileImageURL = userInfo.Picture
	return ui, nil
}

// CreateJWTToken JWTトークンを発行する
func CreateJWTToken(userName string) (string, error) {
	// tokenの作成
	token := jwt.New(jwt.GetSigningMethod("HS256"))

	// claimsの設定
	token.Claims = jwt.MapClaims{
		"user": userName,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	}

	// 署名
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
