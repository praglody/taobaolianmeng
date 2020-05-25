package main

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"io/ioutil"
	"os"
	"taobaolianmeng/ali"
)

func main() {
	app := iris.New()
	app.Use(logger.New())
	app.Use(recover.New())

	app.HandleDir("/js", "./public/js")
	app.HandleDir("/css", "./public/css")
	app.Get("/", func(ctx iris.Context) {
		indexPage, _ := os.Open("./public/index.html")
		s, _ := ioutil.ReadAll(indexPage)
		indexPage.Close()
		ctx.Write(s)
	})

	app.Get("/search", func(ctx iris.Context) {
		code := 200
		q := ctx.URLParam("q")
		p := ctx.URLParam("p")
		if p == "" {
			p = "0"
		}
		resp, err := ali.SearchTaobaoShop(q, p, ctx.RemoteAddr())
		if err != nil {
			if err != nil {
				code = 10005
				resp = []interface{}{}
			}
		}

		ctx.Header("Content-Type", "application/json; charset=utf-8")
		retMsg := map[string]interface{}{
			"code": code,
			"data": map[string]interface{}{
				"result": resp,
			},
		}

		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		b, _ := json.Marshal(&retMsg)
		ctx.Write(b)
	})

	app.Get("/item-info", func(ctx iris.Context) {
		code := 200
		itemId := ctx.URLParam("id")
		ip := ctx.RemoteAddr()
		resp, err := ali.GetItemInfo(itemId, ip)
		if err != nil {
			code = 10005
			resp = map[string]string{}
		}
		ctx.Header("Content-Type", "application/json; charset=utf-8")
		retMsg := map[string]interface{}{
			"code": code,
			"data": map[string]interface{}{
				"result": resp,
			},
		}

		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		b, _ := json.Marshal(&retMsg)
		ctx.Write(b)
	})

	app.Get("/coupon-info", func(ctx iris.Context) {
		code := 200
		itemId := ctx.URLParam("id")
		couponId := ctx.URLParam("coupon_id")
		resp, err := ali.GetCouponInfo(itemId, couponId)
		ctx.Header("Content-Type", "application/json; charset=utf-8")
		if err != nil {
			code = 10005
			resp = map[string]string{}
		}

		retMsg := map[string]interface{}{
			"code": code,
			"data": map[string]interface{}{
				"result": resp,
			},
		}
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		b, _ := json.Marshal(&retMsg)
		ctx.Write(b)
	})

	app.Get("/recommend", func(ctx iris.Context) {
		code := 200
		page := ctx.URLParam("page")
		pageSize := ctx.URLParam("page_size")
		resp, err := ali.GetRecommendList(page, pageSize)
		ctx.Header("Content-Type", "application/json; charset=utf-8")
		if err != nil {
			code = 10005
			resp = map[string]string{}
		}

		retMsg := map[string]interface{}{
			"code": code,
			"data": map[string]interface{}{
				"result": resp,
			},
		}
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		b, _ := json.Marshal(&retMsg)
		ctx.Write(b)
	})

	// 生成口令
	app.Post("/get-share-key", func(ctx iris.Context) {
		type KeyParam struct {
			Title string `json:"title"`
			Url   string `json:"url"`
		}

		code := 200
		retMsg := map[string]interface{}{}
		var share KeyParam
		err := ctx.ReadJSON(&share)
		if err != nil {
			share.Title = ""
			share.Url = ""
		}

		//resp, err := ali.GetShareKey(share.Title, share.Url)
		// {"code":200,"data":{"result":{"model":"￥JbMZ1JpQ3Rq￥"}}}
		resp, err := map[string]string{"model": "￥JbMZ1JpQ3Rq￥"}, nil
		ctx.Header("Content-Type", "application/json; charset=utf-8")
		if err != nil {
			retMsg = map[string]interface{}{
				"code":   10005,
				"errMsg": err.Error(),
			}
		} else {
			retMsg = map[string]interface{}{
				"code": code,
				"data": map[string]interface{}{
					"result": resp,
				},
			}
		}

		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		b, _ := json.Marshal(&retMsg)
		ctx.Write(b)
	})

	app.Run(iris.Addr(fmt.Sprintf(":%d", ali.HttpPort)))
}
