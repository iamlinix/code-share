package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
)

func AddExMatlCfg(emc common.ExMatlCfg) error {

	zaps.Info(">>> add ex-matl config")

	stmt, err := db.Prepare("INSERT INTO ex_matl_cfg(material, " +
		"material_txt, vendor, vendor_name, " +
		"rebate) " +
		"VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		zaps.Error("insert into ex-matl config failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(emc.Material, emc.MaterialTxt, emc.Vendor,
		emc.VendorName, emc.Rebate)
	if err != nil {
		zaps.Error("db exec failed: err")
		return err
	}

	zaps.Info("<<< add ex-matl config done")

	return err
}

func UpdateExMatlCfg(emc common.ExMatlCfg) error {

	zaps.Info(">>> update ex-matl config")

	stmt, err := db.Prepare("UPDATE ex_matl_cfg SET " +
		"material_txt = ?, vendor = ?, " +
		"vendor_name = ?, rebate = ? " +
		"WHERE material = ?")
	if err != nil {
		zaps.Error("update ex-matl config failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(emc.MaterialTxt, emc.Vendor, emc.VendorName,
		emc.Rebate, emc.Material)
	if err != nil {
		zaps.Error("db exec failed: ", err)
		return err
	}

	zaps.Info("<<< update ex-matl config done")

	return err
}

func DelExMatlCfg(code string) {

	zaps.Info(">>> del ex-matl config: ", code)

	_, err := db.Exec("DELETE FROM ex_matl_cfg WHERE material = ?", code)
	if err != nil {
		zaps.Error("delete ex-matl config exec failed: ", err)
	}

	zaps.Info("<<< del ex-matl config done")
}

func ClearExMatlCfg() {

	zaps.Info(">>> clear ex-matl config")

	_, err := db.Exec("DELETE FROM ex_matl_cfg")
	if err != nil {
		zaps.Error("clear ex-matl config exec failed: ", err)
	}

	zaps.Info("<<< clear ex-matl config done")
}

func GetExMatlCfgByCode(code string) (common.ExMatlCfg, bool, error) {

	var emc common.ExMatlCfg
	find := false

	rows, err := db.Query("SELECT material, material_txt, vendor, "+
		"vendor_name, rebate FROM ex_matl_cfg "+
		"WHERE material = ?", code)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return emc, find, err
	}

	defer rows.Close()

	for rows.Next() {
		var cfg common.ExMatlCfg
		var mtxt sql.NullString
		var vname sql.NullString

		err := rows.Scan(&cfg.Material, &mtxt, &cfg.Vendor,
			&vname, &cfg.Rebate)
		if err != nil {
			zaps.Error("query error: ", err)
			return emc, find, err
		} else {
			if mtxt.Valid {
				cfg.MaterialTxt = mtxt.String
			} else {
				cfg.MaterialTxt = "N/A"
			}

			if vname.Valid {
				cfg.VendorName = vname.String
			} else {
				cfg.VendorName = "N/A"
			}

			emc = cfg
			find = true
			break
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return emc, find, err
	}

	return emc, find, err
}

func GetExMatlCfgList() ([]common.ExMatlCfg, int, error) {

	var emList []common.ExMatlCfg
	var rows *sql.Rows
	var err error
	var count int

	rows, err = db.Query("SELECT material, material_txt, " +
		"vendor, vendor_name, rebate " +
		"FROM ex_matl_cfg")
	if err != nil {
		zaps.Error("db query failed: ", err)
		return emList, count, err
	}

	defer rows.Close()

	for rows.Next() {
		var em common.ExMatlCfg
		var mtxt sql.NullString
		var vname sql.NullString

		err := rows.Scan(&em.Material, &mtxt, &em.Vendor,
			&vname, &em.Rebate)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> material: ", em.Material)
			zaps.Debug(">>> vendor: ", em.Vendor)
			zaps.Debug(">>> rebate: ", em.Rebate)

			if mtxt.Valid {
				em.MaterialTxt = mtxt.String
			} else {
				em.MaterialTxt = "N/A"
			}

			if vname.Valid {
				em.VendorName = vname.String
			} else {
				em.VendorName = "N/A"
			}

			emList = append(emList, em)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return emList, count, err
	}

	zaps.Info("<<< get exclude material list done")

	return emList, count, err
}

func GetExMatlVendorList() ([]string, error) {

	var evList []string
	var rows *sql.Rows
	var err error

	rows, err = db.Query("SELECT DISTINCT vendor FROM ex_matl_cfg")
	if err != nil {
		zaps.Error("db query failed: ", err)
		return evList, err
	}

	defer rows.Close()

	for rows.Next() {
		var ev string
		err := rows.Scan(&ev)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> ex-vendor: ", ev)

			evList = append(evList, ev)
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return evList, err
	}

	zaps.Info("<<< get exclude vendor list done")

	return evList, err
}
