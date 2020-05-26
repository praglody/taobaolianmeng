package ali

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"regexp"
	"strings"
	"time"
)

func SearchTaobaoShop(q string, page string, ip string) ([]interface{}, error) {
	retry := 0

	if len(q) == 0 {
		return nil, errors.New("keyword is empty")
	}

	re := regexp.MustCompile("【.*?】")
	find := re.FindString(q)
	if len(find) > 0 {
		find = strings.Trim(find, "【】")
		lindex := strings.Index(find, ":")
		if lindex != -1 {
			lindex++
			find = find[lindex:]
		}
		lindex = strings.Index(find, "：")
		if lindex != -1 {
			lindex++
			find = find[lindex:]
		}
		lindex = strings.Index(find, "(")
		if lindex != -1 {
			find = find[0:lindex]
		}
		lindex = strings.Index(find, "（")
		if lindex != -1 {
			find = find[0:lindex]
		}

		q = find
	}

	if len(q) > 0 {
		return nil, errors.New("keyword is empty")
	}

	p := map[string]string{
		"fields":      "user_id,shop_title,shop_type,seller_nick,pict_url,shop_url",
		"q":           q,
		"page_no":     page,
		"page_size":   "15",
		"adzone_id":   "110280650043",
		"material_id": "17004",
		"ip":          ip,
	}
SearchRequest:
	body, err := SendRequest("taobao.tbk.dg.material.optional", p)
	if err != nil {
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
		"num_iids": itemId,
		"ip":       ip,
	}

ItemInfoRequest:
	body, err := SendRequest("taobao.tbk.item.info.get", p)
	if err != nil {
		return nil, err
	}

	ret := gjson.GetBytes(body, "tbk_item_info_get_response.results.n_tbk_item")
	if ret.IsArray() {
		t := ret.Array()[0].Value().(interface{})
		return t, nil
	}
	ret = gjson.GetBytes(body, "error_response")
	errMsg := ret.Value().(map[string]interface{})
	fmt.Println(string(body))
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
		"item_id":     itemId,
		"activity_id": couponId,
	}

CouponInfoRequest:
	body, err := SendRequest("taobao.tbk.coupon.get", p)
	if err != nil {
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

func GetShareKey(shareTitle, shareUrl string) (interface{}, error) {
	retry := 0

	if shareTitle == "" || shareUrl == "" {
		return nil, errors.New("param error")
	}

	p := map[string]string{
		"url":  shareUrl,
		"text": shareTitle,
	}

	cacheKey := fmt.Sprintf("%X", md5.Sum([]byte("url"+p["url"]+"text"+p["text"])))
	shareKey, err := cache.Get(cacheKey)
	if err == nil {
		return map[string]string{"model": string(shareKey)}, nil
	}

ShareKeyRequest:
	body, err := SendRequest("taobao.tbk.tpwd.create", p)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))

	ret := gjson.GetBytes(body, "tbk_tpwd_create_response.data")
	if ret.Exists() {
		t := ret.Value().(map[string]interface{})
		_ = cache.Set(cacheKey, []byte(t["model"].(string)))
		return t, nil
	}
	ret = gjson.GetBytes(body, "error_response")
	errMsg := ret.Value().(map[string]interface{})
	if errMsg["sub_code"].(string) == "1" && retry < 2 {
		// 服务器错误，重试
		retry++
		time.Sleep(time.Millisecond * 500)
		goto ShareKeyRequest
	}

	return nil, errors.New(errMsg["sub_msg"].(string))
}

func GetRecommendList(page, pageSize string) (interface{}, error) {
	retry := 0
	if pageSize == "" {
		pageSize = "30"
	}
	if page == "" {
		page = "1"
	}
	p := map[string]string{
		"adzone_id":   "110280650043",
		"material_id": "13366",
		"page_no":     page,
		"page_size":   pageSize,
	}

RecommendListRequest:
	body, err := SendRequest("taobao.tbk.dg.optimus.material", p)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))

	ret := gjson.GetBytes(body, "tbk_dg_optimus_material_response.result_list.map_data")
	if ret.IsArray() {
		return ret.Value().([]interface{}), nil
	}

	ret = gjson.GetBytes(body, "error_response")
	errMsg := ret.Value().(map[string]interface{})
	if errMsg["sub_code"].(string) == "40001" && retry < 2 {
		// 服务器错误，重试
		retry++
		time.Sleep(time.Millisecond * 500)
		goto RecommendListRequest
	}

	return map[string]string{}, nil
}

func GetTaoBaoServerTime() (string, error) {
	p := map[string]string{}
	body, err := SendRequest("taobao.time.get", p)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
