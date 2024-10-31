package user

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"user_api/app/v1/contract"
	"user_api/utils/logger"
	redis_utils "user_api/utils/redis"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

const (
	userKey                 = "user:all"
	userDetailsKey          = "user:%d"
	userDetailByUsernameKey = "user:username:%s"
	userDeleteKey           = "user:*"
)

type UserRepositoryImpl struct {
	DB  *gorm.DB
	RDB *redis.Client
}

// FindByUsername implements contract.UserRepository.
func (u *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (contract.User, error) {
	var user contract.User

	// check redis
	key := fmt.Sprintf(userDetailByUsernameKey, username)
	err := u.RDB.Get(ctx, key).Scan(&user)
	if err == nil {
		logger.GetLogger(ctx).Infof("found user in redis: %v", user)
		return user, nil
	}

	err = u.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to find user: %v", err)
		return contract.User{}, err
	}

	// set redis
	err = u.RDB.Set(ctx, key, user, -1).Err()
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to set redis: %v", err)
		return contract.User{}, err
	}

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

	// invalidate redis
	err := redis_utils.DeleteWithPattern(ctx, u.RDB, userDeleteKey)
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to invalidate redis: %v", err)
		return err
	}

	return nil
}

// FindAll implements contract.UserRepository.
func (u *UserRepositoryImpl) FindAll(ctx context.Context) ([]contract.User, error) {
	var users []contract.User

	// check redis
	data, err := u.RDB.Get(ctx, userKey).Result()
	if err == nil {
		json.Unmarshal([]byte(data), &users)
		logger.GetLogger(ctx).Infof("found users in redis: %v", users)
		return users, nil
	}

	err = u.DB.Omit("password").Find(&users).Error
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to find all users: %v", err)
		return []contract.User{}, err
	}

	// set redis
	jsonData, err := json.Marshal(users)
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to marshal users: %v", err)
		return []contract.User{}, err
	}
	err = u.RDB.Set(ctx, userKey, jsonData, -1).Err()
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to set redis: %v", err)
		return []contract.User{}, err
	}

	return users, err
}

// FindById implements contract.UserRepository.
func (u *UserRepositoryImpl) FindById(ctx context.Context, id int) (contract.User, error) {
	var user contract.User

	// check redis
	key := fmt.Sprintf(userDetailsKey, id)
	err := u.RDB.Get(ctx, key).Scan(&user)
	if err == nil {
		logger.GetLogger(ctx).Infof("found user in redis: %v", user)
		return user, nil
	}

	err = u.DB.Omit("password").First(&user, id).Error
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to find user: %v", err)
		return contract.User{}, err
	}

	// set redis
	err = u.RDB.Set(ctx, key, user, -1).Err()
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to set redis: %v", err)
		return contract.User{}, err
	}

	return user, nil
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

	// invalidate redis
	err = redis_utils.DeleteWithPattern(ctx, u.RDB, userDeleteKey)
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to invalidate redis: %v", err)
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

	// invalidate redis
	err = redis_utils.DeleteWithPattern(ctx, u.RDB, userDeleteKey)
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to invalidate redis: %v", err)
		return contract.User{}, err
	}

	return user, err
}

func NewUserRepository(db *gorm.DB, rdb *redis.Client) contract.UserRepository {
	return &UserRepositoryImpl{
		DB:  db,
		RDB: rdb,
	}
}
