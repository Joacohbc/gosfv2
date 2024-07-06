package dtos

import "time"

type UserDTO struct {
	ID       uint    `json:"id"`
	Username *string `json:"username"`
}

// FileDTO: Representa un archivo en el sistema pero sin la información de su Owner
type FileDTO struct {
	ID         uint       `json:"id"`
	OwnerID    uint       `json:"owner_id"`
	Filename   *string    `json:"filename"`
	Shared     *bool      `json:"shared,omitempty"`
	SharedWith []UserDTO  `json:"sharedWith,omitempty"`
	IsDir      *bool      `json:"isDir,omitempty"`
	ParentId   *uint      `json:"parentId,omitempty"`
	Children   []FileDTO  `json:"children,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}
