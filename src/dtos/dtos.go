package dtos

type UserDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

// FileDTO: Representa un archivo en el sistema pero sin la información de su Owner
type FileDTO struct {
	ID         uint      `json:"id"`
	OwnerID    uint      `json:"owner_id,omitempty"`
	Filename   string    `json:"filename"`
	Shared     bool      `json:"shared"`
	SharedWith []UserDTO `json:"shared_with,omitempty"`
}
