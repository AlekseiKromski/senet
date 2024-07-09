package storage

type User struct {
	Id          string       `json:"id"`
	Username    string       `json:"username"`
	Email       string       `json:"email,omitempty"`
	First_name  string       `json:"first_name,omitempty"`
	Second_name string       `json:"second_name,omitempty"`
	Password    string       `json:"password,omitempty"`
	Image       string       `json:"image"`
	Role        string       `json:"role,omitempty"`
	RoleName    *string      `json:"role_name,omitempty"`
	CreateAt    string       `json:"created_at"`
	UpdatedAt   string       `json:"updated_at"`
	DeletedAt   *string      `json:"deleted_at"`
	Endpoints   *[]*Endpoint `json:"endpoints,omitempty"`
}

type StorageUser interface {
	CreateUser(username, email, first_name, second_name, password, role string) error
	DeleteUser(id string) error
	UpdateUser(id, username, email, first_name, second_name, password, role string) error
	GetUserById(id string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	FindUsersByUsername(username string) ([]*User, error)
}
