package MAuth

import (
	"github.com/surgemq/surgemq/auth"
)

type Authenticator interface {
	Authenticate(id string, cred interface{}) error
}

type Manager struct {
}

func (this *Manager) Authenticate(id string, cred interface{}) error {
	if id == "111" {
		return nil
	}
	return auth.ErrAuthFailure
}