package routers

import (
	"api/controllers"
	"api/middleware"
	"github.com/beego/beego/v2/server/web"
)

func init() {
	// Логин
	web.Router("/login", &controllers.AuthController{}, "post:Login")

	// API v1 с защитой JWT
	nameSpace := web.NewNamespace("/v1",
		web.NSBefore(middleware.JWTMiddleware),
		web.NSNamespace("/user",
			web.NSInclude(
				&controllers.UserController{},
			),
		),
	)

	web.AddNamespace(nameSpace)
}
