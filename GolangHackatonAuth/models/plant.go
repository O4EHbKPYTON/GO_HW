// models/plant.go
package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

func init() {
	orm.RegisterModel(new(Plant), new(UserPlant))
}

type Plant struct {
	Id               int    `orm:"column(id);auto"`
	Name             string `orm:"column(name);size(255)"`
	ScientificName   string `orm:"column(scientific_name);size(255);null"`
	WaterFrequency   int    `orm:"column(water_frequency)"`
	LightRequirement string `orm:"column(light_requirement);size(50)"`
	TemperatureRange string `orm:"column(temperature_range);size(50)"`
	HumidityRange    string `orm:"column(humidity_range);size(50)"`
	Description      string `orm:"column(description);type(text)"`
}

type UserPlant struct {
	Id           int       `orm:"auto;column(id)"`
	User         *User     `orm:"rel(fk);column(user_id)"`
	Plant        *Plant    `orm:"rel(fk);column(plant_id)"`
	AcquiredDate time.Time `orm:"column(acquired_date);type(date)"`
	LastWatered  time.Time `orm:"column(last_watered);type(date)"`
	Location     string    `orm:"column(location);size(255)"`
}

func AddPlant(p *Plant) (int64, error) {
	o := orm.NewOrmUsingDB("mydatabase")
	id, err := o.Insert(p)
	if err != nil {
		return 0, fmt.Errorf("error inserting plant: %v", err)
	}
	return id, nil
}

func GetPlant(pid int) (*Plant, error) {
	o := orm.NewOrmUsingDB("mydatabase")
	plant := Plant{Id: pid}
	err := o.Read(&plant)
	if err == orm.ErrNoRows {
		return nil, errors.New("Plant not found")
	}
	return &plant, nil
}

func GetAllPlants() []Plant {
	var plants []Plant
	o := orm.NewOrmUsingDB("mydatabase")
	qb, _ := orm.NewQueryBuilder("postgres")
	qb.Select("*").From("plants")
	_, err := o.Raw(qb.String()).QueryRows(&plants)
	if err != nil {
		fmt.Println("Error getting plants:", err)
	}
	return plants
}

func UpdatePlant(p *Plant) error {
	o := orm.NewOrmUsingDB("mydatabase")
	_, err := o.Update(p)
	if err != nil {
		return errors.New("Error updating plant")
	}
	return nil
}

func DeletePlant(pid int) bool {
	o := orm.NewOrmUsingDB("mydatabase")
	_, err := o.Delete(&Plant{Id: pid})
	return err == nil
}

func AddUserPlant(up *UserPlant) int {
	o := orm.NewOrmUsingDB("mydatabase")
	id, err := o.Insert(up)
	if err != nil {
		fmt.Println("Error adding user plant:", err)
	}
	return int(id)
}

func GetUserPlant(upid int) (*UserPlant, error) {
	o := orm.NewOrmUsingDB("mydatabase")
	userPlant := UserPlant{Id: int(upid)}
	err := o.Read(&userPlant)
	if err == orm.ErrNoRows {
		return nil, errors.New("User's plant not found")
	}
	return &userPlant, nil
}

func GetUserPlants(uid int) []UserPlant {
	var userPlants []UserPlant
	o := orm.NewOrmUsingDB("mydatabase")
	qb, _ := orm.NewQueryBuilder("postgres")
	qb.Select("*").From("user_plants").Where("user_id = ?")
	_, err := o.Raw(qb.String(), uid).QueryRows(&userPlants)
	if err != nil {
		fmt.Println("Error getting user plants:", err)
	}
	return userPlants
}

func UpdateUserPlant(up *UserPlant) error {
	o := orm.NewOrmUsingDB("mydatabase")
	_, err := o.Update(up)
	if err != nil {
		return errors.New("Ошибка обновления растения пользователя")
	}
	return nil
}

func DeleteUserPlant(upid int) bool {
	o := orm.NewOrmUsingDB("mydatabase")
	_, err := o.Delete(&UserPlant{Id: int(upid)})
	return err == nil
}
