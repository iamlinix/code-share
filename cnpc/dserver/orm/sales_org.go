package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"time"
)

///////////////////////////////////////////////////////////////////////////////
//
/*
 * get sales basic metric, from org view
 * input:  head / branch / plant
 * output: sum metric <one item>
 */
func GetSalesMetricByOrg(fr common.FilterReq) (common.SalesMetric,
	bool, error) {

	var sm common.SalesMetric
	var rows *sql.Rows
	var sqlcmd string
	var find bool

	zaps.Info(">>> get sales metric by org")

	if fr.OrgLevel == common.ORG_LEVEL_HEAD {
		sqlcmd = fmt.Sprintf("SELECT sum(inv_qty) AS INVQTY, "+
			"sum(cost) AS COST, sum(netval_inv) AS NETVALINV, "+
			"sum(gross_val) AS GROSSVAL "+
			"FROM bill_zsd "+
			"WHERE calday BETWEEN '%s' AND '%s' "+
			"AND plant != 'A0A1'",
			fr.BeginDate, fr.EndDate)

	} else if fr.OrgLevel == common.ORG_LEVEL_BRANCH {
		sqlcmd = fmt.Sprintf("SELECT sum(a.inv_qty) AS INVQTY, "+
			"sum(a.cost) AS COST, sum(a.netval_inv) AS NETVALINV, "+
			"sum(a.gross_val) AS GROSSVAL "+
			"FROM bill_zsd a FORCE INDEX(idx_calday_plant) "+
			"LEFT JOIN zaplant_xy b ON a.plant = b.bic_zaplant "+
			"WHERE a.calday BETWEEN '%s' AND '%s' "+
			"AND b.bic_zrpa_lcit = '%s'",
			fr.BeginDate, fr.EndDate, fr.OrgCode)

	} else if fr.OrgLevel == common.ORG_LEVEL_PLANT {
		sqlcmd = fmt.Sprintf("SELECT sum(inv_qty) AS INVQTY, "+
			"sum(cost) AS COST, sum(netval_inv) AS NETVALINV, "+
			"sum(gross_val) AS GROSSVAL "+
			"FROM bill_zsd "+
			"WHERE calday BETWEEN '%s' AND '%s' "+
			"AND plant = '%s'",
			fr.BeginDate, fr.EndDate, fr.OrgCode)
	}
	if fr.RateThreshold > 0.0 {
		sqlcmd += fmt.Sprintf(" AND (netval_inv-cost)/netval_inv >= %.3f", fr.RateThreshold)
	}

	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Errorf("(%s)\ndb query failed: ", sqlcmd, err)
		return sm, find, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()
	var nullQty sql.NullInt32
	var nullCost sql.NullFloat64
	var nullNetval sql.NullFloat64
	var nullGross sql.NullFloat64

	for rows.Next() {
		err := rows.Scan(&nullQty, &nullCost, &nullNetval, &nullGross)
		if err != nil {
			zaps.Errorf("(%s)\nrow scan error: ", sqlcmd, err)
			return sm, find, err
		}

		if nullQty.Valid {
			sm.InvQty = int(nullQty.Int32)
		}

		if nullCost.Valid {
			sm.Cost = nullCost.Float64
		}

		if nullNetval.Valid {
			sm.NetvalInv = nullNetval.Float64
		}

		if nullGross.Valid {
			sm.GrossVal = nullGross.Float64
		}

		sm.GrossProfit = sm.NetvalInv - sm.Cost
		sm.GrossMargin = sm.GrossProfit * 100.0 / sm.NetvalInv
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
 * get sales basic metric, group by plant
 * input:  org-level <head/branch>
 * output: metric in org group by plant <list items>
 */
func GetPlantSalesMetricListByOrg(fr common.FilterReq) ([]common.SalesOrgMetric,
	int, error) {

	var soList []common.SalesOrgMetric
	var rows *sql.Rows
	var sqlcmd string
	var count int

	zaps.Info(">>> get plant sales metric list")

	if fr.OrgLevel == common.ORG_LEVEL_HEAD {
		sqlcmd = fmt.Sprintf("SELECT b.bic_zrpa_lcit, a.plant, b.bic_ztxt_jyz, "+
			"b.posx, b.posy, "+
			"sum(a.inv_qty) AS INVQTY, sum(a.cost) AS COST, "+
			"sum(a.netval_inv) AS NETVAL, sum(a.gross_val) AS GVAL, COUNT(DISTINCT a.calday) as SDAY "+
			"FROM bill_zsd a "+
			"LEFT JOIN zaplant_xy b ON a.plant = b.bic_zaplant "+
			"WHERE a.calday BETWEEN '%s' AND '%s' "+
			"AND plant != 'A0A1' ", fr.BeginDate, fr.EndDate)

	} else if fr.OrgLevel == common.ORG_LEVEL_BRANCH {
		sqlcmd = fmt.Sprintf("SELECT b.bic_zrpa_lcit, a.plant, b.bic_ztxt_jyz, "+
			"b.posx, b.posy, "+
			"sum(a.inv_qty) AS INVQTY, sum(a.cost) AS COST, "+
			"sum(a.netval_inv) AS NETVAL, sum(a.gross_val) AS GVAL, COUNT(DISTINCT a.calday) as SDAY "+
			"FROM bill_zsd a "+
			"LEFT JOIN zaplant_xy b ON a.plant = b.bic_zaplant "+
			"WHERE a.calday BETWEEN '%s' AND '%s' "+
			"AND b.bic_zrpa_lcit = '%s' ", fr.BeginDate, fr.EndDate, fr.OrgCode)
	}
	if fr.RateThreshold > 0.0 {
		sqlcmd += fmt.Sprintf(" AND (netval_inv-cost)/netval_inv >= %.3f", fr.RateThreshold)
	}
	sqlcmd += fmt.Sprintf(" GROUP BY a.plant ORDER BY NETVAL %s LIMIT %d", fr.SortBy, fr.Limit)

	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Errorf("(%s)\ndb query failed: ", sqlcmd, err)
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
			zaps.Errorf("(%s)\nrow scan error: ", sqlcmd, err)
			return soList, count, err
		}

		/*
			if so.OrgCode == "A0A1" {
				zaps.Info("skip A0A1 plant")
				continue
			}
		*/

		so.Metric.GrossProfit = so.Metric.NetvalInv - so.Metric.Cost
		if so.Metric.NetvalInv > 0 {
			so.Metric.GrossMargin =
				so.Metric.GrossProfit * 100.0 / so.Metric.NetvalInv
		}
		soList = append(soList, so)
		count++
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return soList, count, err
	}

	zaps.Info("<<< get sales plant metric done with count ", count)

	return soList, count, err
}

