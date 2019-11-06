package main

import (
	"context"
	"encoding/json"
	"log"

	"os"

	"flag"

	"encoding/base64"

	"strings"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/golang/glog"
	"github.com/knative-sample/face-recognition/pkg/kncloudevents"
	"github.com/knative-sample/face-recognition/pkg/manager"
)

const (
	ACCESSKEY_ID     = "ACCESSKEY_ID"
	ACCESSKEY_SECRET = "ACCESSKEY_SECRET"
	UPLOAD_OSS_PATH  = "UPLOAD_OSS_PATH"
)

/*
{"events": [{
            "eventName": "ObjectCreated:PostObject",
            "eventSource": "acs:oss",
            "eventTime": "2019-06-18T06:44:16.000Z",
            "eventVersion": "1.0",
            "oss": {
                "bucket": {
                    "arn": "acs:oss:cn-beijing:1041208914252405:testjian",
                    "name": "testjian",
                    "ownerIdentity": "1041208914252405",
                    "virtualBucket": ""},
                "object": {
                    "deltaSize": 0,
                    "eTag": "137138904F2E18D307D04EB38EA44CDA",
                    "key": "timg.jpg",
                    "size": 12990},
                "ossSchemaVersion": "1.0",
                "ruleId": "demo-image"},
            "region": "cn-beijing",
            "requestParameters": {"sourceIPAddress": "42.120.74.107"},
            "responseElements": {"requestId": "5D08884070BC12B192C65CDF"},
            "userIdentity": {"principalId": "1041208914252405"}}]}

*/

func display(event cloudevents.Event) {
	glog.Infof("cloudevents.Event\n%s", event.String())
	events := &manager.OssEvent{}
	data := event.Data.([]byte)
	tdata := strings.Replace(string(data), "\"", "", -1)
	glog.Infof("TData: %s", tdata)
	msgBytes, err := base64.StdEncoding.DecodeString(tdata)
	if err != nil {
		glog.Errorf(err.Error())
		return
	}
	glog.Infof("Event: %s", msgBytes)
	err = json.Unmarshal(msgBytes, events)
	if err != nil {
		glog.Errorf(err.Error())
		return
	}
	ak, defined := os.LookupEnv(ACCESSKEY_ID)
	if !defined {
		glog.Errorf("required environment variable '%s' not defined", ak)
		return
	}
	sk, defined := os.LookupEnv(ACCESSKEY_SECRET)
	if !defined {
		glog.Errorf("required environment variable '%s' not defined", sk)
		return
	}
	uploadPath, defined := os.LookupEnv(UPLOAD_OSS_PATH)
	if !defined {
		glog.Errorf("required environment variable '%s' not defined", sk)
		return
	}
	cfg := manager.Config{Ak: ak, Sk: sk, OssEvent: events, ConfigPath: *configpath, TargetOssPath: uploadPath}
	manager.DoFace(cfg)

}

var (
	configpath = flag.String("configpath", "", "config path")
)

func main() {
	flag.Set("logtostderr", "true")
	defer glog.Flush()
	flag.Parse()
	c, err := kncloudevents.NewDefaultClient()
	if err != nil {
		log.Fatal("Failed to create client, ", err)
	}
	log.Fatal(c.StartReceiver(context.Background(), display))
}
