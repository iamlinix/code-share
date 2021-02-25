package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"errors"
	"fmt"
	"time"
)

///////////////////////////////////////////////////////////////////////////////
// Materials Purchase

func GetERPMatlZifpurd(fr common.FilterReq) ([]common.ERPMatlZifpurdInfo,
	int, error) {

	var mzList []common.ERPMatlZifpurdInfo
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get material zifpurd list")

	sqlcmd := fmt.Sprintf("SELECT a.material, b.material_txt, "+
		"b.vendor_code, b.vendor_name, "+
		"sum(a.pconf_qty), sum(a.zpurps) "+
		"FROM zifpurd a "+
		"LEFT JOIN (SELECT DISTINCT(t1.material) as material, t4.materialtxt as material_txt, t3.name as vendor_name, "+
		"t1.vendor as vendor_code FROM zifpurd AS t1 INNER JOIN "+
		"(SELECT material, MAX(zgr_date) AS maxdate FROM zifpurd WHERE zgr_date "+
		"<= '%s' AND recordmode = '' GROUP BY material) AS t2 ON "+
		"t1.material = t2.material AND t1.zgr_date = t2.maxdate INNER "+
		"JOIN vendor t3 ON t1.vendor = t3.vendor AND t1.recordmode = '' INNER JOIN material AS t4 ON t1.material = t4.material) AS b "+
		"ON a.material = b.material "+
		"WHERE a.zgr_date >= '%s' and a.zgr_date <= '%s' "+
		"AND a.zpurps > 0 "+
		"GROUP BY a.material, b.vendor_code", fr.EndDate, fr.BeginDate, fr.EndDate)

	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return mzList, 0, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		var mz common.ERPMatlZifpurdInfo
		err := rows.Scan(&mz.Material, &mz.MaterialTxt, &mz.Vendor,
			&mz.VendorName, &mz.PconfQty, &mz.Zpurps)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> material code: ", mz.Material)
			zaps.Debug(">>> pconf qty: ", mz.PconfQty)
			zaps.Debug(">>> zpurps: ", mz.Zpurps)

			mzList = append(mzList, mz)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return mzList, count, err
	}

	zaps.Info("<<< get materials zifpurd done with ", count)

	return mzList, count, err
}

func GetOIEbelnFromZifpurd(begin string, end string) (string, string, error) {

	var minoi sql.NullString
	var maxoi sql.NullString
	var sminoi, smaxoi string

	sql := fmt.Sprintf("SELECT MIN(oi_ebeln), MAX(oi_ebeln) "+
		"FROM zifpurd "+
		"WHERE zgr_date BETWEEN '%s' and '%s'",
		begin, end)

	err := db.QueryRow(sql).Scan(&minoi, &maxoi)
	if err != nil {
		zaps.Error("query zifpurd oi ebeln failed: ", err)
		return "N/A", "N/A", err
	}

	if minoi.Valid != true || maxoi.Valid != true {
		zaps.Error("oi ebeln not valid failed: ", err)
		return "N/A", "N/A", errors.New("not valid")
	}

	sminoi = minoi.String
	smaxoi = maxoi.String

	return sminoi, smaxoi, nil
}

func GetERPVendorFromZifpurd(fr common.FilterReq) ([]common.ERPMatlVendorInfo,
	int, error) {

	var mvList []common.ERPMatlVendorInfo
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get material vendor list")

	sqlcmd := fmt.Sprintf("SELECT t1.material, t1.vendor, t3.name "+
		"FROM zifpurd AS t1 "+
		"INNER JOIN ("+
		"SELECT material, MAX(zgr_date) AS maxdate "+
		"FROM zifpurd "+
		"WHERE zgr_date <= '%s' AND recordmode = '' "+
		"GROUP BY material) AS t2 "+
		"ON t1.material = t2.material AND t1.zgr_date = t2.maxdate "+
		"INNER JOIN vendor t3 ON t1.vendor = t3.vendor "+
		"AND t1.recordmode = '' "+
		"ORDER BY t1.material", fr.EndDate)

	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return mvList, 0, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		var mv common.ERPMatlVendorInfo
		err := rows.Scan(&mv.Material, &mv.Vendor, &mv.VendorName)
		if err != nil {
			zaps.Error("query error: ", err)
			return mvList, count, err
		}

		mvList = append(mvList, mv)
		count++
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return mvList, count, err
	}

	zaps.Info("<<< get materials vendor done with ", count)
	return mvList, count, err
}

