package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func AddPriceZoneCfg(pzc common.PriceZoneCfg) (int64, error) {

	zaps.Info(">>> add price zone cfg")

	stmt, err := db.Prepare("INSERT INTO price_zone_cfg(class_level, class_code, " +
		"zones, user) VALUES(?, ?, ?, ?)")
	if err != nil {
		zaps.Error("insert into price zone cfg failed: ", err)
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(pzc.ClassLevel, pzc.ClassCode, pzc.Zones, pzc.User)
	if err != nil {
		zaps.Error("db exec failed: err")
		return 0, err
	}

	id, err := res.LastInsertId()

	zaps.Info("<<< add user done: ", id)

	return id, err
}

func UpdatePriceZoneCfg(pzc common.PriceZoneCfg) error {

	zaps.Info(">>> update price zone cfg")

	stmt, err := db.Prepare("UPDATE price_zone_cfg SET class_level = ?, zones = ?, " +
		"user = ? WHERE class_code = ?")
	if err != nil {
		zaps.Error("update price zone cfg failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(pzc.ClassLevel, pzc.Zones, pzc.User, pzc.ClassCode)
	if err != nil {
		zaps.Error("db exec failed: err")
		return err
	}

	zaps.Info("<<< update price zone cfg done")

	return err
}

func DelPriceZoneCfg(ccode string) {

	zaps.Info(">>> del price zone cfg")

	_, err := db.Exec("DELETE FROM price_zone_cfg WHERE class_code = ?",
		ccode)
	if err != nil {
		zaps.Error("delete price zone cfg exec failed: ", err)
	}

	zaps.Info("<<< del price zone cfg done")
}

func GetPriceZoneCfgList(all int, page int) ([]common.PriceZoneCfg, int, error) {

	var pzcList []common.PriceZoneCfg
	var rows *sql.Rows
	var err error
	var count int

	zaps.Info(">>> get price zone cfg list with page ", page)

	if all == 1 {
		rows, err = db.Query("SELECT class_level, class_code, zones, user " +
			"FROM price_zone_cfg")
	} else {
		rows, err = db.Query("SELECT class_level, class_code, zones, user "+
			"FROM price_zone_cfg LIMIT ?,10", page*10)
	}

	if err != nil {
		zaps.Error("db query failed: ", err)
		return pzcList, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var pzc common.PriceZoneCfg
		err := rows.Scan(&pzc.ClassLevel, &pzc.ClassCode, &pzc.Zones, &pzc.User)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> class level: ", pzc.ClassLevel)
			zaps.Debug(">>> class code: ", pzc.ClassCode)
			zaps.Debug(">>> zones: ", pzc.Zones)
			zaps.Debug(">>> user: ", pzc.User)

			c := &common.ClassInfo{
				ClassLevel: pzc.ClassLevel,
				ClassCode:  pzc.ClassCode,
			}
			GetClassTextByLevelCode(c)
			pzc.ClassText = c.ClassText

			pzcList = append(pzcList, pzc)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return pzcList, count, err
	}

	zaps.Info("<<< get price zone cfg list done")

	return pzcList, count, err
}

func GetPriceZoneCfgByClass(ccode string) (common.PriceZoneCfg, bool, error) {

	var pzc common.PriceZoneCfg
	find := false

	zaps.Info(">>> get one price zone cfg info: ", ccode)

	rows, err := db.Query("SELECT class_level, class_code, zones, user "+
		"FROM price_zone_cfg WHERE class_code = ?", ccode)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return pzc, find, err
	}

	defer rows.Close()

	for rows.Next() {
		var c common.PriceZoneCfg
		err := rows.Scan(&c.ClassLevel, &c.ClassCode, &c.Zones, &c.User)
		if err != nil {
			zaps.Error("query error: ", err)
			return pzc, find, err
		} else {
			zaps.Debug(">>> class level: ", pzc.ClassLevel)
			zaps.Debug(">>> class code: ", pzc.ClassCode)
			zaps.Debug(">>> zones: ", pzc.Zones)
			zaps.Debug(">>> user: ", pzc.User)

			ci := &common.ClassInfo{
				ClassLevel: c.ClassLevel,
				ClassCode:  c.ClassCode,
			}
			GetClassTextByLevelCode(ci)
			c.ClassText = ci.ClassText

			pzc = c
			find = true
			break
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return pzc, find, err
	}

	zaps.Info("<<< get one price zone cfg info done")

	return pzc, find, err
}
