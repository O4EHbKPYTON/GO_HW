package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/golang-jwt/jwt/v5"
)

type Session struct {
	ID        int       `orm:"auto"`
	UserID    int       `orm:"index"`
	Token     string    `orm:"size(255)"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(User))
	orm.RegisterModel(new(Session))
}

type PostUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PostUserResponse struct {
	Id int64 `json:"id"`
}

type User struct {
	Id       int64  `orm:"auto;column(id)"`
	Username string `orm:"column(username)"`
	Email    string `orm:"column(email)"`
	Password string `orm:"column(password_hash)"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Создание секретного ключа
var SecretKey = []byte("your-secret-key")

// лучше timestamp использовать
// в сессию добавить data update data start
func CreateToken(u User) (string, error) {
	// создаем заявку
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   u.Id,
		"name": u.Username,
	})
	// генерируем токен
	tokenString, err := claims.SignedString(SecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	// Проверка на ошибки
	if err != nil {
		return nil, err
	}

	// Проверка валидности токена
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Возврат данных токена
	return token.Claims.(jwt.MapClaims), nil
}

func AddUser(u User) int64 {
	o := orm.NewOrmUsingDB("mydatabase")
	id, err := o.Insert(&u)
	if err != nil {
		fmt.Println(err)
	}
	return id
}

func GetUser(uid int64) (u *User, err error) {
	o := orm.NewOrmUsingDB("mydatabase")
	user := User{Id: uid}
	err = o.Read(&user)
	if err == orm.ErrNoRows {
		return nil, errors.New("user with this id not found")
	}
	return &user, nil
}

func GetAllUsers() *[]User {
	var users []User
	o := orm.NewOrmUsingDB("mydatabase")
	qb, _ := orm.NewQueryBuilder("postgres")
	qb.Select("id", "username", "password_hash").From("user").Where("id > ?").OrderBy("id").Desc().Limit(10)
	o.Raw(qb.String(), 0).QueryRows(&users)
	return &users
}

func UpdateUser(uu *User) (err error) {
	o := orm.NewOrmUsingDB("mydatabase")
	_, err = o.Update(uu, "username", "password_hash")
	if err != nil {
		return errors.New("user not found")
	}
	return nil
}

func Login(req LoginRequest) (string, error) {
	o := orm.NewOrmUsingDB("mydatabase")

	var user User
	err := o.QueryTable("user").Filter("username", req.Username).One(&user)
	if err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return "", errors.New("invalid username or password")
		}
		return "", err
	}

	// Сравниваем пароль
	if req.Password != user.Password {
		return "", errors.New("invalid username or password")
	}

	// Генерируем токен, если логин и пароль верны
	return CreateToken(user)
}

func DeleteUser(uid int64) bool {
	o := orm.NewOrmUsingDB("mydatabase")
	user := User{Id: uid}
	_, err := o.Delete(&user)
	return err == nil
}
