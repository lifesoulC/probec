package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"probec/prober"
)

func icmpPing(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	req := &icmpReq{}
	err := json.Unmarshal(body, req)
	resp := &icmpResp{}
	if err != nil {
		resp.ErrMsg = "json error"
		resp.ErrCode = errJson
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	}

	opts := &prober.PingOpts{}
	opts.Count = req.Count
	opts.Src = req.Src
	opts.Dest = req.Dest
	opts.Interval = req.Interval
	prob.ICMPPing(opts)

	delays, err := prob.ICMPPing(opts)

	if err != nil {
		resp.ErrMsg = "probe error:" + err.Error()
		resp.ErrCode = errUnkown
	} else {
		resp.Delays = delays
	}
	resp.Token = req.Token
	resp.Src = req.Src
	resp.Dest = req.Dest

	b, _ := json.Marshal(resp)
	w.Write(b)
}

func icmpBroadcast(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	req := &icmpReq{}
	err := json.Unmarshal(body, req)
	resp := &icmpBroadcastResp{}
	if err != nil {
		resp.ErrMsg = "json error"
		resp.ErrCode = errJson
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	}

	opts := &prober.IcmpBroadcastOpts{}
	opts.Src = req.Src
	opts.Dest = req.Dest
	opts.Count = req.Count
	opts.Interval = req.Interval
	delays, e := prob.BroadCastPing(opts)
	if e != nil {
		resp.ErrMsg = "probe error:" + e.Error()
		resp.ErrCode = errUnkown
	} else {
		for _, v := range delays {
			d := &hostDelays{}
			d.Host = v.Dest.String
			d.Delays = append(d.Delays, v.Delays...)
			resp.Delays = append(resp.Delays, d)
		}
	}
	resp.Src = req.Src
	resp.Dest = req.Dest
	resp.Token = req.Token
	b, _ := json.Marshal(resp)
	w.Write(b)
}