func GetERPMatlPurchase2(fr common.FilterReq) ([]common.ERPMatlPurchaseInfo,
	int, error) {

	var mpList []common.ERPMatlPurchaseInfo
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get material purchase list")

	min_oi, max_oi, err := GetOIEbelnFromZifpurd(fr.BeginDate, fr.EndDate)
	if err != nil {
		zaps.Error("get oi ebeln failed: ", err)
		return mpList, 0, err
	}

	sqlcmd := fmt.Sprintf("SELECT a.material, b.material_txt, "+
		"b.vendor, d.name, SUM( "+
		"CASE WHEN movetype = '101' THEN value_lc "+
		"ELSE CASE WHEN movetype = '102' THEN -ABS(value_lc) "+
		"ELSE 0 END END) AS buyin, "+
		"SUM(CASE WHEN movetype = '161' THEN b.zpoamount "+
		"ELSE CASE WHEN movetype = '162' THEN -ABS(b.zpoamount) "+
		"ELSE 0 END END) AS returned "+
		"FROM zinv_d01cg a "+
		"LEFT JOIN (SELECT DISTINCT oi_ebeln, vendor, material, "+
		"MAX(zpoamount) AS zpoamount FROM zifpurd GROUP BY oi_ebeln, vendor, material) b "+
		"ON a.material = b.material AND a.oi_ebeln = b.oi_ebeln "+
		"LEFT JOIN material c ON a.material = c.material "+
		"LEFT JOIN vendor d ON b.vendor = d.vendor "+
		"WHERE pstng_date BETWEEN '%s' AND '%s' "+
		"AND ((dcindic = 'S' AND movetype = '101') OR "+
		"(dcindic = 'H' AND movetype='102') OR "+
		"(movetype IN ('161', '162') AND processkey = '101')) "+
		"AND ((b.oi_ebeln LIKE '46%%' AND b.oi_ebeln BETWEEN '%s' AND '%s') "+
		"OR zbsart = 'Z031') GROUP BY a.material, b.vendor "+
		"ORDER BY a.material",
		fr.BeginDate, fr.EndDate, min_oi, max_oi)

	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err = db.Query(sqlcmd)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return mpList, 0, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		var mp common.ERPMatlPurchaseInfo
		err := rows.Scan(&mp.Material, &mp.MaterialTxt, &mp.Vendor,
			&mp.VendorName, &mp.Buyin, &mp.Return)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> material code: ", mp.Material)
			/*
				zaps.Debug(">>> pconf qty: ", mp.PconfQty)
				zaps.Debug(">>> pconf zje: ", mp.PconfZje)
				zaps.Debug(">>> zpurps: ", mp.Zpurps)
				zaps.Debug(">>> zzje: ", mp.Zzje)
			*/
			if mp.Buyin != 0 || mp.Return != 0 {
				mpList = append(mpList, mp)
				count++
			}
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return mpList, count, err
	}

	zaps.Info("<<< get materials purchase done with ", count)

	return mpList, count, err
}

