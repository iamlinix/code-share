package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"errors"
)

func UpdateVendorRatingCfgTx(cfgs []common.VendorRatingCfg) error {

	var weight float64
	var err error

	zaps.Info(">>> update vendor rating config tx")

	tx, _ := db.Begin()
	defer tx.Rollback()

	for _, cfg := range cfgs {
		err = UpdateVendorRatingCfg(cfg)
		if err != nil {
			zaps.Errorf("update rating cfg failed: %v", err)
			return err
		}

		zaps.Infof("id: %d, weight: %f", cfg.ID, cfg.Weight)
		weight += cfg.Weight
	}

	if weight != 100.0 {
		zaps.Errorf("total weight invalid: %f", weight)
		return errors.New("total weight invalid")
	}

	/* Transaction END */
	err = tx.Commit()
	if err != nil {
		zaps.Errorf("transaction commit failed: %v", err)
		return err
	}

	zaps.Info("<<< update vendor rating config done")

	return err
}

func UpdateVendorRatingCfg(cfg common.VendorRatingCfg) error {

	zaps.Info(">>> update vendor rating config")

	stmt, err := db.Prepare("UPDATE vendor_rating_cfg SET " +
		"weight = ? " +
		"WHERE id = ?")
	if err != nil {
		zaps.Error("update vendor rating config failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cfg.Weight, cfg.ID)
	if err != nil {
		zaps.Error("db exec failed: err")
		return err
	}

	zaps.Info("<<< update vendor rating config done")

	return err
}

func GetVendorRatingCfgList() ([]common.VendorRatingCfg, int, error) {

	var cfgList []common.VendorRatingCfg
	var rows *sql.Rows
	var err error
	var count int

	zaps.Info(">>> get vendor rating config list")

	rows, err = db.Query("SELECT id, name, weight " +
		"FROM vendor_rating_cfg")
	if err != nil {
		zaps.Error("db query failed: ", err)
		return cfgList, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var cfg common.VendorRatingCfg
		err := rows.Scan(&cfg.ID, &cfg.Name, &cfg.Weight)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> id: ", cfg.ID)
			zaps.Debug(">>> name: ", cfg.Name)
			zaps.Debug(">>> weight: ", cfg.Weight)

			cfgList = append(cfgList, cfg)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return cfgList, 0, err
	}

	zaps.Info("<<< get vendor rating config list done")

	return cfgList, count, err
}
