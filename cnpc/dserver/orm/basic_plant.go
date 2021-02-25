package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

func GetPlantNameByCode(code string) (string, bool, error) {

	var name string
	find := false

	zaps.Info(">>> get plant name by code: ", code)

	rows, err := db.Query("SELECT bic_ztxt_jyz FROM zaplant_xy "+
		"WHERE plant = ?", code)
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

	zaps.Info("<<< get plant name by code done")

	return name, find, err
}

func GetPlantByCode(code string) (common.Plant, bool, error) {

	var res common.Plant
	var rows *sql.Rows
	find := false

	rows, err := db.Query("SELECT bic_ztxt_jyz, bic_ztxt_dms, "+
		"bic_ztxt_zcxz, bic_ztxt_dfl, "+
		"bic_zt_type, bic_ztxt_bld, posx, posy "+
		"FROM zaplant_xy WHERE bic_zaplant = ?", code)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return res, find, err
	}

	defer rows.Close()

	for rows.Next() {
		var p common.Plant
		err := rows.Scan(&p.PlantName, &p.BranchName, &p.ZtxtZcxz,
			&p.ZtxtDfl, &p.ZtType,
			&p.ZtxtBld, &p.PosX, &p.PosY)
		if err != nil {
			zaps.Error("query error: ", err)
			return p, find, err
		}

		p.Plant = code

		res = p
		find = true
		break
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return res, find, err
	}

	zaps.Info("<<< get plant by code done")

	return res, find, err
}
