package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"fmt"
)

func GetVendorNameByCode(code string) (string, bool, error) {

	var name string
	find := false

	zaps.Info(">>> get vendor name by code: ", code)

	rows, err := db.Query("SELECT name FROM vendor "+
		"WHERE vendor = ?", code)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return name, find, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			zaps.Error("query error: ", err)
			return name, find, err
		}

		find = true
		break
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return name, find, err
	}

	zaps.Info("<<< get vendor name by code done")

	return name, find, err
}

func GetVendorByCode(code string) (common.Vendor, bool, error) {

	var res common.Vendor
	var rows *sql.Rows
	find := false

	rows, err := db.Query("SELECT vendor, name "+
		"FROM vendor WHERE vendor = ?", code)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return res, find, err
	}

	defer rows.Close()

	for rows.Next() {
		var v common.Vendor
		err := rows.Scan(&v.LongCode, &v.Name)
		if err != nil {
			zaps.Error("query error: ", err)
			return v, find, err
		}

		v.ShortCode = common.GetShortVendorCode(v.LongCode)
		res = v
		find = true
		break
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return res, find, err
	}

	zaps.Info("<<< get vendor by code done")

	return res, find, err
}

func GetVendorByMatlCode(material string, day string) (string, string, bool, error) {

	var code sql.NullString
	var name sql.NullString
	var vcode, vname string
	find := false

	cmd := fmt.Sprintf("SELECT a.vendor, b.name FROM zifpurd a "+
		"LEFT JOIN vendor b ON a.vendor = b.vendor "+
		"WHERE material = '%s' AND zgr_date < '%s' "+
		"ORDER BY zgr_date DESC limit 1 ", material, day)

	err := db.QueryRow(cmd).Scan(&code, &name)
	if err != nil {
		if err == sql.ErrNoRows {
			zaps.Error("query zifpurd vendor no rows")
			return "N/A", "N/A", find, nil
		} else {
			zaps.Error("query zifpurd vendor failed: ", err)
			return "N/A", "N/A", find, err
		}
	}

	find = true
	if code.Valid == true {
		vcode = code.String
	} else {
		vcode = "N/A"
	}

	if name.Valid == true {
		vname = name.String
	} else {
		vname = "N/A"
	}

	return vcode, vname, find, nil
}
