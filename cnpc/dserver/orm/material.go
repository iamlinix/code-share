package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"errors"
	"fmt"
)


///////////////////////////////////////////////////////////////////////////////
// SKU

func GetClassSKU(fr common.FilterReq, all int, page int) ([]common.SKUClassInfo,
	int, error) {

	var classList []common.SKUClassInfo
	var sqld, sqlcmd string
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get class sku list with page ", page)

	code, text, err := common.GetClassInfo(fr.ClassLevel)
	if err != nil {
		zaps.Errorf("get class info failed: %v", err)
		return classList, 0, err
	}

	sqla := fmt.Sprintf("SELECT t.%s, any_value(t.%s), " +
		"count(0) as SKUCOUNT FROM ", code, text)
	sqlb := "SELECT DISTINCT(a.material), b.bic_zklad2, b.zklad2txt, " +
		"b.bic_zklasse_d, b.zklasse_dtxt, b.bic_zrpa_mtl, " +
		"b.zrpa_mtltxt " +
		"FROM zinv a " +
		"LEFT OUTER JOIN material b ON a.material = b.material "
	sqlc := fmt.Sprintf("WHERE a.calday >= '%s' AND a.calday <= '%s' ",
			fr.BeginDate, fr.EndDate)
	if fr.ClassLevel == common.CLASS_LEVEL_MAIN {
		sqld = ""
	} else if fr.ClassLevel == common.CLASS_LEVEL_MID {
		sqld = fmt.Sprintf("AND b.bic_zklad2 = %s", fr.ClassCode)
	} else if fr.ClassLevel == common.CLASS_LEVEL_SUB {
		sqld = fmt.Sprintf("AND b.bic_zklasse_d = %s", fr.ClassCode)
	}
	sqle := fmt.Sprintf("AND a.plant = '%s'", fr.OrgCode)
	sqlf := fmt.Sprintf("GROUP BY t.%s ORDER BY SKUCOUNT DESC", code)

	if all == 1 {
		if fr.OrgLevel == common.ORG_LEVEL_HEAD {
			sqlcmd = fmt.Sprintf("%s(%s %s %s) as t %s",
					sqla, sqlb, sqlc, sqld, sqlf)
		} else if fr.OrgLevel == common.ORG_LEVEL_BRANCH {

		} else if fr.OrgLevel == common.ORG_LEVEL_PLANT {
			sqlcmd = fmt.Sprintf("%s(%s %s %s %s) as t %s",
					sqla, sqlb, sqlc, sqld, sqle, sqlf)
		} else {
			zaps.Error("invalid org level: ", fr.OrgLevel)
			return classList, 0, errors.New(common.ERR_MSG_INVALID_ORG)
		}

		zaps.Info(">>> sql cmd: ", sqlcmd)
		rows, err = db.Query(sqlcmd)
	} else {

	}

	if err != nil {
		zaps.Error("db query failed: ", err)
		return classList, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var sc common.SKUClassInfo
		err := rows.Scan(&sc.ClassCode, &sc.ClassName, &sc.SKUCount)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> ClassCode: ", sc.ClassCode)
			zaps.Debug(">>> ClassName: ", sc.ClassName)
			zaps.Debug(">>> SKUCount: ", sc.SKUCount)

			classList = append(classList, sc)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return classList, count, err
	}

	zaps.Info("<<< get class sku done with ", count)

	return classList, count, err
}


///////////////////////////////////////////////////////////////////////////////
// inventory

func GetClassInv(fr common.FilterReq, all int,
	page int) ([]common.InvClassInfo, int, error) {

	var classList []common.InvClassInfo
	var sqld, sqlcmd string
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get class inv list with page ", page)

	code, text, err := common.GetClassInfo(fr.ClassLevel)
	if err != nil {
		zaps.Errorf("get class info failed: %v", err)
		return classList, 0, err
	}

	sqla := fmt.Sprintf("SELECT t.%s, any_value(t.%s), " +
		"SUM(t.bic_zinvsl) as SUMSL, SUM(t.bic_zinvcost) as SUMCOST " +
		"FROM ", code, text)
	sqlb := "SELECT a.material, a.bic_zinvsl, a.bic_zinvcost, " +
		"b.bic_zklad2, b.zklad2txt, b.bic_zklasse_d, b.zklasse_dtxt, " +
		"b.bic_zrpa_mtl, b.zrpa_mtltxt " +
		"FROM zinv a " +
		"LEFT OUTER JOIN material b ON a.material = b.material "
	sqlc := fmt.Sprintf("WHERE a.calday = '%s'", fr.EndDate)
	if fr.ClassLevel == common.CLASS_LEVEL_MAIN {
		sqld = ""
	} else if fr.ClassLevel == common.CLASS_LEVEL_MID {
		sqld = fmt.Sprintf("AND b.bic_zklad2 = %s", fr.ClassCode)
	} else if fr.ClassLevel == common.CLASS_LEVEL_SUB {
		sqld = fmt.Sprintf("AND b.bic_zklasse_d = %s", fr.ClassCode)
	}
	sqle := fmt.Sprintf("AND a.plant = '%s'", fr.OrgCode)
	sqlf := fmt.Sprintf("GROUP BY t.%s ORDER BY SUMSL DESC", code)

	if all == 1 {
		if fr.OrgLevel == common.ORG_LEVEL_HEAD {
			sqlcmd = fmt.Sprintf("%s(%s %s %s) as t %s",
					sqla, sqlb, sqlc, sqld, sqlf)
		} else if fr.OrgLevel == common.ORG_LEVEL_BRANCH {

		} else if fr.OrgLevel == common.ORG_LEVEL_PLANT {
			sqlcmd = fmt.Sprintf("%s(%s %s %s %s) as t %s",
					sqla, sqlb, sqlc, sqld, sqle, sqlf)
		} else {
			zaps.Error("invalid org level: ", fr.OrgLevel)
			return classList, 0, errors.New(common.ERR_MSG_INVALID_ORG)
		}

		zaps.Info("sql cmd: ", sqlcmd)
		rows, err = db.Query(sqlcmd)
	} else {

	}

	if err != nil {
		zaps.Error("db query failed: ", err)
		return classList, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var sc common.InvClassInfo
		err := rows.Scan(&sc.ClassCode, &sc.ClassName, &sc.InvCount,
				&sc.InvCost)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {

			classList = append(classList, sc)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return classList, count, err
	}

	zaps.Info("<<< get class inv done with count %d", count)

	return classList, count, err
}


