package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"fmt"
)


///////////////////////////////////////////////////////////////////////////////
// specs-zone

func GetSpecsZone(fr common.FilterReq) ([]common.SpecsZoneInfo, int, error) {

	var szList []common.SpecsZoneInfo
	var sqlcmd string
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get specs-zone list")

	code, _, err := common.GetClassInfo(fr.ClassLevel)
	if err != nil {
		zaps.Errorf("get class info failed: %v", err)
		return szList, 0, err
	}

	sqla := "SELECT t.size_dim, SUM(t.inv_qty) as SALECOUNT, " +
		"SUM(t.gross_val) AS SALEINCOME, SUM(t.netval_inv) AS SUMNET, "+
		"SUM(t.cost) AS SUMCOST FROM "
	sqlb := "SELECT a.material, a.gross_val, a.inv_qty, " +
		"a.gross_val/a.inv_qty AS unit_price, a.netval_inv, a.cost, " +
		"a.zlsjg, b.bic_zklad2, b.zklad2txt, b.bic_zklasse_d, " +
		"b.zklasse_dtxt, b.bic_zrpa_mtl, b.zrpa_mtltxt, b.size_dim "
	sqld := "FROM bill_zsd a " +
		"LEFT OUTER JOIN material b ON a.material = b.material "
	sqle := fmt.Sprintf("WHERE a.calday >= '%s' AND a.calday <= '%s'",
			fr.BeginDate, fr.EndDate)
	sqlf := fmt.Sprintf("AND a.plant = '%s'", fr.OrgCode)
	sqlg := fmt.Sprintf("AND b.%s = '%s'", code, fr.ClassCode)
	sqlh := "GROUP BY t.size_dim ORDER BY SALECOUNT DESC"

	if fr.OrgLevel == common.ORG_LEVEL_HEAD {
		sqlcmd = fmt.Sprintf("%s(%s %s %s %s) as t %s",
				sqla, sqlb, sqld, sqle, sqlg, sqlh)
	} else if fr.OrgLevel == common.ORG_LEVEL_BRANCH {

	} else if fr.OrgLevel == common.ORG_LEVEL_PLANT {
		sqlcmd = fmt.Sprintf("%s(%s %s %s %s %s) as t %s",
				sqla, sqlb, sqld, sqle, sqlf, sqlg, sqlh)
	}

	zaps.Info("sql cmd: ", sqlcmd)
	rows, err = db.Query(sqlcmd)

	if err != nil {
		zaps.Error("db query failed: ", err)
		return szList, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var sz common.SpecsZoneInfo
		err := rows.Scan(&sz.SpecsZone, &sz.SalesCount, &sz.SalesIncome,
				&sz.SalesNet, &sz.SalesCost)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> SpecsZone: ", sz.SpecsZone)
			zaps.Debug(">>> SalesCount: ", sz.SalesCount)
			zaps.Debug(">>> SalesIncome: ", sz.SalesIncome)
			zaps.Debug(">>> SalesNet: ", sz.SalesNet)
			zaps.Debug(">>> SalesCost: ", sz.SalesCost)

			szList = append(szList, sz)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return szList, count, err
	}

	zaps.Info("<<< get specs-zone done with count(%d)", count)

	return szList, count, err
}


///////////////////////////////////////////////////////////////////////////////
// specs zone rank

func GetSpecsZoneRank(fr common.FilterReq, zone string) ([]common.MatlRankInfo, int, error) {

	var mList []common.MatlRankInfo
	var sqlcmd string
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get specs-zone rank list")

	code, _, err := common.GetClassInfo(fr.ClassLevel)
	if err != nil {
		zaps.Errorf("get class info failed: %v", err)
		return mList, 0, err
	}

	sqla := "SELECT t.material, t.materialtxt, SUM(t.inv_qty) as SALECOUNT, " +
		"SUM(t.gross_val) AS SALEINCOME, SUM(t.netval_inv) AS SUMNET, " +
		"SUM(t.cost) AS SUMCOST FROM "
	sqlb := "SELECT a.material, b.materialtxt, a.gross_val, a.inv_qty, " +
		"a.gross_val/a.inv_qty AS unit_price, a.netval_inv, a.cost, " +
		"a.zlsjg, b.bic_zklad2, b.zklad2txt, b.bic_zklasse_d, " +
		"b.zklasse_dtxt, b.bic_zrpa_mtl, b.zrpa_mtltxt, b.size_dim "
	sqld := "FROM bill_zsd a " +
		"LEFT OUTER JOIN material b ON a.material = b.material "
	sqle := fmt.Sprintf("WHERE a.calday >= '%s' AND a.calday <= '%s'",
			fr.BeginDate, fr.EndDate)
	sqlf := fmt.Sprintf("AND a.plant = '%s'", fr.OrgCode)
	sqlg := fmt.Sprintf("AND b.%s = '%s'", code, fr.ClassCode)
	sqlh := fmt.Sprintf("WHERE t.size_dim = '%s' ", zone)
	sqli := "GROUP BY t.material ORDER BY SALECOUNT DESC"

	if fr.OrgLevel == common.ORG_LEVEL_HEAD {
		sqlcmd = fmt.Sprintf("%s(%s %s %s %s) as t %s %s",
			sqla, sqlb, sqld, sqle, sqlg, sqlh, sqli)
	} else if fr.OrgLevel == common.ORG_LEVEL_BRANCH {

	} else if fr.OrgLevel == common.ORG_LEVEL_PLANT {
		sqlcmd = fmt.Sprintf("%s(%s %s %s %s %s) as t %s %s",
			sqla, sqlb, sqld, sqle, sqlf, sqlg, sqlh, sqli)
	}

	zaps.Info("sql cmd: ", sqlcmd)
	rows, err = db.Query(sqlcmd)

	if err != nil {
		zaps.Error("db query failed: ", err)
		return mList, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var mr common.MatlRankInfo
		err := rows.Scan(&mr.Material, &mr.MaterialTxt, &mr.SalesCount,
				&mr.SalesIncome, &mr.SalesNet, &mr.SalesCost)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> Material: ", mr.Material)
			zaps.Debug(">>> MaterialTxt: ", mr.MaterialTxt)
			zaps.Debug(">>> SalesCount: ", mr.SalesCount)
			zaps.Debug(">>> SalesIncome: ", mr.SalesIncome)
			zaps.Debug(">>> SalesNet: ", mr.SalesNet)
			zaps.Debug(">>> SalesCost: ", mr.SalesCost)

			mList = append(mList, mr)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return mList, 0, err
	}

	zaps.Info("<<< get specs-zone rank done with count(%d)", count)

	return mList, count, err
}


