package epp

import (
	"github.com/domainr/epp2/schema/epp"
)

func Greeting(cfg *Config) (*epp.Greeting, error) {
	return &epp.Greeting{}, nil // TODO
}

func Command(cfg *Config, action epp.Action, extensions ...epp.Extension) (*epp.Command, error) {
	return &epp.Command{
		Action:              action,
		Extensions:          extensions,
		ClientTransactionID: cfg.TransactionID(),
	}, nil
}

func LoginCommand(cfg *Config, clientID, password string, newPassword *string) (*epp.Command, error) {
	return Command(cfg, &epp.Login{
		ClientID:    clientID,
		Password:    password,
		NewPassword: newPassword,
		Options: epp.Options{
			Version: epp.Version,
		},
	})
}

func LogoutCommand(cfg *Config) (*epp.Command, error) {
	return Command(cfg, &epp.Logout{})
}

func ErrorResponse(cfg *Config, err error) *epp.Response {
	return nil // TODO
}
