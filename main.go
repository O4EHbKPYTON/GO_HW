package main

import (
	_ "api/controllers"
	_ "api/routers"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
	"github.com/beego/beego/v2/server/web/session"
	_ "github.com/beego/beego/v2/server/web/session/postgres"
	_ "github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	_ "net/http"
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
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Sec-WebSocket-Key", "Sec-WebSocket-Version", "Connection", "Upgrade"},
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
	beego.BConfig.WebConfig.StaticDir["/static"] = "static"
	beego.BConfig.Log.AccessLogs = true
	beego.BConfig.Log.Outputs["console"] = ""
	beego.Run()
}
