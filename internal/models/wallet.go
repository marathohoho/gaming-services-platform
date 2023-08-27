package models

import "errors"

type UserWalletRequest struct {
	UserID string   `json:"userId"`
	Amount *float64 `json:"amount"`
}

type UserWallet struct {
	Amount float64 `json:"amount"`
}

const (
	USER_DEPOSIT_MESSAGE      = "User %s added %.2f amount to wallet. Current wallet balance is %.2f"
	USER_DEPOSIT_INIT_MESSAGE = "User %s initialized a wallet with %.2f amount. Current wallet balance is %.2f"
	USER_WITHDRAW_MESSAGE     = "User %s withdrew %.2f amount to wallet. Current wallet balance is %.2f"

	CURRENT_LEADER_MESSAGE                = "User %s is leading the game with %.2f amount in the wallet."
	CURRENT_LEADER_MESSAGE_USER_UNCHANGED = "User %s is still leading the game with %.2f amount in the wallet. Well done!"
)

func (uwr *UserWalletRequest) ValidateUserWalletRequest() error {
	if len(uwr.UserID) <= 0 {
		return errors.New("user id must not be empty")
	}

	if uwr.Amount == nil {
		return errors.New("request must contain wallet amount")
	}

	if uwr.Amount != nil && *uwr.Amount <= 0.0 {
		return errors.New("request amount must be greated than zero")
	}

	return nil
}
