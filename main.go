package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
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
		index, _ := os.Open("./public/index.html")
		s, _ := ioutil.ReadAll(index)
		index.Close()
		ctx.Write(s)
	})

	app.Get("/search", func(ctx iris.Context) {
		q := ctx.URLParam("q")
		resp := SearchTaobaoShop(q)
		ctx.WriteString(resp)
	})

	app.Run(iris.Addr(":8080"))
	//GetTaoBaoServerTime()

}

func SearchTaobaoShop(q string) string {
	p := map[string]string{
		"method":    "taobao.tbk.dg.material.optional",
		"fields":    "user_id,shop_title,shop_type,seller_nick,pict_url,shop_url",
		"q":         q,
		"adzone_id": "110280650043",
	}
	p = ali.GenParameter(p)
	form := url.Values{}
	for k, v := range p {
		form[k] = []string{v}
	}
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.PostForm("http://gw.api.taobao.com/router/rest", form)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	resp.Body.Close()
	fmt.Printf("%s\n", body)
	ret := gjson.GetBytes(body, "tbk_dg_material_optional_response.result_list.map_data")
	//total := gjson.GetBytes(body, "tbk_dg_material_optional_response.total_results")
	for _, v := range ret.Array() {
		t := v.Value()
		for key, value := range t.(map[string]interface{}) {
			fmt.Println(key, ":", value)
		}
		break
	}
	return ret.String()
}

func GetTaoBaoServerTime() {
	p := map[string]string{
		"method": "taobao.time.get",
	}
	p = ali.GenParameter(p)
	form := url.Values{}
	for k, v := range p {
		form[k] = []string{v}
	}
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.PostForm("http://gw.api.taobao.com/router/rest", form)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	resp.Body.Close()
	for k, v := range p {
		fmt.Println(k, v)
	}

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Request.URL)
	fmt.Println(string(body))
}
