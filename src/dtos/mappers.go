package dtos

import "gosfV2/src/models"

func ToUserDTO(user models.User) UserDTO {
	return UserDTO{
		ID:       user.ID,
		Username: user.Username,
	}
}

func ToUserListDTO(users []models.User) []UserDTO {
	var usersDTO []UserDTO
	for _, user := range users {
		usersDTO = append(usersDTO, ToUserDTO(user))
	}
	return usersDTO
}

func ToFileDTO(file models.File) FileDTO {
	return FileDTO{
		ID:         file.ID,
		Filename:   file.Filename,
		Shared:     file.Shared,
		SharedWith: ToUserListDTO(file.SharedWith),
	}
}

func ToFileListDTO(files []models.File) []FileDTO {
	var dtos []FileDTO
	for _, file := range files {
		dtos = append(dtos, ToFileDTO(file))
	}
	return dtos
}
