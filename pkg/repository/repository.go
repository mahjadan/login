package repository

import (
	"context"
	"errors"
)

var ErrAlreadyExists = errors.New("user already exists")
var ErrNotFound = errors.New("user not found")

type UsersRepository interface {
	Get(ctx context.Context, username string) (string, error)
	Set(ctx context.Context, username, password string) error
}