func GetPlantSalesMetricByOrg(fr common.FilterReq) (common.SalesOrgMetric, error) {

	var so common.SalesOrgMetric
	var rows *sql.Rows
	var sqlcmd string

	zaps.Info(">>> get plant sales metric")

	sqlcmd = fmt.Sprintf("SELECT b.bic_zrpa_lcit, a.plant, b.bic_ztxt_jyz, "+
		"b.posx, b.posy, "+
		"sum(a.inv_qty) AS INVQTY, sum(a.cost) AS COST, "+
		"sum(a.netval_inv) AS NETVAL, sum(a.gross_val) AS GVAL, COUNT(DISTINCT a.calday) as SDAY "+
		"FROM bill_zsd a "+
		"LEFT JOIN zaplant_xy b ON a.plant = b.bic_zaplant "+
		"WHERE a.calday BETWEEN '%s' AND '%s' "+
		"AND a.plant = '%s' ", fr.BeginDate, fr.EndDate, fr.OrgCode)
	if fr.RateThreshold > 0.0 {
		sqlcmd += fmt.Sprintf(" AND (netval_inv-cost)/netval_inv >= %.3f", fr.RateThreshold)
	}
	sqlcmd += fmt.Sprintf(" ORDER BY NETVAL %s LIMIT %d", fr.SortBy, fr.Limit)

	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Errorf("(%s)\ndb query failed: ", sqlcmd, err)
		return so, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&so.ParentOrg, &so.OrgCode, &so.OrgText, &so.PosX, &so.PosY,
			&so.Metric.InvQty, &so.Metric.Cost,
			&so.Metric.NetvalInv, &so.Metric.GrossVal, &so.SaleDays)
		if err != nil {
			zaps.Errorf("(%s)\nrow scan error: ", sqlcmd, err)
			return so, err
		}

		so.Metric.GrossProfit = so.Metric.NetvalInv - so.Metric.Cost
		if so.Metric.NetvalInv > 0 {
			so.Metric.GrossMargin =
				so.Metric.GrossProfit * 100.0 / so.Metric.NetvalInv
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return so, err
	}

	return so, nil
}

