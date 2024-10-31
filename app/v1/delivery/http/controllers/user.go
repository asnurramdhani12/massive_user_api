package controllers

import (
	"net/http"
	"user_api/app/errors"
	"user_api/app/middleware"
	"user_api/app/v1/contract"
	"user_api/utils/logger"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	UserUsecase contract.UserUsecase
}

func NewUserController(userUsecase contract.UserUsecase) *UserController {
	return &UserController{
		UserUsecase: userUsecase,
	}
}

func (u *UserController) FindAll(c *fiber.Ctx) error {
	users, err := u.UserUsecase.FindAll(c.UserContext())
	if err != nil {
		logger.GetLogger(c.UserContext()).Errorf("failed to find all users: %v", err)
		return errors.RenderErrorResponse(c, err)
	}

	return c.Status(http.StatusOK).JSON(contract.Response{
		StatusCode: http.StatusOK,
		Message:    "success get all users",
		Data:       users,
		MetaData: contract.MetaData{
			RequestID: c.UserContext().Value(middleware.RequestIDKey).(string),
		},
	})
}

func (u *UserController) FindById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		logger.GetLogger(c.UserContext()).Errorf("failed to parse id: %v", err)
		return c.Status(http.StatusBadRequest).JSON(contract.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "failed to parse id",
			Data:       nil,
			MetaData: contract.MetaData{
				RequestID: c.UserContext().Value(middleware.RequestIDKey).(string),
			},
		})
	}

	result, err := u.UserUsecase.FindById(c.UserContext(), id)
	if err != nil {
		logger.GetLogger(c.UserContext()).Errorf("failed to find user by id: %v", err)
		return errors.RenderErrorResponse(c, err)
	}

	return c.Status(http.StatusOK).JSON(contract.Response{
		StatusCode: http.StatusOK,
		Message:    "success get user by id",
		Data:       result,
		MetaData: contract.MetaData{
			RequestID: c.UserContext().Value(middleware.RequestIDKey).(string),
		},
	})
}

func (u *UserController) Save(c *fiber.Ctx) error {
	var user contract.User
	if err := c.BodyParser(&user); err != nil {
		logger.GetLogger(c.UserContext()).Errorf("failed to parse user: %v", err)
		return c.Status(http.StatusBadRequest).JSON(contract.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "failed to parse user",
			Data:       nil,
			MetaData: contract.MetaData{
				RequestID: c.UserContext().Value(middleware.RequestIDKey).(string),
			},
		})
	}

	// Validate
	if err := user.ValidateInsertOrRegister(); err != nil {
		logger.GetLogger(c.UserContext()).Errorf("failed to validate user: %v", err)
		return c.Status(http.StatusBadRequest).JSON(contract.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
			MetaData: contract.MetaData{
				RequestID: c.UserContext().Value(middleware.RequestIDKey).(string),
			},
		})
	}

	result, err := u.UserUsecase.Save(c.UserContext(), user)
	if err != nil {
		logger.GetLogger(c.UserContext()).Errorf("failed to save user: %v", err)
		return errors.RenderErrorResponse(c, err)
	}

	return c.Status(http.StatusOK).JSON(contract.Response{
		StatusCode: http.StatusOK,
		Message:    "success save user",
		Data:       result,
		MetaData: contract.MetaData{
			RequestID: c.UserContext().Value(middleware.RequestIDKey).(string),
		},
	})
}

func (u *UserController) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		logger.GetLogger(c.UserContext()).Errorf("failed to parse id: %v", err)
		return c.Status(http.StatusBadRequest).JSON(contract.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "failed to parse id",
			Data:       nil,
			MetaData: contract.MetaData{
				RequestID: c.UserContext().Value(middleware.RequestIDKey).(string),
			},
		})
	}
	var user contract.User
	if err := c.BodyParser(&user); err != nil {
		logger.GetLogger(c.UserContext()).Errorf("failed to parse user: %v", err)
		return c.Status(http.StatusBadRequest).JSON(contract.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "failed to parse user",
			Data:       nil,
			MetaData: contract.MetaData{
				RequestID: c.UserContext().Value(middleware.RequestIDKey).(string),
			},
		})
	}

	// Validate
	if err := user.ValidateUpdate(); err != nil {
		logger.GetLogger(c.UserContext()).Errorf("failed to validate user: %v", err)
		return c.Status(http.StatusBadRequest).JSON(contract.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
			MetaData: contract.MetaData{
				RequestID: c.UserContext().Value(middleware.RequestIDKey).(string),
			},
		})
	}

	// set id
	user.ID = id
	result, err := u.UserUsecase.Update(c.UserContext(), user)
	if err != nil {
		logger.GetLogger(c.UserContext()).Errorf("failed to update user: %v", err)
		return errors.RenderErrorResponse(c, err)
	}

	return c.Status(http.StatusOK).JSON(contract.Response{
		StatusCode: http.StatusOK,
		Message:    "success update user",
		Data:       result,
		MetaData: contract.MetaData{
			RequestID: c.UserContext().Value(middleware.RequestIDKey).(string),
		},
	})
}

func (u *UserController) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		logger.GetLogger(c.UserContext()).Errorf("failed to parse id: %v", err)
		return c.Status(http.StatusBadRequest).JSON(contract.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "failed to parse id",
			Data:       nil,
		})
	}

	err = u.UserUsecase.Delete(c.UserContext(), id)
	if err != nil {
		logger.GetLogger(c.UserContext()).Errorf("failed to delete user: %v", err)
		return errors.RenderErrorResponse(c, err)
	}

	return c.Status(http.StatusOK).JSON(contract.Response{
		StatusCode: http.StatusOK,
		Message:    "success delete user",
		Data:       nil,
	})
}
