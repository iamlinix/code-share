package common

import (
	"database/sql"
)

/* for erp vendors module */

///////////////////////////////////////////////////////////////////////////////
// BIBLE CONFIG
type BibleCfg struct {
	Material      string  `json:"material"`
	MaterialTxt   string  `json:"materialTxt"`
	SubClass      string  `json:"subClass"`
	SubClassTxt   string  `json:"subClassTxt"`
	MainClass     string  `json:"mainClass"`
	MainClassTxt  string  `json:"mainClassTxt"`
	ValidSDate    string  `json:"validSDate"`
	Groes         string  `json:"groes"`
	BaseUom       string  `json:"baseUom"`
	PurchasePrice float64 `json:"purchasePrice"`
	VendorCode    string  `json:"vendorCode"`
	VendorName    string  `json:"vendorName"`
	SalePrice     float64 `json:"salePrice"`
	Status        int     `json:"status"`
	Inventory     float64 `json:"stock"` //XXX: inventory
	LatestDate    string  `json:"latestDate"`
	CreatedTime   string  `json:"createdTime"`
}

type BibleCfgListResp struct {
	Code  int        `json:"code"`
	Total int        `json:"total"`
	Cfgs  []BibleCfg `json:"biList"` //XXX: cfgs
}

type BibleCfgDetailResp struct {
	Code int      `json:"code"`
	Cfg  BibleCfg `json:"item"` //XXX: cfg
}

///////////////////////////////////////////////////////////////////////////////
//

/* vendor config info struct */
type VendorCfg struct {
	VendorCode string `json:"vendorCode"`
	VendorName string `json:"vendorName"`

	//MonthUnpaid	float64	`json:"monthUnpaid"`
	//MonthSurplus	float64	`json:"monthSurplus"`

	YearPurchase float64 `json:"yearPurchase"`
	YearSales    float64 `json:"yearSales"`
	YearPaid     float64 `json:"yearPaid"`
	YearUnpaid   float64 `json:"yearUnpaid"`
	YearSurplus  float64 `json:"yearSurplus"`

	Rebate      float64 `json:"rebate"`
	Inactive    bool    `json:"inactive"`
	IsNew       bool    `json:"isNew"`
	CreatedTime string  `json:"createdTime"`
}

/* vendor config add response data */
type VendorCfgAddResp struct {
	Code int `json:"code"`
}

/* vendor config list response data */
type VendorCfgListResp struct {
	Code  int         `json:"code"`
	Total int         `json:"total"`
	Cfgs  []VendorCfg `json:"cfgs"`
}

/* vendor configl detail response data */
type VendorCfgDetailResp struct {
	Code int       `json:"code"`
	Cfg  VendorCfg `json:"cfg"`
}

type ExMatlCfg struct {
	Material    string  `json:"material"`
	MaterialTxt string  `json:"materialTxt"`
	Vendor      string  `json:"vendor"`
	VendorName  string  `json:"vendorName"`
	Rebate      float64 `json:"rebate"`
}

type ExMatlCfgListResp struct {
	Code  int         `json:"code"`
	Total int         `json:"total"`
	Cfgs  []ExMatlCfg `json:"cfgs"`
}

///////////////////////////////////////////////////////////////////////////////
// ERP VENDOR BASIC STRUCT

/* erp materials purchase info */
type ERPMatlPurchaseInfo struct {
	Material    string         `json:"material"`
	MaterialTxt sql.NullString `json:"materialTxt"`
	Vendor      sql.NullString `json:"vendor"`
	VendorName  sql.NullString `json:"vendorName"`

	BuyinCnt   float64 `json:"buyinCnt"`
	Buyin      float64 `json:"buyin"`
	BuyinWtax  float64 `json:"buyinWtax"`
	ReturnCnt  float64 `json:"returnCnt"`
	Return     float64 `json:"return"`
	ReturnWtax float64 `json:"returnWtax"`

	/* TODO: a new struct */
	PconfQty int     `json:"pconfQty"`
	PconfZje float64 `json:"pconfZje"`
	Zpurps   int     `json:"zpurps"`
	Zzje     float64 `json:"zzje"`
}

