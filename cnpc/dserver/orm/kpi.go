package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

func GetKpiByCode(code string) (common.KPI, error) {
	var kpi common.KPI

	zaps.Info(">>> get kpi by code: ", code)

	rows, err := db.Query("SELECT org_code, org_name, month, income_type, value "+
		"FROM kpi WHERE org_code = ? ORDER BY org_code", code)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return kpi, err
	}

	defer rows.Close()

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return kpi, err
	}

	for rows.Next() {
		err := rows.Scan(&kpi.OrgCode, &kpi.OrgCode, &kpi.Month, &kpi.IncomeType, &kpi.Value)
		if err != nil {
			zaps.Error("query error: ", err)
			return kpi, err
		}
	}

	zaps.Info("<<< get kpi by code done")

	return kpi, err
}

func GetQuarterRange(year, quarter int) (string, string, error) {
	if quarter < 0 || quarter > 3 {
		return "", "", errors.New("Invalid Quarter")
	}

	dates, err := GetSettingListByKeyRange("fymon00", "fymon11")
	if err != nil {
		return "", "", err
	}

	year1 := year
	year2 := year
	d1 := dates[quarter*3]
	p := strings.Split(d1.Value, "~")
	m1 := p[0]
	if strings.Compare(p[0], p[1]) > 0 {
		year1--
	}

	d2 := dates[quarter*3+2]
	p = strings.Split(d2.Value, "~")
	m2 := p[1]
	if strings.Compare(p[0], p[1]) > 0 {
		year2++
	}

	start := fmt.Sprintf("%d-%s", year1, m1)
	end := fmt.Sprintf("%d-%s", year2, m2)
	zaps.Info(">>> GetQuarterRange:", year, quarter, start, end)
	return start, end, nil
}

func FindDatePeriod(date string) (string, string, int, int) {
	zaps.Info(">>> FindDatePeriod:", date)
	dates, err := GetSettingListByKeyRange("fymon00", "fymon11")
	if err != nil {
		return "", "", -1, -1
	}

	total := len(dates)
	year, err := strconv.Atoi(date[:4])
	if err != nil {
		zaps.Error(">>> failed to parse date:", date)
		return "", "", -1, -1
	}

	for i := 0; i < total; i++ {
		d := dates[i]
		p := strings.Split(d.Value, "~")
		m1 := fmt.Sprintf("%d-%s", year, p[0])
		m2 := fmt.Sprintf("%d-%s", year, p[1])
		if strings.Compare(p[0], p[1]) > 0 {
			if i < 2 {
				m1 = fmt.Sprintf("%d-%s", year-1, p[0])
			} else {
				m1 = fmt.Sprintf("%d-%s", year+1, p[1])
			}
		}

		if strings.Compare(date, m1) >= 0 && strings.Compare(date, m2) <= 0 {
			mNum := i + 1
			quarter := int(i / 3)
			zaps.Info(">>> FindDatePeriod: ", date, m1, m2, mNum, quarter)
			return m1, m2, mNum, quarter
		}
	}

	zaps.Error(">>> cannot find period for date:", date)
	return "", "", -1, -1
}

func QuarterOfMonth(month time.Month) (int, int, int) {
	if month >= time.January && month <= time.March {
		return 13, 1, 3
	}

	if month >= time.April && month <= time.June {
		return 14, 4, 6
	}

	if month >= time.July && month <= time.September {
		return 15, 7, 9
	}

	if month >= time.October && month <= time.December {
		return 16, 10, 12
	}

	return 0, 0, 0
}

