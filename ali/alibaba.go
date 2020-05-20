package ali

import (
	"crypto/md5"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"time"
)

var appkey = ""
var secret = ""
var HttpPort = 8080

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
