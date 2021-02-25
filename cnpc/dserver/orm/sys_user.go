package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"time"
)

func AddUser(u common.User) (int64, error) {

	zaps.Info(">>> add user")

	stmt, err := db.Prepare("INSERT INTO user(username, password, role, create_time) " +
		"VALUES(?, ?, ?, ?)")
	if err != nil {
		zaps.Error("insert into user failed: ", err)
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(u.Username, u.Password, u.Role, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		zaps.Error("db exec failed: err")
		return 0, err
	}

	id, err := res.LastInsertId()

	zaps.Info("<<< add user done: ", id)

	return id, err
}

func UpdateUser(u common.User) error {

	zaps.Info(">>> update user")

	stmt, err := db.Prepare("UPDATE user SET password = ?, create_time = ? " +
		"WHERE username = ?")
	if err != nil {
		zaps.Error("update user failed: ", err)
		return err
	}
	defer stmt.Close()

	zaps.Info("update new password : ", u.NewPassword)
	_, err = stmt.Exec(u.NewPassword, time.Now().Format("2006-01-02 15:04:05"), u.Username)
	if err != nil {
		zaps.Error("db exec failed: err")
		return err
	}

	zaps.Info("<<< update user done")

	return err
}

func UpdateUserRole(username string, role int) error {
	zaps.Info(">>> update user role")
	stmt, err := db.Prepare("UPDATE user SET role = ? WHERE username = ?")
	if err != nil {
		zaps.Error("update user failed: ", err)
		return err
	}
	defer stmt.Close()

	zaps.Info("update new role : ", role)
	_, err = stmt.Exec(role, username)
	if err != nil {
		zaps.Error("db exec failed: err")
		return err
	}

	zaps.Info("<<< update user role done")

	return nil
}

func DelUser(id int64) {

	zaps.Info(">>> del user")

	_, err := db.Exec("DELETE FROM user WHERE id = ?", id)
	if err != nil {
		zaps.Error("delete user exec failed: ", err)
	}

	zaps.Info("<<< del user done")
}

func GetUserTotal() (int, error) {

	var count int

	zaps.Info(">>> get user total")

	err := db.QueryRow("SELECT count(*) FROM user").Scan(&count)
	if err != nil {
		zaps.Error("db query count failed: ", err)
		return 0, err
	}

	zaps.Info("<<< get user total done")

	return count, nil
}

func GetUserList(all int, page int) ([]common.User, error) {

	var userList []common.User
	var rows *sql.Rows
	var err error

	zaps.Info(">>> get user list with page ", page)

	if all == 1 {
		rows, err = db.Query("SELECT id, username, password, role, IFNULL(view_perm, ''), IFNULL(org_perm, ''), create_time " +
			"FROM user")
	} else {
		rows, err = db.Query("SELECT id, username, password, role, IFNULL(view_perm, ''), IFNULL(org_perm, ''), create_time "+
			"FROM user LIMIT ?,10", page*10)
	}

	if err != nil {
		zaps.Error("db query failed: ", err)
		return userList, err
	}

	defer rows.Close()

	for rows.Next() {
		var u common.User
		err := rows.Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.ViewPerm, &u.OrgPerm, &u.CreateTime)
		if err != nil {
			zaps.Error("query error: ", err)
		} else {
			zaps.Debug(">>> id: ", u.ID)
			zaps.Debug(">>> username: ", u.Username)
			zaps.Debug(">>> passwrod: ", u.Password)
			zaps.Debug(">>> create time: ", u.CreateTime)

			userList = append(userList, u)
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return userList, err
	}

	zaps.Info("<<< get user list done")

	return userList, err
}

func GetUserViewPerm(username string) string {
	var perm sql.NullString
	zaps.Info(">>> get user view perm")
	if err := db.QueryRow("SELECT IFNULL(view_perm, '') FROM user WHERE username = ?",
		username).Scan(&perm); err != nil {
		zaps.Error("failed to query user view perm:", err)
	}
	zaps.Info("<<< get user view perm done")
	return perm.String
}

func GetUserOrgPerm(username string) string {
	var perm sql.NullString
	zaps.Info(">>> get user org perm")
	if err := db.QueryRow("SELECT IFNULL(org_perm, '') FROM user WHERE username = ?",
		username).Scan(&perm); err != nil {
		zaps.Error("failed to query user org perm:", err)
	}
	zaps.Info("<<< get user org perm done")
	return perm.String
}

func UpdateUserViewPerm(username, perm string) error {
	zaps.Info(">>> update user view perm")
	rows, err := db.Query("UPDATE user SET view_perm = ? WHERE username = ?", perm, username)
	if err != nil {
		zaps.Error("failed to update user view perm")
		return err
	}

	rows.Close()
	zaps.Info("<<< update user view perm done")
	return nil
}

func UpdateUserOrgPerm(username, perm string) error {
	zaps.Info(">>> update user org perm")
	rows, err := db.Query("UPDATE user SET org_perm = ? WHERE username = ?", perm, username)
	if err != nil {
		zaps.Error("failed to update user org perm")
		return err
	}

	rows.Close()
	zaps.Info("<<< update user view org done")
	return nil
}

func GetUserByUsername(username string) (common.User, bool, error) {

	var resUser common.User
	findFlag := false

	zaps.Info(">>> get one user info: ", username)

	rows, err := db.Query("SELECT id, username, password, role, IFNULL(view_perm, ''), IFNULL(org_perm, ''), create_time "+
		"FROM user WHERE username = ?", username)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return resUser, findFlag, err
	}

	defer rows.Close()

	for rows.Next() {
		var u common.User
		err := rows.Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.ViewPerm, &u.OrgPerm, &u.CreateTime)
		if err != nil {
			zaps.Error("query error: ", err)
			return resUser, findFlag, err
		} else {
			zaps.Debug(">>> id: ", u.ID)
			zaps.Debug(">>> username: ", u.Username)
			zaps.Debug(">>> password: ", u.Password)
			zaps.Debug(">>> create time: ", u.CreateTime)
			resUser = u
			findFlag = true
			break
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return resUser, findFlag, err
	}

	zaps.Info("<<< get one user info done")

	return resUser, findFlag, err
}
