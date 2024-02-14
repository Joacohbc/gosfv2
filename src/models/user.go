package models

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"gosfV2/src/ent"
	"gosfV2/src/ent/user"
	"gosfV2/src/models/database"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

func Validate(u *ent.User) error {

	if len(u.Username) > 255 {
		return errors.New("invalid username, the name is too long (max 255)")
	}

	matched, err := regexp.MatchString("[a-zA-Z0-9]+", u.Username)
	if !matched || err != nil {
		return errors.New("invalid username, the name must be alphanumeric")
	}

	matched, err = regexp.MatchString(".*[a-z]+.*.*[A-Z]+.*.*[0-9]+.*.*[#!?{}!?*^%@#$]+.*", u.Password)
	if !matched || err != nil {
		return errors.New(`invalid password, the password must have at least one lowercase, one uppercase, one number and one special character ("#", "!", "?", "{", "}", "!", "?", "*", "^", "%", "@", "#", "$")`)
	}

	return nil
}

var (
	ErrUserNotFound           = errors.New("user/s not found")
	ErrIconTooLarge           = errors.New("icon too large (max 512 x 512)")
	ErrIconFormatNotSupported = errors.New("invalid image format, only jpeg/png/gif and png is supported")
)

type UserInterface interface {
	GetAllUsers() ([]*ent.User, error)
	NewUser(user *ent.User) error
	FindUserByName(username string) (*ent.User, error)
	FindUserById(id uint) (*ent.User, error)
	ExistUserByName(username string) (bool, error)
	Rename(id uint, newName string) error
	ChangePassword(id uint, newPassword string) error
	Delete(id uint) error
	UploadIcon(id uint, src io.Reader) error
	DeleteIcon(id uint) error
	GetIcon(id uint) string
	ManageError(err error) error
}

type userService struct {
	Client  *ent.Client
	Context context.Context
}

func Users() UserInterface {
	return userService{Client: database.GetMySQL(), Context: context.Background()}
}

func (u userService) GetAllUsers() ([]*ent.User, error) {
	users, err := u.Client.User.Query().All(u.Context)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return users, nil
}

func (u userService) ManageError(err error) error {
	if ent.IsNotFound(err) {
		return ErrUserNotFound
	}
	return err
}

// Create a new user
func (u userService) NewUser(user *ent.User) error {
	_, err := u.Client.User.Create().SetUsername(user.Username).SetPassword(user.Password).Save(u.Context)
	return err
}

// Find the user by username
// If not found, return ErrUserNotFound
func (u userService) FindUserByName(username string) (*ent.User, error) {
	user, err := u.Client.User.Query().Where(user.UsernameEQ(username)).Only(u.Context)
	if err != nil {
		return nil, u.ManageError(err)
	}

	return user, nil
}

// Find the user by ID
// If not found, return ErrUserNotFound
func (u userService) FindUserById(id uint) (*ent.User, error) {
	user, err := u.Client.User.Get(u.Context, id)
	if err != nil {
		return nil, u.ManageError(err)
	}

	return user, nil
}

// Check if a user with the given username exists
func (u userService) ExistUserByName(username string) (bool, error) {
	_, err := u.FindUserByName(username)
	if err != nil {
		if err == ErrUserNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (u userService) Rename(id uint, newName string) error {
	_, err := u.Client.User.UpdateOneID(id).SetUsername(newName).Save(u.Context)
	return err
}

func (u userService) ChangePassword(id uint, newPassword string) error {
	_, err := u.Client.User.UpdateOneID(id).SetPassword(newPassword).Save(u.Context)
	return err
}

func (u userService) Delete(id uint) error {
	err := u.Client.User.DeleteOneID(id).Exec(u.Context)
	return err
}

func (u userService) UploadIcon(id uint, src io.Reader) error {
	// Read the file (to be able to use the content of the Reader multiple times)
	blob, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	img, _, err := image.DecodeConfig(bytes.NewReader(blob))
	if err != nil {
		return ErrIconFormatNotSupported
	}

	if img.Width > 512 || img.Height > 512 {
		return ErrIconTooLarge
	}

	// Save the file
	file, err := os.Create(u.GetIcon(id))
	if err != nil {
		return err
	}

	// Close the file
	if _, err = io.Copy(file, bytes.NewReader(blob)); err != nil {
		return err
	}

	if err = file.Close(); err != nil {
		return err
	}

	return nil
}

func (u userService) GetIcon(id uint) string {
	return filepath.Join(UserIconDir, fmt.Sprint(id)+"icon")
}

func (u userService) DeleteIcon(id uint) error {
	err := os.Remove(u.GetIcon(id))
	if err != nil {
		return err
	}

	return nil
}
