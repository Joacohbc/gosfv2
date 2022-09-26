package users

import (
	"encoding/json"
	"gosfV2/src/models/env"
	"log"
	"os"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetAllUsers() ([]User, error) {
	jsonb, err := os.ReadFile(env.Config.UsersFilePath)
	if err != nil {
		return nil, err
	}

	var users []User
	err = json.Unmarshal(jsonb, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func NewUser(user User) error {

	users, err := GetAllUsers()
	if err != nil {
		return err
	}

	users = append(users, user)
	jsonb, err := json.Marshal(users)
	if err != nil {
		return err
	}

	err = os.WriteFile(env.Config.UsersFilePath, jsonb, 0644)
	if err != nil {
		return err
	}

	return nil
}

func FindUser(username string, password string) (bool, error) {
	users, err := GetAllUsers()
	if err != nil {
		return false, err
	}

	for _, user := range users {
		if user.Username == username && user.Password == password {
			return true, nil
		}
	}

	return false, nil
}

func FindUserByName(username string) (bool, error) {
	users, err := GetAllUsers()
	if err != nil {
		return false, err
	}

	for _, user := range users {
		if user.Username == username {
			return true, nil
		}
	}

	return false, nil
}

func init() {
	if _, err := os.Stat(env.Config.UsersFilePath); err == nil {
		return
	}

	jsonb, err := json.Marshal([]User{})
	if err != nil {
		log.Fatal("Error creating users.json:", err)
	}

	err = os.WriteFile(env.Config.UsersFilePath, jsonb, 0644)
	if err != nil {
		log.Fatal("Error creating users.json:", err)
	}

	log.Println("users.json created")
}
