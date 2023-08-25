package grpc

import (
	"gaming-services-platform/internal/models"
	"gaming-services-platform/proto"

	"github.com/gofiber/fiber/v2"
)

func GetBalanceHandler(client proto.BalanceServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Params("userId")

		balance, err := client.Get(c.Context(), &proto.RequestBalance{UserId: userId})
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.NewGeneratedError(err.Error()))
		}

		balanceResponse := models.UserWallet{
			Amount: balance.Amount,
		}
		return c.Status(fiber.StatusOK).JSON(balanceResponse)
	}
}
