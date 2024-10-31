package user

import (
	"context"
	"errors"
	"time"
	"user_api/app/v1/contract"
	"user_api/utils/bcrypt"
	"user_api/utils/jwt"
	"user_api/utils/logger"

	appErr "user_api/app/errors"

	"gorm.io/gorm"
)

type UserUsecaseImpl struct {
	UserRepo contract.UserRepository
}

// Delete implements contract.UserUsecase.
func (u *UserUsecaseImpl) Delete(ctx context.Context, id int) error {
	err := u.UserRepo.Delete(ctx, id)
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to delete user: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return appErr.ErrNotFound
		}
		return err
	}

	return err
}

// FindAll implements contract.UserUsecase.
func (u *UserUsecaseImpl) FindAll(ctx context.Context) ([]contract.User, error) {
	result, err := u.UserRepo.FindAll(ctx)

	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to find all users: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []contract.User{}, appErr.ErrNotFound
		}
		return []contract.User{}, err
	}

	if result == nil {
		logger.GetLogger(ctx).Errorf("failed to find all users: %v", "record not found")
		return []contract.User{}, appErr.ErrNotFound
	}

	return result, err
}

// FindById implements contract.UserUsecase.
func (u *UserUsecaseImpl) FindById(ctx context.Context, id int) (contract.User, error) {
	result, err := u.UserRepo.FindById(ctx, id)

	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to find user: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return contract.User{}, appErr.ErrNotFound
		}
		return contract.User{}, err
	}

	return result, err
}

// Login implements contract.UserUsecase.
func (u *UserUsecaseImpl) Login(ctx context.Context, user contract.User) (contract.Token, error) {
	data, err := u.UserRepo.FindByUsername(ctx, user.Username)
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to find user: %v", err)
		return contract.Token{}, err
	}

	if user == (contract.User{}) {
		return contract.Token{}, errors.New("user not found")
	}

	// Check Password
	if !bcrypt.CheckPasswordHash(user.Password, data.Password) {
		return contract.Token{}, errors.New("invalid password")
	}

	// Generate Token
	token, err := jwt.GenerateToken(map[string]any{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	if err != nil {
		return contract.Token{}, errors.New("cannot generate token")
	}

	// Generate Refresh Token
	refreshToken, err := jwt.GenerateToken(map[string]any{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	if err != nil {
		return contract.Token{}, errors.New("cannot generate refresh token")
	}

	return contract.Token{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}

// Register implements contract.UserUsecase.
func (u *UserUsecaseImpl) Register(ctx context.Context, user contract.User) (contract.User, error) {
	//Check Email Exist
	users, err := u.UserRepo.FindByUsername(ctx, user.Username)
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to find user: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return contract.User{}, errors.New("cannot find user")
		}
		return contract.User{}, errors.New("cannot find user")
	}

	if users != (contract.User{}) {
		return contract.User{}, appErr.ErrConflict
	}

	// Hash Password
	hashedPassword, err := bcrypt.HashPassword(user.Password)
	if err != nil {
		return contract.User{}, errors.New("cannot hash password")
	}

	user.Password = hashedPassword

	// Insert User
	user, err = u.UserRepo.Save(ctx, user)
	if err != nil {
		return contract.User{}, errors.New("cannot insert user")
	}

	return user, nil
}

// Save implements contract.UserUsecase.
func (u *UserUsecaseImpl) Save(ctx context.Context, user contract.User) (contract.User, error) {
	// Hash Password
	hashedPassword, err := bcrypt.HashPassword(user.Password)
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to hash password: %v", err)
		return contract.User{}, err
	}

	user.Password = hashedPassword

	// Insert User
	user, err = u.UserRepo.Save(ctx, user)
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to insert user: %v", err)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return contract.User{}, appErr.ErrConflict
		}

		return contract.User{}, err
	}

	return user, nil
}

// Update implements contract.UserUsecase.
func (u *UserUsecaseImpl) Update(ctx context.Context, user contract.User) (contract.User, error) {
	// Hash Password
	if user.Password != "" {
		hashedPassword, err := bcrypt.HashPassword(user.Password)
		if err != nil {
			logger.GetLogger(ctx).Errorf("failed to hash password: %v", err)
			return contract.User{}, err
		}

		user.Password = hashedPassword
	}

	// Update User
	user, err := u.UserRepo.Update(ctx, user)
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to update user: %v", err)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return contract.User{}, appErr.ErrConflict
		}
		return contract.User{}, err
	}

	return user, nil
}

func NewUserUsecase(userRepo contract.UserRepository) contract.UserUsecase {
	return &UserUsecaseImpl{
		UserRepo: userRepo,
	}
}
