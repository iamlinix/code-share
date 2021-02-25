package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"time"
)

func AddVendorCfg(vc common.VendorCfg) error {

	zaps.Info(">>> add vendor config")

	stmt, err := db.Prepare("INSERT INTO vendor_cfg(vendor_code, " +
		"vendor_name, year_purchase, year_sales, " +
		"year_paid, year_unpaid, year_surplus, " +
		"rebate, inactive, is_new, created_time) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		zaps.Error("insert into vendor config failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(vc.VendorCode, vc.VendorName, vc.YearPurchase,
		vc.YearSales, vc.YearPaid, vc.YearUnpaid,
		vc.YearSurplus, vc.Rebate, vc.Inactive,
		vc.IsNew, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		zaps.Error("db exec failed: err")
		return err
	}

	zaps.Info("<<< add vendor config done")

	return err
}

func UpdateVendorCfg(vc common.VendorCfg) error {

	zaps.Info(">>> update vendor config")

	stmt, err := db.Prepare("UPDATE vendor_cfg SET vendor_name = ?, " +
		"year_purchase = ?, year_sales = ?, " +
		"year_paid = ?, year_unpaid = ?, " +
		"year_surplus = ?, rebate = ?, " +
		"inactive = ?, is_new = ?, created_time = ? " +
		"WHERE vendor_code = ?")
	if err != nil {
		zaps.Error("update vendor config failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(vc.VendorName, vc.YearPurchase, vc.YearSales,
		vc.YearPaid, vc.YearUnpaid, vc.YearSurplus, vc.Rebate,
		vc.Inactive, vc.IsNew,
		time.Now().Format("2006-01-02 15:04:05"),
		vc.VendorCode)
	if err != nil {
		zaps.Error("db exec failed: err")
		return err
	}

	zaps.Info("<<< update vendor config done")

	return err
}

func DelVendorCfg(code string) {

	zaps.Info(">>> del vendor config: ", code)

	_, err := db.Exec("DELETE FROM vendor_cfg WHERE vendor_code = ?", code)
	if err != nil {
		zaps.Error("delete vendor config exec failed: ", err)
	}

	zaps.Info("<<< del vendor config done")
}

func ClearVendorCfg() {

	zaps.Info(">>> clear vendor config")

	_, err := db.Exec("DELETE FROM vendor_cfg")
	if err != nil {
		zaps.Error("clear vendor config exec failed: ", err)
	}

	zaps.Info("<<< clear vendor config done")
}

func GetVendorCfgList(all int, page int) ([]common.VendorCfg, int, error) {

	var vcList []common.VendorCfg
	var rows *sql.Rows
	var err error
	var count int

	zaps.Info(">>> get vendor config list with page ", page)

	if all == 1 {
		rows, err = db.Query("SELECT vendor_code, vendor_name, " +
			"year_purchase, year_sales, year_paid, " +
			"year_unpaid, year_surplus, rebate, " +
			"inactive, is_new, created_time " +
			"FROM vendor_cfg")
	} else {
		rows, err = db.Query("SELECT vendor_code, vendor_name, "+
			"year_purchase, year_sales, year_paid, "+
			"year_unpaid, year_surplus, rebate, "+
			"inactive, is_new, created_time "+
			"FROM vendor_cfg LIMIT ?,10", page*10)
	}
	if err != nil {
		zaps.Error("db query failed: ", err)
		return vcList, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var vc common.VendorCfg
		err := rows.Scan(&vc.VendorCode, &vc.VendorName,
			&vc.YearPurchase, &vc.YearSales,
			&vc.YearPaid, &vc.YearUnpaid,
			&vc.YearSurplus, &vc.Rebate,
			&vc.Inactive, &vc.IsNew, &vc.CreatedTime)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> vendor code: ", vc.VendorCode)
			zaps.Debug(">>> vendor name: ", vc.VendorName)
			zaps.Debug(">>> rebate: ", vc.Rebate)
			zaps.Debug(">>> inactive: ", vc.Inactive)
			zaps.Debug(">>> is new: ", vc.IsNew)
			zaps.Debug(">>> created time: ", vc.CreatedTime)

			vcList = append(vcList, vc)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return vcList, 0, err
	}

	zaps.Info("<<< get vendor config list done")

	return vcList, count, err
}

func GetInactiveVendorCfgList() ([]common.VendorCfg, int, error) {

	var vcList []common.VendorCfg
	var rows *sql.Rows
	var err error
	var count int

	zaps.Info(">>> get inactive vendor config list")

	rows, err = db.Query("SELECT vendor_code, vendor_name, " +
		"year_purchase, year_sales, year_paid, " +
		"year_unpaid, year_surplus, rebate, " +
		"inactive, is_new, created_time " +
		"FROM vendor_cfg WHERE inactive = 1")
	if err != nil {
		zaps.Error("db query failed: ", err)
		return vcList, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var vc common.VendorCfg
		err := rows.Scan(&vc.VendorCode, &vc.VendorName,
			&vc.YearPurchase, &vc.YearSales,
			&vc.YearPaid, &vc.YearUnpaid,
			&vc.YearSurplus, &vc.Rebate,
			&vc.Inactive, &vc.IsNew, &vc.CreatedTime)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> vendor code: ", vc.VendorCode)
			zaps.Debug(">>> vendor name: ", vc.VendorName)
			zaps.Debug(">>> rebate: ", vc.Rebate)
			zaps.Debug(">>> inactive: ", vc.Inactive)
			zaps.Debug(">>> is new: ", vc.IsNew)
			zaps.Debug(">>> created time: ", vc.CreatedTime)

			vcList = append(vcList, vc)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return vcList, 0, err
	}

	zaps.Info("<<< get inactive vendor config list done: ", count)

	return vcList, count, err
}

func GetVendorCfgByCode(code string) (common.VendorCfg, bool, error) {

	var cfg common.VendorCfg
	find := false

	rows, err := db.Query("SELECT vendor_code, vendor_name, "+
		"year_purchase, year_sales, year_paid, "+
		"year_unpaid, year_surplus, rebate, "+
		"inactive, is_new, created_time "+
		"FROM vendor_cfg WHERE vendor_code = ?", code)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return cfg, find, err
	}

	defer rows.Close()

	for rows.Next() {
		var vc common.VendorCfg
		err := rows.Scan(&vc.VendorCode, &vc.VendorName,
			&vc.YearPurchase, &vc.YearSales,
			&vc.YearPaid, &vc.YearUnpaid,
			&vc.YearSurplus, &vc.Rebate,
			&vc.Inactive, &vc.IsNew, &vc.CreatedTime)
		if err != nil {
			zaps.Error("query error: ", err)
			return cfg, find, err
		} else {
			zaps.Debug(">>> vendor code: ", vc.VendorCode)
			zaps.Debug(">>> vendor name: ", vc.VendorName)
			zaps.Debug(">>> rebate: ", vc.Rebate)
			zaps.Debug(">>> inactive: ", vc.Inactive)
			zaps.Debug(">>> is new: ", vc.IsNew)
			zaps.Debug(">>> created time: ", vc.CreatedTime)
			cfg = vc
			find = true
			break
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return cfg, find, err
	}

	return cfg, find, err
}

func IsInactiveVendor(vcList []common.VendorCfg, vcode string) bool {

	for _, vc := range vcList {

		if vcode == vc.VendorCode {
			return true
		}
	}

	return false
}
