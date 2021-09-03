package token

import (
	"errors"
	"time"
)

type Payload struct {
	Identity int   `json:"identity"`
	Scope    string `json:"scope"`
	Type     string `json:"type"`
	Exp      int64  `json:"exp"`
}

func (p Payload) Valid() error {
	if time.Now().Unix() >= p.Exp {
		return errors.New("jwt is expired")
	}
	return nil
}

func NewPayload(identity int, scope string, jwtType string, exp int64) *Payload {
	return &Payload{
		Identity: identity,
		Scope:    scope,
		Type:     jwtType,
		Exp:      exp,
	}
}
