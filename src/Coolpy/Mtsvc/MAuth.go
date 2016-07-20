package Mtsvc

import (
	"github.com/surgemq/surgemq/auth"
	"Coolpy/Account"
)

type Authenticator interface {
	Authenticate(id string, cred interface{}) error
}

type Manager struct {
}

func (this *Manager) Authenticate(id string, cred interface{}) error {
	u, err := Account.GetByUkey(id)
	if err != nil {
		return auth.ErrAuthFailure
	}
	if u.Ukey != id {
		return auth.ErrAuthFailure
	}
	return nil
}