func GetERPMatlPurchase(fr common.FilterReq) ([]common.ERPMatlPurchaseInfo,
	int, error) {

	var mpList []common.ERPMatlPurchaseInfo
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get material purchase list")

	sqlcmd := fmt.Sprintf("SELECT a.material, c.materialtxt, "+
		"b.vendor, d.name, "+
		"SUM(CASE WHEN movetype = '101' THEN quant_b "+
		"ELSE CASE WHEN movetype = '102' "+
		"THEN -quant_b ELSE 0 END END) AS buy_count, "+
		"SUM(CASE WHEN movetype = '101' THEN value_lc "+
		"ELSE CASE WHEN movetype = '102' "+
		"THEN -ABS(value_lc) ELSE 0 END END) AS buy_money, "+
		"SUM(CASE WHEN movetype = '101' THEN "+
		"ROUND(value_lc*(1+ztaxrate/100), 2) "+
		"ELSE CASE WHEN movetype = '102' "+
		"THEN ROUND(-ABS(value_lc)*(1+ztaxrate/100),2) "+
		"ELSE 0 END END) AS buy_money_wtax, "+
		"SUM(CASE WHEN movetype = '161' THEN quant_b "+
		"ELSE CASE WHEN movetype = '162' THEN -quant_b "+
		"ELSE 0 END END) AS rt_count, "+
		"SUM(CASE WHEN movetype = '161' THEN b.zpoamount "+
		"ELSE CASE WHEN movetype = '162' "+
		"THEN -ABS(b.zpoamount) ELSE 0 END END) AS rt_money, "+
		"SUM(CASE WHEN movetype = '161' "+
		"THEN b.zpoamount*(1+ztaxrate/100) "+
		"ELSE CASE WHEN movetype = '162' "+
		"THEN -ABS(b.zpoamount)*(1+ztaxrate/100) "+
		"ELSE 0 END END) AS rt_money_wtax "+
		"FROM zinv_d01cg a "+
		"LEFT JOIN (SELECT DISTINCT oi_ebeln, vendor, material, "+
		"MAX(zpoamount) AS zpoamount "+
		"FROM zifpurd WHERE zgr_date BETWEEN '%s' AND '%s' GROUP BY oi_ebeln, "+
		"vendor, material) b ON a.oi_ebeln = b.oi_ebeln AND a.material = b.material "+
		"LEFT JOIN material c ON a.material = c.material "+
		"LEFT JOIN vendor d ON b.vendor = d.vendor "+
		"WHERE a.pstng_date BETWEEN '%s' AND '%s' "+
		"AND (movetype IN ('101', '102', '161', '162') AND b.oi_ebeln LIKE '46%%') "+
		"GROUP BY a.material, b.vendor ORDER BY a.material",
		fr.BeginDate, fr.EndDate, fr.BeginDate, fr.EndDate)

	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return mpList, 0, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		var mp common.ERPMatlPurchaseInfo
		err := rows.Scan(&mp.Material, &mp.MaterialTxt, &mp.Vendor,
			&mp.VendorName, &mp.BuyinCnt, &mp.Buyin,
			&mp.BuyinWtax, &mp.ReturnCnt, &mp.Return,
			&mp.ReturnWtax)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> material code: ", mp.Material)
			/*
				zaps.Debug(">>> pconf qty: ", mp.PconfQty)
				zaps.Debug(">>> pconf zje: ", mp.PconfZje)
				zaps.Debug(">>> zpurps: ", mp.Zpurps)
				zaps.Debug(">>> zzje: ", mp.Zzje)
			*/
			if mp.Buyin != 0 || mp.Return != 0 {
				mpList = append(mpList, mp)
				count++
			}
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return mpList, count, err
	}

	zaps.Info("<<< get materials purchase done with ", count)

	return mpList, count, err
}

