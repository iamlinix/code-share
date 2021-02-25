package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func AddSalesMatlNoSalesCfgTx(cfg common.SalesMatlNoSalesCfg) error {

	tx, _ := db.Begin()
	defer tx.Rollback()

	cfg.Material = common.GetLongMatlCode(cfg.Material)

	for _, c := range cfg.Plants {
		err := AddSalesMatlNoSalesCfg(tx, cfg, c)
		if err != nil {
			zaps.Errorf("add nosales cfg failed: %v", err)
			return err
		}
	}

	/* Transaction END */
	err := tx.Commit()
	if err != nil {
		zaps.Errorf("transaction commit failed: %v", err)
		return err
	}

	zaps.Info("<<< add sales matl no-sales cfg done")

	return nil
}

func AddSalesMatlNoSalesCfg(tx *sql.Tx, cfg common.SalesMatlNoSalesCfg,
	m common.SalesPlantMeta) error {

	stmt, err := db.Prepare("INSERT INTO sales_matl_nosales_cfg( " +
		"material, material_txt, plant, " +
		"plant_name, status) VALUES(?,?,?,?,?)")
	if err != nil {
		zaps.Error("insert into matl no-sales cfg failed:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cfg.Material, cfg.MaterialTxt, m.Plant,
		m.PlantName, m.Status)
	if err != nil {
		zaps.Error("db exec failed: ", err)
		return err
	}

	return err
}

/*
func UpdateSalesMatlNoSalesCfgTx(cfg common.SalesMatlNoSalesCfg) error {

	tx, _ := db.Begin()
	defer tx.Rollback()

	cfg.Material = common.GetLongMatlCode(cfg.Material)

	for _, c := range cfg.Plants {
		err := UpdateSalesMatlNoSalesCfg(tx, cfg, c)
		if err != nil {
			zaps.Error("<<< update matl no-sales cfg failed")
			return err
		}
	}

	err := tx.Commit()
	if err != nil {
		zaps.Errorf("transaction commit failed: %v", err)
		return err
	}

	return nil
}


func UpdateSalesMatlNoSalesCfg(tx *sql.Tx, cfg common.SalesMatlNoSalesCfg,
	m common.SalesPlantMeta) error {

	stmt, err := db.Prepare("UPDATE sales_matl_nosales_cfg SET " +
				"status = ? " +
				"WHERE material = ? AND plant = ?")
	if err != nil {
		zaps.Error("update cfg failed: ", err)
		return err
	}

	_, err = stmt.Exec(m.Status, cfg.Material, m.Plant)
	if err != nil {
		zaps.Error("db exec failed: ", err)
		return err
	}

	zaps.Info("<<< update sales matl no-sales cfg done")

	return err
}
*/

func DelSalesMatlNoSalesCfgByMatl(material string) error {

	material = common.GetLongMatlCode(material)

	_, err := db.Exec("DELETE FROM sales_matl_nosales_cfg "+
		"WHERE material = ?", material)
	if err != nil {
		zaps.Error("delete cfg exec failed: ", err)
		return err
	}

	zaps.Info("<<< del sales material no-sales cfg done")

	return nil
}

func DelSalesMatlNoSalesCfgByMatlPlant(material string, plant string) {

	material = common.GetLongMatlCode(material)

	_, err := db.Exec("DELETE FROM sales_matl_nosales_cfg "+
		"WHERE material = ? AND plant = ?",
		material, plant)
	if err != nil {
		zaps.Error("delete cfg exec failed: ", err)
	}

	zaps.Info("<<< del sales material no-sales cfg done")
}

func GetSalesMatlNoSalesCfgList() ([]common.SalesMatlPlantMeta, int, error) {

	var mList []common.SalesMatlPlantMeta
	var mtxt sql.NullString
	var rows *sql.Rows
	var count int

	rows, err := db.Query("SELECT a.material, b.materialtxt, a.plant, " +
		"c.bic_ztxt_jyz, a.status " +
		"FROM sales_matl_nosales_cfg a " +
		"LEFT JOIN material b ON a.material = b.material " +
		"LEFT JOIN zaplant_xy c ON a.plant = c.bic_zaplant")

	if err != nil {
		zaps.Error("db query failed: ", err)
		return mList, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var m common.SalesMatlPlantMeta
		err := rows.Scan(&m.Material, &mtxt, /*&m.MaterialTxt*/
			&m.Plant, &m.PlantName, &m.Status)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> material: ", m.Material)
			zaps.Debug(">>> material name: ", m.MaterialTxt)
			zaps.Debug(">>> plant: ", m.Plant)
			zaps.Debug(">>> plant name: ", m.PlantName)
			zaps.Debug(">>> status: ", m.Status)

			m.Material = common.GetShortMatlCode(m.Material)
			if mtxt.Valid {
				m.MaterialTxt = mtxt.String
			} else {
				m.MaterialTxt = "N/A"
			}

			mList = append(mList, m)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return mList, count, err
	}

	zaps.Info("<<< get sales matl no-sales cfg list done")

	return mList, count, err
}

func GetSalesMatlNoSalesCfg(material string,
	plant string) (common.SalesMatlPlantMeta, bool, error) {

	var m common.SalesMatlPlantMeta
	var mtxt sql.NullString
	var rows *sql.Rows
	var find bool

	material = common.GetLongMatlCode(material)

	rows, err := db.Query("SELECT material, material_txt, plant, "+
		"plant_name, status "+
		"FROM sales_matl_nosales_cfg "+
		"WHERE material = ? AND plant = ?",
		material, plant)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return m, find, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&m.Material, &mtxt, &m.Plant,
			&m.PlantName, &m.Status)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> material: ", m.Material)
			zaps.Debug(">>> plant: ", m.Plant)
			zaps.Debug(">>> status: ", m.Status)

			m.Material = common.GetShortMatlCode(m.Material)
			if mtxt.Valid {
				m.MaterialTxt = mtxt.String
			} else {
				m.MaterialTxt = "N/A"
			}

			find = true
			break
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return m, find, err
	}

	zaps.Info("<<< get sales matl no-sales cfg done")

	return m, find, err
}
