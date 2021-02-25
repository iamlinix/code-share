package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"time"
)

///////////////////////////////////////////////////////////////////////////////
// sales plant car service

func GetPlantCarServiceList(fr common.FilterReq) ([]common.SalesPlantCarService,
	int, error) {

	var pcsList []common.SalesPlantCarService
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get plant car service list")

	sqlcmd := fmt.Sprintf("SELECT a.plant, b.bic_ztxt_jyz, "+
		"SUM(a.gross_val) as CSALES "+
		"FROM bill_zsd a "+
		"LEFT JOIN zaplant_xy b ON a.plant = b.bic_zaplant "+
		"WHERE a.zklad2 = '3001' "+
		"AND a.plant != 'A0A1' "+
		"AND a.calday BETWEEN '%s' AND '%s' ", fr.BeginDate, fr.EndDate)
	if fr.RateThreshold > 0.0 {
		sqlcmd += fmt.Sprintf(" AND (netval_inv-cost)/netval_inv >= %.3f", fr.RateThreshold)
	}
	sqlcmd += " GROUP BY a.plant"

	zaps.Info("sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return pcsList, 0, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		var pcs common.SalesPlantCarService
		err := rows.Scan(&pcs.Plant, &pcs.PlantName, &pcs.NetSales)
		if err != nil {
			zaps.Error("query error: ", err)
			return pcsList, count, err
		}

		pcsList = append(pcsList, pcs)
		count++
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return pcsList, count, err
	}

	zaps.Info("<<< get plant car service done with count ", count)

	return pcsList, count, err
}

func GetPlantCarService(fr common.FilterReq) (common.SalesPlantCarService, error) {

	var spcs common.SalesPlantCarService
	var rows *sql.Rows

	zaps.Info(">>> get plant car service")

	sqlcmd := fmt.Sprintf("SELECT a.plant, b.bic_ztxt_jyz, "+
		"SUM(a.gross_val) as CSALES "+
		"FROM bill_zsd a "+
		"LEFT JOIN zaplant_xy b ON a.plant = b.bic_zaplant "+
		"WHERE a.zklad2 = '3001' "+
		"AND a.plant = '%s' "+
		"AND a.calday BETWEEN '%s' AND '%s'", fr.OrgCode, fr.BeginDate, fr.EndDate)
	if fr.RateThreshold > 0.0 {
		sqlcmd += fmt.Sprintf(" AND (netval_inv-cost)/netval_inv >= %.3f", fr.RateThreshold)
	}
	zaps.Info("sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return spcs, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return spcs, err
	}

	if rows.Next() {
		err := rows.Scan(&spcs.Plant, &spcs.PlantName, &spcs.NetSales)
		if err != nil {
			zaps.Error("query error: ", err)
			return spcs, err
		}
	}

	return spcs, err
}
