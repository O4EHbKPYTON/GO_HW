package controllers

import (
	"api/models"
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	_ "strings"
)

// Operations about plants
type PlantController struct {
	beego.Controller
}

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func (p *PlantController) checkAuth() bool {
	auth := AuthController{Controller: p.Controller}
	if auth.CheckAuth() {
		p.Data["json"] = Respons{Err: true, Data: "Unauthorized"}
		p.Ctx.Output.SetStatus(401)
		p.ServeJSON()
		return false
	}
	return true
}

// @Title CreatePlant
// @Description create plants
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param	body		body 	models.Plant	true		"body for plant content"
// @Success 200 {int} models.Plant.Id
// @Failure 403 body is empty
// @router / [post]
func (p *PlantController) Post() {
	if !p.checkAuth() {
		return
	}

	var plant models.Plant
	if err := json.Unmarshal(p.Ctx.Input.RequestBody, &plant); err != nil {
		p.Data["json"] = Respons{Err: true, Data: "Invalid JSON format"}
		p.ServeJSON()
		return
	}

	pid, err := models.AddPlant(&plant)
	if err != nil {
		p.Data["json"] = Respons{Err: true, Data: err.Error()}
	} else {
		p.Data["json"] = Respons{Err: false, Data: pid}
	}
	p.ServeJSON()
}

// @Title GetAll
// @Description get all plants
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} models.Plant
// @router / [get]
func (p *PlantController) GetAll() {
	if !p.checkAuth() {
		return
	}
	plants := models.GetAllPlants()
	p.Data["json"] = Respons{Err: false, Data: plants}
	p.ServeJSON()
}

// @Title Get
// @Description get plant by pid
// @Security BearerAuth
// @Param	pid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Plant
// @Failure 403 :pid is empty
// @router /:pid [get]
func (p *PlantController) Get() {
	if !p.checkAuth() {
		return
	}
	pid, err := p.GetInt(":pid")
	if err == nil {
		plant, err := models.GetPlant(pid)
		if err != nil {
			p.Data["json"] = Respons{Err: true, Data: err.Error()}
		} else {
			p.Data["json"] = Respons{Err: false, Data: plant}
		}
	}
	p.ServeJSON()
}

// @Title Update
// @Description update the plant
// @Security BearerAuth
// @Param	body		body 	models.Plant	true		"body for plant content"
// @Success 200 {object} models.Plant
// @Failure 403 body is empty
// @router / [put]
func (p *PlantController) Put() {
	if !p.checkAuth() {
		return
	}
	var plant models.Plant
	json.Unmarshal(p.Ctx.Input.RequestBody, &plant)
	err := models.UpdatePlant(&plant)
	if err != nil {
		p.Data["json"] = Respons{Err: true, Data: err.Error()}
	} else {
		p.Data["json"] = Respons{Err: false, Data: plant}
	}
	p.ServeJSON()
}

// @Title Delete
// @Description delete the plant
// @Security BearerAuth
// @Param	pid		path 	string	true		"The pid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 pid is empty
// @router /:pid [delete]
func (p *PlantController) Delete() {
	if !p.checkAuth() {
		return
	}
	pid, err := p.GetInt(":pid")
	if err == nil {
		del := models.DeletePlant(pid)
		if del {
			p.Data["json"] = Respons{Err: false, Data: "Растение удалено"}
		} else {
			p.Data["json"] = Respons{Err: true, Data: "Растение  не найдено"}
		}
	}
	p.ServeJSON()
}
