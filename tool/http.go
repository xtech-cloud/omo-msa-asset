package tool

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
)

func GetJson(req *http.Request) (string, error) {
	result, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", err
	} else {
		return bytes.NewBuffer(result).String(), nil
	}
}

func GetEncodeJson(req *http.Request) (string, error) {
	bytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", err
	} else {
		tmp := string(bytes)
		data, err := base64.StdEncoding.DecodeString(tmp)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
}
