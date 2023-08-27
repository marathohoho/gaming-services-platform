package user

import (
	"gaming-services-platform/internal"
	"gaming-services-platform/internal/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UserHost string
}

var generatedError *models.GeneratedError

func (h UserHandler) Register() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Print("received request to create a new user")
		reqPayload := c.Body()
		url := h.UserHost + "/users"
		response, err := internal.SendJsonRequest(fiber.MethodPost, url, reqPayload, nil)
		if err != nil {
			generatedError = models.NewGeneratedError(err.Error())
			return c.Status(fiber.ErrInternalServerError.Code).JSON(generatedError)
		}

		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		if response.StatusCode != fiber.StatusCreated {
			return c.Status(response.StatusCode).Send(response.Body)
		}

		log.Print("created a new user: ", response.Body)
		return c.Status(fiber.StatusCreated).Send(response.Body)
	}
}

func (h UserHandler) Get() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Params("id")
		url := h.UserHost + "/user/" + userId
		response, err := internal.SendJsonRequest(fiber.MethodGet, url, nil, nil)
		if err != nil {
			generatedError = models.NewGeneratedError(err.Error())
			return c.Status(fiber.ErrInternalServerError.Code).JSON(generatedError)
		}

		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		if response.StatusCode != fiber.StatusOK {
			return c.Status(response.StatusCode).Send(response.Body)
		}
		return c.Status(fiber.StatusOK).Send(response.Body)
	}
}

func NewUserHandler(host string) *UserHandler {
	return &UserHandler{
		UserHost: host,
	}
}
