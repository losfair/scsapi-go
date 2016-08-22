package scsapi

import (
	"fmt"
	"strings"
	"errors"
	"bytes"
	"time"
	"net/http"
	"io/ioutil"
)

var scsAccessKey string
var scsBucketName string

func SetAccessKey(accessKey string) {
	scsAccessKey = accessKey
}

func SetBucketName(bucketName string) error {
	if strings.Contains(bucketName,"/") {
		return errors.New("Bad bucket name format")
	}
	scsBucketName = bucketName
	return nil
}

func doUpload(sigStr, aclStr, dataType, path string, currentTime time.Time, data []byte) error {
	client := &http.Client{}
	req,err := http.NewRequest("PUT",fmt.Sprintf("http://up.sinacloud.net%s?formatter=json",path),bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization",fmt.Sprintf("SINA %s:%s",scsAccessKey,sigStr))
	req.Header.Set("Date",currentTime.Format("Mon, 02 Jan 2006 15:04:05 GMT"))

	if aclStr != "" {
		req.Header.Set("x-amz-acl",aclStr)
	}

	if dataType != "" {
		req.Header.Set("Content-Type",dataType)
	}

	resp,err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	_,err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func Upload(path, dataType, aclStr string, data []byte) error {
	if !strings.HasPrefix(path,"/") {
		return errors.New("Bad path prefix")
	}
	currentTime := time.Now().UTC()
	amzHeaders := ""
	if aclStr != "" {
		amzHeaders = fmt.Sprintf("x-amz-acl: %s\n",aclStr)
	}
	sigStr := Sign("PUT",dataType,amzHeaders,"/"+scsBucketName+path,currentTime,data)
	if sigStr=="" {
		return errors.New("Signing failed")
	}
	return doUpload(sigStr,aclStr,dataType,"/"+scsBucketName+path,currentTime,data)
}

