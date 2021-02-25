package common

//"database/sql"

/* for sales kpi module */

///////////////////////////////////////////////////////////////////////////////
//

/* posdm zof count info struct */
type PosdmMetric struct {
	BicZplant string

	BicZofType  string
	BicZofCount int
}

///////////////////////////////////////////////////////////////////////////////
// SALES METRIC

/* bill zsd metric info struct */
type SalesMetric struct {
	InvQty    int
	Cost      float64
	NetvalInv float64
	GrossVal  float64

	GrossProfit float64
	GrossMargin float64
}

type SalesOrgMetric struct {
	SaleDays  int
	ParentOrg string
	OrgCode   string
	OrgText   string
	PosX      float64
	PosY      float64
	Metric    SalesMetric
}

type SalesClassMetric struct {
	ClassCode string
	ClassText string
	Metric    SalesMetric
}

type SalesClassResp struct {
	Code      int                `json:"code"`
	Total     int                `json:"total"`
	ClassList []SalesClassMetric `json:"classList"`
}

type SalesMatlMetric struct {
	Material    string
	MaterialTxt string
	Metric      SalesMetric
}

type SalesDateMetric struct {
	Date   string
	Metric SalesMetric
}

type SalesDateResp struct {
	Code     int               `json:"code"`
	Total    int               `json:"total"`
	DateList []SalesDateMetric `json:"dateList"`
}

///////////////////////////////////////////////////////////////////////////////
// SALES KPI

type KPIMetric struct {
	Value float64 `json:"value"`
	YoY   float64 `json:"yoy"`
	MoM   float64 `json:"mom"`
}

// KPI STRUCT
type SalesKPI struct {
	NetSales    KPIMetric `json:"netIncome"` //XXX TODO
	GrossProfit KPIMetric `json:"grossProfit"`
	GrossMargin KPIMetric `json:"grossMargin"`

	NonFuleCount     KPIMetric `json:"nonFuelCount"`
	FuleCount        KPIMetric `json:"fuelCount"`
	FNConversionRate KPIMetric `json:"fnConversionRate"`
	AvgTrxValue      KPIMetric `json:"avgTrxValue"`
	Y2nIncome        KPIMetric `json:"y2nIncome"`
	Q2nIncome        KPIMetric `json:"q2nIncome"`
	M2nIncome        KPIMetric `json:"m2nIncome"`
	Y2nProfit        KPIMetric `json:"y2nProfit"`
	Q2nProfit        KPIMetric `json:"q2nProfit"`
	M2nProfit        KPIMetric `json:"m2nProfit"`
}

type SalesKPIResp struct {
	Code int      `json:"code"`
	KPI  SalesKPI `json:"kpi"`
}

type SalesKPIAnalyseResp struct {
	Code       int        `json:"code"`
	CurrentKPI SalesKPI   `json:"currentKpi"`
	PastYears  []SalesKPI `json:"pastYears"`
}

///////////////////////////////////////////////////////////////////////////////
//

type SalesPlantRankResp struct {
	Code      int              `json:"code"`
	Total     int              `json:"total"`
	PlantList []SalesOrgMetric `json:"plantList"`
}

type SalesMaterialRankResp struct {
	Code     int               `json:"code"`
	Total    int               `json:"total"`
	MatlList []SalesMatlMetric `json:"matlList"`
}

///////////////////////////////////////////////////////////////////////////////
// SALES MATERIAL MODULE STRUCT

/* sales material config struct */
type SalesMaterialCfg struct {
	Material    string `json:"material"`
	MaterialTxt string `json:"materialTxt"`
	Level       int    `json:"level"`
	Tag         string `json:"tag"`
	CreatedTime string `json:"createdTime"`
}

type SalesMaterialCfgAddResp struct {
	Code int              `json:"code"`
	Cfg  SalesMaterialCfg `json:"cfg"`
}

/* sales material config list response data */
type SalesMaterialCfgListResp struct {
	Code  int                `json:"code"`
	Total int                `json:"total"`
	Cfgs  []SalesMaterialCfg `json:"cfgs"`
}

///////////////////////////////////////////////////////////////////////////////
//

/* sales material kpi */
type SalesMaterialKPI struct {
	NetSales    KPIMetric `json:"netIncome"` //XXX TODO
	GrossMargin KPIMetric `json:"grossMargin"`
	GrossProfit KPIMetric `json:"grossProfit"`
}

type SalesMaterialKPIResp struct {
	Code int              `json:"code"`
	KPI  SalesMaterialKPI `json:"kpi"`
}

///////////////////////////////////////////////////////////////////////////////
// SALES PLANT MODULE STRUCT

/* sales plant config struct */
type SalesPlantCfg struct {
	Plant       string `json:"plant"`
	PlantName   string `json:"plantName"`
	Tag         string `json:"tag"`
	CreatedTime string `json:"createdTime"`
}

type SalesPlantCfgAddResp struct {
	Code int           `json:"code"`
	Cfg  SalesPlantCfg `json:"cfg"`
}

type SalesPlantCfgListResp struct {
	Code  int             `json:"code"`
	Total int             `json:"total"`
	Cfgs  []SalesPlantCfg `json:"cfgs"`
}

