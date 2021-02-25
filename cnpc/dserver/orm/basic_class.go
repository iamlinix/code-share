package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"

	"errors"
	"fmt"
)

func DropBasicClassTable() error {

	zaps.Info(">>> drop basic class table")

	stmt, err := db.Prepare("DROP TABLE IF EXISTS basic_class")
	if err != nil {
		zaps.Error("drop basic class db prepare failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		zaps.Error("exec db failed: ", err)
		return err
	}

	zaps.Info("<<< drop basic class table done")

	return nil
}

func InitBasicClassTable() error {

	zaps.Info(">>> init basic class table")

	stmt, err := db.Prepare("CREATE TABLE basic_class(" +
		"SELECT any_value(m.bic_zrpa_mtl) AS sub_class_code, " +
		"any_value(m.zrpa_mtltxt) AS sub_class_txt, " +
		"any_value(m.bic_zklasse_d) AS mid_class_code, " +
		"any_value(m.zklasse_dtxt) AS mid_class_txt, " +
		"any_value(m.bic_zklad2) AS main_class_code, " +
		"any_value(m.zklad2txt) AS main_class_txt " +
		"FROM material m " +
		"GROUP BY m.bic_zrpa_mtl)")
	if err != nil {
		zaps.Error("init basic class table prepare failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		zaps.Error("exec db failed: ", err)
		return err
	}

	zaps.Info("<<< init basic class table done")

	return nil
}

func DropBasicClassDictTable() error {

	zaps.Info(">>> drop basic class dict table")

	stmt, err := db.Prepare("DROP TABLE IF EXISTS basic_class_dict")
	if err != nil {
		zaps.Error("drop basic class dict db prepare failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		zaps.Error("exec db failed: ", err)
		return err
	}

	zaps.Info("<<< drop basic class dict table done")

	return nil
}

func InitBasicClassDictTable() error {

	zaps.Info(">>> init basic class dict table")

	stmt, err := db.Prepare("CREATE TABLE basic_class_dict(" +
		"class_code varchar(32), class_txt varchar(128), " +
		"class_level int)")
	if err != nil {
		zaps.Error("init class dict table create prepare failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		zaps.Error("exec db failed: ", err)
		return err
	}

	/* insert main class */
	stmt, err = db.Prepare("INSERT INTO basic_class_dict(class_code, " +
		"class_txt, class_level) " +
		"SELECT DISTINCT bic_zklad2, zklad2txt, 0 " +
		"FROM material m")
	if err != nil {
		zaps.Error("init class dict main table prepare failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		zaps.Error("exec db failed: ", err)
		return err
	}

	/* insert mid class */
	stmt, err = db.Prepare("INSERT INTO basic_class_dict(class_code, " +
		"class_txt, class_level) " +
		"SELECT DISTINCT bic_zklasse_d, zklasse_dtxt, 1 " +
		"FROM material m")
	if err != nil {
		zaps.Error("init class dict mid table prepare failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		zaps.Error("exec db failed: ", err)
		return err
	}

	/* insert sub class */
	stmt, err = db.Prepare("INSERT INTO basic_class_dict(class_code, " +
		"class_txt, class_level) " +
		"SELECT DISTINCT bic_zrpa_mtl, zrpa_mtltxt, 2 " +
		"FROM material m")
	if err != nil {
		zaps.Error("init class dict sub table prepare failed: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		zaps.Error("exec db failed: ", err)
		return err
	}

	zaps.Info("<<< init basic class dict table done")

	return nil
}

///////////////////////////////////////////////////////////////////////////////
//
//main class
func GetMainClassList() ([]common.ClassInfo, int, error) {

	var mcList []common.ClassInfo
	var rows *sql.Rows
	var err error
	var count int

	zaps.Info(">>> get main class list")

	rows, err = db.Query("SELECT DISTINCT main_class_code, main_class_txt " +
		"FROM basic_class")

	if err != nil {
		zaps.Error("GetMainClassList db query failed: ", err)
		return mcList, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var mc common.ClassInfo
		var text string
		err := rows.Scan(&mc.ClassCode, &text)
		if err != nil {
			zaps.Error("GetMainClassList query error: ", err)
			return mcList, count, err
		}

		mc.ClassText = fmt.Sprintf("%s-%s", mc.ClassCode, text)

		mcList = append(mcList, mc)
		count++
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return mcList, count, err
	}

	zaps.Info("<<< get main class list done")

	return mcList, count, err
}

//mid class
func GetMidClassList() ([]common.ClassInfo, int, error) {

	var mcList []common.ClassInfo
	var rows *sql.Rows
	var err error
	var count int

	zaps.Info(">>> get mid class list")

	rows, err = db.Query("SELECT mid_class_code, any_value(mid_class_txt), " +
		"any_value(main_class_txt) FROM basic_class " +
		"GROUP BY mid_class_code ")

	if err != nil {
		zaps.Error("GetMidClassList db query failed: ", err)
		return mcList, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var mc common.ClassInfo
		var midtxt, maintxt string
		err := rows.Scan(&mc.ClassCode, &midtxt, &maintxt)
		if err != nil {
			zaps.Error("GetMidClassList query error: ", err)
			return mcList, count, err
		}

		mc.ClassText = fmt.Sprintf("%s-%s-%s", mc.ClassCode,
			maintxt, midtxt)

		mcList = append(mcList, mc)
		count++
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return mcList, count, err
	}

	zaps.Info("<<< get mid class list done")

	return mcList, count, err
}

func GetMidClassListByMain(ccode string) ([]common.ClassInfo, int, error) {

	var mcList []common.ClassInfo
	var rows *sql.Rows
	var err error
	var count int

	zaps.Infof(">>> get mid class list with code", ccode)

	sql := fmt.Sprintf("SELECT DISTINCT mid_class_code, mid_class_txt "+
		"FROM basic_class WHERE main_class_code = %s", ccode)
	rows, err = db.Query(sql)

	if err != nil {
		zaps.Error("GetMidClassList db query failed: ", err)
		return mcList, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var mc common.ClassInfo
		err := rows.Scan(&mc.ClassCode, &mc.ClassText)
		if err != nil {
			zaps.Error("GetMidClassList query error: ", err)
			return mcList, 0, err
		}

		mcList = append(mcList, mc)
		count++
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return mcList, count, err
	}

	zaps.Info("<<< get mid class list done")

	return mcList, count, err
}

//sub class
func GetSubClassList() ([]common.ClassInfo, int, error) {

	var scList []common.ClassInfo
	var rows *sql.Rows
	var count int
	var err error

	zaps.Info(">>> get sub class list")

	rows, err = db.Query("SELECT sub_class_code, any_value(sub_class_txt), " +
		"any_value(mid_class_txt), any_value(main_class_txt) " +
		"FROM basic_class GROUP BY sub_class_code ")

	if err != nil {
		zaps.Error("GetSubClassList db query failed: ", err)
		return scList, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var sc common.ClassInfo
		var maintxt, midtxt, subtxt string
		err := rows.Scan(&sc.ClassCode, &subtxt, &midtxt, &maintxt)
		if err != nil {
			zaps.Error("GetSiubClassList query error: ", err)
			return scList, 0, err
		}

		sc.ClassText = fmt.Sprintf("%s-%s-%s-%s", sc.ClassCode,
			maintxt, midtxt, subtxt)

		scList = append(scList, sc)
		count++
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return scList, count, err
	}

	zaps.Info("<<< get sub class list done")

	return scList, count, err
}

func GetSubClassListByMid(ccode string) ([]common.ClassInfo, int, error) {

	var scList []common.ClassInfo
	var rows *sql.Rows
	var err error
	var count int

	zaps.Infof(">>> get sub class list with code: %s", ccode)

	sql := fmt.Sprintf("SELECT DISTINCT sub_class_code, sub_class_txt "+
		"FROM basic_class WHERE mid_class_code = %s", ccode)
	rows, err = db.Query(sql)

	if err != nil {
		zaps.Error("GetSubClassList db query failed: ", err)
		return scList, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var sc common.ClassInfo
		err := rows.Scan(&sc.ClassCode, &sc.ClassText)
		if err != nil {
			zaps.Error("GetSubClassList query error: ", err)
			return scList, 0, err
		}

		scList = append(scList, sc)
		count++
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return scList, count, err
	}

	zaps.Info("<<< get sub class list done")

	return scList, count, err
}

func GetClassTextByLevelCode(class *common.ClassInfo) (bool, error) {

	var sql string
	find := false

	zaps.Info(">>> get class text by level code: ", class.ClassCode)

	level := class.ClassLevel
	code := class.ClassCode

	if level == common.CLASS_LEVEL_MAIN {
		sql = fmt.Sprintf("SELECT DISTINCT main_class_txt "+
			"FROM basic_class "+
			"WHERE main_class_code = '%s'", code)
	} else if level == common.CLASS_LEVEL_MID {
		sql = fmt.Sprintf("SELECT DISTINCT mid_class_txt "+
			"FROM basic_class "+
			"WHERE mid_class_code = '%s'", code)
	} else if level == common.CLASS_LEVEL_SUB {
		sql = fmt.Sprintf("SELECT DISTINCT sub_class_txt "+
			"FROM basic_class "+
			"WHERE sub_class_code = '%s'", code)
	} else {
		zaps.Error("invalid class level: ", level)
		return find, errors.New(common.ERR_MSG_INVALID_CLASS)
	}

	rows, err := db.Query(sql)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return find, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&class.ClassText)
		if err != nil {
			zaps.Error("query error: ", err)
			return find, err
		}

		find = true
		break
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return find, err
	}

	zaps.Info("<<< get class text by level code done")

	return find, err
}

func GetClassFullText(level int, code string) (string, error) {

	var text string
	var err error

	zaps.Info(">>> get class full text: ", code)

	if level == common.CLASS_LEVEL_MAIN {
		text, err = GetMainClassText(code)
	} else if level == common.CLASS_LEVEL_MID {
		text, err = GetMidClassText(code)
	} else if level == common.CLASS_LEVEL_SUB {
		text, err = GetSubClassText(code)
	} else {
		zaps.Error("invalid class level: ", level)
		return "", errors.New(common.ERR_MSG_INVALID_CLASS)
	}

	if err != nil {
		zaps.Error("invalid class text")
		return "", err
	}

	return text, nil
}

func GetMainClassText(code string) (string, error) {

	var rows *sql.Rows
	var err error
	var text string

	sql := fmt.Sprintf("SELECT DISTINCT main_class_txt "+
		"FROM basic_class "+
		"WHERE main_class_code = '%s'", code)

	rows, err = db.Query(sql)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return "", err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&text)
		if err != nil {
			zaps.Error("query error: ", err)
			return "", err
		}

		break
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return "", err
	}

	zaps.Infof("<<< get main class text (%s) done", text)

	return text, err
}

func GetMidClassText(code string) (string, error) {

	var rows *sql.Rows
	var err error
	var text string

	sql := fmt.Sprintf("SELECT DISTINCT main_class_txt, mid_class_txt "+
		"FROM basic_class WHERE mid_class_code = '%s'", code)

	rows, err = db.Query(sql)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return "", err
	}

	defer rows.Close()

	for rows.Next() {
		var main, mid string
		err := rows.Scan(&main, &mid)
		if err != nil {
			zaps.Error("query error: ", err)
			return "", err
		}

		text = fmt.Sprintf("%s-%s", main, mid)
		break
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return "", err
	}

	zaps.Infof("<<< get mid class text (%s) done", text)

	return text, nil
}

func GetSubClassText(code string) (string, error) {

	var rows *sql.Rows
	var err error
	var text string

	sql := fmt.Sprintf("SELECT DISTINCT main_class_txt, mid_class_txt, "+
		"sub_class_txt FROM basic_class "+
		"WHERE sub_class_code = '%s'", code)

	rows, err = db.Query(sql)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return "", err
	}

	defer rows.Close()

	for rows.Next() {
		var main, mid, sub string
		err := rows.Scan(&main, &mid, &sub)
		if err != nil {
			zaps.Error("query error: ", err)
			return "", err
		}

		text = fmt.Sprintf("%s-%s-%s", main, mid, sub)
		break
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return "", err
	}

	zaps.Infof("<<< get sub class text (%s) done", text)

	return text, nil
}

func GetClassNameByCode(code string) (string, bool, error) {

	var rows *sql.Rows
	var name string
	var find bool

	sql := fmt.Sprintf("SELECT class_txt FROM basic_class_dict "+
		"WHERE class_code = '%s'", code)

	rows, err := db.Query(sql)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return name, find, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			zaps.Error("query error: ", err)
			return name, find, err
		}

		find = true
		break
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return name, find, err
	}

	zaps.Infof("<<< get class name (%s) done", name)

	return name, find, nil
}
