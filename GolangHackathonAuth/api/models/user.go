package models

import (
	"errors"
	"fmt"
	_ "fmt"
	_ "strconv"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	orm.RegisterModel(new(User))
}

// Структура пользователя
type User struct {
	Id       int64
	Username string
	Password string
	Email    string
}

// Секретный ключ для JWT
var SecretKey = []byte("your-secret-key")

// Хеширование пароля
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Проверка пароля
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Создание JWT-токена
func CreateToken(u User) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   u.Id,
		"name": u.Username,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // Токен действует 24 часа
	})

	tokenString, err := claims.SignedString(SecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Проверка JWT-токена
func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token.Claims.(jwt.MapClaims), nil
}

// Добавление пользователя
func AddUser(u User) (int64, error) {
	o := orm.NewOrmUsingDB("mydatabase")

	// Хешируем пароль перед сохранением
	hashedPass, err := HashPassword(u.Password)
	if err != nil {
		return 0, err
	}
	u.Password = hashedPass

	id, err := o.Insert(&u)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Получение пользователя по ID
func GetUser(uid int64) (*User, error) {
	o := orm.NewOrmUsingDB("mydatabase")
	user := User{Id: uid}
	err := o.Read(&user)
	if err == orm.ErrNoRows {
		return nil, errors.New("пользователь не найден")
	}
	return &user, err
}

// Получение всех пользователей
func GetAllUsers() ([]User, error) {
	var users []User
	o := orm.NewOrmUsingDB("mydatabase")
	_, err := o.QueryTable("user").All(&users)
	return users, err
}

// Обновление пользователя
func UpdateUser(u *User) error {
	o := orm.NewOrmUsingDB("mydatabase")

	// Проверяем, существует ли пользователь
	existingUser := User{Id: u.Id}
	if err := o.Read(&existingUser); err != nil {
		return errors.New("пользователь не найден")
	}

	// Если передан новый пароль, хешируем его
	if u.Password != "" {
		hashedPass, err := HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashedPass
	}

	_, err := o.Update(u, "Username", "Password", "Email")
	return err
}

// Авторизация пользователя
func Login(username, password string) (string, error) {
	o := orm.NewOrmUsingDB("mydatabase")
	var user User
	err := o.QueryTable("user").Filter("Username", username).One(&user)

	if err == orm.ErrNoRows {
		return "", errors.New("неверный логин или пароль")
	}

	// Проверяем пароль
	if !CheckPasswordHash(password, user.Password) {
		return "", errors.New("неверный логин или пароль")
	}

	// Генерируем токен
	token, err := CreateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func DeleteUser(uid int64) error {
	o := orm.NewOrmUsingDB("mydatabase")
	_, err := o.Delete(&User{Id: uid})
	if err != nil {
		return fmt.Errorf("ошибка при удалении пользователя: %w", err)
	}
	return nil
}
