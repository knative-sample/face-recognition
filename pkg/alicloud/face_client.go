package alicloud

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/knative-sample/face-recognition/pkg/utils"
)

const (
	FaceAttribute = "https://dtplus-cn-shanghai.data.aliyuncs.com/face/attribute"
)

func SendFaceRequest(ak, sk, requestBody string) (body []byte, err error) {
	glog.Infof(ak)
	glog.Infof(sk)
	glog.Infof(requestBody)
	method := "POST"
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, _ := http.NewRequest(method, FaceAttribute, strings.NewReader(requestBody))

	date := time.Now().UTC().String()
	// 1.对body做MD5+BASE64加密
	bodyMd5 := utils.Md5Encrypt(requestBody)
	stringToSign := method + "\napplication/json\n" + bodyMd5 + "\napplication/json\n" + date + "\n" + req.URL.Path
	// 2.计算 HMAC-SHA1
	signature := utils.HmacSha1Base64(sk, stringToSign)

	authHeader := "Dataplus " + ak + ":" + signature

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Date", date)
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Connection", "keep-alive")
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("http send request url %s fails -- %v ", FaceAttribute, err)
		return
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	log.Printf(string(body))
	//status code not in [200, 300) fail
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = fmt.Errorf("response status code %d, error messge: %s", resp.StatusCode, string(body))
		return
	}

	if err != nil {
		log.Printf("read the result of get url %s fails, response status code %d -- %v", FaceAttribute, resp.StatusCode, err)
		return
	}

	return
}
