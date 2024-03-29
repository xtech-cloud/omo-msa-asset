package cache

import (
	"encoding/json"
	"errors"
	"github.com/Baidu-AIP/golang-sdk/aip/censor"
	"github.com/micro/go-micro/v2/logger"
	"omo.msa.asset/config"
	"omo.msa.asset/proxy/nosql"
)

const (
	BD_Conclusion    int = 101 //合规
	BD_NonConclusion int = 102 //不合规
	BD_Uncertain     int = 103 //疑似
	BD_Failed        int = 104 //失败
)

type ExamineResult struct {
	UID        string `json:"-"`
	ErrorCode  int    `json:"error_code"`
	ErrorMsg   string `json:"error_msg"`
	Conclusion string `json:"conclusion"`
	Log        int64  `json:"log_id"`
	HitMd5     bool   `json:"isHitMd5"`
	Kind       int    `json:"conclusionType"`
}

func (mine *ExamineResult) GetStatus() int {
	if config.Schema.Examine.Type == "baidu" {
		return 100 + mine.Kind
	} else {
		return 0
	}
}

func ValidateAssetUrl(uid, url string) (*ExamineResult, error) {
	//通过access_token
	client := censor.NewClient(config.Schema.Examine.AccessKey, config.Schema.Examine.SecretKey)
	//图片url
	msg := client.ImgCensorUrl(url, nil)
	result := new(ExamineResult)
	err := json.Unmarshal([]byte(msg), result)
	if err != nil {
		return nil, err
	}
	if result.ErrorCode > 0 {
		return nil, errors.New(result.ErrorMsg)
	}
	result.UID = uid
	return result, nil
}

func validateAsset(uid, url string) {
	result, err := ValidateAssetUrl(uid, url)
	if err != nil {
		logger.Warn("validate asset error that url = " + url + " and msg = " + err.Error())
		return
	}
	code := result.GetStatus()
	er := nosql.UpdateAssetCode(uid, code)
	if er != nil {
		logger.Warn("set asset code failed that uid = " + uid + " and msg = " + err.Error())
	}
}
