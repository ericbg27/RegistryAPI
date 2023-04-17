package token

import "time"

type Maker interface {
	CreateToken(username string, duration time.Duration) (string, error)
	VerifyToken(tokenToVerify string) (*Payload, error)
}
