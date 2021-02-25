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
// Vendors Purchase

func GetVendorsFirstPurchaseDate(fr common.FilterReq) ([]common.ERPVendorFPDate, int, error) {

	var vdList []common.ERPVendorFPDate
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get vendors first purchase date list")

	sqlcmd := fmt.Sprintf("SELECT vendor, MIN(t2.zgr_date), " +
			"FROM zifpurd " +
			"WHERE zgr_date BETWEEN '%s' AND '%s' " +
			"GROUP BY vendor", fr.BeginDate, fr.EndDate)

	zaps.Info(">>> sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err := db.Query(sqlcmd)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return vdList, 0, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		var vd common.ERPVendorFPDate
		err := rows.Scan(&vd.Vendor, &vd.ZgrDate)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> vendor: ", vd.Vendor)
			zaps.Debug(">>> zgr date: ", vd.ZgrDate)

			vdList = append(vdList, vd)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return vdList, count, err
	}

	zaps.Info("<<< get vendors first purchase date done with ", count)

	return vdList, count, err
}