type ERPMatlPurchaseResp struct {
	Code   int                   `json:"code"`
	Total  int                   `json:"total"`
	MPList []ERPMatlPurchaseInfo `json:"mpList"`
}

type ERPMatlZifpurdInfo struct {
	Material    string         `json:"material"`
	MaterialTxt sql.NullString `json:"materialTxt"`
	Vendor      sql.NullString `json:"vendor"`
	VendorName  sql.NullString `json:"vendorName"`

	PconfQty int `json:"pconfQty"`
	Zpurps   int `json:"zpurps"`
}

type ERPMatlZifpurdResp struct {
	Code   int                  `json:"code"`
	Total  int                  `json:"total"`
	MZList []ERPMatlZifpurdInfo `json:"mzList"`
}

type ERPMatlVendorInfo struct {
	Material    string         `json:"material"`
	MaterialTxt sql.NullString `json:"materialTxt"`
	Vendor      sql.NullString `json:"vendor"`
	VendorName  sql.NullString `json:"vendorName"`
}

type ERPMatlVendorResp struct {
	Code   int                 `json:"code"`
	Total  int                 `json:"total"`
	MVList []ERPMatlVendorInfo `json:"mvList"`
}

/* erp material sales info */
type ERPMatlSalesInfo struct {
	Material    string         `json:"material"`
	MaterialTxt sql.NullString `json:"materialTxt"`
	Vendor      sql.NullString `json:"vendor"`
	VendorName  sql.NullString `json:"vendorName"`

	InvQty    int     `json:"invQty"`
	Cost      float64 `json:"cost"`
	CostWtax  float64 `json:"costWtax"`
	NetvalInv float64 `json:"netvalInv"`
	GrossVal  float64 `json:"grossVal"`
}

type ERPMatlSalesResp struct {
	Code   int                `json:"code"`
	Total  int                `json:"total"`
	MSList []ERPMatlSalesInfo `json:"msList"`
}

/* erp materials inventory info */
type ERPInvSnapshot struct {
	PstngDate string
	Material  string
	Zinvsl    float64
	ZinvCost  float64
}

type ERPMatlInvInfo struct {
	Material    string         `json:"material"`
	MaterialTxt sql.NullString `json:"materialTxt"`
	Vendor      sql.NullString `json:"vendor"`
	VendorName  sql.NullString `json:"vendorName"`

	Zinvsl   float64 `json:"zinvsl"`
	ZinvCost float64 `json:"zinvCost"`
}

type ERPMatlInvResp struct {
	Code   int              `json:"code"`
	Total  int              `json:"total"`
	MSList []ERPMatlInvInfo `json:"msList"`
}


type ERPVendorAccountStatement struct {
	VendorName string
	PstngDate  string
	BuyCount   float64
	BuyMoney   float64
	OiEbeln    string
	PeriodNum  string
}


/* erp materials info */
type ERPMaterialInfo struct {
	Material    string `json:"material"`
	MaterialTxt string `json:"materialTxt"`
	Vendor      string `json:"vendor"`
	VendorName  string `json:"vendorName"`
	VendorFlag  int    `json:"vendorFlag"` /*1 from bible, 2 from zifpurd*/

	/* basic metrics */
	/* purchase metrics */
	Buyin      float64 `json:"buyin"`
	BuyinWtax  float64 `json:"buyinWtax"`
	Return     float64 `json:"return"`
	ReturnWtax float64 `json:"returnWtax"`
	Rebate     float64 `json:"rebate"`

	/* zifpurd metrics */
	PconfQty int `json:"pconfQty"`
	Zpurps   int `json:"zpurps"`

	/* sales metrics */
	InvQty    int     `json:"invQty"`
	Cost      float64 `json:"cost"`
	CostWtax  float64 `json:"costWtax"`
	NetvalInv float64 `json:"netvalInv"`
	GrossVal  float64 `json:"grossVal"`

	/* inventory metrics */
	OpenZinvsl   float64 `json:"openZinvsl"`
	OpenZinvCost float64 `json:"openZinvCost"`

	CloseZinvsl   float64 `json:"closeZinvsl"`
	CloseZinvCost float64 `json:"closeZinvCost"`
}

