package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"

	"errors"
	"fmt"
)


//branch org
func GetBranchOrgList() ([]common.OrgInfo, int, error) {

	var oList []common.OrgInfo
	var rows *sql.Rows
	var err  error
	var count int

	zaps.Info(">>> get branch org list")

	rows, err = db.Query("SELECT bic_zrpa_lcit, any_value(bic_ztxt_dms) "+
			"FROM zaplant_xy " +
			"WHERE bic_zrpa_lcit >= 'A13' AND bic_zrpa_lcit <= 'A16'" +
			"GROUP BY bic_zrpa_lcit")

	if err != nil {
		zaps.Error("GetBranchOrgList db query failed: ", err)
		return oList, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var org common.OrgInfo
		err := rows.Scan(&org.OrgCode, &org.OrgText)
		if err != nil {
			zaps.Error("GetBranchOrgList query error: ", err)
			return oList, count, err
		}

		oList = append(oList, org)
		count++
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return oList, count, err
	}

	zaps.Info("<<< get branch org list done")

	return oList, count, err
}


//palnt org
func GetPlantOrgList() ([]common.OrgInfo, int, error) {

	var oList []common.OrgInfo
	var rows *sql.Rows
	var err  error
	var count int

	zaps.Info(">>> get plant org list")

	rows, err = db.Query("SELECT bic_zaplant, any_value(bic_ztxt_jyz) " +
			"FROM zaplant_xy " +
			"WHERE bic_zrpa_lcit >= 'A13' AND bic_zrpa_lcit <= 'A16'" +
			"GROUP BY bic_zaplant")

	if err != nil {
		zaps.Error("GetPlantOrgList db query failed: ", err)
		return oList, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var org common.OrgInfo
		err := rows.Scan(&org.OrgCode, &org.OrgText)
		if err != nil {
			zaps.Error("GetPlantOrgList query error: ", err)
			return oList, count, err
		}

		oList = append(oList, org)
		count++
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return oList, count, err
	}

	zaps.Info("<<< get plant org list done")

	return oList, count, err
}


func GetOrgTextByLevelCode(org *common.OrgInfo) (bool, error) {

	var sql string
	find := false

	zaps.Info(">>> get org text by level code: ", org.OrgCode)

	level := org.OrgLevel
	if level == common.ORG_LEVEL_BRANCH {
		sql = fmt.Sprintf("SELECT DISTINCT bic_ztxt_dms " +
				"FROM zaplant_xy " +
				"WHERE bic_zrpa_lcit = '%s'", org.OrgCode)
	} else if level == common.ORG_LEVEL_PLANT {
		sql = fmt.Sprintf("SELECT DISTINCT bic_ztxt_jyz " +
				"FROM zaplant_xy " +
				"WHERE bic_zaplant = '%s'", org.OrgCode)
	} else if level == common.ORG_LEVEL_HEAD {
		return true, nil
	} else {
		zaps.Error("invalid org level: ", org.OrgLevel)
		return find, errors.New(common.ERR_MSG_INVALID_ORG)
	}

	rows, err := db.Query(sql)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return find, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&org.OrgText)
		if err != nil {
			zaps.Error("query error: ", err)
			return find, err
		}

		find = true;
		break;
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return find, err
	}

	zaps.Info("<<< get org text by level code done")

	return find, err
}



