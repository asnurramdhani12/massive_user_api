package user

import (
	"context"
	"errors"
	"time"
	"user_api/app/v1/contract"
	"user_api/utils/bcrypt"
	"user_api/utils/jwt"
	"user_api/utils/logger"
)

type UserUsecaseImpl struct {
	UserRepo contract.UserRepository
}

// Delete implements contract.UserUsecase.
func (u *UserUsecaseImpl) Delete(ctx context.Context, id int) error {
	err := u.UserRepo.Delete(ctx, id)
	return err
}

// FindAll implements contract.UserUsecase.
func (u *UserUsecaseImpl) FindAll(ctx context.Context) ([]contract.User, error) {
	result, err := u.UserRepo.FindAll(ctx)
	return result, err
}

// FindById implements contract.UserUsecase.
func (u *UserUsecaseImpl) FindById(ctx context.Context, id int) (contract.User, error) {
	result, err := u.UserRepo.FindById(ctx, id)
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
	if err != nil && err.Error() != "record not found" {
		return contract.User{}, errors.New("cannot find user")
	}

	if users != (contract.User{}) {
		return contract.User{}, errors.New("email already exist")
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
	return u.UserRepo.Save(ctx, user)
}

// Update implements contract.UserUsecase.
func (u *UserUsecaseImpl) Update(ctx context.Context, user contract.User) (contract.User, error) {
	return u.UserRepo.Update(ctx, user)
}

func NewUserUsecase(userRepo contract.UserRepository) contract.UserUsecase {
	return &UserUsecaseImpl{
		UserRepo: userRepo,
	}
}
