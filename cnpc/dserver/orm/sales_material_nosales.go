package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"time"
	"fmt"
)


///////////////////////////////////////////////////////////////////////////////
// sales material no-sales

func GetSalesMatlNoSales(fr common.FilterReq) ([]common.SalesMatlPlantMeta,
	int, error) {

	var mList []common.SalesMatlPlantMeta
	var mtxt sql.NullString
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get sales matl no-sales meta list")

	/*
	sqlcmd := fmt.Sprintf("SELECT a.material, a.material_txt, " +
		"a.plant, a.plant_name " +
		"FROM sales_matl_nosales_cfg a " +
		"WHERE (a.material, a.plant) NOT IN (" +
		"SELECT material, plant FROM bill_zsd " +
		"WHERE calday BETWEEN '%s' AND '%s' " +
		"GROUP BY material, plant)", fr.BeginDate, fr.EndDate)
	*/
	sqlcmd := fmt.Sprintf("SELECT a.material, b.materialtxt, " +
		"a.plant, c.bic_ztxt_jyz " +
		"FROM sales_matl_nosales_cfg a " +
		"LEFT JOIN material b ON a.material = b.material " +
		"LEFT JOIN zaplant_xy c ON a.plant = c.bic_zaplant " +
		"WHERE (a.material, a.plant) NOT IN (" +
		"SELECT material, plant FROM bill_zsd " +
		"WHERE calday BETWEEN '%s' AND '%s' " +
		"GROUP BY material, plant)", fr.BeginDate, fr.EndDate)

	zaps.Info("sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return mList, 0, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		var m common.SalesMatlPlantMeta
		err := rows.Scan(&m.Material, &mtxt, /*&p.MaterialName,*/
				&m.Plant, &m.PlantName)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> Material: ", m.Material)
			zaps.Debug(">>> Plant: ", m.Plant)

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

	zaps.Info("<<< get sales matl no-sales done with count ", count)

	return mList, count, err
}


