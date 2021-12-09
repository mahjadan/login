package service

import (
	"context"
	"errors"
	enc "github.com/mahjadan/login/pkg/encrypt"
	"github.com/mahjadan/login/pkg/repository"
)

var ErrInvalidUserOrPassword = errors.New("invalid username or password")

type UserService interface {
	Login(ctx context.Context, username, password string) error
	Register(ctx context.Context, username, password string) error
}

type service struct {
	repo repository.UsersRepository
}

// todo generate token on successfull login
func (s service) Login(ctx context.Context, username, password string) error {
	hash, err := s.repo.Get(ctx, username)

	if errors.Is(err, repository.ErrNotFound) {
		return ErrInvalidUserOrPassword
	}
	if err != nil {
		return err
	}

	if !enc.PasswordMatch(password, hash) {
		return ErrInvalidUserOrPassword
	}
	return nil
}

func (s service) Register(ctx context.Context, username, password string) error {
	err := s.repo.Set(ctx, username, password)
	if err != nil {
		return err
	}
	return nil
}

func NewService(usersRepository repository.UsersRepository) UserService {
	return &service{
		repo: usersRepository,
	}
}
