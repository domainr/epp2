package epp

import (
	"github.com/domainr/epp2/schema/epp"
)

func Greeting(cfg *Config) (epp.Body, error) {
	return &epp.Greeting{}, nil // TODO
}

func Command(cfg *Config, action epp.Action, extensions ...epp.Extension) (epp.Body, error) {
	return &epp.Command{
		Action:              action,
		Extensions:          extensions,
		ClientTransactionID: cfg.TransactionID(),
	}, nil
}

func Login(cfg *Config, clientID, password string, newPassword *string) (epp.Body, error) {
	return Command(cfg, &epp.Login{
		ClientID:    clientID,
		Password:    password,
		NewPassword: newPassword,
		Options: epp.Options{
			Version: epp.Version,
		},
	})
}

func Logout(cfg *Config) (epp.Body, error) {
	return Command(cfg, &epp.Logout{})
}

func ErrorResponse(cfg *Config, err error) epp.Body {
	return nil // TODO
}
