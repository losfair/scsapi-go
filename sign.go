package scsapi

import (
	"fmt"
	"time"
	"crypto/hmac"
	"crypto/sha1"
//	"crypto/md5"
	"encoding/base64"
)

var scsSecretKey string

func SignString(secretKey, str string) string {
	hm := hmac.New(sha1.New, []byte(secretKey))
	hm.Write([]byte(str))
	return base64.StdEncoding.EncodeToString(hm.Sum(nil))[5:15]
}

func SignRequest(secretKey, httpVerb, contentType string, amzHeaders, targetResource string, currentTime time.Time, targetData []byte) string {
	dataMD5 := ""
/*	if targetData != nil {
		md5Context := md5.New()
		md5Context.Write(targetData)
		dataMD5 = base64.StdEncoding.EncodeToString(md5Context.Sum(nil))
	}*/
	requestString := fmt.Sprintf("%s\n%s\n%s\n%s\n%s%s",httpVerb,dataMD5,contentType,currentTime.Format("Mon, 02 Jan 2006 15:04:05 GMT"),amzHeaders,targetResource)
	sigStr := SignString(secretKey, requestString)
	return sigStr
}

func Sign(httpVerb, contentType string, amzHeaders, targetResource string, currentTime time.Time, targetData []byte) string {
	return SignRequest(scsSecretKey, httpVerb, contentType, amzHeaders, targetResource, currentTime, targetData)
}

func SetSecretKey(secretKey string) {
	scsSecretKey = secretKey
}
