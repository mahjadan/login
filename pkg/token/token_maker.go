package token

import (
	"time"
)

type Maker interface {
	Create(userToken UserToken, duration time.Duration) (string, error)
	Verify(tokenStr string) (*JWT, error)
}
