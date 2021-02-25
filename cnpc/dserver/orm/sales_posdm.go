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
// 

/* group by plant */
func GetOrgPosdmMetricList(fr common.FilterReq) ([]common.PosdmMetric,
	int, error) {

	var pmList []common.PosdmMetric
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get plant posdm metric list")

	sqlcmd := fmt.Sprintf("SELECT bic_zplant, bic_zof_type, " +
			"sum(bic_zof_count) " +
			"FROM posdm " +
			"WHERE rpa_bdd BETWEEN '%s' AND '%s' " +
			"GROUP BY bic_zplant, bic_zof_type",
			fr.BeginDate, fr.EndDate)

	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return pmList, 0, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		var pm common.PosdmMetric
		err := rows.Scan(&pm.BicZplant, &pm.BicZofType, &pm.BicZofCount)
		if err != nil {
			zaps.Error("query error: ", err)
			return pmList, 0, err
		}

		pmList = append(pmList, pm)
		count++
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return pmList, count, err
	}

	zaps.Info("<<< get plant posdm metric done with count ", count)

	return pmList, count, err
}


func GetPosdmMetricByOrg(fr common.FilterReq) ([]common.PosdmMetric,
	int, error) {

	var pmList []common.PosdmMetric
	var rows *sql.Rows
	var sqlcmd string
	var count int

	zaps.Info(">>> get posdm metric list")

	if fr.OrgLevel == common.ORG_LEVEL_HEAD {
		sqlcmd = fmt.Sprintf("SELECT bic_zof_type, sum(bic_zof_count) " +
			"FROM posdm " +
			"WHERE rpa_bdd BETWEEN '%s' AND '%s' " +
			"GROUP BY bic_zof_type",
			fr.BeginDate, fr.EndDate)

	} else if fr.OrgLevel == common.ORG_LEVEL_BRANCH {
		sqlcmd = fmt.Sprintf("SELECT a.bic_zof_type, sum(a.bic_zof_count) " +
			"FROM posdm a " +
			"LEFT JOIN zaplant_xy b ON a.bic_zplant = b.bic_zaplant " +
			"WHERE a.rpa_bdd BETWEEN '%s' AND '%s' " +
			"AND b.bic_zrpa_lcit = '%s' " +
			"GROUP BY a.bic_zof_type",
			fr.BeginDate, fr.EndDate, fr.OrgCode)

	} else if fr.OrgLevel == common.ORG_LEVEL_PLANT {
		sqlcmd = fmt.Sprintf("SELECT bic_zof_type, sum(bic_zof_count) " +
			"FROM posdm " +
			"WHERE rpa_bdd BETWEEN '%s' AND '%s' " +
			"AND bic_zplant = '%s' " +
			"GROUP BY bic_zof_type",
			fr.BeginDate, fr.EndDate, fr.OrgCode)
	}

	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return pmList, 0, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		var pm common.PosdmMetric
		err := rows.Scan(&pm.BicZofType, &pm.BicZofCount)
		if err != nil {
			zaps.Error("query error: ", err)
			return pmList, 0, err
		}

		pmList = append(pmList, pm)
		count++
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return pmList, count, err
	}

	zaps.Info("<<< get posdm metric done")

	return pmList, count, err
}


