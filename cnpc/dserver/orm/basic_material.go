package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
)


func GetMatlNameByCode(code string) (string, bool, error) {

	var name string
	find := false

	zaps.Info(">>> get material name by code: ", code)

	rows, err := db.Query("SELECT materialtxt FROM material " +
			"WHERE material = ?", code)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return name, find, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			zaps.Error("query error: ", err)
			return name, find, err
		}

		find = true
		break
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return name, find, err
	}

	zaps.Info("<<< get material name by code done")

	return name, find, err
}


func GetMaterialByCode(code string) (common.Material, bool, error) {

	var res common.Material
	var rows *sql.Rows
	find := false

	rows, err := db.Query("SELECT material, materialtxt " +
			"FROM material WHERE material = ?", code)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return res, find, err
	}

	defer rows.Close()

	for rows.Next() {
		var m common.Material
		err := rows.Scan(&m.LongCode, &m.Name)
		if err != nil {
			zaps.Error("query error: ", err)
			return m, find, err
		}

		m.ShortCode = common.GetShortMatlCode(m.LongCode)
		res = m
		find = true
		break
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return res, find, err
	}

	zaps.Info("<<< get material by code done")

	return res, find, err
}


