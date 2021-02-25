package db

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"iamlinix.com/partay/logger"
)

type Mysql struct {
	db *gorm.DB
}

type MysqlStatement struct {
	db          *gorm.DB
	queryFormat string
}

func (self *Mysql) DB() *gorm.DB {
	return self.db
}

func (self *Mysql) Connect(driver string, username string, password string, database string,
	args map[string]string, pool map[string]string) error {
	if len(driver) == 0 || len(username) == 0 || len(password) == 0 || len(database) == 0 {
		logger.Errorf("invalid connection parameters: %s, %s, %s, %s", driver, username, password, database)
		return errors.New("Invalid connection parameters")
	}

	connStr := fmt.Sprintf("%s:%s@/%s?parseTime=True", username, password, database)
	if len(args) > 0 {
		for k, v := range args {
			if k == "parseTime" {
				continue
			}

			if k == "charset" && v == "utf8" {
				v = "utf8mb4"
			}

			connStr += fmt.Sprintf("&%s=%s", k, v)
		}
	}

	db, err := gorm.Open(driver, connStr)
	if err != nil {
		logger.Errorf("failed to connect to database: %v", err)
		return err
	}
	self.db = db

	if len(pool) > 0 {
		sqlDb := db.DB()
		if val, ok := pool["maxOpen"]; ok {
			if n, err := strconv.Atoi(val); err == nil {
				sqlDb.SetMaxOpenConns(n)
			} else {
				logger.Errorf("invalid maxOpen value: %s", val)
			}
		}

		if val, ok := pool["maxIdle"]; ok {
			if n, err := strconv.Atoi(val); err == nil {
				sqlDb.SetMaxIdleConns(n)
			} else {
				logger.Errorf("invalid maxIdle value: %s", val)
			}
		}

		if val, ok := pool["maxLife"]; ok {
			if n, err := strconv.Atoi(val); err == nil {
				sqlDb.SetConnMaxLifetime(time.Duration(n) * time.Second)
			} else {
				logger.Errorf("invalid maxLife value: %s", val)
			}
		}
	}

	if err = db.DB().Ping(); err != nil {
		logger.Errorf("database connection ping error: %v", err)
		db.Close()
		return err
	}

	logger.Info("database connected")

	return nil
}

func (self *Mysql) Disconnect(params interface{}) error {
	var err error
	if err = self.db.Close(); err != nil {
		logger.Errorf("error closing database: %v", err)
	}
	return err
}

func (self *Mysql) PrepareStatement(query string) (StatementInterface, error) {
	if len(query) == 0 {
		logger.Errorf("invalid prepared mysql statement: %s", query)
		return nil, errors.New("Invalid prepared statement")
	}

	return &MysqlStatement{
		queryFormat: query,
		db:          self.db,
	}, nil
}

func (self *MysqlStatement) Execute(args ...interface{}) ([]map[string]interface{}, error) {
	if self.db == nil {
		logger.Errorf("db not ready")
		return nil, errors.New("Db not ready")
	}

	if len(self.queryFormat) == 0 {
		logger.Errorf("query not ready")
		return nil, errors.New("Query not ready")
	}

	rows, err := self.db.Raw(self.queryFormat, args...).Rows()
	if err != nil {
		logger.Errorf("error executing prepared query: %v", err)
		return nil, err
	}

	defer rows.Close()
	columns, err := rows.ColumnTypes()
	if err != nil {
		logger.Errorf("error getting result columns: %v", err)
		return nil, err
	}

	rowValues := make([]interface{}, len(columns))
	var results []map[string]interface{}

	for rows.Next() {
		row := map[string]interface{}{}
		for i, col := range columns {
			row[col.Name()] = reflect.New(col.ScanType()).Interface()
			rowValues[i] = row[col.Name()]
		}

		err = rows.Scan(rowValues...)
		if err != nil {
			logger.Errorf("error scaning results: %v", err)
			return nil, err
		}

		results = append(results, row)
	}

	return results, nil
}

func (*MysqlStatement) Close() error {
	logger.Infof("mysql prepared statement closed")
	return nil
}

func (self *Mysql) Execute(query string, args ...interface{}) ([]map[string]interface{}, error) {
	if self.db == nil {
		logger.Errorf("db not ready")
		return nil, errors.New("Db not ready")
	}

	rows, err := self.db.Raw(query, args...).Rows()
	if err != nil {
		logger.Errorf("error executing query: %v", err)
		return nil, err
	}

	defer rows.Close()
	columns, err := rows.ColumnTypes()
	if err != nil {
		logger.Errorf("error getting result columns: %v", err)
		return nil, err
	}

	rowValues := make([]interface{}, len(columns))
	var results []map[string]interface{}

	for rows.Next() {
		row := map[string]interface{}{}
		for i, col := range columns {
			row[col.Name()] = reflect.New(col.ScanType()).Interface()
			rowValues[i] = row[col.Name()]
		}

		err = rows.Scan(rowValues...)
		if err != nil {
			logger.Errorf("error scaning results: %v", err)
			return nil, err
		}

		results = append(results, row)
	}

	return results, nil
}

func (self *Mysql) ExecuteStruct(generator GeneratorFunc, query string, args ...interface{}) ([]interface{}, error) {
	if self.db == nil {
		logger.Errorf("db not ready")
		return nil, errors.New("Db not ready")
	}

	rows, err := self.db.Raw(query, args...).Rows()
	if err != nil {
		logger.Errorf("error executing query: %v", err)
		return nil, err
	}

	defer rows.Close()
	var results []interface{}

	for rows.Next() {
		dst := generator()
		err = self.db.ScanRows(rows, dst)
		if err != nil {
			logger.Errorf("error scan rows: %v", err)
			return nil, err
		}
		logger.Info(dst)
		results = append(results, dst)
	}

	return results, nil
}

func (self *Mysql) Create(model ModelInterface) error {
	self.db.Set("gorm:table_options", "CHARSET=utf8mb4").Table(model.TableName()).Create(model)
	return nil
}

func (self *Mysql) Update(model ModelInterface) error {
	self.db.Set("gorm:table_options", "CHARSET=utf8mb4").Table(model.TableName()).Save(model)
	return nil
}

func (self *Mysql) Delete(model ModelInterface) error {
	self.db.Set("gorm:table_options", "CHARSET=utf8mb4").Table(model.TableName()).Delete(model)
	return nil
}
