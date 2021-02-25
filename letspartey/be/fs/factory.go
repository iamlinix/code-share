package fs

import (
	"errors"
	"sync"

	"iamlinix.com/partay/logger"
)

var _backend BackendType = BackendUnknown
var _inst FSInterface = nil
var _baseDir string = ""

func Get() FSInterface {
	return _inst
}

func SetFsBackend(backend BackendType, baseDir string) error {
	_backend = backend
	_baseDir = baseDir
	var err error
	if _inst, err = ProduceFileSystem(); err != nil {
		logger.Errorf("error inner producing file system: %v", err)
		return err
	}

	return nil
}

func ProduceFileSystem() (FSInterface, error) {
	return InstantiateFileSystem(_backend)
}

func InstantiateFileSystem(backend BackendType) (FSInterface, error) {
	switch backend {
	case BackendLocal:
		return &LocalFS{
			Cache: make(map[string]*CacheFile),
			Mux:   sync.Mutex{},
		}, nil
	default:
		logger.Errorf("unimplemented file system: %d", backend)
		break
	}
	return nil, errors.New("Not implemented")
}