/*
 * get sales basic metric, group by material
 * input:  org-level <head/branch>
 * output: metric in org group by matl <list items>
 */
func GetMatlSalesMetricListByOrg(fr common.FilterReq) ([]common.SalesMatlMetric,
	int, error) {

	var smList []common.SalesMatlMetric
	var rows *sql.Rows
	var sqlcmd string
	var count int

	zaps.Info(">>> get material sales metric list")

	if fr.OrgLevel == common.ORG_LEVEL_HEAD {
		sqlcmd = fmt.Sprintf("SELECT a.material, b.materialtxt, "+
			"sum(a.inv_qty) AS INVQTY, sum(a.cost) AS COST, "+
			"sum(a.netval_inv) AS NETVAL, sum(a.gross_val) AS GVAL "+
			"FROM bill_zsd a "+
			"LEFT JOIN material b ON a.material = b.material "+
			"WHERE a.calday BETWEEN '%s' AND '%s' "+
			"AND plant != 'A0A1' ", fr.BeginDate, fr.EndDate)

	} else if fr.OrgLevel == common.ORG_LEVEL_BRANCH {
		sqlcmd = fmt.Sprintf("SELECT a.material, b.materialtxt, "+
			"sum(a.inv_qty) AS INVQTY, sum(a.cost) AS COST, "+
			"sum(a.netval_inv) AS NETVAL, sum(a.gross_val) AS GVAL "+
			"FROM bill_zsd a "+
			"LEFT JOIN material b ON a.material = b.material "+
			"LEFT JOIN zaplant_xy c ON a.plant = c.bic_zaplant "+
			"WHERE a.calday BETWEEN '%s' AND '%s' "+
			"AND c.bic_zrpa_lcit = '%s'", fr.BeginDate, fr.EndDate, fr.OrgCode)

	} else if fr.OrgLevel == common.ORG_LEVEL_PLANT {
		sqlcmd = fmt.Sprintf("SELECT a.material, b.materialtxt, "+
			"sum(a.inv_qty) AS INVQTY, sum(a.cost) AS COST, "+
			"sum(a.netval_inv) AS NETVAL, sum(a.gross_val) AS GVAL "+
			"FROM bill_zsd a "+
			"LEFT JOIN material b ON a.material = b.material "+
			"WHERE a.calday BETWEEN '%s' AND '%s' "+
			"AND a.plant = '%s' ", fr.BeginDate, fr.EndDate, fr.OrgCode)

	}
	if fr.RateThreshold > 0.0 {
		sqlcmd += fmt.Sprintf(" AND (netval_inv-cost)/netval_inv >= %.3f", fr.RateThreshold)
	}
	sqlcmd += fmt.Sprintf(" GROUP BY a.material ORDER BY NETVAL %s LIMIT %d", fr.SortBy, fr.Limit)

	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Errorf("(%s)\ndb query failed: ", sqlcmd, err)
		return smList, 0, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		var sm common.SalesMatlMetric
		err := rows.Scan(&sm.Material, &sm.MaterialTxt,
			&sm.Metric.InvQty, &sm.Metric.Cost,
			&sm.Metric.NetvalInv, &sm.Metric.GrossVal)
		if err != nil {
			zaps.Errorf("(%s)\nrow scan error: ", sqlcmd, err)
			return smList, count, err
		}

		sm.Metric.GrossProfit =
			sm.Metric.NetvalInv - sm.Metric.Cost
		if sm.Metric.NetvalInv > 0 {
			sm.Metric.GrossMargin =
				sm.Metric.GrossProfit * 100.0 / sm.Metric.NetvalInv
		}
		smList = append(smList, sm)
		count++
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return smList, count, err
	}

	zaps.Info("<<< get sales material metric done with count ", count)

	return smList, count, err
}

