package orm

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/xid"

	"database/sql"
	"time"
)

func AddVendorRatingReportTx(req common.VendorRatingReport) error {

	var wSUB float64

	t1 := time.Now().UnixNano() / 1e6

	cfgs, _, err := GetVendorRatingCfgList()
	if err != nil {
		zaps.Errorf("get vendor rating cfg list failed: %v", err)
		return err
	}

	for _, c := range cfgs {
		switch c.ID {
		case 5:
			wSUB = c.Weight / 100.0
		default:

		}
	}

	req.UUID = xid.New().String()

	tx, _ := db.Begin()
	defer tx.Rollback()

	for _, v := range req.VendorList {

		/* update total score */
		s := v.Metric.SubjectiveScore * wSUB
		v.Metric.SubjectiveScoreW = s
		v.Metric.TotalScore += s

		AddVendorRatingReport(tx, req, v)
	}

	/* Transaction END */
	err = tx.Commit()
	if err != nil {
		zaps.Errorf("transaction commit failed: %v", err)
		return err
	}

	t2 := time.Now().UnixNano() / 1e6
	zaps.Infof("<<< add vendor rating report tx done using %d ms", t2-t1)

	return nil
}

func AddVendorRatingReport(tx *sql.Tx, req common.VendorRatingReport,
	vr common.VendorRating) error {

	zaps.Info(">>> add vendor rating report")

	stmt, err := tx.Prepare("INSERT INTO vendor_rating_report(uuid, " +
		"name, begin_date, end_date, vendor, " +
		"vendor_name, net_sales, net_sales_pct, " +
		"net_sales_score, net_sales_score_w, " +
		"gross_profit, gross_profit_pct, " +
		"gross_profit_score, gross_profit_score_w, " +
		"gross_margin, gross_margin_score, " +
		"gross_margin_score_w, order_fill_rate, " +
		"order_fill_rate_score, order_fill_rate_score_w, " +
		"subjective_score, subjective_score_w, " +
		"total_socre, created_time) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, " +
		"?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		zaps.Error("insert into vendor rating report failed: ", err)
		return err
	}
	defer stmt.Close()

	m := vr.Metric
	_, err = stmt.Exec(req.UUID, req.Name, req.BeginDate, req.EndDate,
		vr.Vendor, vr.VendorName, m.NetSales, m.NetSalesPct,
		m.NetSalesScore, m.NetSalesScoreW, m.GrossProfit,
		m.GrossProfitPct, m.GrossProfitScore, m.GrossProfitScoreW,
		m.GrossMargin, m.GrossMarginScore,
		m.GrossMarginScoreW, m.OrderFillRate,
		m.OrderFillRateScore, m.OrderFillRateScoreW,
		m.SubjectiveScore, m.SubjectiveScoreW,
		m.TotalScore, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		zaps.Error("db exec failed: ", err)
		return err
	}

	zaps.Info("<<< add vendor rating report done")

	return err
}

func DelVendorRatingReport(uuid string) {

	zaps.Info(">>> del vendor rating report")

	_, err := db.Exec("DELETE FROM vendor_rating_report WHERE uuid = ?", uuid)
	if err != nil {
		zaps.Error("delete vendor rating report exec failed: ", err)
	}

	zaps.Info("<<< del vendor rating report done")
}

func GetVendorRatingReportList() ([]common.VendorRatingReportHdr, int, error) {

	var hdrList []common.VendorRatingReportHdr
	var count int

	zaps.Info(">>> get vendor rating report list")

	rows, err := db.Query("SELECT DISTINCT uuid, name, begin_date, " +
		"end_date FROM vendor_rating_report")
	if err != nil {
		zaps.Error("db query failed: ", err)
		return hdrList, count, err
	}

	defer rows.Close()

	for rows.Next() {
		var hdr common.VendorRatingReportHdr
		err := rows.Scan(&hdr.UUID, &hdr.Name, &hdr.BeginDate, &hdr.EndDate)
		if err != nil {
			zaps.Error("query error: ", err)
			return hdrList, count, err
		} else {
			zaps.Debug(">>> id: ", hdr.UUID)
			zaps.Debug(">>> name: ", hdr.Name)
			zaps.Debug(">>> start date: ", hdr.BeginDate)
			zaps.Debug(">>> end date: ", hdr.EndDate)

			hdrList = append(hdrList, hdr)
			count++
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return hdrList, count, err
	}

	zaps.Info("<<< get vendor rating report list done")

	return hdrList, count, err
}

func GetVendorRatingReportByUUID(uuid string) (common.VendorRatingReport,
	bool, error) {

	var report common.VendorRatingReport
	var vrList []common.VendorRating
	var find bool

	zaps.Info(">>> get vendor rating report by uuid: ", uuid)

	rows, err := db.Query("SELECT uuid, name, begin_date, end_date, "+
		"vendor, vendor_name, net_sales, net_sales_pct, "+
		"net_sales_score, net_sales_score_w, "+
		"gross_profit, gross_profit_pct, "+
		"gross_profit_score, gross_profit_score_w, "+
		"gross_margin, gross_margin_score, "+
		"gross_margin_score_w, order_fill_rate, "+
		"order_fill_rate_score, order_fill_rate_score_w, "+
		"subjective_score, subjective_score_w, "+
		"total_socre, created_time "+
		"FROM vendor_rating_report WHERE uuid = ?", uuid)
	if err != nil {
		zaps.Error("db query failed: ", err)
		return report, find, err
	}

	defer rows.Close()

	for rows.Next() {
		var vr common.VendorRating
		err := rows.Scan(&report.UUID, &report.Name, &report.BeginDate,
			&report.EndDate, &vr.Vendor, &vr.VendorName,
			&vr.Metric.NetSales, &vr.Metric.NetSalesPct,
			&vr.Metric.NetSalesScore, &vr.Metric.NetSalesScoreW,
			&vr.Metric.GrossProfit, &vr.Metric.GrossProfitPct,
			&vr.Metric.GrossProfitScore, &vr.Metric.GrossProfitScoreW,
			&vr.Metric.GrossMargin, &vr.Metric.GrossMarginScore,
			&vr.Metric.GrossMarginScoreW, &vr.Metric.OrderFillRate,
			&vr.Metric.OrderFillRateScore, &vr.Metric.OrderFillRateScoreW,
			&vr.Metric.SubjectiveScore, &vr.Metric.SubjectiveScoreW,
			&vr.Metric.TotalScore, &report.CreatedTime)
		if err != nil {
			zaps.Error("query error: ", err)
			return report, find, err
		} else {
			vrList = append(vrList, vr)
			find = true
		}
	}

	err = rows.Err()
	if err != nil {
		zaps.Error("query row error: ", err)
		return report, find, err
	}

	report.VendorList = vrList

	zaps.Info("<<< get vendors report list done")

	return report, find, err
}
