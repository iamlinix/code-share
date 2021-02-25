package orm

import (
	"database/sql"
	"time"

	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"
	_ "github.com/go-sql-driver/mysql"
)

func AddSalesMaterialCfg(cfg common.SalesMaterialCfg) error {

	zaps.Info(">>> add sales material config")

	stmt, err := db.Prepare("INSERT INTO sales_material_cfg(material, " +
		"material_txt, level, tag, created_time) " +
		"VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		zaps.Error("insert into sales material cfg failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cfg.Material, cfg.MaterialTxt, cfg.Level, cfg.Tag,
		time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		zaps.Error("db exec failed: ", err)
		return err
	}

	zaps.Info("<<< add sales material config done")

	return err
}

func UpdateSalesMaterialCfg(cfg common.SalesMaterialCfg) error {

	zaps.Info(">>> update sales material config")

	stmt, err := db.Prepare("UPDATE sales_material_cfg SET " +
		"tag=?, created_time=? " +
		"WHERE material = ?")
	if err != nil {
		zaps.Error("update sales material config failed: ", err)
		return err
	}
	defer stmt.Close()

	t := time.Now().Format("2006-01-02 15:04:05")
	_, err = stmt.Exec(cfg.Tag, t, cfg.Material)
	if err != nil {
		zaps.Error("db exec failed: ", err)
		return err
	}

	zaps.Info("<<< update sales material config done")

	return err
}

func DelSalesMaterialCfg(code string) error {

	zaps.Info(">>> del sales material config: ", code)

	_, err := db.Exec("DELETE FROM sales_material_cfg "+
		"WHERE material = ?", code)
	if err != nil {
		zaps.Error("delete exec failed: ", err)
		return err
	}

	zaps.Info("<<< del sales material config done")

	return err
}

func ClearSalesMaterialCfg() error {

	zaps.Info(">>> clear sales material config")

	_, err := db.Exec("DELETE FROM sales_material_cfg")
	if err != nil {
		zaps.Error("clear config exec failed: ", err)
		return err
	}

	zaps.Info("<<< clear sales material config done")

	return err
}

func GetSalesMaterialCfgList(all int, page int) ([]common.SalesMaterialCfg,
	int, error) {

	var smList []common.SalesMaterialCfg
	var rows *sql.Rows
	var err error
	var count int

	zaps.Info(">>> get sales material config list with page ", page)

	if all == 1 {
		rows, err = db.Query("SELECT material, material_txt, " +
			"level, tag, created_time " +
			"FROM sales_material_cfg " +
			"ORDER BY created_time DESC")
	} else {
		rows, err = db.Query("SELECT material, material_txt, "+
			"level, tag, created_time "+
			"FROM sales_material_cfg "+
			"ORDER BY created_time DESC LIMIT ?,10",
			page*10)
	}

	if err != nil {
		zaps.Error("db query failed: ", err)
		return smList, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var sm common.SalesMaterialCfg
		err := rows.Scan(&sm.Material, &sm.MaterialTxt, &sm.Level,
			&sm.Tag, &sm.CreatedTime)
		if err != nil {
			zaps.Error("query error: ", err)
			return smList, count, err
		}

		sm.Material = common.GetShortMatlCode(sm.Material)
		smList = append(smList, sm)
		count++
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return smList, 0, err
	}

	zaps.Info("<<< get sales material config list done")

	return smList, count, err
}
