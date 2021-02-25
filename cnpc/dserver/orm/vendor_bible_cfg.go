package orm

import (
	"database/sql"
	"time"

	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"
	_ "github.com/go-sql-driver/mysql"
)

func ImportBibleCfgTx(cfgs []common.BibleCfg) error {

	tx, _ := db.Begin()
	defer tx.Rollback()
	var err error

	for _, cfg := range cfgs {

		/* check bible item exist first */
		_, find, err := GetBibleCfgByMatl(cfg.Material)
		if err != nil {
			zaps.Error("get bible cfg by matl failed: ", err)
			return err

		} else if find == true {
			err = UpdateBibleCfgTx(tx, cfg)
			if err != nil {
				zaps.Error("update bible cfg failed: ", err)
				return err
			}
		} else {
			err = AddBibleCfgTx(tx, cfg)
			if err != nil {
				zaps.Error("add bible cfg tx failed: ", err)
				return err
			}
		}
	}

	/* Transaction END */
	err = tx.Commit()
	if err != nil {
		zaps.Error("transaction commit failed: ", err)
		return err
	}

	return nil
}

func AddBibleCfg(cfg common.BibleCfg) error {

	zaps.Info(">>> add bible cfg")

	stmt, err := db.Prepare("INSERT INTO vendor_bible_cfg(material, " +
		"material_txt, subclass, subclass_txt, " +
		"mainclass, mainclass_txt, valid_sdate, " +
		"groes, base_uom, purchase_price, " +
		"vendor_code, vendor_name, sale_price, " +
		"status, stock, created_time) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, " +
		"?, ?, ?, ?)")
	if err != nil {
		zaps.Error("insert into bile cfg failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cfg.Material, cfg.MaterialTxt, cfg.SubClass,
		cfg.SubClassTxt, cfg.MainClass, cfg.MainClassTxt,
		cfg.ValidSDate, cfg.Groes, cfg.BaseUom,
		cfg.PurchasePrice, cfg.VendorCode, cfg.VendorName,
		cfg.SalePrice, cfg.Status, cfg.Inventory,
		time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		zaps.Error("db exec failed: ", err)
		return err
	}

	zaps.Info("<<< add bible cfg done")

	return nil
}

func AddBibleCfgTx(tx *sql.Tx, cfg common.BibleCfg) error {

	zaps.Info(">>> add bible cfg tx")

	stmt, err := tx.Prepare("INSERT INTO vendor_bible_cfg(material, " +
		"material_txt, subclass, subclass_txt, " +
		"mainclass, mainclass_txt, valid_sdate, " +
		"groes, base_uom, purchase_price, " +
		"vendor_code, vendor_name, sale_price, " +
		"status, stock, created_time) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, " +
		"?, ?, ?, ?)")
	if err != nil {
		zaps.Error("insert into bile cfg failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cfg.Material, cfg.MaterialTxt, cfg.SubClass,
		cfg.SubClassTxt, cfg.MainClass, cfg.MainClassTxt,
		cfg.ValidSDate, cfg.Groes, cfg.BaseUom,
		cfg.PurchasePrice, cfg.VendorCode, cfg.VendorName,
		cfg.SalePrice, cfg.Status, cfg.Inventory,
		time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		zaps.Error("db exec failed: ", err)
		return err
	}

	zaps.Info("<<< add bible cfg done")

	return nil
}

func UpdateBibleCfg(cfg common.BibleCfg) error {

	zaps.Info(">>> update bible cfg")

	stmt, err := db.Prepare("UPDATE vendor_bible_cfg SET " +
		"material_txt = ?, " +
		"subclass = ?, subclass_txt = ?, " +
		"mainclass = ?, mainclass_txt = ?, " +
		"valid_sdate = ?, groes = ?, base_uom = ?, " +
		"purchase_price = ?, vendor_code = ?, " +
		"vendor_name = ?, sale_price = ?, " +
		"status = ?, stock = ?, created_time = ? " +
		"WHERE material = ?")
	if err != nil {
		zaps.Error("update bible cfg failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cfg.MaterialTxt, cfg.SubClass, cfg.SubClassTxt,
		cfg.MainClass, cfg.MainClassTxt, cfg.ValidSDate,
		cfg.Groes, cfg.BaseUom, cfg.PurchasePrice,
		cfg.VendorCode, cfg.VendorName, cfg.SalePrice,
		cfg.Status, cfg.Inventory,
		time.Now().Format("2006-01-02 15:04:05"),
		cfg.Material)
	if err != nil {
		zaps.Error("db exec failed: ", err)
		return err
	}

	zaps.Info("<<< update bible cfg done")

	return nil
}

func UpdateBibleCfgTx(tx *sql.Tx, cfg common.BibleCfg) error {

	zaps.Info(">>> update bible cfg tx")

	stmt, err := tx.Prepare("UPDATE vendor_bible_cfg SET " +
		"material_txt = ?, " +
		"subclass = ?, subclass_txt = ?, " +
		"mainclass = ?, mainclass_txt = ?, " +
		"valid_sdate = ?, groes = ?, base_uom = ?, " +
		"purchase_price = ?, vendor_code = ?, " +
		"vendor_name = ?, sale_price = ?, " +
		"status = ?, stock = ?, created_time = ? " +
		"WHERE material = ?")
	if err != nil {
		zaps.Error("update bible cfg failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cfg.MaterialTxt, cfg.SubClass, cfg.SubClassTxt,
		cfg.MainClass, cfg.MainClassTxt, cfg.ValidSDate,
		cfg.Groes, cfg.BaseUom, cfg.PurchasePrice,
		cfg.VendorCode, cfg.VendorName, cfg.SalePrice,
		cfg.Status, cfg.Inventory,
		time.Now().Format("2006-01-02 15:04:05"),
		cfg.Material)
	if err != nil {
		zaps.Error("db exec failed: ", err)
		return err
	}

	zaps.Info("<<< update bible cfg done")

	return nil
}

func DelBibleCfg(code string) error {

	zaps.Info(">>> del bible cfg: ", code)

	_, err := db.Exec("DELETE FROM vendor_bible_cfg "+
		"WHERE material = ?", code)
	if err != nil {
		zaps.Error("delete bible cfg failed: ", err)
		return err
	}

	zaps.Info("<<< del bible cfg done")

	return nil
}

func ClearBibleCfg() error {

	zaps.Info(">>> clear bible cfg")

	_, err := db.Exec("DELETE FROM vendor_bible_cfg")
	if err != nil {
		zaps.Error("clear bible cfg failed: ", err)
		return err
	}

	zaps.Info("<<< clear bible cfg done")

	return nil
}

func GetBibleCfgList(all int, page int) ([]common.BibleCfg, int, error) {

	var cfgs []common.BibleCfg
	var rows *sql.Rows
	var count int
	var err error

	zaps.Infof(">>> get bible cfg list with page (%d), all (%d)",
		page, all)

	if all == 1 {
		rows, err = db.Query("SELECT material, material_txt, " +
			"subclass, subclass_txt, mainclass, " +
			"mainclass_txt, valid_sdate, groes, " +
			"base_uom, purchase_price, vendor_code, " +
			"vendor_name, sale_price, status, stock, latest_date, " +
			"created_time FROM vendor_bible_cfg")
	} else {
		rows, err = db.Query("SELECT material, material_txt, "+
			"subclass, subclass_txt, mainclass, "+
			"mainclass_txt, valid_sdate, groes, "+
			"base_uom, purchase_price, vendor_code, "+
			"vendor_name, sale_price, status, stock, latest_date, "+
			"created_time FROM vendor_bible_cfg "+
			"LIMIT ?,10", page*20)
	}
	if err != nil {
		zaps.Error("db query failed: ", err)
		return cfgs, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var cfg common.BibleCfg
		err := rows.Scan(&cfg.Material, &cfg.MaterialTxt, &cfg.SubClass,
			&cfg.SubClassTxt, &cfg.MainClass, &cfg.MainClassTxt,
			&cfg.ValidSDate, &cfg.Groes, &cfg.BaseUom,
			&cfg.PurchasePrice, &cfg.VendorCode, &cfg.VendorName,
			&cfg.SalePrice, &cfg.Status, &cfg.Inventory, &cfg.LatestDate,
			&cfg.CreatedTime)
		if err != nil {
			zaps.Error("query error: ", err)
			return cfgs, count, err
		}

		cfg.Material = common.GetShortMatlCode(cfg.Material)
		cfg.VendorCode = common.GetShortVendorCode(cfg.VendorCode)

		cfgs = append(cfgs, cfg)
		count++
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return cfgs, 0, err
	}

	zaps.Info("<<< get bible cfg list done")

	return cfgs, count, err
}

func GetBibleCfgByMatl(code string) (common.BibleCfg, bool, error) {

	var cfg common.BibleCfg
	var find bool

	rows, err := db.Query("SELECT material, material_txt, "+
		"subclass, subclass_txt, mainclass, "+
		"mainclass_txt, valid_sdate, groes, "+
		"base_uom, purchase_price, vendor_code, "+
		"vendor_name, sale_price, status, stock, "+
		"created_time FROM vendor_bible_cfg "+
		"WHERE material = ?", code)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return cfg, find, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&cfg.Material, &cfg.MaterialTxt,
			&cfg.SubClass, &cfg.SubClassTxt, &cfg.MainClass,
			&cfg.MainClassTxt, &cfg.ValidSDate, &cfg.Groes,
			&cfg.BaseUom, &cfg.PurchasePrice, &cfg.VendorCode,
			&cfg.VendorName, &cfg.SalePrice, &cfg.Status,
			&cfg.Inventory, &cfg.CreatedTime)
		if err != nil {
			zaps.Error("query error: ", err)
			return cfg, find, err
		}

		cfg.Material = common.GetShortMatlCode(cfg.Material)
		cfg.VendorCode = common.GetShortMatlCode(cfg.VendorCode)

		find = true
		break
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return cfg, find, err
	}

	return cfg, find, err
}

func GetVendorMatlCountList() ([]common.VendorMatlCount, int, error) {

	var vmcList []common.VendorMatlCount
	var rows *sql.Rows
	var err error
	var count int

	rows, err = db.Query("SELECT vendor_code, count(0) AS CNT " +
		"FROM vendor_bible_cfg GROUP BY vendor_code")
	if err != nil {
		zaps.Error("db query failed: ", err)
		return vmcList, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var vmc common.VendorMatlCount
		err := rows.Scan(&vmc.Vendor, &vmc.MatlCount)
		if err != nil {
			zaps.Error("query error: ", err)
			return vmcList, count, err
		}

		vmcList = append(vmcList, vmc)
		count++
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return vmcList, 0, err
	}

	zaps.Info("<<< get vendor matl count list done")

	return vmcList, count, err
}
