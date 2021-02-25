package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	_ "github.com/go-sql-driver/mysql"

	"encoding/base64"
	"encoding/json"
	"database/sql"
	"strings"
	"errors"
	"time"
	"fmt"
)


func ParsePriceZoneCfg(ccode string) (string, error) {

	var pzc common.PriceZoneCfg
	var zi common.ZoneInfo
	var cmd, sqlA, sqlB string

	//XXX now we ONLY support main-class price zone
	mcode := ccode[0:4]
	pzc, find, err := GetPriceZoneCfgByClass(mcode)
	if err != nil {
		zaps.Error("get price zone cfg error: ", err)
		return "", err
	}

	if find != true {
		pzc, find, err = GetPriceZoneCfgByClass(common.PRICE_ZONE_DEFAULT)
		if err != nil || find != true {
			cmd = "ELT(INTERVAL(any_value(a.unit_price), " +
			"0.00,5.00,10.00,15.00,20.00,25.00,30.00,35.00,40.00,45.00,50.00), " +
			"'0.00-5.00','5.00-10.00','10.00-15.00','15.00-20.00', " +
			"'20.00-25.00','25.00-30.00','30.00-35.00','35.00-40.00', " +
			"'40.00-45.00','45.00-50.00','>=50.00')  SECTION,"
			zaps.Warn("using interval default config")
			return cmd, nil
		}

		zaps.Warn("using db default config")
	}

	jsonStr, err := base64.StdEncoding.DecodeString(pzc.Zones)

	json.Unmarshal(jsonStr, &zi)

	zaps.Info("price zone count: ", zi.Count)

	size := len(zi.Values)
	for i, v := range zi.Values {

		var s1, s2 string
		if v.End == common.PRICE_ZONE_MAX {
			s1 = fmt.Sprintf("%.2f", v.Begin)
			s2 = fmt.Sprintf("'>=%.2f' ", v.Begin)
		} else {
			s1 = fmt.Sprintf("%.2f", v.Begin)
			s2 = fmt.Sprintf("'%.2f-%.2f'", v.Begin, v.End)
		}

		sqlA = sqlA + s1
		sqlB = sqlB + s2

		if (i != size - 1) {
			sqlA = sqlA + ","
			sqlB = sqlB + ","
		}
	}

	cmd = fmt.Sprintf("ELT(INTERVAL(any_value(a.unit_price), %s), %s)", sqlA, sqlB)
	zaps.Debug("zone cmd: ", cmd)

	return cmd, nil
}


func ParsePriceZoneString(zone string) (string, error) {

	var cmd string

	if zone[0:2] == ">=" {
		cmd = fmt.Sprintf("AND a.unit_price%s ", zone)
	} else {
		idx := strings.Index(zone, "-")
		if idx <= 0 {
			zaps.Info("invalid zone string: ", zone)
			return "", errors.New("invalid zone")
		}

		cmd = fmt.Sprintf("AND a.unit_price >= %s AND a.unit_price < %s",
				zone[0:idx], zone[idx+1:len(zone)-1])
	}

	return cmd, nil
}


///////////////////////////////////////////////////////////////////////////////
// price-zone


func GetPriceZone(fr common.FilterReq) ([]common.PriceZoneInfo, int, error) {

	var pzList []common.PriceZoneInfo
	var sqlcmd string
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get price-zone list")

	code, _, err := common.GetClassInfo(fr.ClassLevel)
	if err != nil {
		zaps.Errorf("get class info failed: %v", err)
		return pzList, 0, err
	}

	sqla, err := ParsePriceZoneCfg(fr.ClassCode)
	if err != nil {
		zaps.Error("parse price zone error: ", err)
		return pzList, 0, err
	}
	sqla0:= "SECTION,"
	sqlb := "SUM(a.inv_qty) as SALECOUNT, SUM(a.gross_val) AS SALEINCOME, " +
		"SUM(a.netval_inv) AS SUMNET, SUM(a.cost) AS SUMCOST "
	sqlc := "bill_zsd a force index(idx_calday) " +
		"LEFT OUTER JOIN zaplant b ON a.plant = b.bic_zaplant " +
		"LEFT OUTER JOIN material c ON a.material = c.material "
	sqld := fmt.Sprintf("a.calday >= '%s' AND a.calday <= '%s' AND a.no_inv_it > 0 ",
			fr.BeginDate, fr.EndDate)
	sqle1:= fmt.Sprintf("AND b.bic_zrpa_lcit = '%s' ", fr.OrgCode)
	sqle2:= fmt.Sprintf("AND a.plant = '%s' ", fr.OrgCode)
	sqlf := fmt.Sprintf("AND c.%s = %s ", code, fr.ClassCode)
	sqlg := "SALECOUNT DESC"

	if fr.OrgLevel == common.ORG_LEVEL_HEAD {
		sqlcmd = fmt.Sprintf("SELECT %s %s %s FROM %s WHERE %s %s GROUP BY %s ORDER BY %s",
				sqla, sqla0, sqlb, sqlc, sqld, sqlf, sqla, sqlg)
	} else if fr.OrgLevel == common.ORG_LEVEL_BRANCH {
		sqlcmd = fmt.Sprintf("SELECT %s %s %s FROM %s WHERE %s %s %s GROUP BY %s ORDER BY %s",
				sqla, sqla0, sqlb, sqlc, sqld, sqle1, sqlf, sqla, sqlg)
	} else if fr.OrgLevel == common.ORG_LEVEL_PLANT {
		sqlcmd = fmt.Sprintf("SELECT %s %s %s FROM %s WHERE %s %s %s GROUP BY %s ORDER BY %s",
				sqla, sqla0, sqlb, sqlc, sqld, sqle2, sqlf, sqla, sqlg)
	}

	zaps.Info("sql cmd: ", sqlcmd)
	t1 := time.Now().UnixNano() / 1e6

	rows, err = db.Query(sqlcmd)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return pzList, 0, err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("db query using %d ms", (t2 - t1))

	defer rows.Close()

	for rows.Next() {
		var zp common.PriceZoneInfo
		err := rows.Scan(&zp.PriceZone, &zp.SalesCount, &zp.SalesIncome,
				&zp.SalesNet, &zp.SalesCost)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> PriceZone: ", zp.PriceZone)
			zaps.Debug(">>> SalesCount: ", zp.SalesCount)
			zaps.Debug(">>> SalesIncome: ", zp.SalesIncome)
			zaps.Debug(">>> SalesNet: ", zp.SalesNet)
			zaps.Debug(">>> SalesCost: ", zp.SalesCost)

			pzList = append(pzList, zp)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return pzList, count, err
	}

	zaps.Info("<<< get price-zone done with count ", count)

	return pzList, count, err
}


