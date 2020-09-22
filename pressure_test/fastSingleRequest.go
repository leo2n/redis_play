package main

import (
	//"encoding/json"
	"go_redis/pressure_test/jsonStruct"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

//
var client2 = &fasthttp.Client{
	MaxConnsPerHost: 60000, // 一个fasthttp.Client客户端的最大TCP数量, 一般达不到65535就不会报错
	Dial: func(addr string) (conn net.Conn, err error) {
		//return connLocal, err
		return fasthttp.DialTimeout(addr, 30*time.Second) // tcp 层
	},
	ReadTimeout: 60 * time.Second, // http 应用层, 如果tcp建立起来, 但是服务器不给你回应||回应的时间太久, 难道你要一直耗着吗?  当然是关闭http链接啊
}

func fastSingleRequest(userID string, productID string, w *sync.WaitGroup, timeStatistics chan float64) (bool, error) {
	// 首先, 构造client
	client := client2
	// 构造request body里面的值
	r := jsonStruct.ReqBuy{
		UserId:      userID,
		ProductId:   productID,
		PurchaseNum: 1,
	}
	reqBody, err := r.MarshalJSON()
	if err != nil {
		w.Done()
		log.Println(err)
		return false, err
	}
	req := &fasthttp.Request{}
	//req := fasthttp.AcquireRequest()
	req.Header.Set("Authorization", token)
	req.Header.SetContentType("application/json")
	req.Header.SetMethod(http.MethodPost)
	req.SetRequestURI(url)
	req.SetBody(reqBody)

	resp := &fasthttp.Response{}
	//resp := fasthttp.AcquireResponse()
	// 开始发送请求
	t0 := time.Now() // 客户端开始发起请求的时间

	err = client.Do(req, resp)
	if err != nil {
		w.Done()
		log.Println("发送请求时:", err)
		return false, err
	}
	if resp.StatusCode() != 200 {
		w.Done()
		return false, err
	}
	t1 := time.Since(t0)           // 客户端结束发起请求的时间
	timeStatistics <- t1.Seconds() // 将客户端发起请求的时间发送给timeStatistics
	// 请求结束了
	w.Done()
	return true, nil
}
