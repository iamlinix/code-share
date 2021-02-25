package orm

import (
	"fmt"
	"time"

	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

///////////////////////////////////////////////////////////////////////////////
//
/*
 * get sales basic metric, from material view
 * input:  material
 * output: sum metric <one item>
 */
func GetSalesMetricByMaterial(fr common.FilterReq) (common.SalesMetric,
	bool, error) {

	var sm common.SalesMetric
	var rows *sql.Rows
	var find bool

	zaps.Info(">>> get sales metric by material")

	sqlcmd := fmt.Sprintf("SELECT sum(a.inv_qty) AS INVQTY, "+
		"sum(a.cost) AS COST, sum(a.netval_inv) AS NETVALINV, "+
		"sum(a.gross_val) AS GROSSVAL "+
		"FROM bill_zsd a "+
		"WHERE a.calday BETWEEN '%s' AND '%s' "+
		"AND a.plant != 'A0A1' "+
		"AND a.material = '%s'",
		fr.BeginDate, fr.EndDate, fr.Material)

	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return sm, find, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&sm.InvQty, &sm.Cost, &sm.NetvalInv,
			&sm.GrossVal)
		if err != nil {
			zaps.Error("query error: ", err)

		} else {
			sm.GrossProfit = sm.NetvalInv - sm.Cost
			if sm.NetvalInv > 0 {
				sm.GrossMargin = sm.GrossProfit * 100.0 / sm.NetvalInv
			}
			find = true
			break
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return sm, find, err
	}

	zaps.Info("<<< get sales metric by org done")

	return sm, find, err
}

/*
 * get plant sales metric list, by material
 * input:  material code
 * output: metric group by plant <list items>
 */
func GetPlantSalesMetricListByMatl(fr common.FilterReq) ([]common.SalesOrgMetric,
	int, error) {

	var soList []common.SalesOrgMetric
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get plant sales metric list by material")

	sqlcmd := fmt.Sprintf("SELECT b.bic_zrpa_lcit, a.plant, b.bic_ztxt_jyz, b.posx, b.posy, "+
		"sum(a.inv_qty) AS INVQTY, sum(a.cost) AS COST, "+
		"sum(a.netval_inv) AS NETVAL, sum(a.gross_val) AS GVAL, COUNT(DISTINCT a.calday) as SDAY "+
		"FROM bill_zsd a "+
		"LEFT JOIN zaplant_xy b ON a.plant = b.bic_zaplant "+
		"WHERE a.calday BETWEEN '%s' AND '%s' "+
		"AND a.plant != 'A0A1' "+
		"AND a.material = '%s' "+
		"GROUP BY a.plant ORDER BY NETVAL %s LIMIT %d",
		fr.BeginDate, fr.EndDate, fr.Material, fr.SortBy, fr.Limit)

	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return soList, 0, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		var so common.SalesOrgMetric
		err := rows.Scan(&so.ParentOrg, &so.OrgCode, &so.OrgText, &so.PosX, &so.PosY,
			&so.Metric.InvQty, &so.Metric.Cost,
			&so.Metric.NetvalInv, &so.Metric.GrossVal, &so.SaleDays)
		if err != nil {
			zaps.Error("query error: ", err)

		} else {
			so.Metric.GrossProfit =
				so.Metric.NetvalInv - so.Metric.Cost
			if so.Metric.NetvalInv > 0 {
				so.Metric.GrossMargin =
					so.Metric.GrossProfit * 100.0 / so.Metric.NetvalInv
			}
			soList = append(soList, so)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return soList, count, err
	}

	zaps.Info("<<< get sales plant metric by matl done with ", count)

	return soList, count, err
}

/*
 * get date sales metric list, by material
 * input:  material code
 * output: metric group by date <list items>
 */
func GetDateSalesMetricListByMatl(fr common.FilterReq) ([]common.SalesDateMetric,
	int, error) {

	var sdList []common.SalesDateMetric
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get date sales metric list by material")

	sqlcmd := fmt.Sprintf("SELECT a.calday, "+
		"sum(a.inv_qty) AS INVQTY, sum(a.cost) AS COST, "+
		"sum(a.netval_inv) AS NETVAL, sum(a.gross_val) AS GVAL "+
		"FROM bill_zsd a "+
		"WHERE a.calday BETWEEN '%s' AND '%s' "+
		"AND a.plant != 'A0A1' "+
		"AND a.material = '%s' "+
		"GROUP BY a.calday ORDER BY a.calday",
		fr.BeginDate, fr.EndDate, fr.Material)

	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return sdList, 0, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		var sd common.SalesDateMetric
		err := rows.Scan(&sd.Date, &sd.Metric.InvQty, &sd.Metric.Cost,
			&sd.Metric.NetvalInv, &sd.Metric.GrossVal)
		if err != nil {
			zaps.Error("query error: ", err)

		} else {
			sd.Metric.GrossProfit =
				sd.Metric.NetvalInv - sd.Metric.Cost
			if sd.Metric.NetvalInv > 0 {
				sd.Metric.GrossMargin =
					sd.Metric.GrossProfit * 100.0 / sd.Metric.NetvalInv
			}
			sdList = append(sdList, sd)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return sdList, count, err
	}

	zaps.Info("<<< get sales date metric by matl done with count ", count)

	return sdList, count, err
}