///////////////////////////////////////////////////////////////////////////////
// SALES PLANT Metric

/* sales plant kpi */
type SalesPlantKPI struct {
	NetSales    KPIMetric `json:"netIncome"` //XXX TODO
	GrossMargin KPIMetric `json:"grossMargin"`
	GrossProfit KPIMetric `json:"grossProfit"`

	NonFuleCount     KPIMetric `json:"nonFuelCount"`
	FNConversionRate KPIMetric `json:"fnConversionRate"`
	AvgTrxValue      KPIMetric `json:"avgTrxValue"`
	SalesPerM2       KPIMetric `json:"salesPerM2"`
	Y2nIncome        KPIMetric `json:"y2nIncome"`
	Q2nIncome        KPIMetric `json:"q2nIncome"`
	M2nIncome        KPIMetric `json:"m2nIncome"`
}

type SalesPlantKPIResp struct {
	Code int           `json:"code"`
	KPI  SalesPlantKPI `json:"kpi"`
}

///////////////////////////////////////////////////////////////////////////////
// SALES PLANT LEVEL CFG

/* sales plant level config */
type SalesPlantLevelCfg struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Begin float64 `json:"begin"`
	End   float64 `json:"end"`
}

/* sales plant level config add response data */
type SalesPlantLevelCfgAddResp struct {
	Code int   `json:"code"`
	ID   int64 `json:"id"`
}

/* sales plant level cfg list response data */
type SalesPlantLevelCfgListResp struct {
	Code  int                  `json:"code"`
	Total int                  `json:"total"`
	Cfgs  []SalesPlantLevelCfg `json:"pszList"` //XXX cfgs
}

/* sales plant level cfg detail response data */
type SalesPlantLevelCfgDetailResp struct {
	Code int                `json:"code"`
	Cfg  SalesPlantLevelCfg `json:"pszCfg"` //XXX cfg
}

///////////////////////////////////////////////////////////////////////////////
// SALES PLANT LEVEL

/* sales plant level */
type SalesPlantLevel struct {
	ParentOrg string `json:"parentOrg"`
	Plant     string `json:"plant"`
	PlantName string `json:"plantName"`
	SaleDays  int    `json:"saleDays"`

	LevelID      int64   `json:"zoneID"`   //XXX levelID
	LevelName    string  `json:"zoneName"` //XXX levelName
	NetSalesDAvg float64 `json:"sales"`    //XXX netSalesDAvg
	NetSales     float64 `json:"netSales"` //XXX netSales

	/* car service */
	CarServiceSR bool `json:"carServSR"`
	CarServiceLR bool `json:"carServLR"`
}

/* sales plant level response data */
type SalesPlantLevelResp struct {
	Code      int               `json:"code"`
	Total     int               `json:"total"`
	PlantList []SalesPlantLevel `json:"pszList"` //XXX plantList
}

type SalesSinglePlantLevelResp struct {
	Code  int             `json:"code"`
	Plant SalesPlantLevel `json:"psz"` //XXX plantList
}

///////////////////////////////////////////////////////////////////////////////
// SALES PLANT CAR SEVICE

/* sales plant car service */
type SalesPlantCarService struct {
	Plant     string  `json:"plant"`
	PlantName string  `json:"plantName"`
	NetSales  float64 `json:"netSales"`
}

/* sales plant car service response data */
type SalesPlantCarServiceResp struct {
	Code      int                    `json:"code"`
	Total     int                    `json:"total"`
	PlantList []SalesPlantCarService `json:"plantList"`
}

///////////////////////////////////////////////////////////////////////////////
// SALES MATERIAL PLANT CFG

/* sales material no-sales config */
type SalesMatlPlantMeta struct {
	Material    string `json:"material"`
	MaterialTxt string `json:"materialName"` //XXX materialTxt
	Plant       string `json:"plant"`
	PlantName   string `json:"plantName"`
	Status      int    `json:"status"`
}

type SalesPlantMeta struct {
	Plant     string `json:"plant"`
	PlantName string `json:"plantName"`
	Status    int    `json:"status"`
}

type SalesMatlNoSalesCfg struct {
	Material    string           `json:"material"`
	MaterialTxt string           `json:"materialName"` //XXX materialTxt
	All         bool             `json:"all"`
	Plants      []SalesPlantMeta `json:"plants"`
}

type SalesMatlNoSalesCfgListResp struct {
	Code    int                   `json:"code"`
	Total   int                   `json:"total"`
	CfgList []SalesMatlNoSalesCfg `json:"pmmcList"` //XXX cfgList
}

///////////////////////////////////////////////////////////////////////////////
// SALES MATL PLANT NO-SALES

type SalesMatlNoSales struct {
	Material    string           `json:"material"`
	MaterialTxt string           `json:"materialName"` //XXX materialTxt
	Plants      []SalesPlantMeta `json:"plants"`
}

type SalesMatlNoSalesResp struct {
	Code     int                `json:"code"`
	Total    int                `json:"total"`
	MatlList []SalesMatlNoSales `json:"pmnsList"` //XXX MatlList
}
