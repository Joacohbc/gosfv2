package auth

import (
	"context"
	"errors"
	"fmt"
	"gosfV2/src/models/database"
	"gosfV2/src/models/env"
	"net/http"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/labstack/echo"
)

var (
	ErrTokenNotFound               = fmt.Errorf("token not found")
	ErrMaxTokens                   = fmt.Errorf("max tokens reached")
	TokenDuration    time.Duration = time.Minute * time.Duration(env.Config.JWTMinutes)
)

type TokenManager interface {
	AddToken(id uint, token string) error
	RemoveToken(id uint, token string) error
	GetTokens(id uint) ([]string, error)
	RemoveUserTokens(id uint) error
	ExistsToken(id uint, token string) (bool, error)
	ManageError(err error) error
}

func NewTokenManager() TokenManager {
	return &tokenRedis{db: database.GetRedis(), ctx: context.Background()}
}

type tokenRedis struct {
	db  *redis.Client
	ctx context.Context
}

func (t *tokenRedis) ManageError(err error) error {
	if errors.Is(err, redis.Nil) {
		return ErrTokenNotFound
	}
	return err
}

func (t *tokenRedis) AddToken(id uint, token string) error {
	l, err := t.db.SCard(t.ctx, fmt.Sprint(id)).Result()
	if err != nil {
		return t.ManageError(err)
	}

	// Si hay un limite de tokens por usuario
	// que verifique si se ha alcanzado
	if env.Config.MaxTokenPerUser != -1 {
		if l >= int64(env.Config.MaxTokenPerUser) {
			return ErrMaxTokens
		}
	}

	err = t.db.SAdd(t.ctx, fmt.Sprint(id), token).Err()
	return t.ManageError(err)
}

func (t *tokenRedis) RemoveToken(id uint, token string) error {
	err := t.db.SRem(t.ctx, fmt.Sprint(id), token).Err()
	return t.ManageError(err)
}

func (t *tokenRedis) GetTokens(id uint) ([]string, error) {
	tokens, err := t.db.SMembers(t.ctx, fmt.Sprint(id)).Result()
	if err != nil {
		return []string{}, t.ManageError(err)
	}

	return tokens, nil
}

func (t *tokenRedis) ExistsToken(id uint, token string) (bool, error) {
	ok, err := t.db.SIsMember(t.ctx, fmt.Sprint(id), token).Result()
	if err != nil {
		return false, t.ManageError(err)
	}
	return ok, nil
}

func (t *tokenRedis) RemoveUserTokens(id uint) error {
	err := t.db.Del(t.ctx, fmt.Sprint(id)).Err()
	return t.ManageError(err)
}

func HandlerTokenError(err error) error {
	if err == ErrTokenNotFound {
		return echo.NewHTTPError(http.StatusNotFound, ErrTokenNotFound.Error())
	}

	if err == ErrMaxTokens {
		return echo.NewHTTPError(http.StatusForbidden, ErrMaxTokens.Error())
	}

	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}
