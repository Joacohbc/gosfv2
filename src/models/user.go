package models

import (
	"context"
	"errors"
	"gosfV2/src/models/db"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique; not null; type: varchar(30)"`
	Password string `json:"password" gorm:"not null; type: longtext"`
}

type UserDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

func init() {
	if err := db.GetBd().AutoMigrate(&User{}); err != nil {
		log.Fatal("Error to create User table:", err)
	}
}

var (
	Users           users
	ErrUserNotFound = errors.New("user/s not found")
)

type users struct{}

// Inscrita el password con AES y retorna la cadena encriptada
func generatePassowrd(password *string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	*password = string(hash)
	return nil
}

func checkPassword(passoword, bdHash string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(bdHash), []byte(passoword)); err != nil {
		return false, err
	}
	return true, nil
}

func (u users) GetAllUsers(ctx context.Context) ([]User, error) {
	var users []User
	if err := db.GetBdCtx(ctx).Find(&users).Error; err != nil {
		if db.IsNotFound(err) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return users, nil
}

// Crea un nuevo usuario
func (u users) NewUser(ctx context.Context, user User) error {
	if err := generatePassowrd(&user.Password); err != nil {
		return err
	}
	return db.GetBd().Create(&user).Error
}

// Devuelve el primer usuario que encuentre con el nombre de usuario y la contraseña
// Si no encuentra ninguno, no devuelve error (solo retorna false)
func (u users) ExistUser(ctx context.Context, username string, password string) (bool, error) {

	user, err := u.FindUserByName(ctx, username)
	if err != nil {
		if err == ErrUserNotFound {
			return false, nil
		}
		return false, err
	}

	ok, err := checkPassword(password, user.Password)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, nil
	}

	return true, nil
}

// Devuelve el usuario con el nombre de usuario
// Si no encuentra ninguno, devuelve ErrUserNotFound
func (u users) FindUserByName(ctx context.Context, username string) (User, error) {
	bd := db.GetBdCtx(ctx)

	var user User
	if err := bd.Where(User{Username: username}).First(&user).Error; err != nil {
		if db.IsNotFound(err) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}

	return user, nil
}

// Devuelve el usuario con el nombre de usuario
// Si no encuentra ninguno, no devuelve error (solo retorna false)
func (u users) ExistUserByName(ctx context.Context, username string) (bool, error) {
	_, err := u.FindUserByName(ctx, username)
	if err != nil {
		if err == ErrUserNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
