package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"probec/prober"
)

func icmpPing(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body) //将http中的body读出来
	req := &icmpReq{}                 //初始化icmpreq结构体
	err := json.Unmarshal(body, req)  //解析json文件 body放入结构体 icmpReq 中
	resp := &icmpResp{}               //初始化icmpresp结构体
	if err != nil {
		resp.ErrMsg = "json error" //如果错误 返回错误信息
		resp.ErrCode = errJson
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	}

	//	err = checkSrcIP(req.Src) //检查ip的合法性
	//	if err != nil {
	//		resp.ErrMsg = err.Error()
	//		resp.ErrCode = errSrcIP
	//		b, _ := json.Marshal(resp)
	//		w.Write(b)
	//		return
	//	}

	opts := &prober.PingOpts{} //将opts初始化
	opts.Count = req.Count
	opts.Src = req.Src
	opts.Dest = req.Dest
	opts.Interval = req.Interval
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
	resp.Count = len(resp.Delays)

	b, _ := json.Marshal(resp)
	w.Write(b)
	//return
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

	//	err = checkSrcIP(req.Src)
	//	if err != nil {
	//		resp.ErrMsg = err.Error()
	//		resp.ErrCode = errSrcIP
	//		b, _ := json.Marshal(resp)
	//		w.Write(b)
	//		return
	//	}

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
			if v.Dest == nil {
				continue
			}
			d := &hostDelays{}
			d.Host = v.Dest.String
			d.Delays = append(d.Delays, v.Delays...)
			d.Count = len(d.Delays)
			resp.Delays = append(resp.Delays, d)
		}
	}
	resp.Src = req.Src
	resp.Dest = req.Dest
	resp.Token = req.Token
	b, _ := json.Marshal(resp)
	w.Write(b)
}

func udpTrace(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	req := &icmpReq{}
	err := json.Unmarshal(body, req)
	resp := &traceResp{}
	if err != nil {
		resp.ErrMsg = "json error"
		resp.ErrCode = errJson
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	}

	//	err = checkSrcIP(req.Src)
	//	if err != nil {
	//		resp.ErrMsg = err.Error()
	//		resp.ErrCode = errSrcIP
	//		b, _ := json.Marshal(resp)
	//		w.Write(b)
	//		return
	//	}

	opts := &prober.TraceOpts{}
	opts.Src = req.Src
	opts.Dest = req.Dest
	opts.Count = req.Count
	opts.Interval = req.Interval
	delays, e := prob.Trace(opts)

	if e != nil {
		resp.ErrMsg = "probe error:" + e.Error()
		resp.ErrCode = errUnkown
	} else {

		for _, v := range delays {
			if v.Host == nil {
				continue
			}
			d := &hostTraceDelays{}
			d.Host = v.Host.String
			d.TTL = v.TTL
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
