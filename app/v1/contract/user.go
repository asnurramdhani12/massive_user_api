package contract

import "context"

type User struct {
	ID       int    `json:"id" gorm:"primary_key"`
	Username string `json:"username" gorm:"unique"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password" gorm:"not null"`
}

func (u User) TableName() string {
	return "users"
}

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