// Materials Sales
func GetERPMaterialSales(fr common.FilterReq) ([]common.ERPMatlSalesInfo,
	int, error) {

	var msList []common.ERPMatlSalesInfo
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get material sales list")

	sqlcmd := fmt.Sprintf("SELECT a.material, b.material_txt, "+
		"b.vendor_code, b.vendor_name, "+
		"sum(a.inv_qty) AS INVQTY, sum(a.cost) AS COST, "+
		"ROUND(sum(a.cost*(1+c.bic_ztax_rat/100)), 2) "+
		"AS COSTWTAX, "+
		"sum(a.netval_inv) AS NETVALINV, "+
		"sum(a.gross_val) AS GROSSVAL "+
		"FROM bill_zsd a "+
		"LEFT JOIN (SELECT DISTINCT(t1.material) as material, t4.materialtxt as material_txt, t3.name as vendor_name, "+
		"t1.vendor as vendor_code FROM zifpurd AS t1 INNER JOIN "+
		"(SELECT material, MAX(zgr_date) AS maxdate FROM zifpurd WHERE zgr_date "+
		"<= '%s' AND recordmode = '' GROUP BY material) AS t2 ON "+
		"t1.material = t2.material AND t1.zgr_date = t2.maxdate INNER "+
		"JOIN vendor t3 ON t1.vendor = t3.vendor AND t1.recordmode = '' INNER JOIN material AS t4 ON t1.material = t4.material) AS b "+
		"ON a.material = b.material "+
		"LEFT JOIN material c ON a.material = c.material "+
		"WHERE calday BETWEEN '%s' AND '%s' "+
		"GROUP BY a.material, b.vendor_code",
		fr.EndDate, fr.BeginDate, fr.EndDate)

	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return msList, 0, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		var ms common.ERPMatlSalesInfo
		err := rows.Scan(&ms.Material, &ms.MaterialTxt, &ms.Vendor,
			&ms.VendorName, &ms.InvQty, &ms.Cost,
			&ms.CostWtax, &ms.NetvalInv, &ms.GrossVal)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> Material: ", ms.Material)

			msList = append(msList, ms)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return msList, count, err
	}

	zaps.Info("<<< get materials sales done with count ", count)

	return msList, count, err
}

// Materials Open Inv
/*
func GetERPMatlOpenInv(fr common.FilterReq) ([]common.ERPMatlInvInfo,
	int, error) {

	var msList []common.ERPMatlInvInfo
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get material open inv")

	sqlcmd := fmt.Sprintf("SELECT a.material, b.material_txt, " +
		"b.vendor_code, b.vendor_name, " +
		"SUM(bic_zinvsl) as BICZINVSL, " +
		"SUM(bic_zinvcost) as BICZINVCOST " +
		"FROM zinv a " +
		"LEFT JOIN bible b ON a.material = b.material " +
		"WHERE calday = '%s' " +
		"GROUP BY a.material",
		fr.BeginDate)

	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return msList, 0, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		var ms common.ERPMatlInvInfo
		err := rows.Scan(&ms.Material, &ms.MaterialTxt, &ms.Vendor,
				&ms.VendorName, &ms.Zinvsl, &ms.ZinvCost)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> Material: ", ms.Material)
			zaps.Debug(">>> Zinvsl: ", ms.Zinvsl)
			zaps.Debug(">>> ZinvCost: ", ms.ZinvCost)

			msList = append(msList, ms)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return msList, count, err
	}

	zaps.Info("<<< get materials open inv done with count ", count)

	return msList, count, err
}


func GetERPMatlCloseInv(fr common.FilterReq) ([]common.ERPMatlInvInfo,
	int, error) {

	var msList []common.ERPMatlInvInfo
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get material close inv")

	sqlcmd := fmt.Sprintf("SELECT a.material, b.material_txt, " +
		"b.vendor_code, b.vendor_name, " +
		"SUM(bic_zinvsl) as BICZINVSL, " +
		"SUM(bic_zinvcost) as BICZINVCOST " +
		"FROM zinv a " +
		"LEFT JOIN bible b ON a.material = b.material " +
		"WHERE calday = '%s' " +
		"GROUP BY a.material",
		fr.EndDate)

	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return msList, 0, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		var ms common.ERPMatlInvInfo
		err := rows.Scan(&ms.Material, &ms.MaterialTxt, &ms.Vendor,
				&ms.VendorName, &ms.Zinvsl, &ms.ZinvCost)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> Material: ", ms.Material)
			zaps.Debug(">>> Zinvsl: ", ms.Zinvsl)
			zaps.Debug(">>> ZinvCost: ", ms.ZinvCost)

			msList = append(msList, ms)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return msList, count, err
	}

	zaps.Info("<<< get materials close inv done with count ", count)

	return msList, count, err
}

*/