/*
 * get sales basic metric, group by main-class
 * input:  org-level <head/branch>
 * output: metric in org group by main-class <list items>
 */
func GetClassSalesMetricListByOrg(fr common.FilterReq) ([]common.SalesClassMetric,
	int, error) {

	var scList []common.SalesClassMetric
	var rows *sql.Rows
	var sqlcmd string
	var count int

	zaps.Info(">>> get class sales metric list by org")

	if fr.OrgLevel == common.ORG_LEVEL_HEAD {
		sqlcmd = fmt.Sprintf("SELECT a.zklad2, "+
			"any_value(b.class_txt), "+
			"sum(a.inv_qty) AS INVQTY, sum(a.cost) AS COST, "+
			"sum(a.netval_inv) AS NETVAL, "+
			"sum(a.gross_val) AS GROSSVAL "+
			"FROM bill_zsd a "+
			"LEFT JOIN basic_class_dict b "+
			"ON a.zklad2 = b.class_code "+
			"WHERE a.calday BETWEEN '%s' AND '%s' "+
			"AND a.plant != 'A0A1' ", fr.BeginDate, fr.EndDate)

	} else if fr.OrgLevel == common.ORG_LEVEL_BRANCH {
		sqlcmd = fmt.Sprintf("SELECT a.zklad2, "+
			"any_value(b.class_txt), "+
			"sum(a.inv_qty) AS INVQTY, sum(a.cost) AS COST, "+
			"sum(a.netval_inv) AS NETVAL, "+
			"sum(a.gross_val) AS GROSSVAL "+
			"FROM bill_zsd a "+
			"LEFT JOIN basic_class_dict b "+
			"ON a.zklad2 = b.class_code "+
			"LEFT JOIN zaplant_xy c "+
			"ON a.plant = c.bic_zaplant "+
			"WHERE a.calday BETWEEN '%s' AND '%s' "+
			"AND c.bic_zrpa_lcit = '%s' ", fr.BeginDate, fr.EndDate, fr.OrgCode)

	} else if fr.OrgLevel == common.ORG_LEVEL_PLANT {
		sqlcmd = fmt.Sprintf("SELECT a.zklad2, "+
			"any_value(b.class_txt), "+
			"sum(a.inv_qty) AS INVQTY, sum(a.cost) AS COST, "+
			"sum(a.netval_inv) AS NETVAL, "+
			"sum(a.gross_val) AS GROSSVAL "+
			"FROM bill_zsd a "+
			"LEFT JOIN basic_class_dict b "+
			"ON a.zklad2 = b.class_code "+
			"WHERE a.calday BETWEEN '%s' AND '%s' "+
			"AND a.plant = '%s' ", fr.BeginDate, fr.EndDate, fr.OrgCode)
	}
	if fr.RateThreshold > 0.0 {
		sqlcmd += fmt.Sprintf(" AND (netval_inv-cost)/netval_inv >= %.3f", fr.RateThreshold)
	}
	sqlcmd += " GROUP BY a.zklad2 ORDER BY NETVAL"
	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Errorf("(%s)\ndb query failed: ", sqlcmd, err)
		return scList, 0, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		var sc common.SalesClassMetric
		err := rows.Scan(&sc.ClassCode, &sc.ClassText,
			&sc.Metric.InvQty, &sc.Metric.Cost,
			&sc.Metric.NetvalInv, &sc.Metric.GrossVal)
		if err != nil {
			zaps.Errorf("(%s)\nrow scan error: ", sqlcmd, err)
			return scList, count, err
		}

		sc.Metric.GrossProfit = sc.Metric.NetvalInv - sc.Metric.Cost
		if sc.Metric.NetvalInv > 0 {
			sc.Metric.GrossMargin =
				sc.Metric.GrossProfit * 100.0 / sc.Metric.NetvalInv
		}
		scList = append(scList, sc)
		count++
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return scList, count, err
	}

	zaps.Info("<<< get sales class metric done with count ", count)

	return scList, count, err
}

