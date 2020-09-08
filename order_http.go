package main

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"github.com/valyala/fasthttp"
	"go_redis/jsonStruct"
	shop "go_redis/mysql/shop/goods"
	"log"
	"net/http"
	"strconv"
)

func errorHandle(w http.ResponseWriter, err error, code int) {
	log.Println(err)
	http.Error(w, err.Error(), code)
}

//var cancelBuyLock sync.Mutex

// 处理用户要购买某种商品时, 提交的参数: userId, productId, productNum 的参数的处理呀
// 使用application/json的方式
func test(w http.ResponseWriter, r *http.Request) {
}

func buy(ctx *fasthttp.RequestCtx) {
	//// 请求方法限定为post
	//if ctx.Request.Header.IsPost() == false {
	//	ctx.Response.Header.Set("Allow", fasthttp.MethodPost)
	//	ctx.Error("request method must be post", 405)
	//	return
	//}
	//if r.Method != http.MethodPost {
	//	w.Header().Set("Allow", http.MethodPost)
	//	errorHandle(w, errors.New("请求方法不合法!"), 405)
	//	return
	//}

	buyReqPointer := new(jsonStruct.BuyReq)
	err := buyReqPointer.UnmarshalJSON(ctx.PostBody())
	//err := json.Unmarshal(ctx.PostBody(), buyReqPointer)
	if err != nil {
		log.Printf("%v", err)
		ctx.Error("decode json body error", 500)
		return
	}

	// 一些数据校验部分, 校验用户id, productId, productNum
	u := new(User)
	u.userID = buyReqPointer.UserId
	// 判断productId和productNum是否合法
	ok, err := u.CanBuyIt(buyReqPointer.ProductId, buyReqPointer.PurchaseNum)
	if err != nil {
		c := jsonStruct.CommonResponse{
			Code: 8005,
			Msg:  "您购买的商品数量已达到上限或者缺货",
			Data: nil,
		}
		content, err := c.MarshalJSON()
		//content, err := jsonStruct.CommonResp(c)
		if err != nil {
			ctx.Error("response info error", 500)
			return
		}
		ctx.SetContentType("application/json")
		ctx.Response.SetBody(content)
		return
	}
	if ok {
		// 生成订单信息
		_, err := u.orderGenerator(buyReqPointer.ProductId, buyReqPointer.PurchaseNum)
		if err != nil {
			c := jsonStruct.CommonResponse{
				Code: 8002,
				Msg:  "库存数量不足呀~",
				Data: nil,
			}
			content, err := c.MarshalJSON()
			//content, err := jsonStruct.CommonResp(c)
			if err != nil {
				ctx.Error("store num is not enough", 500)
				return
				//errorHandle(w, errors.New(err.Error()), 500)
			}
			ctx.SetContentType("application/json")
			ctx.SetBody(content)
			//w.Header().Set("Content-Type", "application/json")
			//w.Write(content)
			return
		}

		// 给用户的已经购买的商品hash表里面的值添加数量
		err = u.Bought(buyReqPointer.ProductId, buyReqPointer.PurchaseNum)
		if err != nil {
			c := jsonStruct.CommonResponse{
				Code: 8004,
				Msg:  "给用户的已经购买的商品hash表单productId添加数量时发生错误!",
				Data: nil,
			}
			content, err := c.MarshalJSON()
			//content, err := jsonStruct.CommonResp(c)
			if err != nil {
				//errorHandle(w, errors.New(err.Error()), 500)
				ctx.Error("add bought list error", 500)
				return
			}
			ctx.SetContentType("application/json")
			ctx.SetBody(content)
			//w.Header().Set("Content-Type", "application/json")
			//w.Write(content)
			return
		}

		//w.Header().Set("application/json", "json")
		c := jsonStruct.CommonResponse{
			Code: 8001,
			Msg:  "操作成功",
			Data: nil,
		}
		content, err := c.MarshalJSON()
		//content, err := jsonStruct.CommonResp(c)
		if err != nil {
			ctx.Error("json marshal error", 500)
			//errorHandle(w, errors.New(err.Error()), 500)
		}
		ctx.SetContentType("application/json")
		//w.Header().Set("Content-Type", "application/json")
		//w.Write(content)
		ctx.SetBody(content)
		return
	}
}

