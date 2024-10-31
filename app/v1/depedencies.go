package v1

import (
	"context"
	"user_api/app/v1/contract"
	userRepo "user_api/app/v1/repository/user"
	userUseCase "user_api/app/v1/usecase/user"
	"user_api/config"
	"user_api/utils/logger"
)

type UseCase struct {
	userUsecase contract.UserUsecase
}

type Repository struct {
	userRepo contract.UserRepository
}

type Dependencies struct {
	UseCase
	Repository
}

func NewDependencies(ctx context.Context, cfg *config.Configuration) *Dependencies {
	db, err := config.NewPostgreSql(ctx, *cfg)
	if err != nil {
		logger.GetLogger(ctx).Panicf("failed to connect database: %v", err)
		return nil
	}

	rdb, err := config.NewCache(ctx, *cfg)
	if err != nil {
		logger.GetLogger(ctx).Panicf("failed to connect redis: %v", err)
		return nil
	}

	userRepo := userRepo.NewUserRepository(db, rdb)
	userUsecase := userUseCase.NewUserUsecase(userRepo)

	return &Dependencies{
		UseCase: UseCase{
			userUsecase: userUsecase,
		},
		Repository: Repository{
			userRepo: userRepo,
		},
	}
}
