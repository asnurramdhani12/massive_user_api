package controllers

import (
	"net/http"
	"user_api/app/errors"
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
	users, err := u.UserUsecase.FindAll(c.Context())
	if err != nil {
		logger.GetLogger(c.Context()).Errorf("failed to find all users: %v", err)
		return errors.RenderErrorResponse(c, err)
	}

	return c.Status(http.StatusOK).JSON(contract.Response{
		StatusCode: http.StatusOK,
		Message:    "success get all users",
		Data:       users,
	})
}

func (u *UserController) FindById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		logger.GetLogger(c.Context()).Errorf("failed to parse id: %v", err)
		return c.Status(http.StatusBadRequest).JSON(contract.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "failed to parse id",
			Data:       nil,
		})
	}

	result, err := u.UserUsecase.FindById(c.Context(), id)
	if err != nil {
		logger.GetLogger(c.Context()).Errorf("failed to find user by id: %v", err)
		return errors.RenderErrorResponse(c, err)
	}

	return c.Status(http.StatusOK).JSON(contract.Response{
		StatusCode: http.StatusOK,
		Message:    "success get user by id",
		Data:       result,
	})
}

func (u *UserController) Save(c *fiber.Ctx) error {
	var user contract.User
	if err := c.BodyParser(&user); err != nil {
		logger.GetLogger(c.Context()).Errorf("failed to parse user: %v", err)
		return c.Status(http.StatusBadRequest).JSON(contract.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "failed to parse user",
			Data:       nil,
		})
	}

	result, err := u.UserUsecase.Save(c.Context(), user)
	if err != nil {
		logger.GetLogger(c.Context()).Errorf("failed to save user: %v", err)
		return errors.RenderErrorResponse(c, err)
	}

	return c.Status(http.StatusOK).JSON(contract.Response{
		StatusCode: http.StatusOK,
		Message:    "success save user",
		Data:       result,
	})
}

func (u *UserController) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		logger.GetLogger(c.Context()).Errorf("failed to parse id: %v", err)
		return c.Status(http.StatusBadRequest).JSON(contract.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "failed to parse id",
			Data:       nil,
		})
	}
	var user contract.User
	if err := c.BodyParser(&user); err != nil {
		logger.GetLogger(c.Context()).Errorf("failed to parse user: %v", err)
		return c.Status(http.StatusBadRequest).JSON(contract.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "failed to parse user",
			Data:       nil,
		})
	}

	// set id
	user.ID = id
	result, err := u.UserUsecase.Update(c.Context(), user)
	if err != nil {
		logger.GetLogger(c.Context()).Errorf("failed to update user: %v", err)
		return errors.RenderErrorResponse(c, err)
	}

	return c.Status(http.StatusOK).JSON(contract.Response{
		StatusCode: http.StatusOK,
		Message:    "success update user",
		Data:       result,
	})
}

func (u *UserController) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		logger.GetLogger(c.Context()).Errorf("failed to parse id: %v", err)
		return c.Status(http.StatusBadRequest).JSON(contract.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "failed to parse id",
			Data:       nil,
		})
	}

	err = u.UserUsecase.Delete(c.Context(), id)
	if err != nil {
		logger.GetLogger(c.Context()).Errorf("failed to delete user: %v", err)
		return errors.RenderErrorResponse(c, err)
	}

	return c.Status(http.StatusOK).JSON(contract.Response{
		StatusCode: http.StatusOK,
		Message:    "success delete user",
		Data:       nil,
	})
}