func GetKpiByMonthV2(startDate, endDate, orgCode string, incomeType int) (common.KPIMonthResp, error) {
	var resp common.KPIMonthResp
	var monthValue float64 = 0.0
	var quarterValue float64 = 0.0
	var yearValue float64 = 0.0
	layout := "2006-01-02"

	s, err := time.Parse(layout, startDate)
	if err != nil {
		zaps.Error("failed to parse start date: ", startDate, err)
		return resp, err
	}

	e, err := time.Parse(layout, endDate)
	if err != nil {
		zaps.Error("failed to parse end date: ", endDate, err)
		return resp, err
	}

	var val sql.NullInt64
	// s1 := s
	// sy := s1.Year()
	// ey := e.Year()
	// sm := s1.Month()
	// em := e.Month()
	// ms := s1.AddDate(0, 0, -s1.Day()+1)
	// me := ms.AddDate(0, 1, 0)
	// monthDays := me.Sub(ms).Hours() / 24
	// days := me.Sub(s1).Hours() / 24

	// get month values
	ms, _ := time.Parse(layout, fmt.Sprintf("%d-%02d-01", s.Year(), int(s.Month())))
	for ms.Before(e) {
		sm := ms.Month()
		me := ms.AddDate(0, 1, 0)
		monthDays := me.Sub(ms).Hours() / 24

		if ms.Before(s) {
			ms = s
		}
		if me.After(e) {
			me = e
		}
		days := me.Sub(ms).Hours() / 24

		if err = db.QueryRow("SELECT value FROM kpi WHERE org_code = ? AND month = ? AND income_type = ?",
			orgCode, int(sm), incomeType).Scan(&val); err != nil {
			zaps.Error("failed to get month kpi value: ", err)
		} else {
			monthValue += float64(val.Int64) * days / monthDays
		}
		zaps.Infof("by month get month value: %f, %d, %d", monthValue, days, monthDays)

		ms = me
	}

	// for true {
	// 	if err = db.QueryRow("SELECT value FROM kpi WHERE org_code = ? AND month = ? AND income_type = ?",
	// 		orgCode, int(sm), incomeType).Scan(&val); err != nil {
	// 		zaps.Error("failed to get month kpi value: ", err)
	// 	} else {
	// 		monthValue += float64(val.Int64) * days / monthDays
	// 	}
	// 	zaps.Infof("by month get month value: %f, %d, %d", monthValue, days, monthDays)

	// 	if sm == em && sy == ey {
	// 		break
	// 	}

	// 	ms = me
	// 	s1 = me
	// 	sy = s1.Year()
	// 	sm = s1.Month()
	// 	me = me.AddDate(0, 1, 0)
	// 	if sm == em && sy == ey {
	// 		// s1.AddDate(0, 0, e.Day()-1)
	// 		me = e
	// 	}
	// 	monthDays = me.Sub(ms).Hours() / 24
	// 	days = me.Sub(s1).Hours() / 24
	// }
	resp.MonthValue = monthValue

	// get quarter value
	s1 := s
	sy := s1.Year()
	sm := s1.Month()
	jobDone := false
	for true {
		quarter, qs, _ := QuarterOfMonth(sm)
		sq, _ := time.Parse(layout, fmt.Sprintf("%d-%02d-01", sy, qs))
		eq := sq.AddDate(0, 3, -1)
		ee := eq
		quarterDays := (eq.Sub(sq).Hours() / 24) + 1
		if sq.Before(s) {
			sq = s
		}

		if eq.After(e) {
			jobDone = true
			ee = e
		}
		days := (ee.Sub(sq).Hours() / 24) + 1

		if err = db.QueryRow("SELECT value FROM kpi WHERE org_code = ? AND month = ? AND income_type = ?",
			orgCode, quarter, incomeType).Scan(&val); err != nil {
			zaps.Error("failed to get quarter kpi value: ", err)
		} else {
			quarterValue += float64(val.Int64) * days / quarterDays
		}
		zaps.Infof("by month get quarter value: %f, %d, %d", quarterValue, days, quarterDays)

		if jobDone {
			break
		}

		s1 = s1.AddDate(0, 3, 0)
		sy = s1.Year()
		sm = s1.Month()
	}
	resp.QuarterValue = quarterValue

	// get year value
	year := e.Year()
	ys, _ := time.Parse(layout, fmt.Sprintf("%d-01-01", year))
	ye := ys.AddDate(1, 0, 0)
	yearDays := ye.Sub(ys).Hours() / 24
	days := (e.Sub(s).Hours() / 24) + 1
	if err = db.QueryRow("SELECT value FROM kpi WHERE org_code = ? AND month = 0 AND income_type = ?",
		orgCode, incomeType).Scan(&val); err != nil {
		zaps.Error("failed to get year kpi value: ", err)
	} else {
		yearValue += float64(val.Int64) * days / yearDays
	}
	zaps.Infof("by month get year value: %f, %d, %d", yearValue, days, yearDays)
	resp.YearValue = yearValue

	return resp, nil
}

