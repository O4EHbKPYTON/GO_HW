package controllers

import (
	"strings"

	"api/models"
	"github.com/beego/beego/v2/server/web"
)

type AuthController struct {
	web.Controller
}

// CheckAuth проверяет авторизацию пользователя
func (a *AuthController) CheckAuth() bool {
	authHeader := a.Ctx.Input.Header("Authorization")
	if authHeader == "" {
		a.Abort("401")
		return true
	}

	accessToken := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := models.VerifyToken(accessToken)
	if err != nil || claims == nil {
		a.Abort("401")
		return true
	}

	return false
}
