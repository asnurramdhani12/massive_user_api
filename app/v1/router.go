package v1

import (
	"net/http"
	"user_api/app/v1/contract"
	"user_api/app/v1/delivery/http/controllers"

	appConfig "user_api/config"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func NewRouter(f *fiber.App, deps *Dependencies) {
	authControllers := controllers.NewAuthController(deps.userUsecase)
	userContrllers := controllers.NewUserController(deps.userUsecase)

	f.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	v1 := f.Group("/api/v1")

	auth := v1.Group("/auth")
	auth.Post("/login", authControllers.Login)
	auth.Post("/register", authControllers.Register)

	v1.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(appConfig.Config().JWTSecret)},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(http.StatusUnauthorized).JSON(contract.Response{
				StatusCode: http.StatusUnauthorized,
				Message:    "unauthorized",
				Data:       nil,
			})
		},
	}))

	admin := v1.Group("/admin")

	user := admin.Group("/user")
	user.Get("/", userContrllers.FindAll)
	user.Get("/:id", userContrllers.FindById)
	user.Post("/", userContrllers.Save)
	user.Put("/:id", userContrllers.Update)
	user.Delete("/:id", userContrllers.Delete)
}
