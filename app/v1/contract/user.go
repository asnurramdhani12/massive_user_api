package contract

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type User struct {
	ID       int    `json:"id" gorm:"primary_key"`
	Username string `json:"username" gorm:"unique"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password,omitempty" gorm:"not null"`
}

// GORM
func (u User) TableName() string {
	return "users"
}

// VALIDATION
func (u User) ValidateInsertOrRegister() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Username, validation.Required),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required),
	)
}

func (u User) ValidateUpdate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Username, validation.Required),
		validation.Field(&u.Email, validation.Required, is.Email),
	)
}

func (u User) ValidateLogin() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Username, validation.Required),
		validation.Field(&u.Password, validation.Required),
	)
}

// CONTRACT
type UserRepository interface {
	FindAll(ctx context.Context) ([]User, error)
	FindById(ctx context.Context, id int) (User, error)
	Save(ctx context.Context, user User) (User, error)
	Update(ctx context.Context, user User) (User, error)
	Delete(ctx context.Context, id int) error
	FindByUsername(ctx context.Context, username string) (User, error)
}

type UserUsecase interface {
	FindAll(ctx context.Context) ([]User, error)
	FindById(ctx context.Context, id int) (User, error)
	Save(ctx context.Context, user User) (User, error)
	Update(ctx context.Context, user User) (User, error)
	Delete(ctx context.Context, id int) error
	Login(ctx context.Context, user User) (Token, error)
	Register(ctx context.Context, user User) (User, error)
}
