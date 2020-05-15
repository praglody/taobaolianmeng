package ali

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func SearchTaobaoShop(q string) (string, error) {
	p := map[string]string{
		"method":    "taobao.tbk.dg.material.optional",
		"fields":    "user_id,shop_title,shop_type,seller_nick,pict_url,shop_url",
		"q":         q,
		"adzone_id": "110280650043",
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
		return "", errors.New("request error")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("io error")
	}
	if err = resp.Body.Close(); err != nil {
		return "", err
	}
	fmt.Printf("%s\n", body)

	ret := gjson.GetBytes(body, "tbk_dg_material_optional_response.result_list.map_data")
	if ret.IsArray() {
		for _, v := range ret.Array() {
			t := v.Value()
			for key, value := range t.(map[string]interface{}) {
				fmt.Println(key, ":", value)
			}
			break
		}
		fmt.Printf("%s\n", ret.String())
		return ret.String(), nil
	}

	return "", nil
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