// redis收到后台的请求, 用户取消了订单, 需要用到的参数有: userId, productId, purchaseNum,  redis直接操作用户的: user:[userId]:bought 里面key为productId的, 赋值为0
// 这个接口必须由后台调用, 因为我没有做数据校验
func cancelBuy(ctx *fasthttp.RequestCtx) {
	//if ctx.Request.Header.IsPost() == false {
	//	ctx.Request.Header.Set("Allow", http.MethodPost)
	//	ctx.Error("request method is not supported", 405)
	//	return
	//}
	//if r.Method != http.MethodPost {
	//	w.Header().Set("Allow", http.MethodPost)
	//	errorHandle(w, errors.New("请求方式不合法!"), 405)
	//	return
	//}

	// 解析: /cancelBuy接口传过来的四个参数, userId, productId, purchaseNum, orderId
	cancelBuyReqPointer := new(jsonStruct.CancelBuyReq)
	err := json.Unmarshal(ctx.Request.Body(), cancelBuyReqPointer)
	if err != nil {
		log.Println(err)
		ctx.Error("decode request body error", 500)
		return
	}
	//cancelBuyReqPointer, err := decodeCancelBuyReq(r.Body)
	//if err!=nil {
	//	errorHandle(w, errors.New("reqBody解析到struct时出错!"), 500)
	//	return
	//}
	u := new(User)
	u.userID = cancelBuyReqPointer.UserId
	err = u.CancelBuy(cancelBuyReqPointer.OrderNum)
	if err != nil {
		c := jsonStruct.CommonResponse{
			Code: 8006,
			Msg:  "取消订单时失败!",
			Data: nil,
		}
		content, err := jsonStruct.CommonResp(c)
		if err != nil {
			ctx.Error("encode resp body to []byte error", 500)
			return
		}
		ctx.SetContentType("application/json")
		ctx.SetBody(content)
		//w.Header().Set("Content-Type", "application/json")
		//w.Write(content)
		return
	}
	c := jsonStruct.CommonResponse{
		Code: 8007,
		Msg:  "取消订单成功!",
		Data: nil,
	}
	content, err := jsonStruct.CommonResp(c)
	if err != nil {
		ctx.Error("encode resp body to []byte error", 500)
		return
		//errorHandle(w, errors.New(err.Error()), 500)
	}
	ctx.SetContentType("application/json")
	ctx.SetBody(content)
	//w.Header().Set("Content-Type", "application/json")
	//w.Write(content)
	return
}

// 调用这个函数, 立刻同步(mysql中存在 && redis中不存在)的商品数据到redis
func syncRedis(ctx *fasthttp.RequestCtx) {
	redisconn := pool.Get()
	defer redisconn.Close()

	storeList, err := redis.Strings(redisconn.Do("keys", "store:*"))
	if err!=nil {
		log.Println(err)
	}
	// 分离商品的ID出来, 到单独的store id list
	storeIDlist := make([]string, 0, 128)
	for _, v := range storeList {
		storeIDlist = append(storeIDlist, v[6:])
	}
	log.Println(storeIDlist)
	// 从现有的MySQL表格中找到具体数据
	goodsList, err := shop.QueryGoods()
	if err!=nil {
		log.Println(err)
	}
	for _, v := range goodsList {
		_, ok := FindElement(storeIDlist, strconv.Itoa(v.ProductId))
		if !ok {
			// 给redis中添加相关商品数据
			err = redisconn.Send("hmset", "store:"+strconv.Itoa(v.ProductId), "productName", v.ProductName, "productId", v.ProductId, "storeNum", v.Inventory)
			if err != nil {
				log.Printf("%+v创建hash `store:%s`失败", err, v.ProductId)
			}
		}
	}
	if err!=nil {
		ctx.Error("内部处理错误", fasthttp.StatusInternalServerError)
	}
	ctx.Response.SetStatusCode(200)
	respJson, err := jsonStruct.CommonResp(jsonStruct.CommonResponse{
		Code: 8001,
		Msg:  "处理成功",
		Data: nil,
	})
	if err!=nil {
		ctx.Error("内部处理错误", fasthttp.StatusInternalServerError)
	}
	ctx.Response.SetBody(respJson)
	ctx.Response.Header.Set("Content-Type", "application/json")
}

// 展示商品清单
func goodsList(ctx *fasthttp.RequestCtx) {
	redisconn := pool.Get()
	defer redisconn.Close()

	reply, err := redis.Strings(redisconn.Do("keys", "store:*"))
	if err!=nil {
		log.Println(err)
		ctx.Error("内部处理错误", fasthttp.StatusInternalServerError)
	}
	type goods struct {
		ProductName string `redis:"productName"`
		ProductId int `redis:"productId"`
		StoreNum int `redis:"storeNum"`
	}
	goodsList := make([]*goods, 0)
	for _, v := range reply {
		log.Println(v)
		goodsMap, err := redis.Values(redisconn.Do("hgetall", v))
		if err!=nil {
			log.Println(err)
		}
		//log.Println(goodsMap)
		g := new(goods)
		err = redis.ScanStruct(goodsMap, g)
		if err!=nil {
			log.Println("redis scanStruct error: ", err)
		}
		log.Println(g)
		goodsList = append(goodsList, g)
	}
	response := jsonStruct.CommonResponse{
		Code: 8001,
		Msg:  "success",
		Data: goodsList,
	}
	err = json.NewEncoder(ctx.Response.BodyWriter()).Encode(response)
	if err!=nil {
		ctx.Error("internel error", fasthttp.StatusInternalServerError)
	}
	ctx.Response.Header.Set("Content-Type", "application/json")
}

// 更新商品限制计划
// 例如, 在更新MySQL后, 若要将商品购买限制同步到mysql中, 只需要调用goodsLimit这个接口就可以
func syncGoodsLimit(ctx *fasthttp.RequestCtx) {
	// 加载limit限制计划
	err := loadLimit()
	if err != nil {
		log.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
	response := jsonStruct.CommonResponse{
		Code: 8001,
		Msg:  "success",
		Data: nil,
	}
	err = json.NewEncoder(ctx.Response.BodyWriter()).Encode(response)
	if err!=nil {
		ctx.Error("internel error", fasthttp.StatusInternalServerError)
	}
	ctx.Response.Header.Set("Content-Type", "application/json")
}