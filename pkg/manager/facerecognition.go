package manager

import (
	"encoding/json"
	"log"

	"fmt"

	"github.com/knative-sample/face-recognition/pkg/alicloud"
)

const ImageReq = "{\"image_url\":\"%s\"}"

type FaceAttribute struct {
	FaceNum  int       `json:"face_num"`
	FaceRect []int     `json:"face_rect"`
	FaceProb []float32 `json:"face_prob"`
	Gender   []int     `json:"gender"`
	Age      []int     `json:"age"`
}

func DoFaceAttribute(ak, sk, imageUrl string) (*FaceAttribute, error) {

	resp, err := alicloud.SendFaceRequest(ak, sk, fmt.Sprintf(ImageReq, imageUrl))
	if err != nil {
		log.Printf("SendFaceRequest fails -- %v ", err)
		return nil, err
	}
	fa := &FaceAttribute{}
	err = json.Unmarshal(resp, fa)
	if err != nil {
		log.Printf("Unmarshal error: %v ", err)
		return nil, err
	}
	log.Println(fa.FaceRect)
	log.Println(fa.Gender)
	log.Println(fa.Age)
	return fa, nil
}
