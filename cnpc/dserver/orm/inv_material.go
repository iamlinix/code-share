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
 * get inv basic metric, from material view
 * input:  material
 * output: sum metric <one item>
 */
 /*
func GetInvMetricByMaterial(fr common.FilterReq) (common.SalesMetric,
	bool, error) {

	var sm common.SalesMetric
	var rows *sql.Rows
	var find bool

	zaps.Info(">>> get sales metric by material")

	sqlcmd := fmt.Sprintf("SELECT sum(a.inv_qty) AS INVQTY, " +
			"sum(a.cost) AS COST, sum(a.netval_inv) AS NETVALINV, " +
			"sum(a.gross_val) AS GROSSVAL " +
			"FROM bill_zsd a " +
			"WHERE a.calday BETWEEN '%s' AND '%s' " +
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
				sm.GrossMargin = sm.GrossProfit*100.0/sm.NetvalInv
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
 * get inv basic metric, group by material
 * input:  head / branch
 * output: metric group by org <list items>
 */
func GetMatlInvMetricListByOrg(fr common.FilterReq,
	dtype int) ([]common.InvMatlMetric, int, error) {

	var smList []common.InvMatlMetric
	var rows *sql.Rows
	var sqlcmd string
	var ssdate string
	var qdate string
	var count int
	var err error

	zaps.Info(">>> get material inv metric list")

	if dtype == 0 {
		ssdate, err = GetLatestSSDate(fr.BeginDate)
		if err != nil {
			zaps.Errorf("get latest ss date failed: ", err)
			return smList, 0, err
		}
		qdate = common.GetSubDay(fr.BeginDate, 1)

	} else {
		ssdate, err = GetLatestSSDate(fr.EndDate)
		if err != nil {
			zaps.Errorf("get latest ss date failed: ", err)
			return smList, 0, err
		}
		qdate = fr.EndDate
	}

	sss, err := GetInvSnapshot(ssdate)
	if err != nil {
		zaps.Errorf("get latest inv ss failed: ", err)
		return smList, 0, err
	}

	if fr.OrgLevel == common.ORG_LEVEL_HEAD {
		sqlcmd = fmt.Sprintf("SELECT a.material, b.materialtxt, " +
			"SUM(CASE WHEN dcindic = 'S' THEN quant_b ELSE -quant_b END) AS ZINVSL, " +
			"ROUND(SUM(CASE WHEN dcindic = 'S' THEN value_lc ELSE -value_lc END),2) AS ZINVCOST " +
			"FROM zinv_d01cg a " +
			"LEFT JOIN material b ON a.material = b.material " +
			"WHERE a.pstng_date BETWEEN '%s' AND '%s' " +
			"AND a.plant NOT LIKE 'AA%%' " +
			"GROUP BY a.material",
			ssdate, qdate)

	} else if fr.OrgLevel == common.ORG_LEVEL_BRANCH {
		sqlcmd = fmt.Sprintf("SELECT a.material, b.materialtxt, " +
			"SUM(CASE WHEN dcindic = 'S' THEN quant_b ELSE -quant_b END) AS ZINVSL, " +
			"ROUND(SUM(CASE WHEN dcindic = 'S' THEN value_lc ELSE -value_lc END),2) AS ZINVCOST " +
			"FROM zinv_d01cg a " +
			"LEFT JOIN material b ON a.material = b.material " +
			"LEFT JOIN zaplant_xy c ON a.plant = c.bic_zaplant " +
			"WHERE a.pstng_date BETWEEN '%s' AND '%s' " +
			"AND a.plant NOT LIKE 'AA%%' " +
			"AND c.bic_zrpa_lcit = '%s'" +
			"GROUP BY a.material",
			ssdate, qdate, fr.OrgCode)
	}

	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err = db.Query(sqlcmd)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return smList, 0, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		var sm common.InvMatlMetric
		err := rows.Scan(&sm.Material, &sm.MaterialTxt,
				&sm.Metric.Zinvsl, &sm.Metric.ZinvCost)
		if err != nil {
			zaps.Error("query error: ", err)

		} else {
			//XXX TODO: add history snapshot
			ssm := LookupSSMaterial(sss, sm.Material)
			sm.Metric.Zinvsl += ssm.Zinvsl
			sm.Metric.ZinvCost += ssm.ZinvCost

			smList = append(smList, sm)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return smList, count, err
	}

	zaps.Info("<<< get inv material metric done with count ", count)

	return smList, count, err
}


/* 
 * get plant sales metric list, by material
 * input:  material code
 * output: metric group by plant <list items>
 */
 /*
func GetPlantSalesMetricListByMatl(fr common.FilterReq) ([]common.SalesOrgMetric,
	int, error) {

	var soList []common.SalesOrgMetric
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get plant sales metric list by material")

	sqlcmd := fmt.Sprintf("SELECT a.plant, b.bic_ztxt_jyz, b.posx, b.posy, " +
		"sum(a.inv_qty) AS INVQTY, sum(a.cost) AS COST, " +
		"sum(a.netval_inv) AS NETVAL, sum(a.gross_val) AS GVAL " +
		"FROM bill_zsd a " +
		"LEFT JOIN zaplant_xy b ON a.plant = b.bic_zaplant " +
		"WHERE a.calday BETWEEN '%s' AND '%s' " +
		"AND a.material = '%s' " +
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
		err := rows.Scan(&so.OrgCode, &so.OrgText, &so.PosX, &so.PosY,
				&so.Metric.InvQty, &so.Metric.Cost,
				&so.Metric.NetvalInv, &so.Metric.GrossVal)
		if err != nil {
			zaps.Error("query error: ", err)

		} else {
			so.Metric.GrossProfit =
				so.Metric.NetvalInv - so.Metric.Cost
			if so.Metric.NetvalInv > 0 {
				so.Metric.GrossMargin =
				so.Metric.GrossProfit*100.0/so.Metric.NetvalInv
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
/*
func GetDateSalesMetricListByMatl(fr common.FilterReq) ([]common.SalesDateMetric,
	int, error) {

	var sdList []common.SalesDateMetric
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get date sales metric list by material")

	sqlcmd := fmt.Sprintf("SELECT a.calday, " +
		"sum(a.inv_qty) AS INVQTY, sum(a.cost) AS COST, " +
		"sum(a.netval_inv) AS NETVAL, sum(a.gross_val) AS GVAL " +
		"FROM bill_zsd a " +
		"WHERE a.calday BETWEEN '%s' AND '%s' " +
		"AND a.material = '%s' " +
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
				sd.Metric.GrossProfit*100.0/sd.Metric.NetvalInv
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

*/
