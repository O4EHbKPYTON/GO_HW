package main

import (
	_ "api/controllers"
	_ "api/routers"
	"github.com/beego/beego/v2/server/web/session"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
	_ "github.com/beego/beego/v2/server/web/session/postgres"
	_ "github.com/lib/pq"
)

func main() {
	sqlconn, _ := beego.AppConfig.String("sqlconn")
	orm.RegisterDataBase("mydatabase", "postgres", sqlconn)
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	//Настройка CORS
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: []string{"*"},                                                //Доступ с любого адреса
		AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"}, //Допустимые запросы
		AllowHeaders: []string{"Origin", "x-requested-with",
			"content-type",
			"accept",
			"origin",
			"authorization",
			"x-csrftoken"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	sessionConfig := &session.ManagerConfig{
		CookieName:      beego.BConfig.WebConfig.Session.SessionName,
		EnableSetCookie: true,
		Gclifetime:      int64(beego.BConfig.WebConfig.Session.SessionGCMaxLifetime),
		ProviderConfig:  beego.BConfig.WebConfig.Session.SessionProviderConfig,
	}
	globalSessions, _ := session.NewManager("postgresql", sessionConfig)
	go globalSessions.GC()

	beego.Run()
}
