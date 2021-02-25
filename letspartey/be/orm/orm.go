package orm

import (
	"errors"

	"iamlinix.com/partay/db"
	"iamlinix.com/partay/logger"
	"iamlinix.com/partay/orm/models"
)

var _inst db.DBInterface = nil

func InitOrmWithDbInst(inst db.DBInterface) error {
	if _inst != nil {
		logger.Errorf("Orm db instance can only be initiated once")
		return errors.New("Orm db instance already initiated")
	}

	_inst = inst
	return nil
}

func InitOrm(backend db.BackendType, driver string, username string, password string, database string,
	args map[string]string, pool map[string]string) error {
	if _inst != nil {
		logger.Errorf("Orm db instance can only be initiated once")
		return errors.New("Orm db instance already initiated")
	}

	var err error
	_inst, err = db.InstantiateDatabase(backend)
	if err != nil {
		logger.Errorf("error instantiating db backend: %v", err)
		return err
	}

	if err = _inst.Connect(driver, username, password, database, args, pool); err != nil {
		logger.Errorf("error connecting database: %v", err)
		_inst = nil
		return err
	}

	return nil
}

func InitAllTables() error {
	switch db.Get().(type) {
	case *db.Mysql:
		inst := db.Get().(*db.Mysql)
		inst.DB().Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&models.WxOpen{})
		inst.DB().Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&models.WxUser{})
		inst.DB().Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&models.WxDevice{})
		inst.DB().Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&models.User{})
		inst.DB().Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&models.Activity{})
		inst.DB().Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&models.ActivityImage{})
		inst.DB().Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&models.Post{})
		inst.DB().Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&models.PostImage{})
		break
	default:
		return errors.New("Unimplemented database backend")
	}
	return nil
}
