package controllers

//import (
//	"api/models"
//	"github.com/beego/beego/v2/server/web"
//	"github.com/golang-jwt/jwt/v5"
//	"time"
//)
//
//type AuthController struct {
//	web.Controller
//}
//
//type LoginRequest struct {
//	Username string `json:"username"`
//	Password string `json:"password"`
//}
//
//func (c *AuthController) Login() {
//	var req LoginRequest
//	if err := c.BindJSON(&req); err != nil {
//		c.Ctx.Output.SetStatus(400)
//		_ = c.Ctx.Output.Body([]byte("Invalid request"))
//		return
//	}
//
//	user, err := models.Login(req.Username, req.Password)
//	if err != nil {
//		c.Ctx.Output.SetStatus(401)
//		_ = c.Ctx.Output.Body([]byte("Invalid credentials"))
//		return
//	}
//
//	claims := jwt.MapClaims{
//		"user_id":  user.Id,
//		"username": user.Username,
//		"exp":      time.Now().Add(2 * time.Hour).Unix(),
//	}
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	tokenString, err := token.SignedString(secretKey)
//	if err != nil {
//		c.Ctx.Output.SetStatus(500)
//		_ = c.Ctx.Output.Body([]byte("Error generating token"))
//		return
//	}
//
//	c.Data["json"] = map[string]string{
//		"token":   tokenString,
//		"expires": time.Unix(claims["exp"].(int64), 0).Format(time.RFC3339),
//	}
//	_ = c.ServeJSON()
//}
