package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"time"
)

func AddMaterialReportTx(fr common.FilterReq, emis []common.ERPMaterialInfo) error {

	t1 := time.Now().UnixNano() / 1e6

	tx, _ := db.Begin()
	defer tx.Rollback()

	for _, v := range emis {
		err := AddMaterialReport(tx, fr.Month, v)
		if err != nil {
			zaps.Errorf("add material report failed: %v", err)
			return err
		}
	}

	/* Transaction END */
	err := tx.Commit()
	if err != nil {
		zaps.Errorf("transaction commit failed: %v", err)
		return err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("<<< add material report tx done using %d ms", t2-t1)

	return nil
}

func AddMaterialReport(tx *sql.Tx, month string, v common.ERPMaterialInfo) error {

	stmt, err := tx.Prepare("INSERT INTO materials_report(month, " +
		"material, material_txt, vendor, vendor_name, " +
		"vendor_flag, buyin, buyin_wtax, returned, " +
		"returned_wtax, rebate, pconf_qty, zpurps, " +
		"inv_qty, cost, cost_wtax, netval_inv, gross_val, " +
		"open_zinvsl, open_zinvcost, close_zinvsl, " +
		"close_zinvcost, created_time) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, " +
		"?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		zaps.Error("insert into materials report failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(month, v.Material, v.MaterialTxt, v.Vendor,
		v.VendorName, v.VendorFlag, v.Buyin, v.BuyinWtax,
		v.Return, v.ReturnWtax, v.Rebate, v.PconfQty,
		v.Zpurps, v.InvQty, v.Cost, v.CostWtax, v.NetvalInv,
		v.GrossVal, v.OpenZinvsl, v.OpenZinvCost,
		v.CloseZinvsl, v.CloseZinvCost,
		time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		zaps.Error("db exec failed: ", err)
		return err
	}

	return err
}

func DelMaterialReport(month string) {

	zaps.Info(">>> del material report")

	_, err := db.Exec("DELETE FROM materials_report WHERE month = ?", month)
	if err != nil {
		zaps.Error("delete materials report exec failed: ", err)
	}

	zaps.Info("<<< del materials report done")
}

func GetMaterialReportListByMonthAll(month string) ([]common.ERPMaterialInfo, int, error) {

	var emiList []common.ERPMaterialInfo
	var count int

	zaps.Info(">>> get materials report by month list")

	rows, err := db.Query("SELECT material, material_txt, vendor, "+
		"vendor_name, vendor_flag, buyin, buyin_wtax, "+
		"returned, returned_wtax, rebate, pconf_qty, zpurps, "+
		"inv_qty, cost, cost_wtax, netval_inv, gross_val,"+
		"open_zinvsl, open_zinvcost, close_zinvsl, "+
		"close_zinvcost "+
		"FROM materials_report WHERE month = ? "+
		"ORDER BY vendor", month)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return emiList, count, err
	}

	defer rows.Close()

	for rows.Next() {
		var v common.ERPMaterialInfo
		err := rows.Scan(&v.Material, &v.MaterialTxt, &v.Vendor,
			&v.VendorName, &v.VendorFlag, &v.Buyin, &v.BuyinWtax,
			&v.Return, &v.ReturnWtax, &v.Rebate, &v.PconfQty,
			&v.Zpurps, &v.InvQty, &v.Cost, &v.CostWtax, &v.NetvalInv,
			&v.GrossVal, &v.OpenZinvsl, &v.OpenZinvCost,
			&v.CloseZinvsl, &v.CloseZinvCost)
		if err != nil {
			zaps.Error("query error: ", err)
			return emiList, count, err
		} else {
			zaps.Debug(">>> month: ", month)

			emiList = append(emiList, v)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return emiList, count, err
	}

	zaps.Info("<<< get materials report list done")

	return emiList, count, err
}

func GetMaterialReportListByMonthActive(month string) ([]common.ERPMaterialInfo, int, error) {

	var emiList []common.ERPMaterialInfo
	var count int

	zaps.Info(">>> get materials report by month list")

	rows, err := db.Query("SELECT a.material, a.material_txt, a.vendor, "+
		"a.vendor_name, a.vendor_flag, a.buyin, a.buyin_wtax, "+
		"a.returned, a.returned_wtax, a.rebate, a.pconf_qty, "+
		"a.zpurps, a.inv_qty, a.cost, a.cost_wtax, a.netval_inv, "+
		"a.gross_val, a.open_zinvsl, a.open_zinvcost, "+
		"a.close_zinvsl, a.close_zinvcost "+
		"FROM materials_report a "+
		"LEFT JOIN vendor_cfg b ON a.vendor = b.vendor_code "+
		"WHERE a.month = ? AND b.inactive = 0 "+
		"ORDER BY vendor", month)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return emiList, count, err
	}

	defer rows.Close()

	for rows.Next() {
		var v common.ERPMaterialInfo
		err := rows.Scan(&v.Material, &v.MaterialTxt, &v.Vendor,
			&v.VendorName, &v.VendorFlag, &v.Buyin, &v.BuyinWtax,
			&v.Return, &v.ReturnWtax, &v.Rebate, &v.PconfQty,
			&v.Zpurps, &v.InvQty, &v.Cost, &v.CostWtax, &v.NetvalInv,
			&v.GrossVal, &v.OpenZinvsl, &v.OpenZinvCost,
			&v.CloseZinvsl, &v.CloseZinvCost)
		if err != nil {
			zaps.Error("query error: ", err)
			return emiList, count, err
		} else {
			zaps.Debug(">>> month: ", month)

			emiList = append(emiList, v)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return emiList, count, err
	}

	zaps.Info("<<< get materials report list done")

	return emiList, count, err
}

func GetMaterialReportListByMonthVendor(month string, vendor string) ([]common.ERPMaterialInfo, int, error) {

	var emiList []common.ERPMaterialInfo
	var count int

	zaps.Infof(">>> get materials report by m(%s) v(%s)", month, vendor)

	vcode := common.GetLongVendorCode(vendor)

	rows, err := db.Query("SELECT material, material_txt, vendor, "+
		"vendor_name, vendor_flag, buyin, buyin_wtax, "+
		"returned, returned_wtax, rebate, pconf_qty, zpurps, "+
		"inv_qty, cost, cost_wtax, netval_inv, gross_val,"+
		"open_zinvsl, open_zinvcost, close_zinvsl, "+
		"close_zinvcost "+
		"FROM materials_report "+
		"WHERE month = ? AND vendor = ?", month, vcode)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return emiList, count, err
	}

	defer rows.Close()

	for rows.Next() {
		var v common.ERPMaterialInfo
		err := rows.Scan(&v.Material, &v.MaterialTxt, &v.Vendor,
			&v.VendorName, &v.VendorFlag, &v.Buyin, &v.BuyinWtax,
			&v.Return, &v.ReturnWtax, &v.Rebate, &v.PconfQty,
			&v.Zpurps, &v.InvQty, &v.Cost, &v.CostWtax, &v.NetvalInv,
			&v.GrossVal, &v.OpenZinvsl, &v.OpenZinvCost,
			&v.CloseZinvsl, &v.CloseZinvCost)
		if err != nil {
			zaps.Error("query error: ", err)
			return emiList, count, err
		} else {
			zaps.Debug(">>> month: ", month)

			emiList = append(emiList, v)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return emiList, count, err
	}

	zaps.Info("<<< get materials report list done")

	return emiList, count, err
}

func GetPurchaseReportGroupByOiEbeln(beginDate string, endDate string) ([]common.ERPVendorAccountStatement, int, error) {
	var emiList []common.ERPVendorAccountStatement
	var count int

	zaps.Infof(">>> get vendor account statement from day(%s) to day(%s)", beginDate, endDate)

	rows, err := db.Query("SELECT d.NAME, a.pstng_date,"+
		"SUM( CASE movetype WHEN '101' THEN quant_b WHEN '102' THEN - quant_b "+
		"WHEN '161' THEN - quant_b WHEN '162' THEN quant_b END ) AS buy_count,"+
		"SUM( CASE movetype  WHEN '101' THEN ROUND(( 1+a.ztaxrate * 0.01 )* value_lc, 2 )"+
		"WHEN '102' THEN - ROUND( ABS( ( 1+a.ztaxrate * 0.01 )* value_lc ), 2 )"+
		"WHEN '161' THEN - ROUND( ABS( ( 1+a.ztaxrate * 0.01 )* b.zpoamount ), 2 )"+
		"WHEN '162' THEN Round( ABS( ( 1+a.ztaxrate * 0.01 )* b.zpoamount ), 2 )"+
		"END  ) AS buy_money,a.oi_ebeln  FROM zinv_d01cg a "+
		"LEFT JOIN ( SELECT DISTINCT oi_ebeln, vendor, material, MAX( zpoamount )"+
		"AS zpoamount FROM zifpurd GROUP BY oi_ebeln, vendor, material )"+
		"b ON a.oi_ebeln = b.oi_ebeln  AND a.material = b.material "+
		"LEFT JOIN material c ON a.material = c.material "+
		"LEFT JOIN vendor d ON b.vendor = d.vendor  WHERE pstng_date BETWEEN ? AND ? "+
		"AND ( movetype IN ( '101', '102', '161', '162' ) AND b.oi_ebeln LIKE '46%%' ) "+
		"GROUP BY a.oi_ebeln, d.NAME, a.pstng_date ORDER BY d.NAME DESC", beginDate, endDate)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return emiList, count, err
	}

	defer rows.Close()

	for rows.Next() {
		var v common.ERPVendorAccountStatement
		err := rows.Scan(&v.VendorName, &v.PstngDate, &v.BuyCount,
			&v.BuyMoney, &v.OiEbeln)
		if err != nil {
			zaps.Error("query error: ", err)
			return emiList, count, err
		} else {
			zaps.Debug(">>> beginDate: ", beginDate)
			emiList = append(emiList, v)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return emiList, count, err
	}

	zaps.Info("<<< get vendor account statement done")

	return emiList, count, err
}
