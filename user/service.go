package user

import (
	"gaming-services-platform/internal/models"
	"gaming-services-platform/internal/repositories"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

type UserService struct {
	rdb        *redis.Client
	repository *repositories.UserRepository
}

func (us UserService) Create() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(models.User)

		// parse the request body
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.ErrBadRequest.Code).JSON(models.NewGeneratedError(err.Error()))
		}

		// validate the request body
		if err := req.ValidateUserRequest(); err != nil {
			return c.Status(fiber.ErrBadRequest.Code).JSON(models.NewGeneratedError(err.Error()))
		}

		// store our new user to inmemory storage
		userId, err := us.repository.Create(req)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.NewGeneratedError(err.Error()))
		}

		createdUser := models.UserResponse{
			ID:    userId,
			Email: req.Email,
		}
		return c.Status(fiber.StatusOK).JSON(createdUser)
	}
}

func (us UserService) Get() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Params("id")

		user, err := us.repository.Get(userId)
		if err != nil {
			var generatedError *models.GeneratedError
			if user == nil {
				generatedError = models.NewGeneratedError("user was not found")
				return c.Status(fiber.ErrBadRequest.Code).JSON(generatedError)
			}
			generatedError = models.NewGeneratedError(err.Error())
			return c.Status(fiber.ErrBadRequest.Code).JSON(generatedError)
		}

		return c.Status(fiber.StatusOK).JSON(user)
	}
}

func NewUserService(redis *redis.Client, userRepository *repositories.UserRepository) *UserService {
	return &UserService{
		rdb:        redis,
		repository: userRepository,
	}
}
