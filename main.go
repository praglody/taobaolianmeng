package main

import (
	"fmt"
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
		q := ctx.URLParam("q")
		resp, err := ali.SearchTaobaoShop(q)
		if err != nil {
			ctx.StatusCode(403)
			return
		}
		ctx.WriteString(resp)
	})

	app.Run(iris.Addr(fmt.Sprintf(":%d", ali.HttpPort)))
}
