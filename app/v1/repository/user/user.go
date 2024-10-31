package user

import (
	"context"
	"fmt"
	"user_api/app/v1/contract"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

// FindByUsername implements contract.UserRepository.
func (u *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (contract.User, error) {
	var user contract.User
	err := u.DB.Where("username = ?", username).First(&user).Error
	fmt.Println(err)
	return user, err
}

// Delete implements contract.UserRepository.
func (u *UserRepositoryImpl) Delete(ctx context.Context, id int) error {
	err := u.DB.Delete(&contract.User{}, id).Error
	return err
}

// FindAll implements contract.UserRepository.
func (u *UserRepositoryImpl) FindAll(ctx context.Context) ([]contract.User, error) {
	var users []contract.User
	err := u.DB.Find(&users).Error
	return users, err
}

// FindById implements contract.UserRepository.
func (u *UserRepositoryImpl) FindById(ctx context.Context, id int) (contract.User, error) {
	var user contract.User
	err := u.DB.First(&user, id).Error
	return user, err
}

// Save implements contract.UserRepository.
func (u *UserRepositoryImpl) Save(ctx context.Context, user contract.User) (contract.User, error) {
	err := u.DB.Create(&user).Error
	return user, err
}

// Update implements contract.UserRepository.
func (u *UserRepositoryImpl) Update(ctx context.Context, user contract.User) (contract.User, error) {
	err := u.DB.Save(&user).Error
	return user, err
}

func NewUserRepository(db *gorm.DB) contract.UserRepository {
	return &UserRepositoryImpl{
		DB: db,
	}
}
