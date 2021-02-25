package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func UpdateEffMetricCfg(cfg common.EffMetricCfg) error {

	zaps.Info(">>> update eff metric config")

	stmt, err := db.Prepare("UPDATE eff_metric_cfg SET " +
		"value=? " +
		"WHERE name = ?")
	if err != nil {
		zaps.Error("update eff metric config failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cfg.Value, cfg.Name)
	if err != nil {
		zaps.Error("db exec failed: err")
		return err
	}

	zaps.Info("<<< update eff metric config done")

	return err
}

func GetEffMetricCfgListByModule(mid string) ([]common.EffMetricCfg,
	int, error) {

	var cfgList []common.EffMetricCfg
	var rows *sql.Rows
	var err error
	var count int

	zaps.Info(">>> get eff metric config list by module ", mid)

	rows, err = db.Query("SELECT name, value, module "+
		"FROM eff_metric_cfg "+
		"WHERE module = ?", mid)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return cfgList, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var cfg common.EffMetricCfg
		err := rows.Scan(&cfg.Name, &cfg.Value, &cfg.Module)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> cfg name: ", cfg.Name)
			zaps.Debug(">>> cfg value: ", cfg.Value)

			cfgList = append(cfgList, cfg)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return cfgList, 0, err
	}

	zaps.Info("<<< get eff metric config list done")

	return cfgList, count, err
}
