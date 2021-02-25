package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
)

func AddSalesPlantLevelCfg(cfg common.SalesPlantLevelCfg) (int64, error) {

	zaps.Info(">>> add sales plant level cfg")

	stmt, err := db.Prepare("INSERT INTO sales_plant_level_cfg(name, " +
		"begin, end) VALUES(?, ?, ?)")
	if err != nil {
		zaps.Error("insert into sales plant level cfg failed: ", err)
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(cfg.Name, cfg.Begin, cfg.End)
	if err != nil {
		zaps.Error("db exec failed: ", err)
		return 0, err
	}

	id, err := res.LastInsertId()

	zaps.Info("<<< add sales plant level cfg done: ", id)

	return id, err
}

func UpdateSalesPlantLevelCfg(cfg common.SalesPlantLevelCfg) error {

	zaps.Info(">>> update sales plant level cfg")

	stmt, err := db.Prepare("UPDATE sales_plant_level_cfg SET " +
		"name = ?, begin = ?, end = ? " +
		"WHERE id = ?")
	if err != nil {
		zaps.Error("update level cfg failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cfg.Name, cfg.Begin, cfg.End, cfg.ID)
	if err != nil {
		zaps.Error("db exec failed: ", err)
		return err
	}

	zaps.Info("<<< update sales plant level cfg done")

	return err
}

func DelSalesPlantLevelCfg(id int64) {

	zaps.Info(">>> del sales plant level cfg")

	_, err := db.Exec("DELETE FROM sales_plant_level_cfg WHERE id = ?", id)
	if err != nil {
		zaps.Error("delete sales plant level cfg exec failed: ", err)
	}

	zaps.Info("<<< del sales plant level cfg done")
}

func GetSalesPlantLevelCfgList(all int, page int) ([]common.SalesPlantLevelCfg, int, error) {

	var cfgList []common.SalesPlantLevelCfg
	var rows *sql.Rows
	var err error
	var count int

	zaps.Info(">>> get sales plant level cfg list with page ", page)

	if all == 1 {
		rows, err = db.Query("SELECT id, name, begin, end " +
			"FROM sales_plant_level_cfg")
	} else {
		rows, err = db.Query("SELECT id, name, begin, end "+
			"FROM sales_plant_level_cfg LIMIT ?,10", page*10)
	}

	if err != nil {
		zaps.Error("db query failed: ", err)
		return cfgList, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var cfg common.SalesPlantLevelCfg
		err := rows.Scan(&cfg.ID, &cfg.Name, &cfg.Begin, &cfg.End)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> id: ", cfg.ID)
			zaps.Debug(">>> name: ", cfg.Name)
			zaps.Debug(">>> begin: ", cfg.Begin)
			zaps.Debug(">>> end: ", cfg.End)

			cfgList = append(cfgList, cfg)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return cfgList, count, err
	}

	zaps.Info("<<< get sales plant level cfg list done")

	return cfgList, count, err
}

func GetSalesPlantLevelCfgByID(id int64) (common.SalesPlantLevelCfg, bool, error) {

	var cfg common.SalesPlantLevelCfg
	find := false

	zaps.Info(">>> get one sales plant level cfg info: ", id)

	rows, err := db.Query("SELECT id, name, begin, end "+
		"FROM sales_plant_level_cfg "+
		"WHERE id = ?", id)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return cfg, find, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&cfg.ID, &cfg.Name, &cfg.Begin, &cfg.End)
		if err != nil {
			zaps.Error("query error: ", err)
			return cfg, find, err
		} else {
			zaps.Debug(">>> id: ", cfg.ID)
			zaps.Debug(">>> name: ", cfg.Name)
			zaps.Debug(">>> begin: ", cfg.Begin)
			zaps.Debug(">>> end: ", cfg.End)

			find = true
			break
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return cfg, find, err
	}

	zaps.Info("<<< get one sales plant level cfg info done")

	return cfg, find, err
}
