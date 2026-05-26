package templatestring

import (
	"errors"
	"os"
	"strings"
)

type Plugin interface {
	ProcessToken(token string) (rv string, isProcessed bool, err error)
}

type delegatePlugin struct {
	cb func(token string) (rv string, isProcessed bool, err error)
}

func NewDelegatePlugin(cb func(token string) (string, bool, error)) *delegatePlugin {
	return &delegatePlugin{
		cb: cb,
	}
}

func (d *delegatePlugin) ProcessToken(token string) (rv string, isProcessed bool, err error) {
	return d.cb(token)
}

type envPlugin struct{}

func NewEnvPlugin() *envPlugin {
	return &envPlugin{}
}

func (d *envPlugin) ProcessToken(token string) (string, bool, error) {
	if strings.HasPrefix(token, "env:") {
		token = token[4:]
		if len(token) == 0 {
			return "", true, errors.New("env variable should have name")
		}
		token = os.Getenv(token)
		return token, true, nil
	}

	return token, false, nil
}
