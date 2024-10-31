package errors

import (
	"errors"
	"user_api/app/middleware"
	"user_api/app/v1/contract"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrNotFound = errors.New("record not found")
	ErrConflict = errors.New("record already exist")
)

func RenderErrorResponse(c *fiber.Ctx, err error) error {
	switch err {
	case ErrNotFound:
		return c.Status(fiber.StatusNotFound).JSON(contract.Response{
			StatusCode: fiber.StatusNotFound,
			Message:    err.Error(),
			Data:       nil,
			MetaData: contract.MetaData{
				RequestID: c.UserContext().Value(middleware.RequestIDKey).(string),
			},
		})
	case ErrConflict:
		return c.Status(fiber.StatusConflict).JSON(contract.Response{
			StatusCode: fiber.StatusConflict,
			Message:    err.Error(),
			Data:       nil,
			MetaData: contract.MetaData{
				RequestID: c.UserContext().Value(middleware.RequestIDKey).(string),
			},
		})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(contract.Response{
			StatusCode: fiber.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
			MetaData: contract.MetaData{
				RequestID: c.UserContext().Value(middleware.RequestIDKey).(string),
			},
		})
	}
}