func GetERPMatlOpenInv(fr common.FilterReq) ([]common.ERPMatlInvInfo,
	int, error) {

	var msList []common.ERPMatlInvInfo
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get material open inv")

	ssdate, err := GetLatestSSDate(fr.BeginDate)
	if err != nil {
		zaps.Errorf("get latest ss date failed: ", err)
		return msList, 0, err
	}

	sss, err := GetInvSnapshot(ssdate)
	if err != nil {
		zaps.Errorf("get latest inv ss failed: ", err)
		return msList, 0, err
	}

	day := common.GetSubDay(fr.BeginDate, 1)

	sqlcmd := fmt.Sprintf("SELECT a.material, b.material_txt, "+
		"b.vendor_code, b.vendor_name, "+
		"SUM(CASE WHEN dcindic = 'S' THEN quant_b ELSE -quant_b END) "+
		"AS ZINVSL, "+
		"ROUND(SUM(CASE WHEN dcindic = 'S' THEN "+
		"value_lc ELSE -value_lc END),2) AS ZINVCOST "+
		"FROM zinv_d01cg a "+
		"LEFT JOIN (SELECT DISTINCT(t1.material) as material, t4.materialtxt as material_txt, t3.name as vendor_name, "+
		"t1.vendor as vendor_code FROM zifpurd AS t1 INNER JOIN "+
		"(SELECT material, MAX(zgr_date) AS maxdate FROM zifpurd WHERE zgr_date "+
		"<= '%s' AND recordmode = '' GROUP BY material) AS t2 ON "+
		"t1.material = t2.material AND t1.zgr_date = t2.maxdate INNER "+
		"JOIN vendor t3 ON t1.vendor = t3.vendor AND t1.recordmode = '' INNER JOIN material AS t4 ON t1.material = t4.material) AS b ON a.material = b.material "+
		"WHERE a.pstng_date BETWEEN '%s' AND '%s' "+
		"AND a.plant NOT LIKE 'AA%%' "+
		"GROUP BY a.material, b.vendor_code",
		day, ssdate, day)

	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err = db.Query(sqlcmd)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return msList, 0, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		var ms common.ERPMatlInvInfo
		err := rows.Scan(&ms.Material, &ms.MaterialTxt, &ms.Vendor,
			&ms.VendorName, &ms.Zinvsl, &ms.ZinvCost)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> Material: ", ms.Material)
			zaps.Debug(">>> Zinvsl: ", ms.Zinvsl)
			zaps.Debug(">>> ZinvCost: ", ms.ZinvCost)

			//lookup snapshot
			ssm := LookupSSMaterial(sss, ms.Material)
			ms.Zinvsl += ssm.Zinvsl
			ms.ZinvCost += ssm.ZinvCost

			msList = append(msList, ms)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return msList, count, err
	}

	zaps.Info("<<< get materials open inv done with count ", count)

	return msList, count, err
}

