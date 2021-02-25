package common


///////////////////////////////////////////////////////////////////////////////
// scenes info

type Scenes struct {
	ID		int	`json:"id"`
	Name		string	`json:"name"`
	BeginDate	string	`json:"beginDate"`
	EndDate		string	`json:"endDate"`
	OrgInfo		[]OrgInfo	`json:"orgInfo"`	//multiple orgs
	ClassInfo	[]ClassInfo	`json:"classInfo"`	//multiple classes
	User		string	`json:"user"`
	Remark		string	`json:"remark"`
	CreatedTime	string	`json:"createdTime"`
}


/* scenes add response data */
type ScenesAddResp struct {
	Code	int	`json:"code"`
	SID	int64	`json:"SID"`
}

/* scenes list response data */
type ScenesListResp struct {
	Code		int	`json:"code"`
	Total		int	`json:"total"`
	Scenes	[]Scenes	`json:"scenes"`
}

/* scenes one response data */
type ScenesDetailResp struct {
	Code	int	`json:"code"`
	Scenes	Scenes	`json:"scenes"`
}


///////////////////////////////////////////////////////////////////////////////
//

type ScenesOrgReq struct {
	BeginDate	string	`json:"beginDate"`
	EndDate		string	`json:"endDate"`
	OrgLevel	int	`json:"orgLevel"`
	OrgCode		string	`json:"orgCode"`
	ClassInfo []ClassInfo	`json:"classInfo"`
}


///////////////////////////////////////////////////////////////////////////////
//

/* sku info */
type SKUClassInfo struct {
	ClassCode	string	`json:"classCode"`
	ClassName	string	`json:"className"`
	SKUCount	int	`json:"skuCount"`
}

/* sku class response data */
type SKUClassResp struct {
	Code		int	`json:"code"`
	Total		int	`json:"total"`
	ClassList	[]SKUClassInfo	`json:"classList"`
}


/* inventory info */
type InvClassInfo struct {
	ClassCode	string	`json:"classCode"`
	ClassName	string	`json:"className"`
	InvCount	int	`json:"stockCount"`
	InvCost		float64	`json:"stockCost"`
}

/* inv class response data */
type InvClassResp struct {
	Code		int	`json:"code"`
	Total		int	`json:"total"`
	ClassList	[]InvClassInfo	`json:"classList"`
}

///////////////////////////////////////////////////////////////////////////////
//

type Zone struct {
	Begin		float64	`json:"begin"`
	End		float64	`json:"end"`
}

type ZoneInfo struct {
	Count		int	`json:"count"`
	Values		[]Zone	`json:"values"`
}


/* price zone config */
type PriceZoneCfg struct {
	ClassLevel	int	`json:"classLevel"`
	ClassCode	string	`json:"classCode"`
	ClassText	string	`json:"classText"`
	Zones		string	`json:"zones"`
	User		string	`json:"user"`
}

/* price-zone config add response data */
type PriceZoneCfgAddResp struct {
	Code		int	`json:"code"`
	PzcID		int64	`json:"pzcID"`
}

/* price zone cfg list response data */
type PriceZoneCfgListResp struct {
	Code		int	`json:"code"`
	Total		int	`json:"total"`
	PzcList		[]PriceZoneCfg	`json:"pzcList"`
}

/* price zone cfg detail response data */
type PriceZoneCfgDetailResp struct {
	Code	int	`json:"code"`
	PZCfg	PriceZoneCfg	`json:"pzCfg"`
}

///////////////////////////////////////////////////////////////////////////////

/* price zone info */
type PriceZoneInfo struct {
	PriceZone	string	`json:"priceZone"`
	SalesCount	int	`json:"salesCount"`
	SalesIncome	float64	`json:"salesIncome"`
	SalesNet	float64	`json:"salesNet"`
	SalesCost	float64	`json:"salesCost"`
}

/* price-zone response data */
type PriceZoneResp struct {
	Code		int	`json:"code"`
	Total		int	`json:"total"`
	PZList	[]PriceZoneInfo	`json:"pzList"`
}


type PriceZoneClassInfo struct {
	ClassLevel	int	`json:"classLevel"`
	ClassCode	string	`json:"classCode"`
	ClassText	string	`json:"classText"`
	Total		int	`json:"total"`
	PZList	[]PriceZoneInfo	`json:"pzList"`
}


type PriceZoneClassResp struct {
	Code		int	`json:"code"`
	Total		int	`json:"total"`
	ClassList	[]PriceZoneClassInfo	`json:"classList"`
}

///////////////////////////////////////////////////////////////////////////////


/* material rank info */
type MatlRankInfo struct {
	Material	string	`json:"material"`
	MaterialTxt	string	`json:"materialTxt"`
	SalesCount	float64	`json:"salesCount"`
	SalesIncome	float64	`json:"salesIncome"`
	SalesNet	float64	`json:"salesNet"`
	SalesCost	float64	`json:"salesCost"`
}

/* material rank response data */
type MatlRankResp struct {
	Code		int	`json:"code"`
	Total		int	`json:"total"`
	MList	[]MatlRankInfo	`json:"mList"`
}


///////////////////////////////////////////////////////////////////////////////
//

/* specs zone info */
type SpecsZoneInfo struct {
	SpecsZone	string	`json:"specsZone"`
	SalesCount	int	`json:"salesCount"`
	SalesIncome	float64	`json:"salesIncome"`
	SalesNet	float64	`json:"salesNet"`
	SalesCost	float64	`json:"salesCost"`
}

/* specs-zone class response data */
type SpecsZoneResp struct {
	Code		int	`json:"code"`
	Total		int	`json:"total"`
	SZList	[]SpecsZoneInfo	`json:"szList"`
}


