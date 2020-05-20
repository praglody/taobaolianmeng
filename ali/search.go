package ali

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

func SearchTaobaoShop(q string, page string) ([]interface{}, error) {
	retry := 0
	p := map[string]string{
		"method":    "taobao.tbk.dg.material.optional",
		"fields":    "user_id,shop_title,shop_type,seller_nick,pict_url,shop_url",
		"q":         q,
		"page_no":   page,
		"page_size": "30",
		"adzone_id": "110280650043",
	}
	p = GenParameter(p)
	form := url.Values{}
	for k, v := range p {
		form[k] = []string{v}
	}

SearchRequest:
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.PostForm("http://gw.api.taobao.com/router/rest", form)
	if err != nil {
		return nil, errors.New("request error")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("io error")
	}
	if err = resp.Body.Close(); err != nil {
		return nil, err
	}

	ret := gjson.GetBytes(body, "tbk_dg_material_optional_response.result_list.map_data")
	if ret.IsArray() {
		return ret.Value().([]interface{}), nil
	}
	ret = gjson.GetBytes(body, "error_response")
	errMsg := ret.Value().(map[string]interface{})
	if errMsg["code"].(float64) == 15 && retry < 2 {
		// 服务器错误，重试
		retry++
		time.Sleep(time.Millisecond * 500)
		goto SearchRequest
	}

	return []interface{}{}, nil
}

func GetItemInfo(itemId, ip string) (interface{}, error) {
	retry := 0
	p := map[string]string{
		"method":   "taobao.tbk.item.info.get",
		"num_iids": itemId,
		"ip":       ip,
	}
	p = GenParameter(p)
	form := url.Values{}
	for k, v := range p {
		form[k] = []string{v}
	}

ItemInfoRequest:
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.PostForm("http://gw.api.taobao.com/router/rest", form)
	if err != nil {
		return nil, errors.New("request error")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("io error")
	}
	if err = resp.Body.Close(); err != nil {
		return nil, err
	}

	ret := gjson.GetBytes(body, "tbk_item_info_get_response.results.n_tbk_item")
	if ret.IsArray() {
		t := ret.Array()[0].Value().(interface{})
		return t, nil
	}
	ret = gjson.GetBytes(body, "error_response")
	errMsg := ret.Value().(map[string]interface{})
	if errMsg["code"].(float64) == 15 && errMsg["sub_code"].(string) != "50001" && retry < 2 {
		// 服务器错误，重试
		retry++
		time.Sleep(time.Millisecond * 500)
		goto ItemInfoRequest
	}

	return map[string]string{}, nil
}

func GetCouponInfo(itemId, couponId string) (interface{}, error) {
	retry := 0
	p := map[string]string{
		"method":      "taobao.tbk.coupon.get",
		"item_id":     itemId,
		"activity_id": couponId,
	}
	p = GenParameter(p)
	form := url.Values{}
	for k, v := range p {
		form[k] = []string{v}
	}

CouponInfoRequest:
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.PostForm("http://gw.api.taobao.com/router/rest", form)
	if err != nil {
		return nil, errors.New("request error")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("io error")
	}
	if err = resp.Body.Close(); err != nil {
		return nil, err
	}

	fmt.Println(string(body))

	ret := gjson.GetBytes(body, "tbk_coupon_get_response.data")
	if ret.Exists() {
		t := ret.Value().(map[string]interface{})
		return t, nil
	}
	ret = gjson.GetBytes(body, "error_response")
	errMsg := ret.Value().(map[string]interface{})
	if errMsg["code"].(float64) == 15 && errMsg["sub_code"].(string) == "1" && retry < 2 {
		// 服务器错误，重试
		retry++
		time.Sleep(time.Millisecond * 500)
		goto CouponInfoRequest
	}

	return map[string]string{}, nil
}

func GetTaoBaoServerTime() {
	p := map[string]string{
		"method": "taobao.time.get",
	}
	p = GenParameter(p)
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
