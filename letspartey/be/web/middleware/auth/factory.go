package auth

import (
	"errors"

	"iamlinix.com/partay/logger"
)

var _method AuthMethod = AuthNone
var _inst AuthInterface = &BaseAuth{}
var WebAuthEnabled bool = false

func Get() AuthInterface {
	return _inst
}

func SetAuthMethod(method AuthMethod) error {
	_method = method
	if method > AuthNone {
		WebAuthEnabled = true
	}
	var err error
	if _inst, err = ProduceAuth(); err != nil {
		logger.Errorf("error inner producing auth: %v", err)
		return err
	}

	return nil
}

func ProduceAuth() (AuthInterface, error) {
	return InstantiateAuth(_method)
}

func InstantiateAuth(method AuthMethod) (AuthInterface, error) {
	switch method {
	case AuthJwt:
		return &JwtAuth{}, nil
	default:
		logger.Errorf("unimplemented auth: %d", method)
		break
	}
	return nil, errors.New("Not implemented")
}
