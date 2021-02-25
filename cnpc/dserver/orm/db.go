package orm

import (
	"database/sql"
	"fmt"

	"cnpc.com.cn/cnpc/dserver/zaps"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var err int

func DBConnect(driver, addr, port, user, pass, dbname string) error {

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, addr, port, dbname)
	zaps.Info(">>> connecting to db: ", addr)

	_db, err := sql.Open(driver, connStr)
	if err != nil {
		zaps.Fatalf("<<< sql open failed")
		return err
	}

	err = _db.Ping()
	if err != nil {
		zaps.Fatalf("<<< sql open db failed: ", addr)
	}

	zaps.Info("<<< open database success")
	db = _db
	db.SetMaxIdleConns(32)
	db.SetMaxOpenConns(0)

	return nil
}

func GetDB() *sql.DB {
	return db
}
