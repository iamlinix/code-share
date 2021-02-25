package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"time"
)

func ImportTestMarketMatlCfg(cfgList []common.EffTestMktMatlCfg) error {

	tx, _ := db.Begin()
	defer tx.Rollback()
	var err error

	for _, cfg := range cfgList {

		/* get material name first */
		name, _, err := GetMatlNameByCode(cfg.Material)
		if err != nil {
			zaps.Error("get material name by code failed: ", err)
			return err
		}

		cfg.MaterialTxt = name
		err = AddTestMarketMatlCfgTx(tx, cfg)
		if err != nil {
			zaps.Error("add test market matl cfg failed")
			return err
		}
	}

	/* Transaction END */
	err = tx.Commit()
	if err != nil {
		zaps.Errorf("transaction commit failed: %v", err)
		return err
	}

	return nil
}

func AddTestMarketMatlCfg(cfg common.EffTestMktMatlCfg) error {

	zaps.Info(">>> add test market matl cfg")

	name, _, err := GetMatlNameByCode(cfg.Material)
	if err != nil {
		zaps.Error("get material name by code failed: ", err)
		return err
	}

	cfg.MaterialTxt = name
	stmt, err := db.Prepare("INSERT INTO eff_test_market_cfg(material, " +
		"material_txt, status, created_time) " +
		"VALUES(?, ?, ?, ?)")
	if err != nil {
		zaps.Error("insert into test market matl cfg failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cfg.Material, cfg.MaterialTxt, cfg.Status,
		time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		zaps.Error("db exec failed: err")
		return err
	}

	zaps.Info("<<< add test market matl cfg done")

	return err
}

func AddTestMarketMatlCfgTx(tx *sql.Tx, cfg common.EffTestMktMatlCfg) error {

	zaps.Info(">>> add test market matl cfg tx")

	stmt, err := tx.Prepare("INSERT INTO eff_test_market_cfg(material, " +
		"material_txt, status, created_time) " +
		"VALUES(?, ?, ?, ?)")
	if err != nil {
		zaps.Error("insert into test market matl cfg failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cfg.Material, cfg.MaterialTxt, cfg.Status,
		time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		zaps.Error("db exec failed: err")
		return err
	}

	zaps.Info("<<< add test market matl cfg tx done")

	return err
}

func UpdateTestMarketMatlCfg(cfg common.EffTestMktMatlCfg) error {

	zaps.Info(">>> update test market matl cfg")

	stmt, err := db.Prepare("UPDATE eff_test_market_cfg SET " +
		"status = ?, created_time = ? " +
		"WHERE material = ?")
	if err != nil {
		zaps.Error("update test market matl cfg failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cfg.Status, time.Now().Format("2006-01-02 15:04:05"),
		cfg.Material)
	if err != nil {
		zaps.Error("db exec failed: err")
		return err
	}

	zaps.Info("<<< update test market matl cfg done")

	return err
}

func DelTestMarketMatlCfg(code string) error {

	zaps.Info(">>> del test market matl cfg: ", code)

	_, err := db.Exec("DELETE FROM eff_test_market_cfg "+
		"WHERE material = ?", code)
	if err != nil {
		zaps.Error("delete test market matl cfg exec failed: ", err)
		return err
	}

	zaps.Info("<<< del test market matl cfg done")
	return nil
}

func ClearTestMarketMatlCfg() {

	zaps.Info(">>> clear test market matl cfg")

	_, err := db.Exec("DELETE FROM eff_test_market_cfg")
	if err != nil {
		zaps.Error("clear cfg exec failed: ", err)
	}

	zaps.Info("<<< clear test market matl cfg done")
}

func GetTestMarketMatlCfgList(all int, page int) ([]common.EffTestMktMatlCfg,
	int, error) {

	var cfgList []common.EffTestMktMatlCfg
	var rows *sql.Rows
	var err error
	var count int

	zaps.Info(">>> get test market matl cfg list with page ", page)

	if all == 1 {
		rows, err = db.Query("SELECT material, material_txt, status," +
			"created_time FROM eff_test_market_cfg")
	} else {
		rows, err = db.Query("SELECT material, material_txt, status,"+
			"created_time FROM eff_test_market_cfg "+
			"LIMIT ?,10", page*20)
	}
	if err != nil {
		zaps.Error("db query failed: ", err)
		return cfgList, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var cfg common.EffTestMktMatlCfg
		err := rows.Scan(&cfg.Material, &cfg.MaterialTxt,
			&cfg.Status, &cfg.CreatedTime)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> material: ", cfg.Material)

			cfgList = append(cfgList, cfg)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return cfgList, 0, err
	}

	zaps.Info("<<< get test market matl cfg list done")

	return cfgList, count, err
}
