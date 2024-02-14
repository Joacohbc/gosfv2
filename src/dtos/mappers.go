package dtos

import (
	"gosfV2/src/ent"
)

func ToUserDTO(user *ent.User) UserDTO {
	return UserDTO{
		ID:       user.ID,
		Username: &user.Username,
	}
}

func ToUserListDTO(users []*ent.User) []UserDTO {
	var usersDTO []UserDTO
	for _, user := range users {
		usersDTO = append(usersDTO, ToUserDTO(user))
	}
	return usersDTO
}

func ToFileDTO(file *ent.File) FileDTO {
	return FileDTO{
		ID:         file.ID,
		Filename:   &file.Filename,
		Shared:     &file.IsShared,
		SharedWith: ToUserListDTO(file.Edges.SharedWith),
	}
}

func ToFileListDTO(files []*ent.File) []FileDTO {
	var dtos []FileDTO
	for _, file := range files {
		dtos = append(dtos, ToFileDTO(file))
	}
	return dtos
}
