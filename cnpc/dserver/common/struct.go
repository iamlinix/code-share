package common

///////////////////////////////////////////////////////////////////////////////
// common info

/* filter request message */
type FilterReq struct {
	/* common query flter */
	BeginDate string `json:"beginDate"`
	EndDate   string `json:"endDate"`

	OrgLevel      int    `json:"orgLevel"` //0:head | 1:branch | 2:plant
	OrgCode       string `json:"orgCode"`
	GroupByOLevel int

	ClassLevel    int    `json:"classLevel"` //0:main | 1:mid | 2:sub
	ClassCode     string `json:"classCode"`
	GroupByCLevel int

	Material string `json:"material"`

	SortBy        string  `json:"sortBy"`
	Limit         int     `json:"limit"`
	RateThreshold float64 `json:"rateThreshold"`

	/* module  */
	ZoneTxt    string `json:"zoneTxt"`
	Vendor     string `json:"vendor"`
	Month      string `json:"month"`
	ActiveOnly bool   `json:"activeOnly"`
	StartYear  int    `json:"startYear"`
	EndYear    int    `json:"endYear"`
}

/* material data */
type Material struct {
	ShortCode string `json:"shortCode"`
	LongCode  string `json:"longCode"`
	Name      string `json:"name"`
}

type MaterialResp struct {
	Code     int      `json:"code"`
	Total    int      `json:"total"`
	Material Material `json:"material"`
}

/* vendor data */
type Vendor struct {
	ShortCode string `json:"shortCode"`
	LongCode  string `json:"longCode"`
	Name      string `json:"name"`
}

type VendorResp struct {
	Code   int    `json:"code"`
	Total  int    `json:"total"`
	Vendor Vendor `json:"vendor"`
}

/* class data */
type ClassInfo struct {
	ClassLevel int    `json:"classLevel"`
	ClassCode  string `json:"classCode"`
	ClassText  string `json:"classText"`
}

type ClassResp struct {
	Code      int         `json:"code"`
	Total     int         `json:"total"`
	ClassList []ClassInfo `json:"classList"`
}

/* org data */
type OrgInfo struct {
	OrgLevel int    `json:"orgLevel"`
	OrgCode  string `json:"orgCode"`
	OrgText  string `json:"orgText"`
}

type OrgResp struct {
	Code    int       `json:"code"`
	Total   int       `json:"total"`
	OrgList []OrgInfo `json:"orgList"`
}

/* plant data */
type Plant struct {
	Plant      string  `json:"plant"`
	PlantName  string  `json:"plantName"`
	BranchName string  `json:"branchName"`
	ZtxtZcxz   string  `json:"ztxtZcxz"`
	ZtxtXjms   string  `json:"ztxtXjms"`
	ZtxtDfl    string  `json:"ztxtDfl"`
	ZtType     string  `json:"ztType"`
	ZtxtBld    float64 `json:"ztxtBld"`
	PosX       float64 `json:"posX"`
	PosY       float64 `json:"posY"`
}

type PlantResp struct {
	Code  int   `json:"code"`
	Total int   `json:"total"`
	Plant Plant `json:"plant"`
}

///////////////////////////////////////////////////////////////////////////////
//

/* common response data */
type CommonResp struct {
	Code int    `json:"code"`
	Mesg string `json:"message"`
}

///////////////////////////////////////////////////////////////////////////////
//

/* kpi setting */
type KPI struct {
	OrgCode    string `json:"orgCode"`
	OrgName    string `json:"orgName"`
	Month      int    `json:"month"`
	IncomeType int    `json:"incomeType"`
	Category   string `json:"category"`
	Value      int64  `json:"value"`
}

type KPIByMonthReq struct {
	StartDate  string `json:"startDate"`
	EndDate    string `json:"endDate"`
	OrgCode    string `json:"orgCode"`
	IncomeType int    `json:"incomeType"`
}

type KPIListResp struct {
	Code int   `json:"code"`
	Kpis []KPI `json:"kpis"`
}

type KPIResp struct {
	Code int `json:"code"`
	Kpi  KPI `json:"kpi"`
}

type KPIMonthResp struct {
	Code         int     `json:"code"`
	MonthValue   float64 `json:"monthValue"`
	QuarterValue float64 `json:"quarterValue"`
	YearValue    float64 `json:"yearValue"`
}

type PaymentType struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Money float64 `json:"money"`
	Count int     `json:"count"`
}

type PaymentTypeResp struct {
	Code     int           `json:"code"`
	Payments []PaymentType `json:"payments"`
}

///////////////////////////////////////////////////////////////////////////////
//

/* general settings */

type GeneralSetting struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type GeneralSettingResp struct {
	Code  int    `json:"code"`
	Value string `json:"value"`
}

type GeneralListResp struct {
	Code   int              `json:"code"`
	Values []GeneralSetting `json:"values"`
}