func GetKpiByMonth(code, month string, incomeType int) (common.KPIMonthResp, error) {
	var kpi common.KPIMonthResp
	kpi.MonthValue = 0
	kpi.YearValue = 0
	zaps.Info(">>> get kpi by month: ", code, month)

	_, _, mon, quarter := FindDatePeriod(month)
	if mon == -1 {
		zaps.Error(">>> failed to find period for month:", month)
		return kpi, errors.New("Invalid Month")
	}

	rowMonth, err := db.Query("SELECT value from kpi WHERE org_code = ? AND month = ? AND income_type = ?",
		code, mon, incomeType)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return kpi, err
	}
	defer rowMonth.Close()
	err = rowMonth.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return kpi, err
	}
	if rowMonth.Next() {
		rowMonth.Scan(&kpi.MonthValue)
	}

	rowQuarter, err := db.Query("SELECT value from kpi WHERE org_code = ? AND month = ? AND income_type = ?",
		code, quarter+13, incomeType)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return kpi, err
	}
	defer rowQuarter.Close()
	err = rowQuarter.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return kpi, err
	}
	if rowQuarter.Next() {
		rowQuarter.Scan(&kpi.QuarterValue)
	}

	rowYear, err := db.Query("SELECT value from kpi WHERE org_code = ? AND month = 0 AND income_type = ?",
		code, incomeType)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return kpi, err
	}
	defer rowYear.Close()
	err = rowYear.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return kpi, err
	}
	if rowYear.Next() {
		rowYear.Scan(&kpi.YearValue)
	}

	zaps.Info("<<< get kpi by month done")
	return kpi, nil
}

func GetKpis() ([]common.KPI, error) {
	var kpis []common.KPI = nil

	zaps.Info(">>> get kpis")

	rows, err := db.Query("SELECT org_code, org_name, month, income_type, value " +
		"FROM kpi ORDER BY org_code, month")
	if err != nil {
		zaps.Error("db query failed: ", err)
		return kpis, err
	}

	defer rows.Close()

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return kpis, err
	}

	for rows.Next() {
		var kpi common.KPI
		err := rows.Scan(&kpi.OrgCode, &kpi.OrgName, &kpi.Month, &kpi.IncomeType, &kpi.Value)
		if err != nil {
			zaps.Error("query error: ", err)
			return nil, err
		}
		kpis = append(kpis, kpi)
	}

	zaps.Info("<<< get kpis done")

	return kpis, err
}

func UpdateKpi(kpi *common.KPI) bool {
	zaps.Info(">>> update kpi")

	if rows, err := db.Query("INSERT INTO kpi (org_code, org_name, month, income_type, value) "+
		"VALUES (? ,?, ?, ?, ?) ON DUPLICATE KEY UPDATE value = ?",
		kpi.OrgCode, kpi.OrgName, kpi.Month, kpi.IncomeType, kpi.Value, kpi.Value); err != nil {
		zaps.Error("query error: ", err)
		return false
	} else {
		rows.Close()
	}

	return true
}

func DeleteKpi(kpi *common.KPI) bool {
	zaps.Info(">>> delete kpi")

	if rows, err := db.Query("DELETE FROM kpi WHERE org_code = ? AND month = ? AND income_type = ?",
		kpi.OrgCode, kpi.Month, kpi.IncomeType); err != nil {
		zaps.Error("query error: ", err)
		return false
	} else {
		rows.Close()
	}

	return true
}

