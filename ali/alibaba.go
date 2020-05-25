package ali

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/allegro/bigcache/v2"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"time"
)

var appkey = ""
var secret = ""
var HttpPort = 8080
var cache *bigcache.BigCache

func init() {
	var confYaml = "./config.yaml"
	YamlFile, err := os.Open(confYaml)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("File： %s does not exists\n", confYaml)
		} else {
			log.Fatalln(err.Error())
		}
	}
	defer YamlFile.Close()

	conf, _ := ioutil.ReadAll(YamlFile)
	confMap := make(map[interface{}]interface{})
	err = yaml.Unmarshal(conf, &confMap)
	if err != nil {
		log.Fatalln(err.Error())
	}

	if confMap["appkey"] == nil || confMap["appkey"] == "" {
		log.Fatalln("config appkey is empty")
	}

	if confMap["secret"] == nil || confMap["secret"] == "" {
		log.Fatalln("config secret is empty")
	}

	if confMap["http_port"] != nil {
		HttpPort = confMap["http_port"].(int)
	}

	appkey = fmt.Sprintf("%s", confMap["appkey"])
	secret = fmt.Sprintf("%s", confMap["secret"])

	cache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(60 * time.Minute))
}

func GenParameter(Param map[string]string) map[string]string {
	Param["app_key"] = appkey
	Param["sign_method"] = "md5"
	Param["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	Param["format"] = "json"
	Param["v"] = "2.0"

	var keys []string
	for k, _ := range Param {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var tmp = secret
	for _, v := range keys {
		tmp += v + fmt.Sprint(Param[v])
	}
	tmp += secret
	Param["sign"] = fmt.Sprintf("%X", md5.Sum([]byte(tmp)))
	return Param
}

func SendRequest(method string, p map[string]string) ([]byte, error) {
	p["method"] = method
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
		return nil, errors.New("request error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, errors.New("io error")
	}

	return body, nil
}
