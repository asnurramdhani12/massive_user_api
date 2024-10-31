package controllers

import (
	"net/http"
	"user_api/app/v1/contract"
	"user_api/utils/logger"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	UserUsecase contract.UserUsecase
}

func NewAuthController(userUsecase contract.UserUsecase) *AuthController {
	return &AuthController{
		UserUsecase: userUsecase,
	}
}

func (u *AuthController) Login(c *fiber.Ctx) error {
	var user contract.User
	if err := c.BodyParser(&user); err != nil {
		logger.GetLogger(c.Context()).Errorf("failed to parse user: %v", err)
		return c.Status(http.StatusBadRequest).JSON(contract.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "failed to parse user",
			Data:       nil,
		})
	}

	result, err := u.UserUsecase.Login(c.Context(), user)
	if err != nil {
		logger.GetLogger(c.Context()).Errorf("failed to login user: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(contract.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to login user",
			Data:       nil,
		})
	}

	return c.Status(http.StatusOK).JSON(contract.Response{
		StatusCode: http.StatusOK,
		Message:    "success login user",
		Data:       result,
	})
}

func (u *AuthController) Register(c *fiber.Ctx) error {
	var user contract.User
	if err := c.BodyParser(&user); err != nil {
		logger.GetLogger(c.Context()).Errorf("failed to parse user: %v", err)
		return c.Status(http.StatusBadRequest).JSON(contract.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "failed to parse user",
			Data:       nil,
		})
	}

	result, err := u.UserUsecase.Register(c.Context(), user)
	if err != nil {
		logger.GetLogger(c.Context()).Errorf("failed to register user: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(contract.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to register user",
			Data:       nil,
		})
	}

	return c.Status(http.StatusOK).JSON(contract.Response{
		StatusCode: http.StatusOK,
		Message:    "success register user",
		Data:       result,
	})
}
