package repositories

import (
	"context"
	"encoding/json"
	"gaming-services-platform/internal/models"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

var ValueDefatulExpDuration = 120 * time.Minute

type UserRepository struct {
	client *redis.Client
	ctx    context.Context
}

func (u UserRepository) Create(user *models.User) (string, error) {
	userId := uuid.New()

	userMarshalled, err := json.Marshal(user)
	if err != nil {
		return "", nil
	}

	if err := u.client.Set(u.ctx, userId.String(), userMarshalled, ValueDefatulExpDuration).Err(); err != nil {
		return "", err
	}
	return userId.String(), nil
}

func (u UserRepository) Get(userId string) (*models.User, error) {
	userMarshalled, err := u.client.Get(u.ctx, userId).Result()
	if err != nil {
		return nil, err
	}

	user := new(models.User)
	if err := json.Unmarshal([]byte(userMarshalled), &user); err != nil {
		return nil, err
	}

	return user, nil
}

func NewUserRepository(ctx context.Context, client *redis.Client) *UserRepository {
	return &UserRepository{
		client: client,
		ctx:    ctx,
	}
}
