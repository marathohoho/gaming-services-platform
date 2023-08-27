package wallet

import (
	"gaming-services-platform/internal"
	"gaming-services-platform/internal/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

type WalletHandler struct {
	WalletHost string
}

var generatedError *models.GeneratedError

func (w WalletHandler) Deposit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Print("received a request for user balance deposit")
		reqPayload := c.Body()
		url := w.WalletHost + "/deposit"
		response, err := internal.SendJsonRequest(fiber.MethodPost, url, reqPayload, nil)
		if err != nil {
			generatedError = models.NewGeneratedError(err.Error())
			return c.Status(fiber.ErrBadRequest.Code).JSON(generatedError)
		}

		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		if response.StatusCode != fiber.StatusOK {
			return c.Status(response.StatusCode).Send(response.Body)
		}
		return c.Status(fiber.StatusAccepted).Send(response.Body)
	}
}

func (w WalletHandler) Withdraw() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Print("received a request for user balance withdraw")
		reqPayload := c.Body()
		url := w.WalletHost + "/withdraw"
		response, err := internal.SendJsonRequest(fiber.MethodPost, url, reqPayload, nil)
		if err != nil {
			generatedError = models.NewGeneratedError(err.Error())
			return c.Status(fiber.ErrBadRequest.Code).JSON(generatedError)
		}

		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		if response.StatusCode != fiber.StatusOK {
			return c.Status(response.StatusCode).Send(response.Body)
		}
		return c.Status(fiber.StatusAccepted).Send(response.Body)
	}
}

func NewWalletHandler(host string) *WalletHandler {
	return &WalletHandler{
		WalletHost: host,
	}
}
