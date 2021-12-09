package repository

import (
	"context"
	enc "github.com/mahjadan/login/pkg/encrypt"
	"github.com/pkg/errors"
)

func NewInMemory() UsersRepository {
	return InMemory{
		users: make(map[string]string),
	}
}

type InMemory struct {
	users map[string]string
}

func (i InMemory) Get(ctx context.Context, username string) (string, error) {
	if hash, ok := i.users[username]; ok {
		return hash, nil
	}
	return "", ErrNotFound
}

func (i InMemory) Set(ctx context.Context, username, password string) error {
	if _, ok := i.users[username]; ok {
		return ErrAlreadyExists
	}
	hashPassword, err := enc.HashPassword(password)
	if err != nil {
		return errors.Wrap(err, "hash err")
	}
	i.users[username] = hashPassword
	return nil
}
