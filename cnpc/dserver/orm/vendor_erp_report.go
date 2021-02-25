package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"time"
)

func AddVendorReportTx(fr common.FilterReq,
	evis []common.ERPVendorInfo) ([]common.VendorReport,
	common.VendorReport, error) {

	var vrList []common.VendorReport
	var sum common.VendorReport

	t1 := time.Now().UnixNano() / 1e6

	lastm, err := common.GetLastMonth(fr.Month)
	if err != nil {
		zaps.Errorf("get last month failed: %v", err)
		return vrList, sum, err
	}

	zaps.Info("query last month: ", lastm)

	lmvrs, findL, err := GetVendorReportByMonth(lastm)
	if err != nil {
		zaps.Errorf("get vendors report by month failed: %v", err)
		return vrList, sum, err
	}
	if findL != true {
		zaps.Warnf("vendor for last month not exist, using config")

		cfgs, _, err := GetVendorCfgList(1, 0)
		if err != nil {
			zaps.Errorf("get vendors cfg list failed: %v", err)
			return vrList, sum, err
		}

		lmvrs = TransCfgToReport(cfgs)
	}

	/* lookup first month of current year */
	firstm, err := common.GetFirstMonth(fr.Month)
	if err != nil {
		zaps.Errorf("get first month failed: %v", err)
		return vrList, sum, err
	}

	fmvrs, findF, err := GetVendorReportByMonth(firstm)
	if err != nil {
		zaps.Errorf("get vendors report by month failed: %v", err)
		return vrList, sum, err
	}
	if findF != true {
		zaps.Warnf("vendor for first month not exist")
	}

	zaps.Warnf("find F: ", findF)

	/* get vendor material count from bible */
	vmcList, _, err := GetVendorMatlCountList()
	if err != nil {
		zaps.Errorf("get vendors material count list failed: %v", err)
		return vrList, sum, err
	}

	tx, _ := db.Begin()
	defer tx.Rollback()

	for _, v := range evis {

		var vr common.VendorReport

		//looking for vr for last month
		lmvr := LookupVendorReport(lmvrs, v.Vendor)
		zaps.Infof("last month vendor (%s) s: %f, p: %f, up: %f",
			v.Vendor, lmvr.MonthUnpaid, lmvr.MonthSurplus)

		//looking for vr for first month
		fmvr := LookupVendorReport(fmvrs, v.Vendor)
		zaps.Infof("first month vendor (%s) os: %f",
			v.Vendor, fmvr.MonthOpenInv)

		if v.VendorName == "N/A" {
			cfg, _, _ := GetVendorCfgByCode(v.Vendor)
			v.VendorName = cfg.VendorName
		}

		vr.Month = fr.Month
		vr.BeginDate = fr.BeginDate
		vr.EndDate = fr.EndDate
		vr.VendorCode = v.Vendor
		vr.VendorName = v.VendorName

		vr.MonthSalesQty = v.InvQty
		vr.MonthSales = v.CostWtax
		vr.MonthSalesOrig = vr.MonthSales
		vr.MonthSales2 = v.Cost
		vr.MonthOpenInv = v.OpenZinvCost
		vr.MonthCloseInv = v.CloseZinvCost
		vr.MonthPurchase = v.PurWtaxVal //with tax
		vr.MonthPurchaseOrig = vr.MonthPurchase
		vr.MonthPurchase2 = v.PurVal
		vr.MonthRebate = v.RebateVal
		vr.MonthPurFinal = v.PurFinalVal
		vr.MonthReceipt = v.PurFinalVal

		if v.Zpurps > 0 {
			vr.OrderFillRate =
				100.0 * float64(v.PconfQty) / float64(v.Zpurps)
		}
		if v.NetvalInv > 0 {
			vr.GrossMargin =
				100.0 * (v.NetvalInv - v.Cost) / v.NetvalInv
		}

		//LastMonthUnpaid is earlyMonthUnpaid
		vr.LastMonthUnpaid = lmvr.YearUnpaid
		//vr.MonthPaid
		if vr.MonthSales > vr.MonthReceipt {
			if vr.MonthSales > vr.MonthReceipt+vr.LastMonthUnpaid {
				vr.MonthPaidOrig = vr.MonthReceipt + vr.LastMonthUnpaid
			} else {
				vr.MonthPaidOrig = vr.MonthSales
			}
		} else {
			vr.MonthPaidOrig = vr.MonthSales
		}
		vr.MonthPaid = vr.MonthPaidOrig

		//vr.MonthUnpaid
		if vr.MonthReceipt > vr.MonthPaid {
			vr.MonthUnpaid = vr.MonthReceipt - vr.MonthPaid
		} else {
			vr.MonthUnpaid = 0
		}
		vr.MonthUnpaidOrig = vr.MonthUnpaid

		/* calc some metrics */
		/* 2021-01-15 lvxs left these out because a new calculation method is adopted*/
		if vr.MonthSales+lmvr.MonthSurplus >
			vr.MonthPurFinal+lmvr.MonthUnpaid {
			//vr.MonthPaid = vr.MonthPurFinal + lmvr.MonthUnpaid
			//vr.MonthPaidOrig = vr.MonthPaid
			//vr.MonthUnpaid = 0
			//vr.MonthUnpaidOrig = vr.MonthUnpaid
			vr.MonthSurplus = (vr.MonthSales + lmvr.MonthSurplus) -
				vr.MonthPaid
		} else {
			//vr.MonthPaid = vr.MonthSales + lmvr.MonthSurplus
			//vr.MonthPaidOrig = vr.MonthPaid
			//vr.MonthUnpaid = vr.MonthReceipt - vr.MonthPaid
			//vr.MonthUnpaidOrig = vr.MonthUnpaid
		}
		vr.MonthSurplus = 0
		if findF == true {
			vr.YearOpenInv = fmvr.MonthOpenInv
		} else {
			vr.YearOpenInv = vr.MonthOpenInv //XXX TODO: something wrong
		}

		vr.YearPurchase = vr.MonthPurFinal + lmvr.YearPurchase
		vr.YearSalesQty = vr.MonthSalesQty + lmvr.YearSalesQty
		vr.YearSales = vr.MonthSales + lmvr.YearSales
		vr.YearPaid = vr.MonthPaid + lmvr.YearPaid
		vr.YearUnpaidOrig = vr.LastMonthUnpaid + vr.MonthUnpaid
		vr.YearUnpaid = vr.YearUnpaidOrig
		vr.YearSurplus = vr.MonthSurplus

		/* other metrics */
		vr.MonthAvgInv = (vr.MonthOpenInv + vr.MonthCloseInv) / 2.0
		vr.YearAvgInv = (vr.YearOpenInv + vr.MonthCloseInv) / 2.0
		if vr.MonthAvgInv == 0 {
			vr.MonthDaysTo = 0
		} else {
			vr.MonthDaysTo = 31.0 / (vr.MonthSales / vr.MonthAvgInv)
		}
		if vr.YearAvgInv == 0 {
			vr.YearDaysTo = 0
		} else {
			vr.YearDaysTo = 31.0 / (vr.YearSales / vr.YearAvgInv)
		}

		//lookup current vendor matl count
		mcnt := LookupVendorMatlCount(vmcList, v.Vendor)
		if mcnt > 0 {
			zaps.Infof("GrossVal: %f, Mcnt: %d", v.GrossVal, mcnt)
			vr.CommEff = v.GrossVal / float64(mcnt)
		}

		AddVendorReport(tx, vr)

		SumVendorReport(&sum, vr)

		vrList = append(vrList, vr)
	}

	/* Transaction END */
	err = tx.Commit()
	if err != nil {
		zaps.Errorf("transaction commit failed: %v", err)
		return vrList, sum, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("<<< add vendor report tx done using %d ms", t2-t1)

	return vrList, sum, nil
}

func LookupVendorMatlCount(vmcs []common.VendorMatlCount, vendor string) int {

	for _, v := range vmcs {
		if v.Vendor == vendor {
			return v.MatlCount
		}
	}

	return 0
}

func TransCfgToReport(cfg []common.VendorCfg) []common.VendorReport {

	var rlist []common.VendorReport

	for _, c := range cfg {

		var r common.VendorReport
		r.VendorCode = c.VendorCode
		r.VendorName = c.VendorName
		/**/
		r.MonthUnpaid = c.YearUnpaid   ///XXX NOTICE
		r.MonthSurplus = c.YearSurplus ///XXX NOTICE
		/**/
		r.YearPurchase = c.YearPurchase
		r.YearSales = c.YearSales
		r.YearPaid = c.YearPaid
		r.YearUnpaid = c.YearUnpaid
		r.YearUnpaidOrig = r.YearUnpaid
		r.YearSurplus = c.YearSurplus

		rlist = append(rlist, r)
	}

	return rlist
}

func LookupVendorReport(vrs []common.VendorReport,
	code string) common.VendorReport {

	var vr common.VendorReport

	for _, v := range vrs {
		if v.VendorCode == code {
			return v
		}
	}

	return vr
}

func SumVendorReport(sum *common.VendorReport, v common.VendorReport) {

	sum.YearOpenInv += v.YearOpenInv
	sum.MonthOpenInv += v.MonthOpenInv
	sum.MonthPurchase += v.MonthPurchase
	sum.MonthPurchaseOrig += v.MonthPurchaseOrig
	sum.MonthReceipt += v.MonthReceipt
	sum.MonthSalesQty += v.MonthSalesQty
	sum.MonthSales += v.MonthSales
	sum.MonthSalesOrig += v.MonthSalesOrig
	sum.MonthPaid += v.MonthPaid
	sum.MonthPaidOrig += v.MonthPaidOrig
	sum.MonthUnpaid += v.MonthUnpaid
	sum.MonthUnpaidOrig += v.MonthUnpaidOrig
	sum.MonthSurplus += v.MonthSurplus
	sum.MonthCloseInv += v.MonthCloseInv
	sum.YearPurchase += v.YearPurchase
	sum.YearSalesQty += v.YearSalesQty
	sum.YearSales += v.YearSales
	sum.YearPaid += v.YearPaid
	sum.YearUnpaid += v.YearUnpaid
	sum.YearUnpaidOrig += v.YearUnpaidOrig
	sum.YearSurplus += v.YearSurplus
}

func AddVendorReport(tx *sql.Tx, vr common.VendorReport) error {

	zaps.Info(">>> add vendors report")

	stmt, err := tx.Prepare("INSERT INTO vendors_report(month, begin_date, " +
		"end_date, vendor_code, vendor_name, " +
		"month_purchase, month_purchase_orig, month_receipt, month_rebate, month_pur_final, " +
		"month_purchase2, month_sales_qty, month_sales, month_sales_orig, " +
		"month_sales2, month_paid, month_paid_orig, last_month_unpaid, " +
		"month_unpaid, month_unpaid_orig, month_surplus, year_purchase, " +
		"year_sales_qty, year_sales, year_paid, " +
		"year_unpaid, year_unpaid_orig, year_surplus, year_open_stock, " +
		"month_open_stock, month_close_stock, " +
		"month_avg_stock, year_avg_stock, month_days_to, " +
		"year_days_to, comm_eff, gross_margin, " +
		"order_fill_rate, first_order_time, created_time) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, " +
		"?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, " +
		"?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		zaps.Error("insert into vendors report failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(vr.Month, vr.BeginDate, vr.EndDate, vr.VendorCode,
		vr.VendorName, vr.MonthPurchase, vr.MonthPurchaseOrig, vr.MonthReceipt, vr.MonthRebate,
		vr.MonthPurFinal, vr.MonthPurchase2, vr.MonthSalesQty,
		vr.MonthSales, vr.MonthSalesOrig, vr.MonthSales2, vr.MonthPaid, vr.MonthPaidOrig,
		vr.LastMonthUnpaid, vr.MonthUnpaid, vr.MonthUnpaidOrig, vr.MonthSurplus,
		vr.YearPurchase, vr.YearSalesQty, vr.YearSales,
		vr.YearPaid, vr.YearUnpaid, vr.YearUnpaidOrig, vr.YearSurplus,
		vr.YearOpenInv, vr.MonthOpenInv, vr.MonthCloseInv,
		vr.MonthAvgInv, vr.YearAvgInv, vr.MonthDaysTo,
		vr.YearDaysTo, vr.CommEff, vr.GrossMargin, vr.OrderFillRate,
		vr.FirstOrder, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		zaps.Error("db exec failed: ", err)
		return err
	}

	zaps.Info("<<< add vendors report done")

	return err
}

func UpdateVendorReportHistory(vr *common.VendorReportUpdate) error {
	stmt, err := db.Prepare("INSERT INTO vendor_report_history (month, type, vendor, old_value, new_value, user, reason) VALUES " +
		"(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		zaps.Error("update vendors report history failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(vr.Month, vr.Type, vr.Vendor, vr.OldValue, vr.NewValue, vr.User, vr.Reason)
	if err != nil {
		zaps.Error("db exec failed: ", err)
		return err
	}

	return nil
}

func GetVendorReportHistory(month, vendor string, tp int) ([]common.VendorReportUpdate, error) {
	rows, err := db.Query("SELECT month, type, vendor, old_value, new_value, user, reason, mtime "+
		" FROM vendor_report_history WHERE month = ? AND type = ? AND vendor = ? ORDER BY mtime DESC",
		month, tp, vendor)
	if err != nil {
		zaps.Error("failed to get vendor report history:", err)
		return nil, err
	}
	defer rows.Close()

	var history []common.VendorReportUpdate
	for rows.Next() {
		var h common.VendorReportUpdate
		if err = rows.Scan(&h.Month, &h.Type, &h.Vendor, &h.OldValue, &h.NewValue, &h.User, &h.Reason, &h.MTime); err != nil {
			zaps.Error("failed to scan vendor report history: ", err)
			return nil, err
		}
		history = append(history, h)
	}

	return history, nil
}

func UpdateVendorReport(tp int, vr common.VendorReport) error {

	zaps.Info(">>> update vendors report")

	if tp == 0 {
		stmt, err := db.Prepare("UPDATE vendors_report SET " +
			"month_paid = ?, month_unpaid = ?, " +
			"year_paid = ?, year_unpaid = ? " +
			"WHERE month = ? AND vendor_code = ?")
		if err != nil {
			zaps.Error("update vendors report failed: ", err)
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(vr.MonthPaid, vr.MonthUnpaid, vr.YearPaid,
			vr.YearUnpaid, vr.Month, vr.VendorCode)
		if err != nil {
			zaps.Error("db exec failed: ", err)
			return err
		}
	} else if tp == 1 {
		stmt, err := db.Prepare("UPDATE vendors_report SET " +
			"month_receipt = ?, month_unpaid = ?, year_unpaid = ? WHERE month = ? AND vendor_code = ?")
		if err != nil {
			zaps.Error("update vendors report failed: ", err)
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(vr.MonthReceipt, vr.MonthUnpaid, vr.YearUnpaid, vr.Month, vr.VendorCode)
		if err != nil {
			zaps.Error("db exec failed: ", err)
			return err
		}
	} else if tp == 2 {
		stmt, err := db.Prepare("UPDATE vendors_report SET " +
			"month_sales = ?, year_sales = ?, month_surplus = ?, year_surplus = ?, month_days_to = ? WHERE month = ? AND vendor_code = ?")
		if err != nil {
			zaps.Error("update vendors report failed: ", err)
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(vr.MonthSales, vr.YearSales, vr.MonthSurplus, vr.YearSurplus, vr.MonthDaysTo, vr.Month, vr.VendorCode)
		if err != nil {
			zaps.Error("db exec failed: ", err)
			return err
		}
	} else if tp == 3 {
		stmt, err := db.Prepare("UPDATE vendors_report SET " +
			"year_unpaid = ? WHERE month = ? AND vendor_code = ?")
		if err != nil {
			zaps.Error("update vendors report failed: ", err)
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(vr.YearUnpaid, vr.Month, vr.VendorCode)
		if err != nil {
			zaps.Error("db exec failed: ", err)
			return err
		}
	} else if tp == 4 {
		stmt, err := db.Prepare("UPDATE vendors_report SET " +
			"month_purchase = ? WHERE month = ? AND vendor_code = ?")
		if err != nil {
			zaps.Error("update vendors report failed: ", err)
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(vr.MonthPurchase, vr.Month, vr.VendorCode)
		if err != nil {
			zaps.Error("db exec failed: ", err)
			return err
		}
	}

	zaps.Info("<<< update vendors report done")

	return nil
}

func DelVendorReport(month string) {

	zaps.Info(">>> del vendors report")

	_, err := db.Exec("DELETE FROM vendors_report WHERE month = ?", month)
	if err != nil {
		zaps.Error("delete vendors report exec failed: ", err)
	}

	zaps.Info("<<< del vendors report done")
}

func GetVendorReportList() ([]common.VendorReport, int, error) {

	var vrList []common.VendorReport
	var count int

	zaps.Info(">>> get vendors report list")

	rows, err := db.Query("SELECT distinct month, begin_date, end_date " +
		"FROM vendors_report")
	if err != nil {
		zaps.Error("db query failed: ", err)
		return vrList, count, err
	}

	defer rows.Close()

	for rows.Next() {
		var vr common.VendorReport
		err := rows.Scan(&vr.Month, &vr.BeginDate, &vr.EndDate)
		if err != nil {
			zaps.Error("query error: ", err)
			return vrList, count, err
		} else {
			zaps.Debug(">>> month: ", vr.Month)
			zaps.Debug(">>> start date: ", vr.BeginDate)
			zaps.Debug(">>> end date: ", vr.EndDate)

			vrList = append(vrList, vr)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return vrList, count, err
	}

	zaps.Info("<<< get vendors report list done")

	return vrList, count, err
}

func GetVendorReportByMonth(month string) ([]common.VendorReport, bool, error) {

	var vrList []common.VendorReport
	var find bool

	zaps.Info(">>> get vendors report by month: ", month)

	rows, err := db.Query("SELECT begin_date, end_date, vendor_code, "+
		"vendor_name, "+
		"month_purchase, month_purchase_orig, month_receipt, month_rebate, month_pur_final, "+
		"month_purchase2, month_sales_qty, month_sales, month_sales_orig, "+
		"month_sales2, month_paid, month_paid_orig, last_month_unpaid, "+
		"month_unpaid, month_unpaid_orig, month_surplus, year_purchase, "+
		"year_sales_qty, year_sales, year_paid, "+
		"year_unpaid, year_unpaid_orig, year_surplus, year_open_stock, "+
		"month_open_stock, month_close_stock, "+
		"month_avg_stock, year_avg_stock, month_days_to, "+
		"year_days_to, comm_eff, gross_margin, "+
		"order_fill_rate, first_order_time, created_time "+
		"FROM vendors_report WHERE month = ?", month)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return vrList, find, err
	}

	defer rows.Close()

	for rows.Next() {
		var vr common.VendorReport
		err := rows.Scan(&vr.BeginDate, &vr.EndDate, &vr.VendorCode,
			&vr.VendorName,
			&vr.MonthPurchase, &vr.MonthPurchaseOrig, &vr.MonthReceipt, &vr.MonthRebate, &vr.MonthPurFinal,
			&vr.MonthPurchase2, &vr.MonthSalesQty, &vr.MonthSales, &vr.MonthSalesOrig,
			&vr.MonthSales2, &vr.MonthPaid, &vr.MonthPaidOrig, &vr.LastMonthUnpaid,
			&vr.MonthUnpaid, &vr.MonthUnpaidOrig, &vr.MonthSurplus, &vr.YearPurchase,
			&vr.YearSalesQty, &vr.YearSales, &vr.YearPaid,
			&vr.YearUnpaid, &vr.YearUnpaidOrig, &vr.YearSurplus, &vr.YearOpenInv,
			&vr.MonthOpenInv, &vr.MonthCloseInv,
			&vr.MonthAvgInv, &vr.YearAvgInv, &vr.MonthDaysTo,
			&vr.YearDaysTo, &vr.CommEff, &vr.GrossMargin,
			&vr.OrderFillRate, &vr.FirstOrder, &vr.CreatedTime)
		if err != nil {
			zaps.Error("query error: ", err)
			return vrList, find, err
		} else {
			zaps.Debug(">>> month: ", vr.Month)
			zaps.Debug(">>> start date: ", vr.BeginDate)
			zaps.Debug(">>> end date: ", vr.EndDate)
			zaps.Debug(">>> vendor code: ", vr.VendorCode)
			zaps.Debug(">>> vendor name: ", vr.VendorName)

			if vr.VendorCode != "" {
				vrList = append(vrList, vr)
				find = true
			}
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return vrList, find, err
	}

	zaps.Info("<<< get vendors report list done")

	return vrList, find, err
}

func GetVendorReportActiveByMonth(month string) ([]common.VendorReport,
	bool, error) {

	var vrList []common.VendorReport
	var find bool

	zaps.Info(">>> get vendors report by month: ", month)

	rows, err := db.Query("SELECT a.begin_date, a.end_date, "+
		"a.vendor_code, a.vendor_name, "+
		"a.month_purchase, a.month_purchase_orig, a.month_receipt, a.month_rebate, "+
		"a.month_pur_final, a.month_purchase2, "+
		"a.month_sales_qty, a.month_sales, a.month_sales_orig, "+
		"a.month_sales2, a.month_paid, a.month_paid_orig, "+
		"a.last_month_unpaid, a.month_unpaid, a.month_unpaid_orig, "+
		"a.month_surplus, a.year_purchase, "+
		"a.year_sales_qty, a.year_sales, "+
		"a.year_paid, a.year_unpaid, a.year_unpaid_orig, a.year_surplus, "+
		"a.year_open_stock, a.month_open_stock, "+
		"a.month_close_stock, a.month_avg_stock, "+
		"a.year_avg_stock, a.month_days_to, "+
		"a.year_days_to, a.comm_eff, a.gross_margin, "+
		"a.order_fill_rate, a.first_order_time, "+
		"a.created_time "+
		"FROM vendors_report a "+
		"LEFT JOIN vendor_cfg b "+
		"ON a.vendor_code = b.vendor_code "+
		"WHERE month = ? AND b.inactive = 0", month)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return vrList, find, err
	}

	defer rows.Close()

	for rows.Next() {
		var vr common.VendorReport
		err := rows.Scan(&vr.BeginDate, &vr.EndDate, &vr.VendorCode,
			&vr.VendorName,
			&vr.MonthPurchase, &vr.MonthPurchaseOrig, &vr.MonthReceipt, &vr.MonthRebate, &vr.MonthPurFinal,
			&vr.MonthPurchase2, &vr.MonthSalesQty, &vr.MonthSales, &vr.MonthSalesOrig,
			&vr.MonthSales2, &vr.MonthPaid, &vr.MonthPaidOrig, &vr.LastMonthUnpaid,
			&vr.MonthUnpaid, &vr.MonthUnpaidOrig, &vr.MonthSurplus, &vr.YearPurchase,
			&vr.YearSalesQty, &vr.YearSales, &vr.YearPaid,
			&vr.YearUnpaid, &vr.YearUnpaidOrig, &vr.YearSurplus, &vr.YearOpenInv,
			&vr.MonthOpenInv, &vr.MonthCloseInv, &vr.MonthAvgInv,
			&vr.YearAvgInv, &vr.MonthDaysTo, &vr.YearDaysTo,
			&vr.CommEff, &vr.GrossMargin, &vr.OrderFillRate,
			&vr.FirstOrder, &vr.CreatedTime)
		if err != nil {
			zaps.Error("query error: ", err)
			return vrList, find, err
		} else {
			zaps.Debug(">>> month: ", vr.Month)
			zaps.Debug(">>> start date: ", vr.BeginDate)
			zaps.Debug(">>> end date: ", vr.EndDate)
			zaps.Debug(">>> vendor code: ", vr.VendorCode)
			zaps.Debug(">>> vendor name: ", vr.VendorName)

			vrList = append(vrList, vr)
			find = true
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return vrList, find, err
	}

	zaps.Info("<<< get vendors report list done")

	return vrList, find, err
}

func GetVendorReportByMonthVendor(month string,
	vendor string) (common.VendorReport, bool, error) {

	var res common.VendorReport
	var find bool

	zaps.Info(">>> get vendors report by month & vendor: ", month, vendor)

	rows, err := db.Query("SELECT begin_date, end_date, vendor_code, "+
		"vendor_name, month_purchase, month_purchase_orig, month_receipt, month_rebate, "+
		"month_pur_final, month_purchase2, month_sales_qty, "+
		"month_sales, month_sales_orig, month_sales2, month_paid, month_paid_orig, "+
		"last_month_unpaid, month_unpaid, month_unpaid_orig, month_surplus, "+
		"year_purchase, year_sales_qty, year_sales, "+
		"year_paid, year_unpaid, year_unpaid_orig, year_surplus, "+
		"year_open_stock, month_open_stock, "+
		"month_close_stock, month_avg_stock, year_avg_stock, "+
		"month_days_to, year_days_to, comm_eff, gross_margin, "+
		"order_fill_rate, first_order_time, created_time "+
		"FROM vendors_report "+
		"WHERE month = ? AND vendor_code = ?",
		month, vendor)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return res, find, err
	}

	defer rows.Close()

	for rows.Next() {
		var vr common.VendorReport
		err := rows.Scan(&vr.BeginDate, &vr.EndDate, &vr.VendorCode,
			&vr.VendorName,
			&vr.MonthPurchase, &vr.MonthPurchaseOrig, &vr.MonthReceipt, &vr.MonthRebate, &vr.MonthPurFinal,
			&vr.MonthPurchase2, &vr.MonthSalesQty, &vr.MonthSales, &vr.MonthSalesOrig,
			&vr.MonthSales2, &vr.MonthPaid, &vr.MonthPaidOrig, &vr.LastMonthUnpaid,
			&vr.MonthUnpaid, &vr.MonthUnpaidOrig, &vr.MonthSurplus, &vr.YearPurchase,
			&vr.YearSalesQty, &vr.YearSales, &vr.YearPaid,
			&vr.YearUnpaid, &vr.YearUnpaidOrig, &vr.YearSurplus, &vr.YearOpenInv,
			&vr.MonthOpenInv, &vr.MonthCloseInv, &vr.MonthAvgInv,
			&vr.YearAvgInv, &vr.MonthDaysTo, &vr.YearDaysTo,
			&vr.CommEff, &vr.GrossMargin, &vr.OrderFillRate,
			&vr.FirstOrder, &vr.CreatedTime)
		if err != nil {
			zaps.Error("query error: ", err)
			return res, find, err
		} else {
			vr.Month = month

			zaps.Debug(">>> month: ", vr.Month)
			zaps.Debug(">>> start date: ", vr.BeginDate)
			zaps.Debug(">>> end date: ", vr.EndDate)
			zaps.Debug(">>> vendor code: ", vr.VendorCode)
			zaps.Debug(">>> vendor name: ", vr.VendorName)

			res = vr
			find = true
			break
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return res, find, err
	}

	zaps.Info("<<< get vendors report by month and vendor done")

	return res, find, err
}