///////////////////////////////////////////////////////////////////////////////
// price zone rank

func GetPriceZoneRank(fr common.FilterReq) ([]common.MatlRankInfo, int, error) {

	var mList []common.MatlRankInfo
	var sqlcmd string
	var rows *sql.Rows
	var count int

	zaps.Info(">>> get price-zone rank list")

	code, _, err := common.GetClassInfo(fr.ClassLevel)
	if err != nil {
		zaps.Errorf("get class info failed: %v", err)
		return mList, 0, err
	}

	sqla := "c.material, c.materialtxt, SUM(a.inv_qty) AS SALECOUNT, " +
		"SUM(a.gross_val) AS SALEINCOME,  SUM(a.netval_inv) AS SUMNET, " +
		"SUM(a.cost) AS SUMCOST "
	sqlb := "bill_zsd a force index(idx_unit_calday) " +
		"LEFT OUTER JOIN zaplant b ON a.plant = b.bic_zaplant " +
		"LEFT OUTER JOIN material c ON a.material = c.material "
	sqlc := fmt.Sprintf("a.calday >= '%s' AND a.calday <= '%s' AND a.no_inv_it > 0 ",
			fr.BeginDate, fr.EndDate)
	sqld, err := ParsePriceZoneString(fr.ZoneTxt)
	if err != nil {
		zaps.Error("parse price zone string failed: ", err)
		return mList, 0, err
	}
	sqle1:= fmt.Sprintf("AND b.bic_zrpa_lcit = '%s'", fr.OrgCode)
	sqle2:= fmt.Sprintf("AND a.plant = '%s'", fr.OrgCode)
	sqlf := fmt.Sprintf("AND c.%s = %s", code, fr.ClassCode)
	sqlg := "c.material ORDER BY SALECOUNT DESC LIMIT 10"

	if fr.OrgLevel == common.ORG_LEVEL_HEAD {
		sqlcmd = fmt.Sprintf("SELECT %s FROM %s WHERE %s %s %s GROUP BY %s",
			sqla, sqlb, sqlc, sqld, sqlf, sqlg)
	} else if fr.OrgLevel == common.ORG_LEVEL_BRANCH {
		sqlcmd = fmt.Sprintf("SELECT %s FROM %s WHERE %s %s %s %s GROUP BY %s",
			sqla, sqlb, sqlc, sqld, sqle1, sqlf, sqlg)
	} else if fr.OrgLevel == common.ORG_LEVEL_PLANT {
		sqlcmd = fmt.Sprintf("SELECT %s FROM %s WHERE %s %s %s %s GROUP BY %s",
			sqla, sqlb, sqlc, sqld, sqle2, sqlf, sqlg)
	}

	zaps.Info("sql cmd: ", sqlcmd)
	rows, err = db.Query(sqlcmd)

	if err != nil {
		zaps.Error("db query failed: ", err)
		return mList, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var mr common.MatlRankInfo
		err := rows.Scan(&mr.Material, &mr.MaterialTxt, &mr.SalesCount,
				&mr.SalesIncome, &mr.SalesNet, &mr.SalesCost)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> Material: ", mr.Material)
			zaps.Debug(">>> MaterialTxt: ", mr.MaterialTxt)
			zaps.Debug(">>> SalesCount: ", mr.SalesCount)
			zaps.Debug(">>> SalesIncome: ", mr.SalesIncome)
			zaps.Debug(">>> SalesNet: ", mr.SalesNet)
			zaps.Debug(">>> SalesCost: ", mr.SalesCost)

			mList = append(mList, mr)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return mList, 0, err
	}

	zaps.Infof("<<< get price-zone rank done with count (%d)", count)

	return mList, count, err
}


