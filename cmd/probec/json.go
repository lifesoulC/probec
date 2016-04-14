package main

type icmpReq struct {
	Token    string `json:"token"`
	Src      string `json:"src"`
	Dest     string `json:"dest"`
	Count    int    `json:"count"`
	Interval int    `json:"interval"`
	TTL      int    `json:"ttl"`
}

type icmpResp struct {
	ErrCode int    `json:"ErrCode"`
	ErrMsg  string `json:"ErrMsg"`
	Token   string `json:"Token"`
	Src     string `json:"Src"`
	Dest    string `json:"Dest"`
	Delays  []int  `json:"Delays"`
	Count   int    `json:"Count"`
}

type hostDelays struct {
	Host   string `json:"Host"`
	Count  int    `json:"Count"`
	Delays []int  `json:"Delays"`
}

type hostTraceDelays struct {
	Host   string `json:"Host"`
	TTL    int    `json:"TTL"`
	Delays []int  `json:"Delays"`
}

type icmpBroadcastResp struct {
	ErrCode int           `json:"ErrCode"`
	ErrMsg  string        `json:"ErrMsg"`
	Token   string        `json:"Token"`
	Src     string        `json:"Src"`
	Dest    string        `json:"Dest"`
	Delays  []*hostDelays `json:"Delays"`
}

type traceResp struct {
	ErrCode int                `json:"ErrCode"`
	ErrMsg  string             `json:"ErrMsg"`
	Token   string             `json:"Token"`
	Src     string             `json:"Src"`
	Dest    string             `json:"Dest"`
	Delays  []*hostTraceDelays `json:"Delays"`
}