/*
 * get sales basic metric, group by date
 * input:  org-level <head/branch>
 * output: metric in org group by date <list items>
 */
func GetDateSalesMetricListByOrg(fr common.FilterReq) ([]common.SalesDateMetric,
	int, error) {

	var sdList []common.SalesDateMetric
	var rows *sql.Rows
	var sqlcmd string
	var count int

	zaps.Info(">>> get date sales metric list by org")

	if fr.OrgLevel == common.ORG_LEVEL_HEAD {
		sqlcmd = fmt.Sprintf("SELECT a.calday, "+
			"sum(a.inv_qty) AS INVQTY, sum(a.cost) AS COST, "+
			"sum(a.netval_inv) AS NETVAL, "+
			"sum(a.gross_val) AS GROSSVAL "+
			"FROM bill_zsd a "+
			"WHERE a.calday BETWEEN '%s' AND '%s' "+
			"AND a.plant != 'A0A1' ", fr.BeginDate, fr.EndDate)

	} else if fr.OrgLevel == common.ORG_LEVEL_BRANCH {
		sqlcmd = fmt.Sprintf("SELECT a.calday, "+
			"sum(a.inv_qty) AS INVQTY, sum(a.cost) AS COST, "+
			"sum(a.netval_inv) AS NETVAL, "+
			"sum(a.gross_val) AS GROSSVAL "+
			"FROM bill_zsd a "+
			"LEFT JOIN zaplant_xy b "+
			"ON a.plant = b.bic_zaplant "+
			"WHERE a.calday BETWEEN '%s' AND '%s' "+
			"AND b.bic_zrpa_lcit = '%s' ", fr.BeginDate, fr.EndDate, fr.OrgCode)

	} else if fr.OrgLevel == common.ORG_LEVEL_PLANT {
		sqlcmd = fmt.Sprintf("SELECT a.calday, "+
			"sum(a.inv_qty) AS INVQTY, sum(a.cost) AS COST, "+
			"sum(a.netval_inv) AS NETVAL, "+
			"sum(a.gross_val) AS GROSSVAL "+
			"FROM bill_zsd a "+
			"WHERE a.calday BETWEEN '%s' AND '%s' "+
			"AND a.plant = '%s' ", fr.BeginDate, fr.EndDate, fr.OrgCode)
	}
	if fr.RateThreshold > 0.0 {
		sqlcmd += fmt.Sprintf(" AND (netval_inv-cost)/netval_inv >= %.3f", fr.RateThreshold)
	}
	sqlcmd += " GROUP BY a.calday ORDER BY a.calday"
	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Errorf("(%s)\ndb query failed: ", sqlcmd, err)
		return sdList, 0, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		var sd common.SalesDateMetric
		err := rows.Scan(&sd.Date, &sd.Metric.InvQty,
			&sd.Metric.Cost, &sd.Metric.NetvalInv,
			&sd.Metric.GrossVal)
		if err != nil {
			zaps.Errorf("(%s)\nrow scan error: ", sqlcmd, err)
			return sdList, count, err
		}

		sd.Metric.GrossProfit = sd.Metric.NetvalInv - sd.Metric.Cost
		if sd.Metric.NetvalInv > 0 {
			sd.Metric.GrossMargin =
				sd.Metric.GrossProfit * 100.0 / sd.Metric.NetvalInv
		}
		sdList = append(sdList, sd)
		count++
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return sdList, count, err
	}

	zaps.Info("<<< get sales date metric done with count ", count)

	return sdList, count, err
}
