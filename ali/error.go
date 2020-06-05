package ali

import (
	"fmt"
	"github.com/kataras/iris/v12"
)

const ErrorClientCopy = 10060

type errMsg struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func ErrorHandle(ctx iris.Context) {
	var body errMsg
	if err := ctx.ReadJSON(&body); err != nil {
		return
	}

	switch body.Code {
	case ErrorClientCopy:
		fmt.Printf("%s: %s User-Agent:%s\n", ctx.RemoteAddr(), body.Msg, ctx.GetHeader("user-agent"))
	default:
		fmt.Println("unknown error")
	}
}
