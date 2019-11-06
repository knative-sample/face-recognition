package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"os"
)

func Md5Encrypt(data string) string {
	md5Ctx := md5.New()                                           //md5 init
	md5Ctx.Write([]byte(data))                                    //md5 updata
	cipherStr := md5Ctx.Sum(nil)                                  //md5 final
	encryptedData := base64.StdEncoding.EncodeToString(cipherStr) //base64
	return encryptedData
}

//HMAC-SHA1签名算法是一种常用的签名算法，用于对一段信息进行生成签名摘要
func HmacSha1Base64(key, data string) string {
	keyByte := []byte(key)
	dataByte := []byte(data)
	hmacSha1 := hmac.New(sha1.New, keyByte)
	hmacSha1.Write(dataByte)
	return base64.StdEncoding.EncodeToString(hmacSha1.Sum(nil))
}

func RemoveContents(file string) error {
	err := os.RemoveAll(file)
	if err != nil {
		return err
	}
	return nil
}
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
