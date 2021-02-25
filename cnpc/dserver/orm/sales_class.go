package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	_ "github.com/go-sql-driver/mysql"
	"database/sql"

	"time"
	"fmt"
)


///////////////////////////////////////////////////////////////////////////////
// 
/* 
 * get sales basic metric, from class view
 * input:  class code/level
 * output: sum metric <one item>
 */
func GetSalesMetricByClass(fr common.FilterReq) (common.SalesMetric,
	bool, error) {

	var sm common.SalesMetric
	var sqlcmd string
	var rows *sql.Rows
	var find bool

	zaps.Info(">>> get sales metric by class")

	sqla := fmt.Sprintf("SELECT sum(a.inv_qty) AS INVQTY, " +
			"sum(a.cost) AS COST, sum(a.netval_inv) AS NETVAL, " +
			"sum(a.gross_val) AS GROSSVAL " +
			"FROM bill_zsd a " +
			"LEFT JOIN material b ON a.material = b.material "+
			"WHERE a.calday BETWEEN '%s' AND '%s' " +
			"AND a.plant != 'A0A1' ",
			fr.BeginDate, fr.EndDate)

	if fr.ClassLevel == common.CLASS_LEVEL_MAIN {
		sqlcmd = fmt.Sprintf("%s AND b.bic_zklad2 = '%s'",
					sqla, fr.Material)

	} else if fr.ClassLevel == common.CLASS_LEVEL_MID {
		sqlcmd = fmt.Sprintf("%s AND b.bic_zklasse_d = '%s'",
					sqla, fr.Material)

	} else if fr.ClassLevel == common.CLASS_LEVEL_SUB {
		sqlcmd = fmt.Sprintf("%s AND b.bic_zrpa_mat = '%s'",
					sqla, fr.Material)
	}

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
			return sm, find, err
		}

		sm.GrossProfit = sm.NetvalInv - sm.Cost
		if sm.NetvalInv > 0 {
			sm.GrossMargin = sm.GrossProfit*100.0/sm.NetvalInv
		}
		find = true
		break
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
 * get plant sales metric list, by class
 * input:  class level/code
 * output: metric group by plant <list items>
 */
func GetPlantSalesMetricListByClass(fr common.FilterReq) ([]common.SalesOrgMetric,
	int, error) {

	var soList []common.SalesOrgMetric
	var sqlcmd string
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get plant sales metric list by class")

	sqla := fmt.Sprintf("SELECT a.plant, b.bic_ztxt_jyz, b.posx, b.posy, " +
		"sum(a.inv_qty) AS INVQTY, sum(a.cost) AS COST, " +
		"sum(a.netval_inv) AS NETVAL, sum(a.gross_val) AS GVAL " +
		"FROM bill_zsd a " +
		"LEFT JOIN zaplant_xy b ON a.plant = b.bic_zaplant " +
		"LEFT JOIN material c ON a.material = c.material "+
		"WHERE a.calday BETWEEN '%s' AND '%s' "+
		"AND a.plant != 'A0A1' ",
		fr.BeginDate, fr.EndDate)

	sqlb := fmt.Sprintf("GROUP BY a.plant ORDER BY NETVAL %s LIMIT %d",
			fr.SortBy, fr.Limit)

	if fr.ClassLevel == common.CLASS_LEVEL_MAIN {
		sqlcmd = fmt.Sprintf("%s AND c.bic_zklad2 = '%s' %s",
					sqla, fr.Material, sqlb)

	} else if fr.ClassLevel == common.CLASS_LEVEL_MID {
		sqlcmd = fmt.Sprintf("%s AND c.bic_zklasse_d = '%s' %s",
					sqla, fr.Material, sqlb)

	} else if fr.ClassLevel == common.CLASS_LEVEL_SUB {
		sqlcmd = fmt.Sprintf("%s AND c.bic_zrpa_mat = '%s' %s",
					sqla, fr.Material, sqlb)
	}

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
		err := rows.Scan(&so.OrgCode, &so.OrgText, &so.PosX, &so.PosY,
				&so.Metric.InvQty, &so.Metric.Cost,
				&so.Metric.NetvalInv, &so.Metric.GrossVal)
		if err != nil {
			zaps.Error("query error: ", err)
			return soList, count, err
		}

		so.Metric.GrossProfit = so.Metric.NetvalInv - so.Metric.Cost
		if so.Metric.NetvalInv > 0 {
			so.Metric.GrossMargin =
				so.Metric.GrossProfit*100.0/so.Metric.NetvalInv
		}
		soList = append(soList, so)
		count++
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return soList, count, err
	}

	zaps.Info("<<< get sales plant metric by class done with count ", count)

	return soList, count, err
}


/* 
 * get date sales metric list, by class
 * input:  class code/level
 * output: metric group by date <list items>
 */
func GetDateSalesMetricListByClass(fr common.FilterReq) ([]common.SalesDateMetric,
	int, error) {

	var sdList []common.SalesDateMetric
	var rows *sql.Rows
	var sqlcmd string
	var count int

	zaps.Info(">>> get date sales metric list by class")

	sqla := fmt.Sprintf("SELECT a.calday, " +
		"sum(a.inv_qty) AS INVQTY, sum(a.cost) AS COST, " +
		"sum(a.netval_inv) AS NETVAL, sum(a.gross_val) AS GVAL " +
		"FROM bill_zsd a " +
		"LEFT JOIN material b ON a.material = b.material " +
		"WHERE a.calday BETWEEN '%s' AND '%s' " +
		"AND a.plant != 'A0A1' ",
		fr.BeginDate, fr.EndDate)

	sqlb := fmt.Sprintf("GROUP BY a.calday ORDER BY a.calday")

	if fr.ClassLevel == common.CLASS_LEVEL_MAIN {
		sqlcmd = fmt.Sprintf("%s AND b.bic_zklad2 = '%s' %s",
					sqla, fr.Material, sqlb)

	} else if fr.ClassLevel == common.CLASS_LEVEL_MID {
		sqlcmd = fmt.Sprintf("%s AND b.bic_zklasse_d = '%s' %s",
					sqla, fr.Material, sqlb)

	} else if fr.ClassLevel == common.CLASS_LEVEL_SUB {
		sqlcmd = fmt.Sprintf("%s AND b.bic_zrpa_mat = '%s' %s",
					sqla, fr.Material, sqlb)
	}

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
			return sdList, count, err
		}

		sd.Metric.GrossProfit = sd.Metric.NetvalInv - sd.Metric.Cost
		if sd.Metric.NetvalInv > 0 {
			sd.Metric.GrossMargin =
			sd.Metric.GrossProfit*100.0/sd.Metric.NetvalInv
		}
		sdList = append(sdList, sd)
		count++
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return sdList, count, err
	}

	zaps.Info("<<< get sales date metric by class done with count ", count)

	return sdList, count, err
}