type ERPMaterialResp struct {
	Code         int               `json:"code"`
	Total        int               `json:"total"`
	MaterialList []ERPMaterialInfo `json:"materialList"`
}

type ERPMaterialReportListResp struct {
	Code        int               `json:"code"`
	Total       int               `json:"total"`
	MReportList []ERPMaterialInfo `json:"materialReportList"`
}

/* erp vendor info */
type ERPVendorInfo struct {
	Vendor     string `json:"vendor"`
	VendorName string `json:"vendorName"`

	PurVal      float64 `json:"purVal"`
	PurWtaxVal  float64 `json:"purWtaxVal"`
	RebateVal   float64 `json:"rebateVal"`
	PurFinalVal float64 `json:"purFinalVal"`

	PconfQty int `json:"pconfQty"`
	Zpurps   int `json:"zpurps"`

	InvQty    int     `json:"invQty"`
	Cost      float64 `json:"cost"`
	CostWtax  float64 `json:"costWtax"`
	NetvalInv float64 `json:"netvalInv"`
	GrossVal  float64 `json:"grossVal"`

	OpenZinvsl   float64 `json:"openZinvsl"`
	OpenZinvCost float64 `json:"openZinvCost"`

	CloseZinvsl   float64 `json:"closeZinvsl"`
	CloseZinvCost float64 `json:"closeZinvCost"`
}

type ERPVendorResp struct {
	Code       int             `json:"code"`
	Total      int             `json:"total"`
	VendorList []ERPVendorInfo `json:"vendorList"`
}

type ERPVendorFPDate struct {
	Vendor  string
	ZgrDate string
}

///////////////////////////////////////////////////////////////////////////////
// vendors report
type VendorReportHistoryResp struct {
	Code    int                  `json:"code"`
	History []VendorReportUpdate `json:"history"`
}

type VendorReportUpdate struct {
	Type     int     `json:"type"`
	Month    string  `json:"month"`
	Vendor   string  `json:"vendor"`
	User     string  `json:"user"`
	OldValue float64 `json:"oldValue"`
	NewValue float64 `json:"newValue"`
	Reason   string  `json:"reason"`
	MTime    string  `json:"mtime"`
}

type VendorReport struct {
	Month     string `json:"month"`
	BeginDate string `json:"beginDate"`
	EndDate   string `json:"endDate"`

	VendorCode string `json:"vendorCode"`
	VendorName string `json:"vendorName"`

	MonthPurchase     float64 `json:"monthPurchase"`
	MonthPurchaseOrig float64 `json:"monthPurchaseOrig"`
	MonthReceipt      float64 `json:"monthReceipt"`
	MonthRebate       float64 `json:"monthRebate"`
	MonthPurFinal     float64 `json:"monthPurFinal"`
	MonthPurchase2    float64 `json:"monthPurchase2"` //before tax
	MonthSalesQty     int     `json:"monthSalesQty"`
	MonthSales        float64 `json:"monthSales"`
	MonthSalesOrig    float64 `json:"monthSalesOrig"`
	MonthSales2       float64 `json:"monthSales2"` //before tax
	MonthPaidOrig     float64 `json:"monthPaidOrig"`
	MonthPaid         float64 `json:"monthPaid"`
	LastMonthUnpaid   float64 `json:"lastMonthUnpaid"`
	MonthUnpaid       float64 `json:"monthUnpaid"`
	MonthUnpaidOrig   float64 `json:"monthUnpaidOrig"`
	MonthSurplus      float64 `json:"monthSurplus"`

	YearPurchase   float64 `json:"yearPurchase"`
	YearSales      float64 `json:"yearSales"`
	YearSalesQty   int     `json:"yearSalesQty"`
	YearPaid       float64 `json:"yearPaid"`
	YearUnpaid     float64 `json:"yearUnpaid"`
	YearUnpaidOrig float64 `json:"yearUnpaidOrig"`
	YearSurplus    float64 `json:"yearSurplus"`

	YearOpenInv   float64 `json:"yearOpenStock"` //XXX inv
	MonthOpenInv  float64 `json:"monthOpenStock"`
	MonthCloseInv float64 `json:"monthCloseStock"`
	MonthAvgInv   float64 `json:"monthAvgStock"`
	YearAvgInv    float64 `json:"yearAvgStock"`
	MonthDaysTo   float64 `json:"monthDaysTo"`
	YearDaysTo    float64 `json:"yearDayTo"`
	CommEff       float64 `json:"commEff"`
	GrossMargin   float64 `json:"grossMargin"`
	OrderFillRate float64 `json:"orderFillRate"`
	FirstOrder    string  `json:"firstOrder"`

	CreatedTime string `json:"createdTime"`
}

