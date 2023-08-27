package grpc

import (
	"gaming-services-platform/internal/models"
	"gaming-services-platform/proto"
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetBalanceHandler(client proto.BalanceServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Print("received a request to get a user balance via gRPC call")

		userId := c.Params("userId")
		balance, err := client.Get(c.Context(), &proto.RequestBalance{UserId: userId})
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.NewGeneratedError(err.Error()))
		}

		balanceResponse := models.UserWallet{
			Amount: balance.Amount,
		}

		log.Print("successfully receive a user balance via gRPC")
		return c.Status(fiber.StatusOK).JSON(balanceResponse)
	}
}
