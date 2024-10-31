package user

import (
	"context"
	"fmt"
	"strings"
	"user_api/app/v1/contract"
	"user_api/utils/logger"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

// FindByUsername implements contract.UserRepository.
func (u *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (contract.User, error) {
	var user contract.User
	err := u.DB.Where("username = ?", username).First(&user).Error

	return user, err
}

// Delete implements contract.UserRepository.
func (u *UserRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := u.DB.Delete(&contract.User{}, id)

	if query.RowsAffected == 0 {
		logger.GetLogger(ctx).Errorf("failed user not found with id: %d", id)
		return gorm.ErrRecordNotFound
	}

	if query.Error != nil {
		logger.GetLogger(ctx).Errorf("failed to delete user: %v", query.Error)
		return fmt.Errorf("failed to delete user: %v", query.Error)
	}

	return nil
}

// FindAll implements contract.UserRepository.
func (u *UserRepositoryImpl) FindAll(ctx context.Context) ([]contract.User, error) {
	var users []contract.User
	err := u.DB.Omit("password").Find(&users).Error

	return users, err
}

// FindById implements contract.UserRepository.
func (u *UserRepositoryImpl) FindById(ctx context.Context, id int) (contract.User, error) {
	var user contract.User
	err := u.DB.Omit("password").First(&user, id).Error
	return user, err
}

// Save implements contract.UserRepository.
func (u *UserRepositoryImpl) Save(ctx context.Context, user contract.User) (contract.User, error) {
	err := u.DB.Create(&user).Error
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to insert user: %v", err)
		if strings.Contains(err.Error(), "duplicate key") {
			return contract.User{}, gorm.ErrDuplicatedKey
		}

		return contract.User{}, err
	}

	return user, nil
}

// Update implements contract.UserRepository.
func (u *UserRepositoryImpl) Update(ctx context.Context, user contract.User) (contract.User, error) {
	query := u.DB.Where("id = ?", user.ID).Updates(user)

	err := query.Error
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to update user: %v", err)
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "duplicated key") {
			return contract.User{}, gorm.ErrDuplicatedKey
		}
		return contract.User{}, err
	}

	if query.RowsAffected == 0 {
		logger.GetLogger(ctx).Errorf("failed user not found with id: %d", user.ID)
		return contract.User{}, gorm.ErrRecordNotFound
	}

	return user, err
}

func NewUserRepository(db *gorm.DB) contract.UserRepository {
	return &UserRepositoryImpl{
		DB: db,
	}
}