type VendorReportAddResp struct {
	Code       int            `json:"code"`
	Total      int            `json:"total"`
	VendorSum  VendorReport   `json:"sum"`
	VendorList []VendorReport `json:"vendorsList"`
}

type VendorReportListResp struct {
	Code        int            `json:"code"`
	Total       int            `json:"total"`
	VReportList []VendorReport `json:"vendorsReportList"`
}

type VendorReportDetailResp struct {
	Code       int            `json:"code"`
	VendorList []VendorReport `json:"vendorsList"`
}

type VendorMatlCount struct {
	Vendor    string
	MatlCount int
}

///////////////////////////////////////////////////////////////////////////////
// VENDORS RATING

type VendorRatingCfg struct {
	ID     int64   `json:"id"`
	Name   string  `json:"name"`
	Weight float64 `json:"weight"`
}

type VendorRatingCfgUpdate struct {
	Cfgs []VendorRatingCfg `json:"cfgs"`
}

type VendorRatingCfgListResp struct {
	Code  int               `json:"code"`
	Total int               `json:"total"`
	Cfgs  []VendorRatingCfg `json:"cfgs"`
}

type VendorMeta struct {
	Vendor     string `json:"vendor"`
	VendorName string `json:"vendorName"`

	PconfQty int `json:"pconfQty"`
	Zpurps   int `json:"zpurps"`

	InvQty    int     `json:"invQty"`
	Cost      float64 `json:"cost"`
	NetvalInv float64 `json:"netvalInv"`
	GrossVal  float64 `json:"grossVal"`
}

type VendorRatingMetric struct {
	NetSales       float64 `json:"netSales"`
	NetSalesPct    float64 `json:"netSalesPct"`
	NetSalesScore  float64 `json:"netSalesScore"`
	NetSalesScoreW float64 `json:"netSalesScoreW"`

	GrossProfit       float64 `json:"grossProfit"`
	GrossProfitPct    float64 `json:"grossProfitPct"`
	GrossProfitScore  float64 `json:"grossProfitScore"`
	GrossProfitScoreW float64 `json:"grossProfitScoreW"`

	GrossMargin       float64 `json:"grossMargin"`
	GrossMarginScore  float64 `json:"grossMarginScore"`
	GrossMarginScoreW float64 `json:"grossMarginScoreW"`

	OrderFillRate       float64 `json:"orderFillRate"`
	OrderFillRateScore  float64 `json:"orderFillRateScore"`
	OrderFillRateScoreW float64 `json:"orderFillRateScoreW"`

	SubjectiveScore  float64 `json:"subjectiveScore"`
	SubjectiveScoreW float64 `json:"subjectiveScoreW"`

	TotalScore float64 `json:"totalScore"`
}

type VendorRating struct {
	Vendor     string `json:"vendor"`
	VendorName string `json:"vendorName"`

	Metric VendorRatingMetric `json:"metric"`
}

type VendorRatingResp struct {
	Code       int            `json:"code"`
	Total      int            `json:"total"`
	VendorList []VendorRating `json:"vendorList"`
}

type VendorRatingReportHdr struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	BeginDate   string `json:"beginDate"`
	EndDate     string `json:"endDate"`
	CreatedTime string `json:"createdTime"`
}

type VendorRatingReport struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	BeginDate   string `json:"beginDate"`
	EndDate     string `json:"endDate"`
	CreatedTime string `json:"createdTime"`

	VendorList []VendorRating `json:"vendorList"`
}

type VendorRatingReportListResp struct {
	Code  int `json:"code"`
	Total int `json:"total"`

	ReportList []VendorRatingReportHdr `json:"reportList"`
}

type VendorRatingReportDetailResp struct {
	Code   int                `json:"code"`
	Report VendorRatingReport `json:"report"`
}
