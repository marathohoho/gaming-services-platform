package models

import "errors"

type UserWalletRequest struct {
	UserID string   `json:"userId"`
	Amount *float64 `json:"amount"`
}

type UserWallet struct {
	Amount float64 `json:"amount"`
}

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
