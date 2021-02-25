package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"
)

func GetSetting(name string) (string, error) {
	zaps.Info(">>> get setting: ", name)
	var value string = ""
	row := db.QueryRow("SELECT value FROM settings WHERE name = ?", name)
	if err := row.Scan(&value); err != nil {
		zaps.Error(">>> scan setting error:", err)
		return value, err
	}
	return value, nil
}

func GetSettingDefault(name, defaultValue string) string {
	zaps.Info(">>> get setting: ", name)
	var value string
	row := db.QueryRow("SELECT value FROM settings WHERE name = ?", name)
	if err := row.Scan(&value); err != nil {
		zaps.Error(">>> scan setting error:", err)
		value = defaultValue
	}
	return value
}

func UpdateSetting(name, value string) error {
	zaps.Info(">>> set setting: ", name, value)
	rows, err := db.Query("INSERT INTO settings (name, value) VALUES (?, ?) ON DUPLICATE KEY UPDATE value = ?", name, value, value)
	if err != nil {
		zaps.Error(">>> set setting error:", err)
	}
	rows.Close()

	return err
}

func GetSettingListByKeyRange(start, end string) ([]common.GeneralSetting, error) {
	zaps.Info(">>> get setting list by key range: ", start, end)
	rows, err := db.Query("SELECT name, value FROM settings WHERE name BETWEEN ? AND ? ORDER BY name", start, end)
	if err != nil {
		zaps.Error(">>> get setting list error:", err)
		return nil, err
	}
	defer rows.Close()

	var settings []common.GeneralSetting
	for rows.Next() {
		var setting common.GeneralSetting
		rows.Scan(&setting.Name, &setting.Value)
		settings = append(settings, setting)
	}

	return settings, nil
}
