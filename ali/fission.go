package ali

import (
	"crypto/md5"
	"errors"
	"fmt"
	QrCode "github.com/skip2/go-qrcode"
	"os"
	"path/filepath"
	"time"
)

type FissionParam struct {
	TaoKey string `json:"taoKey"`
	Title  string `json:"title"`
	Cover  string `json:"cover"`
}

func GetFissionUrl(fission FissionParam) (string, error) {
	uniqKey := fmt.Sprintf("%X", md5.Sum([]byte(fission.Title+fission.Cover+fission.TaoKey)))
	fileName := fmt.Sprintf("/pic/%s/%s.png", time.Now().Format("2006-01-02"), uniqKey)

	if _, err := os.Stat("./runtime" + fileName); err != nil {
		if !os.IsNotExist(err) {
			return "", errors.New("file status error")
		}

		// 文件不存在
		fileFd, err := os.OpenFile("./runtime"+fileName, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			path := filepath.Dir("./runtime" + fileName)
			_ = os.Mkdir(path, 0755)

			fileFd, err = os.OpenFile("./runtime"+fileName, os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return "", err
			}
		}
		defer fileFd.Close()
		png, err := getFissionQrCode(fission.TaoKey)
		if err != nil {
			return "", err
		}

		_, err = fileFd.Write(png)
		if err != nil {
			return "", nil
		}
	}

	return fileName, nil
}

func getFissionQrCode(shareKey string) ([]byte, error) {
	var png []byte
	png, err := QrCode.Encode(shareKey, QrCode.Medium, 256)
	if err != nil {
		return nil, errors.New("generate qrCode error: " + err.Error())
	}

	return png, nil
}
