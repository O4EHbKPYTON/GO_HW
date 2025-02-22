package middleware

//
//import (
//	"github.com/beego/beego/v2/server/web/context"
//	"github.com/golang-jwt/jwt/v5"
//	"strings"
//)
//
//var secretKey = []byte("your_secret_key")
//
//func JWTMiddleware(ctx *context.Context) {
//	authHeader := ctx.Input.Header("Authorization")
//	if authHeader == "" {
//		ctx.Output.SetStatus(401)
//		_ = ctx.Output.Body([]byte("Missing authorization header"))
//		return
//	}
//
//	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
//	if tokenString == authHeader {
//		ctx.Output.SetStatus(401)
//		_ = ctx.Output.Body([]byte("Invalid token format"))
//		return
//	}
//
//	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, jwt.ErrSignatureInvalid
//		}
//		return secretKey, nil
//	})
//
//	if err != nil || !token.Valid {
//		ctx.Output.SetStatus(401)
//		_ = ctx.Output.Body([]byte("Invalid token"))
//		return
//	}
//
//	// Добавляем claims в контекст
//	if claims, ok := token.Claims.(jwt.MapClaims); ok {
//		ctx.Input.SetData("userID", claims["user_id"])
//		ctx.Input.SetData("username", claims["username"])
//	}
//}
