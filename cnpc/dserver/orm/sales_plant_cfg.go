package orm

import (
	"database/sql"
	"time"

	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"
	_ "github.com/go-sql-driver/mysql"
)

func AddSalesPlantCfg(cfg common.SalesPlantCfg) error {

	zaps.Info(">>> add sales plant config")

	stmt, err := db.Prepare("INSERT INTO sales_plant_cfg(plant, " +
		"plant_name, tag, created_time) " +
		"VALUES(?, ?, ?, ?)")
	if err != nil {
		zaps.Error("insert into sales plant cfg failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cfg.Plant, cfg.PlantName, cfg.Tag,
		time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		zaps.Error("db exec failed: err")
		return err
	}

	zaps.Info("<<< add sales plant config done")

	return err
}

func UpdateSalesPlantCfg(cfg common.SalesPlantCfg) error {

	zaps.Info(">>> update sales plant config")

	stmt, err := db.Prepare("UPDATE sales_plant_cfg SET " +
		"tag=?, created_time=? " +
		"WHERE plant = ?")
	if err != nil {
		zaps.Error("update sales plant config failed: ", err)
		return err
	}
	defer stmt.Close()

	t := time.Now().Format("2006-01-02 15:04:05")
	_, err = stmt.Exec(cfg.Tag, t, cfg.Plant)
	if err != nil {
		zaps.Error("db exec failed: ", err)
		return err
	}

	zaps.Info("<<< update sales plant config done")

	return err
}

func DelSalesPlantCfg(code string) error {

	zaps.Info(">>> del sales plant config: ", code)

	_, err := db.Exec("DELETE FROM sales_plant_cfg "+
		"WHERE plant = ?", code)
	if err != nil {
		zaps.Error("delete exec failed: ", err)
		return err
	}

	zaps.Info("<<< del sales plant config done")

	return err
}

func ClearSalesPlantCfg() error {

	zaps.Info(">>> clear sales plant config")

	_, err := db.Exec("DELETE FROM sales_plant_cfg")
	if err != nil {
		zaps.Error("clear config exec failed: ", err)
		return err
	}

	zaps.Info("<<< clear sales plant config done")

	return err
}

func GetSalesPlantCfgList(all int, page int) ([]common.SalesPlantCfg,
	int, error) {

	var cfgList []common.SalesPlantCfg
	var rows *sql.Rows
	var err error
	var count int

	zaps.Info(">>> get sales material config list with page ", page)

	if all == 1 {
		rows, err = db.Query("SELECT plant, plant_name, " +
			"tag, created_time " +
			"FROM sales_plant_cfg " +
			"ORDER BY created_time DESC")
	} else {
		rows, err = db.Query("SELECT plant, plant_name, "+
			"tag, created_time "+
			"FROM sales_plant_cfg "+
			"ORDER BY created_time DESC LIMIT ?,10",
			page*10)
	}

	if err != nil {
		zaps.Error("db query failed: ", err)
		return cfgList, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var cfg common.SalesPlantCfg
		err := rows.Scan(&cfg.Plant, &cfg.PlantName, &cfg.Tag,
			&cfg.CreatedTime)
		if err != nil {
			zaps.Error("query error: ", err)
			return cfgList, count, err
		}

		cfgList = append(cfgList, cfg)
		count++
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return cfgList, 0, err
	}

	zaps.Info("<<< get sales plant config list done")

	return cfgList, count, err
}
