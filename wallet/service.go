package wallet

import (
	"encoding/json"
	"fmt"
	"gaming-services-platform/internal/models"
	"gaming-services-platform/internal/repositories"

	ws "gaming-services-platform/websocket"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/websocket"
)

type WalletService struct {
	rdb                 *redis.Client
	repository          *repositories.WalletRepository
	userRepository      *repositories.UserRepository
	websocketConnection *websocket.Conn
}

var currentGameLeader string
var currentHighestWallet float64

func (w WalletService) DepositFunds() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(models.UserWalletRequest)

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewGeneratedError(err.Error()))
		}

		if err := req.ValidateUserWalletRequest(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewGeneratedError(err.Error()))
		}

		// retrieve and check existing wallet
		walletRetrieved, err := w.repository.Get(req.UserID)
		if err != nil {
			// check if the wallet does not exist for current user, if so create a new wallet
			if err == redis.Nil {
				createdWallet := models.UserWalletRequest{
					UserID: req.UserID,
					Amount: req.Amount,
				}
				if err := w.repository.CreateOrUpdate(createdWallet); err != nil {
					return c.Status(fiber.StatusUnprocessableEntity).JSON(models.NewGeneratedError(err.Error()))
				}

				_ = w.sendWSMessage(createdWallet, *req.Amount, models.USER_DEPOSIT_INIT_MESSAGE)

				return c.Status(fiber.StatusOK).JSON(createdWallet)
			}
			return c.Status(fiber.StatusBadRequest).JSON(models.NewGeneratedError(err.Error()))
		}

		if req.Amount == nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.NewGeneratedError("unable to process the request"))
		}
		updatedAmount := walletRetrieved.Amount + *req.Amount
		updatedWallet := models.UserWalletRequest{
			UserID: req.UserID,
			Amount: &updatedAmount,
		}
		if err := w.repository.CreateOrUpdate(updatedWallet); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.NewGeneratedError(err.Error()))
		}

		_ = w.sendWSMessage(updatedWallet, *req.Amount, models.USER_DEPOSIT_MESSAGE)
		return c.Status(fiber.StatusOK).JSON(updatedWallet)
	}
}

func (w WalletService) sendWSMessage(wallet models.UserWalletRequest, amountAdded float64, message string) error {
	// send a current game update message
	var newEvent ws.Event
	newEvent.Type = ws.EventSendMessage
	var newMessageEvent ws.SendMessageEvent
	newMessageEvent.Message = fmt.Sprintf(message, wallet.UserID, amountAdded, *wallet.Amount)
	newMessageEvent.GameEvent = ws.GameOutcomes
	t, _ := json.Marshal(newMessageEvent)
	newEvent.Payload = t

	_ = w.websocketConnection.WriteJSON(newEvent)

	// send a leaderboard game message
	if *wallet.Amount > currentHighestWallet {
		var messageTemplate string
		currentHighestWallet = *wallet.Amount
		newMessageEvent.GameEvent = ws.LeaderboardChanges

		if currentGameLeader == wallet.UserID {
			messageTemplate = models.CURRENT_LEADER_MESSAGE_USER_UNCHANGED
		} else {
			messageTemplate = models.CURRENT_LEADER_MESSAGE
		}

		newMessageEvent.Message = fmt.Sprintf(messageTemplate, wallet.UserID, *wallet.Amount)
		t, _ := json.Marshal(newMessageEvent)
		newEvent.Payload = t

		_ = w.websocketConnection.WriteJSON(newEvent)
		currentGameLeader = wallet.UserID
	}

	return nil
}

func (w WalletService) WithdrawFunds() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(models.UserWalletRequest)

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewGeneratedError(err.Error()))
		}

		if err := req.ValidateUserWalletRequest(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewGeneratedError(err.Error()))
		}

		walletRetrieved, err := w.repository.Get(req.UserID)
		if err != nil {
			// check if the wallet does not exist for current user, if so decline the request
			if err == redis.Nil {
				return c.Status(fiber.StatusBadRequest).JSON(models.NewGeneratedError("no wallet exists for the user id provided"))
			}
			return c.Status(fiber.StatusBadRequest).JSON(models.NewGeneratedError(err.Error()))
		}

		if req.Amount == nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.NewGeneratedError("unable to process the request"))
		}
		updatedAmount := walletRetrieved.Amount - *req.Amount
		if updatedAmount < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewGeneratedError("cannot withdraw from walledt due to negative generated balance"))
		}

		updatedWallet := models.UserWalletRequest{
			UserID: req.UserID,
			Amount: &updatedAmount,
		}

		_ = w.sendWSMessage(updatedWallet, *req.Amount, models.USER_WITHDRAW_MESSAGE)
		if err := w.repository.CreateOrUpdate(updatedWallet); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.NewGeneratedError(err.Error()))
		}

		return c.Status(fiber.StatusOK).JSON(updatedWallet)
	}
}

func NewWalletService(redis *redis.Client, walletRepository *repositories.WalletRepository, connection *websocket.Conn) *WalletService {
	return &WalletService{
		rdb:                 redis,
		repository:          walletRepository,
		websocketConnection: connection,
	}
}
