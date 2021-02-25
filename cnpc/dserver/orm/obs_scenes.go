package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"encoding/json"
	"time"
)

func AddScenes(s common.Scenes) (int64, error) {

	zaps.Info(">>> add scenes")

	stmt, err := db.Prepare("INSERT INTO scenes(name, begin_date, end_date, " +
		"org_info, class_info, user, remark, created_time) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		zaps.Error("insert into scenes failed: ", err)
		return 0, err
	}
	defer stmt.Close()

	org, _ := json.Marshal(s.OrgInfo)
	css, _ := json.Marshal(s.ClassInfo)

	res, err := stmt.Exec(s.Name, s.BeginDate, s.EndDate, org, css,
		s.User, s.Remark,
		time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		zaps.Errorf("db exec failed: %v", err)
		return 0, err
	}

	id, err := res.LastInsertId()

	zaps.Info("<<< add scenes done: ", id)

	return id, err
}

func UpdateScenes(s common.Scenes) error {

	zaps.Info(">>> update scenes")

	stmt, err := db.Prepare("UPDATE scenes SET begin_date = ?, end_date = ?, " +
		"org_info = ?, class_info = ?, user = ?, remark = ?, " +
		"created_time = ? WHERE name = ?")
	if err != nil {
		zaps.Error("update scenes failed: ", err)
		return err
	}
	defer stmt.Close()

	org, _ := json.Marshal(s.OrgInfo)
	css, _ := json.Marshal(s.ClassInfo)

	_, err = stmt.Exec(s.BeginDate, s.EndDate, org, css, s.User, s.Remark,
		time.Now().Format("2006-01-02 15:04:05"), s.Name)
	if err != nil {
		zaps.Error("db exec failed: ", err)
		return err
	}

	zaps.Info("<<< update scenes done")

	return err
}

func DelScenes(id int64) {

	zaps.Info(">>> del scenes")

	_, err := db.Exec("DELETE FROM scenes WHERE id = ?", id)
	if err != nil {
		zaps.Error("delete scenes exec failed: ", err)
	}

	zaps.Info("<<< del scenes done")
}

func GetScenesList(all int, page int) ([]common.Scenes, int, error) {

	var sList []common.Scenes
	var rows *sql.Rows
	var err error
	var count int

	zaps.Info(">>> get scenes list with page ", page)

	if all == 1 {
		rows, err = db.Query("SELECT id, name, begin_date, end_date, " +
			"org_info, class_info, user, remark, created_time " +
			"FROM scenes ORDER BY created_time DESC")
	} else {
		rows, err = db.Query("SELECT id, name, begin_date, end_date, "+
			"org_info, class_info, user, remark, created_time "+
			"FROM user LIMIT ?,10", page*10)
	}

	if err != nil {
		zaps.Error("db query failed: ", err)
		return sList, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var s common.Scenes
		var oi, ci string
		err := rows.Scan(&s.ID, &s.Name, &s.BeginDate, &s.EndDate,
			&oi, &ci, &s.User, &s.Remark,
			&s.CreatedTime)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> id: ", s.ID)
			zaps.Debug(">>> name: ", s.Name)

			json.Unmarshal([]byte(oi), &s.OrgInfo)
			json.Unmarshal([]byte(ci), &s.ClassInfo)

			sList = append(sList, s)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return sList, count, err
	}

	zaps.Info("<<< get scenes list done")

	return sList, count, err
}

func GetScenesByName(name string) (common.Scenes, bool, error) {

	var res common.Scenes
	find := false

	zaps.Info(">>> get one scenes info: ", name)

	rows, err := db.Query("SELECT id, name, begin_date, end_date, "+
		"org_info, class_info, user, remark, created_time "+
		"FROM scenes WHERE name = ?", name)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return res, find, err
	}

	defer rows.Close()

	for rows.Next() {
		var s common.Scenes
		var oi, ci string
		err := rows.Scan(&s.ID, &s.Name, &s.BeginDate, &s.EndDate,
			&oi, &ci, &s.User, &s.Remark, &s.CreatedTime)
		if err != nil {
			zaps.Error("query error: ", err)
			return res, find, err
		} else {
			zaps.Debug(">>> id: ", s.ID)
			zaps.Debug(">>> name: ", s.Name)

			json.Unmarshal([]byte(oi), &s.OrgInfo)
			json.Unmarshal([]byte(ci), &s.ClassInfo)

			res = s
			find = true
			break
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return res, find, err
	}

	zaps.Info("<<< get one scenes info done")

	return res, find, err
}
