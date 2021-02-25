package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"fmt"
	"time"
)

///////////////////////////////////////////////////////////////////////////////
// Test-market Materials List: include First Order Date

func GetTestMarketMatlList(fr common.FilterReq) ([]common.TestMktMatl,
	int, error) {

	var mList []common.TestMktMatl
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get material zifpurd list")

	edate := common.GetSubDay(fr.EndDate, 90) //XXX TODO: use dynamic config

	sqlcmd := fmt.Sprintf("SELECT t.material, m.materialtxt, t.FDATE "+
		"FROM "+
		"(SELECT material, MIN(zgr_date) as FDATE "+
		"FROM zifpurd "+
		"WHERE zgr_date <= '%s' "+
		"GROUP BY material) t "+
		"LEFT JOIN material m ON t.material = m.material "+
		"WHERE t.FDATE BETWEEN '%s' AND '%s'",
		edate, fr.BeginDate, fr.EndDate)

	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return mList, 0, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		var m common.TestMktMatl
		err := rows.Scan(&m.Material, &m.MaterialTxt,
			&m.FirstOrderDate)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> material code: ", m.Material)
			zaps.Debug(">>> first order: ", m.FirstOrderDate)

			mList = append(mList, m)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return mList, count, err
	}

	zaps.Info("<<< get test-market materials list with ", count)

	return mList, count, err
}

///////////////////////////////////////////////////////////////////////////////
// Get First-Order-Date For Material

func GetFirstOrderDateByMatl(fr common.FilterReq) (string, error) {

	var fdate string
	var rows *sql.Rows

	zaps.Info(">>> get material first order date")

	sqlcmd := fmt.Sprintf("SELECT MIN(a.zgr_date) as FDATE "+
		"FROM zifpurd a "+
		"WHERE a.zgr_date <= '%s' "+
		"AND material = '%s'",
		fr.EndDate, fr.Material)

	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return fdate, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&fdate)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> fdate: ", fdate)
			break
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return fdate, err
	}

	zaps.Info("<<< get material first-order-date done")

	return fdate, err
}

///////////////////////////////////////////////////////////////////////////////
// Test-market Materials Metric List

func GetTestMarketMetricByMatl(fr common.FilterReq,
	item common.EffTestMktMatl) (common.EffTestMktMatl, error) {

	var rows *sql.Rows

	//XXX TODO: optimize cmd -> oi_ebeln
	sqlcmd := fmt.Sprintf("SELECT t1.material, m.materialtxt, "+
		"t1.buy_count, t1.buy_money, t1.rt_count, "+
		"t1.rt_money, t2.SUMQTY, t2.SUMCOST "+
		"FROM (SELECT a.material, "+
		"SUM(CASE WHEN movetype = '101' THEN quant_b "+
		"ELSE CASE WHEN movetype = '102' "+
		"THEN -quant_b ELSE 0 END END) AS buy_count, "+
		"SUM(CASE WHEN movetype = '101' THEN value_lc "+
		"ELSE CASE WHEN movetype = '102' "+
		"THEN -ABS(value_lc) ELSE 0 END END) AS buy_money, "+
		"SUM(CASE WHEN movetype = '161' THEN quant_b "+
		"ELSE CASE WHEN movetype = '162' THEN -quant_b "+
		"ELSE 0 END END) AS rt_count, "+
		"SUM(CASE WHEN movetype = '161' THEN b.zpoamount "+
		"ELSE CASE WHEN movetype = '162' "+
		"THEN -ABS(b.zpoamount) ELSE 0 END END) AS rt_money "+
		"FROM zinv_d01cg a "+
		"LEFT JOIN ( SELECT DISTINCT oi_ebeln, vendor, material, "+
		"MAX(zpoamount) AS zpoamount FROM zifpurd GROUP BY oi_ebeln, vendor, material ) b ON "+
		"a.oi_ebeln = b.oi_ebeln AND a.material = b.material "+
		"WHERE pstng_date BETWEEN '%s' AND '%s' "+
		"AND a.material = '%s' "+
		"AND (movetype IN ('101', '102', '161', '162') "+
		"AND b.oi_ebeln LIKE '46%%' "+
		"AND ((a.oi_ebeln BETWEEN "+
		"(SELECT min(oi_ebeln) FROM zifpurd WHERE zgr_date "+
		"BETWEEN '%s' AND '%s') "+
		"AND "+
		"(SELECT max(oi_ebeln) FROM zifpurd WHERE zgr_date "+
		"BETWEEN '%s' AND '%s') OR zbsart = 'Z031')) "+
		")) t1 "+
		"LEFT JOIN ("+
		"SELECT material, sum(inv_qty) as SUMQTY, "+
		"sum(cost) as SUMCOST "+
		"FROM bill_zsd "+
		"WHERE material = '%s' "+
		"AND calday BETWEEN '%s' AND '%s' "+
		") t2 "+
		"ON t1.material = t2.material "+
		"LEFT JOIN material m ON t1.material = m.material ",
		fr.BeginDate, fr.EndDate, item.Material,
		fr.BeginDate, fr.EndDate, fr.BeginDate, fr.EndDate,
		item.Material, fr.BeginDate, fr.EndDate)

	//zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return item, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("m (%s) db query using %d ms", item.Material, (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		var buyinQty, buyinCost sql.NullFloat64
		var returnQty, returnCost sql.NullFloat64
		var salesQty, salesCost sql.NullFloat64
		var mcode, mtxt sql.NullString

		err := rows.Scan(&mcode, &mtxt,
			&buyinQty, &buyinCost, &returnQty,
			&returnCost, &salesQty, &salesCost)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> material: ", item.Material)

			if mtxt.Valid {
				item.MaterialTxt = mtxt.String
			} else {
				item.MaterialTxt = "N/A"
			}

			if buyinQty.Valid && returnQty.Valid {
				item.Metric.PurchQty =
					buyinQty.Float64 - returnQty.Float64
			} else {
				item.Metric.PurchQty = 0
			}

			if buyinCost.Valid && returnCost.Valid {
				item.Metric.PurchCost =
					buyinCost.Float64 - returnCost.Float64
			} else {
				item.Metric.PurchCost = 0
			}

			if salesQty.Valid {
				item.Metric.SalesQty = salesQty.Float64
			} else {
				item.Metric.SalesQty = 0
			}

			if salesCost.Valid {
				item.Metric.SalesCost = salesCost.Float64
			} else {
				item.Metric.SalesCost = 0
			}
			break
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return item, err
	}

	zaps.Info("<<< get test-market materials metric done")

	return item, err
}