func UpdateFinanceYearStart(start string) bool {
	zaps.Info(">>> set finace year start", start)
	if err := UpdateSetting("fystart", start); err != nil {
		zaps.Error("query error: ", err)
		return false
	}
	return true
}

func GetFinanceYearStart() string {
	zaps.Info(">>> get finace year start")
	start, err := GetSetting("fystart")
	if err != nil {
		zaps.Error("query error: ", err)
	}
	return start
}

func UpdateFyMonth(month, span string) bool {
	zaps.Info(">>> update fy month:", month, span)
	if err := UpdateSetting(month, span); err != nil {
		zaps.Error("query error: ", err)
		return false
	}
	return true
}

func GetFyMonthList() []common.GeneralSetting {
	zaps.Info(">>> get fy month list")
	if settings, err := GetSettingListByKeyRange("fymon00", "fymon11"); err != nil {
		zaps.Error("query error: ", err)
		return nil
	} else {
		return settings
	}
}

func GetPaymentTypes(beginDate, endDate, orgCode string) ([]common.PaymentType, error) {
	zaps.Info(">>> get payment types")
	var payments []common.PaymentType
	if orgCode == "000" {
		zaps.Info("get payments for company")
		rows, err := db.Query("SELECT zpay_type, SUM(amount), COUNT(*) FROM zpay_d03h_orc WHERE calday "+
			"BETWEEN ? AND ? AND plant IN (SELECT bic_zaplant FROM zaplant_xy WHERE	bic_zrpa_lcit "+
			"BETWEEN 'A13' AND 'A16') AND material NOT IN (SELECT material FROM material WHERE bic_zklad2 "+
			"IN ('5002', '7001', '1001', '1002', '7002', '8001')) GROUP BY zpay_type", beginDate, endDate)
		if err != nil {
			zaps.Error("failed to get payments: ", err)
			return nil, err
		}

		defer rows.Close()
		for rows.Next() {
			var pt common.PaymentType
			if err = rows.Scan(&pt.ID, &pt.Money, &pt.Count); err != nil {
				zaps.Error("failed to scan payments: ", err)
			} else {
				payments = append(payments, pt)
			}
		}
	} else if len(orgCode) == 3 {
		zaps.Info("get payments for branch")
		rows, err := db.Query("SELECT zpay_type, SUM(amount), COUNT(*) FROM zpay_d03h_orc WHERE calday "+
			"BETWEEN ? AND ? AND plant IN (SELECT bic_zaplant FROM zaplant_xy WHERE	bic_zrpa_lcit = ?) "+
			"AND material NOT IN (SELECT material FROM material WHERE bic_zklad2 IN ('5002', '7001', "+
			"'1001', '1002', '7002', '8001')) "+
			"GROUP BY zpay_type", beginDate, endDate, orgCode)
		if err != nil {
			zaps.Error("failed to get payments: ", err)
			return nil, err
		}

		defer rows.Close()
		for rows.Next() {
			var pt common.PaymentType
			if err = rows.Scan(&pt.ID, &pt.Money, &pt.Count); err != nil {
				zaps.Error("failed to scan payments: ", err)
			} else {
				payments = append(payments, pt)
			}
		}
	} else {
		zaps.Info("get payments for plant")
		rows, err := db.Query("SELECT zpay_type, SUM(amount), COUNT(*) FROM zpay_d03h_orc WHERE calday "+
			"BETWEEN ? AND ? AND plant = ? AND material NOT IN (SELECT material FROM material WHERE "+
			"bic_zklad2 IN ('5002', '7001', '1001', '1002', '7002', '8001')) "+
			"GROUP BY zpay_type", beginDate, endDate, orgCode)
		if err != nil {
			zaps.Error("failed to get payments: ", err)
			return nil, err
		}

		defer rows.Close()
		for rows.Next() {
			var pt common.PaymentType
			if err = rows.Scan(&pt.ID, &pt.Money, &pt.Count); err != nil {
				zaps.Error("failed to scan payments: ", err)
			} else {
				payments = append(payments, pt)
			}
		}
	}

	return payments, nil
}
