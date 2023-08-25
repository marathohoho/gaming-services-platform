package wallet

import (
	"gaming-services-platform/internal/models"
	"gaming-services-platform/internal/repositories"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

type WalletService struct {
	rdb        *redis.Client
	repository *repositories.WalletRepository
}

func (ws WalletService) DepositFunds() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(models.UserWalletRequest)

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewGeneratedError(err.Error()))
		}

		if err := req.ValidateUserWalletRequest(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewGeneratedError(err.Error()))
		}

		// retrieve and check existing wallet
		walletRetrieved, err := ws.repository.Get(req.UserID)
		if err != nil {
			// check if the wallet does not exist for current user, if so create a new wallet
			if err == redis.Nil {
				createdWallet := models.UserWalletRequest{
					UserID: req.UserID,
					Amount: req.Amount,
				}
				if err := ws.repository.CreateOrUpdate(createdWallet); err != nil {
					return c.Status(fiber.StatusUnprocessableEntity).JSON(models.NewGeneratedError(err.Error()))
				}

				return c.Status(fiber.StatusOK).JSON(createdWallet)
			}
			return c.Status(fiber.StatusBadRequest).JSON(models.NewGeneratedError(err.Error()))
		}

		if walletRetrieved.Amount == nil || req.Amount == nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.NewGeneratedError("unable to process the request"))
		}
		updatedAmount := *walletRetrieved.Amount + *req.Amount
		updatedWallet := models.UserWalletRequest{
			UserID: req.UserID,
			Amount: &updatedAmount,
		}
		if err := ws.repository.CreateOrUpdate(updatedWallet); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.NewGeneratedError(err.Error()))
		}

		return c.Status(fiber.StatusOK).JSON(updatedWallet)
	}
}

func (ws WalletService) WithdrawFunds() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(models.UserWalletRequest)

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewGeneratedError(err.Error()))
		}

		if err := req.ValidateUserWalletRequest(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewGeneratedError(err.Error()))
		}

		walletRetrieved, err := ws.repository.Get(req.UserID)
		if err != nil {
			// check if the wallet does not exist for current user, if so decline the request
			if err == redis.Nil {
				return c.Status(fiber.StatusBadRequest).JSON(models.NewGeneratedError("no wallet exists for the user id provided"))
			}
			return c.Status(fiber.StatusBadRequest).JSON(models.NewGeneratedError(err.Error()))
		}

		if walletRetrieved.Amount == nil || req.Amount == nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.NewGeneratedError("unable to process the request"))
		}
		updatedAmount := *walletRetrieved.Amount - *req.Amount
		if updatedAmount < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewGeneratedError("cannot withdraw from walledt due to negative generated balance"))
		}

		updatedWallet := models.UserWalletRequest{
			UserID: req.UserID,
			Amount: &updatedAmount,
		}
		if err := ws.repository.CreateOrUpdate(updatedWallet); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.NewGeneratedError(err.Error()))
		}

		return c.Status(fiber.StatusOK).JSON(updatedWallet)
	}
}

func NewWalletService(redis *redis.Client, walletRepository *repositories.WalletRepository) *WalletService {
	return &WalletService{
		rdb:        redis,
		repository: walletRepository,
	}
}
