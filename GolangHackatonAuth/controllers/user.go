package controllers

import (
	"api/models"
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	"strings"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

type Respons struct {
	Err  bool `json:"err"`
	Data any  `json:"data"`
}

func (u *UserController) SessionTest() {
	err := u.SetSession("test_key", "test_value")
	if err != nil {
		u.Data["json"] = Respons{Err: true, Data: err.Error()}
	} else {
		u.Data["json"] = Respons{Err: false, Data: "Session saved"}
	}
	u.ServeJSON()
}

// Проверка заголовков для всех запросов
func (u *UserController) HandlerFunc(rules string) bool {
	switch rules {
	case "GetAll", "Logout":
		// Проверка сессии
		accessToken, ok := u.GetSession("accessToken").(string) //27 строка
		if !ok || accessToken == "" {
			u.Abort("401")
			return true
		}

		// Проверка заголовка
		authHeader := u.Ctx.Input.Header("Authorization")
		if authHeader == "" {
			u.Abort("401")
			return true
		}

		// Извлечение токена
		accessTokenFromHeader := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			accessTokenFromHeader = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			accessTokenFromHeader = authHeader
		}

		// Верификация токена
		claims, err := models.VerifyToken(accessTokenFromHeader)
		if err != nil || claims == nil {
			u.Abort("401")
			return true
		}

		// Сравнение токенов
		if accessTokenFromHeader != accessToken {
			u.Abort("401")
			return true
		}

		return false
	default:
		return false
	}
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router / [post]
func (u *UserController) Post() {
	var user models.User
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	uid := models.AddUser(user)
	u.Data["json"] = Respons{Err: false, Data: uid}
	u.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} models.User
// @router / [get]
func (u *UserController) GetAll() {
	users := models.GetAllUsers()
	u.Data["json"] = Respons{Err: false, Data: users}
	u.ServeJSON()
}

//Сначала нужно ввести токен из Login()

// @Title Get
// @Description get user by uid
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router /:uid [get]
func (u *UserController) Get() {
	uid, err := u.GetInt64(":uid")
	if err == nil {
		user, err := models.GetUser(uid)
		if err != nil {
			u.Data["json"] = Respons{Err: true, Data: err.Error()}
		} else {
			u.Data["json"] = Respons{Err: false, Data: user}
		}
	}
	u.ServeJSON()
}

// @Title Update
// @Description update the user
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 403 body is empty
// @router / [put]
func (u *UserController) Put() {
	var user models.User
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	err := models.UpdateUser(&user)
	if err != nil {
		u.Data["json"] = Respons{Err: true, Data: err.Error()}
	} else {
		u.Data["json"] = Respons{Err: false, Data: user}
	}
	u.ServeJSON()
}

// @Title Delete
// @Description delete the user
// @Param	uid		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /:uid [delete]
func (u *UserController) Delete() {
	uid, err := u.GetInt64(":uid")
	if err == nil {
		del := models.DeleteUser(uid)
		if del {
			u.Data["json"] = Respons{Err: false, Data: "Пользователь удален"}
		} else {
			u.Data["json"] = Respons{Err: true, Data: "Пользователь не найден"}
		}
	}
	u.ServeJSON()
}

// @Title Login
// @Description Logs user into the system
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router /login [post]
func (u *UserController) Login() {
	var user models.User
	if err := json.Unmarshal(u.Ctx.Input.RequestBody, &user); err != nil {
		u.Data["json"] = Respons{Err: true, Data: "Invalid request"}
		u.ServeJSON()
		return
	}

	token, err := models.Login(user)
	if err != nil {
		u.Data["json"] = Respons{Err: true, Data: err.Error()}
	} else {
		u.SetSession("accessToken", token)
		u.Data["json"] = Respons{Err: false, Data: token}
	}
	u.ServeJSON()
}

// @Title logout
// @Description Logs out current logged in user session
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} logout success
// @router /logout [get]
func (u *UserController) Logout() {
	// получение значения сессии
	//u.GetSession("accessToken")
	// удаление значения сессии
	// u.DelSession("accessToken")
	// уничтожение сессии
	u.DestroySession()
	u.Data["json"] = Respons{Err: false, Data: "Вышли из сессии"}
	u.ServeJSON()
}
