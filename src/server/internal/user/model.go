package user

import "errors"

var (
	ErrNotFound           = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PublicUser struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u User) Public() PublicUser {
	return PublicUser{ID: u.ID, Name: u.Name, Email: u.Email}
}

type Store interface {
	Create(name, email, password string) (PublicUser, error)
	Authenticate(email, password string) (PublicUser, error)
	List() ([]PublicUser, error)
	Delete(id int) error
}
