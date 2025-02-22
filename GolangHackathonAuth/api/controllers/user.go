package controllers

import (
	"api/models"
	"encoding/json"
	"fmt"
	"slices"

	beego "github.com/beego/beego/v2/server/web"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

// Единая структура ответа на запросы
type Respons struct {
	Err  bool `json:"err"`
	Data any  `json:"data"`
}

// Проверка заголовков для всех запросов
func (u *UserController) HandlerFunc(rules string) bool {
	fmt.Println(u.Ctx.Request.Header["Authorization"])
	switch rules {
	case "GetAll", "Logout": // rules - в значении имеет название функции, которая выполняется при вызове метода
		if u.GetSession("accessToken") == nil { // GetSession возвращает nil если нет ключа, поэтому приходится проверять
			break
		}
		accessToken := u.GetSession("accessToken").(string)
		arrayToken := u.Ctx.Request.Header["Authorization"]
		token, _ := models.VerifyToken(accessToken)                                    // забираем данные из токена
		fmt.Println(token, token["id"])                                                // Можем делать проверку
		if len(arrayToken) > 0 && slices.Contains(arrayToken, "Bearer "+accessToken) { //проверка токена, тут проверка через сессию
			return false
		}
	default: //все не указанные методы будут выполняться без авторизации
		return false
	}
	u.Abort("401") // выдаем ошибку авторизации
	return true
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
	uid, _ := models.AddUser(user)
	u.Data["json"] = Respons{Err: false, Data: uid}
	u.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} models.User
// @router / [get]
func (u *UserController) GetAll() {
	users, _ := models.GetAllUsers()
	u.Data["json"] = Respons{Err: false, Data: users}
	u.ServeJSON()
}

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
		err := models.DeleteUser(uid)
		if err == nil {
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
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	token, _ := models.Login(user.Username, user.Password)
	u.Data["json"] = Respons{Err: false, Data: token}
	// установка значения сессии
	u.SetSession("accessToken", token)
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
