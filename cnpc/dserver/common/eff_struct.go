package common


///////////////////////////////////////////////////////////////////////////////
// INVENTORY METRIC

/* inv metric info struct */
type InvMetric struct {
	Zinvsl		float64
	ZinvCost	float64
}

type InvOrgMetric struct {
	OrgCode		string
	OrgText		string
	PosX		float64
	PosY		float64
	Metric		InvMetric
}

type InvClassMetric struct {
	ClassCode	string
	ClassText	string
	Metric		InvMetric
}

type InvMatlMetric struct {
	Material	string
	MaterialTxt	string
	Metric		InvMetric
}

type InvDateMetric struct {
	Date		string
	Metric		InvMetric
}


///////////////////////////////////////////////////////////////////////////////
// EFF Metric CFG

type EffMetricCfg struct {
	Name		string	`json:"name"`
	Value		float64	`json:"value"`
	Module		string	`json:"module"`
}

type EffMetricCfgListResp struct {
	Code		int	`json:"code"`
	Total		int	`json:"total"`
	Cfgs	[]EffMetricCfg	`json:"cfgs"`
}


///////////////////////////////////////////////////////////////////////////////
// SALES-ACT-RATE Metric

type SARMetric struct {
	NetSales	float64	`json:"netSales"`
	InvCost		float64	`json:"invCost"`

	SalesActRate	float64	`json:"salesActRate"`
}

// Material SAR STRUCT
type EffMatlSAR struct {
	Material	string	`json:"material"`
	MaterialTxt	string	`json:"materialTxt"`

	Metric		SARMetric `json:"metric"`
}


type EffMatlSARResp struct {
	Code		int	`json:"code"`
	Total		int	`json:"total"`
	MatlList	[]EffMatlSAR `json:"matlList"`
}


///////////////////////////////////////////////////////////////////////////////
// INVENTORY Metric

type InvEffMetric struct {
	Cost		float64	`json:"cost"`
	OpenInvCost	float64	`json:"openInvCost"`
	CloseInvCost	float64	`json:"closeInvCost"`

	DaysSalesInv	float64	`json:"daysSalesInv"`
	AvgInvCost	float64	`json:"avgInvCost"`
}

// Material INV STRUCT
type EffMatlInv struct {
	Material	string	`json:"material"`
	MaterialTxt	string	`json:"materialTxt"`

	Metric		InvEffMetric `json:"metric"`
}


type EffMatlInvResp struct {
	Code		int	`json:"code"`
	Total		int	`json:"total"`
	MatlList	[]EffMatlInv `json:"matlList"`
}


///////////////////////////////////////////////////////////////////////////////
// EFF TEST-MARKET CFG

type EffTestMktMatlCfg struct {
	Material	string	`json:"material"`
	MaterialTxt	string	`json:"materialTxt"`
	SrateTH		float64
	Status		int	`json:"status"`
	CreatedTime	string	`json:"createdTime"`
}

type EffTestMktMatlCfgListResp struct {
	Code		int	`json:"code"`
	Total		int	`json:"total"`
	Cfgs	[]EffTestMktMatlCfg	`json:"cfgs"`
}


///////////////////////////////////////////////////////////////////////////////
// EFF TEST-MARKET

type TestMktMetric struct {
	PurchQty	float64	`json:"purchQty"`
	PurchCost	float64	`json:"purchCost"`
	SalesQty	float64	`json:"salesQty"`
	SalesCost	float64	`json:"salesCost"`

	SurvivalRate	float64	`json:"survivalRate"`
}

type EffTestMktMatl struct {
	Material	string	`json:"material"`
	MaterialTxt	string	`json:"materialTxt"`
	FirstOrderDate	string	`json:"firstOrderDate"`
	OnMarketDays	float64	`json:"onMarketDays"`

	Metric		TestMktMetric `json:"metric"`
}

type EffTestMktMatlSlice []EffTestMktMatl

func (s EffTestMktMatlSlice) Len() int           { return len(s) }
func (s EffTestMktMatlSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s EffTestMktMatlSlice) Less(i, j int) bool {
	return s[i].Metric.SurvivalRate < s[j].Metric.SurvivalRate
}

type EffTestMktMatlFODListResp struct {
	Code		int	`json:"code"`
	Total		int	`json:"total"`
	MatlList	[]TestMktMatl `json:"matlList"`
}


type EffTestMktMatlListResp struct {
	Code		int	`json:"code"`
	Total		int	`json:"total"`
	MatlList	[]EffTestMktMatl `json:"matlList"`
}

type EffTestMktMatlSingleResp struct {
	Code		int	`json:"code"`
	Matl	EffTestMktMatl	`json:"matl"`
}


