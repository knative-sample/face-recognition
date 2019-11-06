package manager

import (
	"log"

	"io"
	"os"

	"path/filepath"

	"fmt"

	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/knative-sample/face-recognition/pkg/utils"
)

type Config struct {
	Ak            string
	Sk            string
	OssEvent      *OssEvent
	ConfigPath    string
	TargetOssPath string
	//OssEndpoint string
	//OssBucket   string
}
type OssEvent struct {
	Events []OssEventInfo `json:"events"`
}
type OssEventInfo struct {
	EventName    string `json:"eventName"`
	EventSource  string `json:"eventSource"`
	EventTime    string `json:"eventTime"`
	EventVersion string `json:"eventVersion"`
	Region       string `json:"region"`
	Oss          Oss    `json:"oss"`
}
type Oss struct {
	Bucket Bucket    `json:"bucket"`
	Object OssObject `json:"object"`
}
type Bucket struct {
	Arn  string `json:"arn"`
	Name string `json:"name"`
}
type OssObject struct {
	Key string `json:"key"`
}

func DoFace(cfg Config) {
	ak := cfg.Ak
	sk := cfg.Sk
	targetPath := cfg.TargetOssPath
	region := cfg.OssEvent.Events[0].Region
	bucketName := cfg.OssEvent.Events[0].Oss.Bucket.Name
	imageName := cfg.OssEvent.Events[0].Oss.Object.Key
	if strings.HasPrefix(imageName, targetPath) {
		return
	}

	endpoint := fmt.Sprintf("oss-%s.aliyuncs.com", region)
	imageUrl := fmt.Sprintf("https://%s.oss-%s.aliyuncs.com/%s", bucketName, region, imageName)
	client, err := oss.New(endpoint, ak, sk)
	if err != nil {
		log.Printf("oss.New error:%s", err.Error())
		return
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		log.Printf("oss connect to bucket error:%s", err.Error())
		return
	}
	if strings.Contains(imageName, "/") {
		imageNames := strings.Split(imageName, "/")
		imageName = imageNames[len(imageNames)-1]
	}
	tmpPath := "/app/tmp/face"
	tmpTargetPath := "/app/tmp/face/target/"
	if !utils.Exists(tmpPath) {
		os.MkdirAll(tmpPath, os.ModePerm)
	}
	if !utils.Exists(tmpTargetPath) {
		os.MkdirAll(tmpTargetPath, os.ModePerm)
	}
	fd, err := os.OpenFile(filepath.Join(tmpPath, imageName), os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		log.Println("occurred error:", err)
		return
	}
	defer fd.Close()

	body, err := bucket.GetObject(cfg.OssEvent.Events[0].Oss.Object.Key)
	if err != nil {
		log.Println("occurred error:", err)
		return
	}
	io.Copy(fd, body)
	body.Close()
	fa, err := DoFaceAttribute(ak, sk, imageUrl)
	if err != nil {
		log.Println("DoFaceAttribute error:", err)
		return
	}
	//
	err = Mark(cfg.ConfigPath, filepath.Join(tmpPath, imageName), filepath.Join(tmpTargetPath, imageName), fa)
	if err != nil {
		log.Printf("Mark %s error:%s", tmpPath+imageName, err.Error())
		return
	}
	//
	err = bucket.PutObjectFromFile(filepath.Join(targetPath, "rec-"+imageName), tmpTargetPath+imageName)
	if err != nil {
		log.Printf("oss put object %s error:%s", filepath.Join(targetPath, imageName), err.Error())
		return
	}
	err = utils.RemoveContents(tmpTargetPath + imageName)
	if err != nil {
		log.Printf("removeContents %s error:%s", tmpTargetPath+imageName, err.Error())
	}
}
