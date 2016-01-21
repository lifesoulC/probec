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
}

type hostDelays struct {
	Host   string `json:"Host"`
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
