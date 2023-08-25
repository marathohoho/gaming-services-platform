package repositories

import (
	"context"
	"encoding/json"
	"gaming-services-platform/internal/models"

	"github.com/go-redis/redis/v8"
)

type WalletRepository struct {
	client *redis.Client
	ctx    context.Context
}

func (w WalletRepository) CreateOrUpdate(userWalletRequest models.UserWalletRequest) error {
	walletRequestMarshalled, err := json.Marshal(models.UserWallet{
		Amount: userWalletRequest.Amount,
	})
	if err != nil {
		return err
	}

	if err := w.client.Set(w.ctx, userWalletRequest.UserID, walletRequestMarshalled, ValueDefatulExpDuration).Err(); err != nil {
		return err
	}

	return nil
}

func (w WalletRepository) Get(userId string) (*models.UserWallet, error) {
	walletMarshalled, err := w.client.Get(w.ctx, userId).Result()
	if err != nil {
		return nil, err
	}

	wallet := new(models.UserWallet)
	if err := json.Unmarshal([]byte(walletMarshalled), &wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

func NewWalletRepository(ctx context.Context, client *redis.Client) *WalletRepository {
	return &WalletRepository{
		client: client,
		ctx:    ctx,
	}
}