func GetERPMatlCloseInv(fr common.FilterReq) ([]common.ERPMatlInvInfo,
	int, error) {

	var msList []common.ERPMatlInvInfo
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get material close inv")

	ssdate, err := GetLatestSSDate(fr.EndDate)
	if err != nil {
		zaps.Errorf("get latest ss date failed: ", err)
		return msList, 0, err
	}

	sss, err := GetInvSnapshot(ssdate)
	if err != nil {
		zaps.Errorf("get latest inv ss failed: ", err)
		return msList, 0, err
	}

	sqlcmd := fmt.Sprintf("SELECT a.material, b.material_txt, "+
		"b.vendor_code, b.vendor_name, "+
		"SUM(CASE WHEN dcindic = 'S' THEN quant_b ELSE -quant_b END) "+
		"AS ZINVSL, "+
		"ROUND(SUM(CASE WHEN dcindic = 'S' THEN "+
		"value_lc ELSE -value_lc END),2) AS ZINVCOST "+
		"FROM zinv_d01cg a "+
		"LEFT JOIN (SELECT DISTINCT(t1.material) as material, t4.materialtxt as material_txt, t3.name as vendor_name, "+
		"t1.vendor as vendor_code FROM zifpurd AS t1 INNER JOIN "+
		"(SELECT material, MAX(zgr_date) AS maxdate FROM zifpurd WHERE zgr_date "+
		"<= '%s' AND recordmode = '' GROUP BY material) AS t2 ON "+
		"t1.material = t2.material AND t1.zgr_date = t2.maxdate INNER "+
		"JOIN vendor t3 ON t1.vendor = t3.vendor AND t1.recordmode = '' INNER JOIN material AS t4 ON t1.material = t4.material) AS b ON a.material = b.material "+
		"WHERE a.pstng_date BETWEEN '%s' AND '%s' "+
		"AND a.plant NOT LIKE 'AA%%' "+
		"GROUP BY a.material, b.vendor_code",
		fr.EndDate, ssdate, fr.EndDate)

	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err = db.Query(sqlcmd)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return msList, 0, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		var ms common.ERPMatlInvInfo
		err := rows.Scan(&ms.Material, &ms.MaterialTxt, &ms.Vendor,
			&ms.VendorName, &ms.Zinvsl, &ms.ZinvCost)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> Material: ", ms.Material)
			zaps.Debug(">>> Zinvsl: ", ms.Zinvsl)
			zaps.Debug(">>> ZinvCost: ", ms.ZinvCost)

			//lookup snapshot
			ssm := LookupSSMaterial(sss, ms.Material)
			ms.Zinvsl += ssm.Zinvsl
			ms.ZinvCost += ssm.ZinvCost

			msList = append(msList, ms)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return msList, count, err
	}

	zaps.Info("<<< get materials close inv done with count ", count)

	return msList, count, err
}

func GetLatestSSDate(date string) (string, error) {

	var pstng string

	sqlcmd := fmt.Sprintf("SELECT MAX(pstng_date) FROM stock_snapshot "+
		"WHERE pstng_date < '%s'", date)

	err := db.QueryRow(sqlcmd).Scan(&pstng)
	if err != nil {
		zaps.Error("query latest pstng failed: ", err)
		return "", err
	}

	return pstng, nil
}

func GetInvSnapshot(pstng string) ([]common.ERPInvSnapshot, error) {

	var ssList []common.ERPInvSnapshot
	var rows *sql.Rows

	sqlcmd := fmt.Sprintf("SELECT material, zinvsl, zinvcost "+
		"FROM stock_snapshot WHERE pstng_date = '%s'", pstng)

	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return ssList, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		var ss common.ERPInvSnapshot
		err := rows.Scan(&ss.Material, &ss.Zinvsl, &ss.ZinvCost)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> Material: ", ss.Material)
			zaps.Debug(">>> Zinvsl: ", ss.Zinvsl)
			zaps.Debug(">>> ZinvCost: ", ss.ZinvCost)

			ssList = append(ssList, ss)
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return ssList, err
	}

	zaps.Info("<<< get inv snapshot done")

	return ssList, err
}

/* XXX TODO: !performance!*/
func LookupSSMaterial(ssList []common.ERPInvSnapshot,
	material string) common.ERPInvSnapshot {

	var res common.ERPInvSnapshot

	for _, s := range ssList {
		if s.Material == material {
			return s
		}
	}

	return res
}
