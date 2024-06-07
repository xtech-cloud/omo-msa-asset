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
		if mine.Kind == 1 {
			return BD_Conclusion
		} else if mine.Kind == 2 {
			return BD_NonConclusion
		} else if mine.Kind == 3 {
			return BD_Uncertain
		}
		return BD_Failed
	} else {
		return -1
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

func validateAsset(info *AssetInfo) {
	key, url := info.getMinURL()
	if len(url) < 1 {
		return
	}
	result, err := ValidateAssetUrl(key, url)
	if err != nil {
		logger.Warn("validate asset error that url = " + url + " and msg = " + err.Error())
		return
	}
	code := result.GetStatus()
	er := nosql.UpdateAssetCode(info.UID, code)
	if er != nil {
		logger.Warn("set asset code failed that uid = " + info.UID + " and msg = " + err.Error())
		return
	}
	info.Code = code
	if code == BD_Conclusion {
		group := FaceGroupDefault
		if info.Scope == AssetScopeOrg {
			group = info.Owner
		}
		_ = CheckFaceGroup(group)
		er = checkFaces(info.UID, info.Owner, url, group, info.Quote, info.Creator)
		if er != nil {
			logger.Warn("check faces failed that uid = " + info.UID + " and msg = " + er.Error())
		}
	}
}

func checkFaces(asset, owner, url, group, quote, operator string) error {
	resp, er := detectFaces(url)
	if er != nil {
		return er
	}
	er = clipFaces(asset, owner, url, group, quote, operator, resp)
	if er != nil {
		return er
	}
	return nil
}
