package db

import (
	"errors"

	"iamlinix.com/partay/logger"
)

var _backend BackendType = BackendUnknown
var _inst DBInterface = nil

func Get() DBInterface {
	return _inst
}

func SetDbBackend(backend BackendType) error {
	_backend = backend
	var err error
	if _inst, err = ProduceDatabase(); err != nil {
		logger.Errorf("error inner producing db backend: %v", err)
		return err
	}

	return nil
}

func ProduceDatabase() (DBInterface, error) {
	return InstantiateDatabase(_backend)
}

func InstantiateDatabase(backend BackendType) (DBInterface, error) {
	switch backend {
	case BackendMysql:
		return &Mysql{}, nil
	default:
		logger.Errorf("Unimplemented database: %d", backend)
		break
	}
	return nil, errors.New("Not implemented")
}
