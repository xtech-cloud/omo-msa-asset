package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"omo.msa.asset/config"
	"omo.msa.asset/proxy"
)

const (
	QualityNone   = "none"
	QualityLow    = "low"
	QualityNormal = "normal"
	QualityHigh   = "high"
)

const (
	LiveNone = "none"
)

const (
	ImageTypeBase64 = "base64"
	ImageTypeFace   = "face_token"
)

type FaceSearchReq struct {
	Image     string `json:"image"`
	Type      string `json:"image_type"`
	Groups    string `json:"group_id_list"` //用户组列表，多个就用逗号分割
	Quality   string `json:"quality_control"`
	User      string `json:"user_id"`
	MaxUser   int    `json:"max_user_num"`    //返回的用户数量[1,50]，默认1
	Threshold int    `json:"match_threshold"` //匹配阈值[0, 100]，默认80
}

type FaceSearchResponse struct {
	Token string        `json:"face_token"`
	Users []*UserResult `json:"user_list"`
}

type FaceMultiSearchResp struct {
	LogID    uint64            `json:"log_id"`
	FaceNum  int               `json:"face_num"`
	FaceList []*UserFaceResult `json:"face_list"`
}

type UserResult struct {
	Group  string  `json:"group_id"`
	ID     string  `json:"user_id"`
	Score  float32 `json:"score"`
	Remark string  `json:"user_info"`
}

type UserFaceResult struct {
	Token    string              `json:"face_token"`
	Location *proxy.LocationInfo `json:"location"`
	Users    []*UserResult       `json:"user_list"`
}

type ImageFaceResult struct {
	LogID    uint64              `json:"log_id"`
	Token    string              `json:"face_token"`
	Location *proxy.LocationInfo `json:"location"`
}

type FaceListResp struct {
	LogID    uint64       `json:"log_id"`
	FaceList []*FaceBrief `json:"face_list"`
}

type UserListResp struct {
	LogID uint64        `json:"log_id"`
	Users []*UserResult `json:"user_list"`
}

type UserAddReq struct {
	Image   string `json:"image"`
	Type    string `json:"image_type"`
	Group   string `json:"group_id"`
	User    string `json:"user_id"`
	Meta    string `json:"user_info"`
	Quality string `json:"quality_control"`
}

func searchByOneFace(info *FaceSearchReq) (*FaceSearchResponse, error) {
	token, er := getDetectAccessToken()
	if er != nil {
		return nil, er
	}
	addr := fmt.Sprintf("%s?access_token=%s", config.Schema.Detection.OneSearch, token)
	data, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}
	msg, er := httpPost(addr, string(data))
	if er != nil {
		return nil, er
	}
	reply := new(FaceSearchResponse)
	er = json.Unmarshal(msg, reply)
	return reply, er
}

func searchByMultiFace(info *FaceSearchReq) (*FaceMultiSearchResp, error) {
	token, er := getDetectAccessToken()
	if er != nil {
		return nil, er
	}
	addr := fmt.Sprintf("%s?access_token=%s", config.Schema.Detection.MultiSearch, token)
	data, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}
	msg, er := httpPost(addr, string(data))
	if er != nil {
		return nil, er
	}
	reply := new(FaceMultiSearchResp)
	er = json.Unmarshal(msg, reply)
	return reply, er
}

func registerUserFace(info *UserAddReq) (*ImageFaceResult, error) {
	token, er := getDetectAccessToken()
	if er != nil {
		return nil, er
	}
	addr := fmt.Sprintf("%s?access_token=%s", config.Schema.Detection.Face.Add, token)
	data, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}
	bts, er := httpPost(addr, string(data))
	if er != nil {
		return nil, er
	}
	result := gjson.ParseBytes(bts)
	code := result.Get("error_code").Int()
	if code != 0 {
		return nil, errors.New(result.Get("error_msg").String())
	}
	reply := new(ImageFaceResult)
	er = json.Unmarshal(bts, reply)
	return reply, er
}

func updateUserFace(info *UserAddReq) (*ImageFaceResult, error) {
	token, er := getDetectAccessToken()
	if er != nil {
		return nil, er
	}
	addr := fmt.Sprintf("%s?access_token=%s", config.Schema.Detection.Face.Update, token)
	data, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}
	bts, er := httpPost(addr, string(data))
	if er != nil {
		return nil, er
	}
	result := gjson.ParseBytes(bts)
	code := result.Get("error_code").Int()
	if code != 0 {
		return nil, errors.New(result.Get("error_msg").String())
	}
	reply := new(ImageFaceResult)
	er = json.Unmarshal(bts, reply)
	return reply, er
}

func createUserGroup(uid string) error {
	token, er := getDetectAccessToken()
	if er != nil {
		return er
	}
	addr := fmt.Sprintf("%s?access_token=%s", config.Schema.Detection.Group.Add, token)
	data := fmt.Sprintf(`{"group_id":"%s"}`, uid)
	bts, er := httpPost(addr, data)
	if er != nil {
		return er
	}
	result := gjson.ParseBytes(bts)
	code := result.Get("error_code").Int()
	if code != 0 {
		return errors.New(result.Get("error_msg").String())
	}
	return nil
}

func removeUserGroup(uid string) error {
	token, er := getDetectAccessToken()
	if er != nil {
		return er
	}
	addr := fmt.Sprintf("%s?access_token=%s", config.Schema.Detection.Group.Delete, token)
	data := fmt.Sprintf(`{"group_id":"%s"}`, uid)
	bts, er := httpPost(addr, data)
	if er != nil {
		return er
	}
	result := gjson.ParseBytes(bts)
	code := result.Get("error_code").Int()
	if code != 0 {
		return errors.New(result.Get("error_msg").String())
	}
	return nil
}

func getUserGroups() ([]string, error) {
	token, er := getDetectAccessToken()
	if er != nil {
		return nil, er
	}
	addr := fmt.Sprintf("%s?access_token=%s", config.Schema.Detection.Group.List, token)
	bts, er := httpPost(addr, "")
	if er != nil {
		return nil, er
	}
	result := gjson.ParseBytes(bts)
	code := result.Get("error_code").Int()
	if code != 0 {
		return nil, errors.New(result.Get("error_msg").String())
	}
	arr := result.Get("result/group_id_list").Array()
	list := make([]string, 0, len(arr))
	for _, item := range arr {
		list = append(list, item.String())
	}
	return list, er
}

func getUsersByGroup(group string) ([]string, error) {
	token, er := getDetectAccessToken()
	if er != nil {
		return nil, er
	}
	addr := fmt.Sprintf("%s?access_token=%s", config.Schema.Detection.User.List, token)
	data := fmt.Sprintf(`{"group_id":"%s"}`, group)
	bts, er := httpPost(addr, data)
	if er != nil {
		return nil, er
	}
	result := gjson.ParseBytes(bts)
	code := result.Get("error_code").Int()
	if code != 0 {
		return nil, errors.New(result.Get("error_msg").String())
	}
	re := result.Get("result")
	arr := re.Get("user_id_list").Array()
	list := make([]string, 0, len(arr))
	for _, item := range arr {
		list = append(list, item.String())
	}
	return list, er
}

func getUserFacesByGroup(group, user string) (*FaceListResp, error) {
	token, er := getDetectAccessToken()
	if er != nil {
		return nil, er
	}
	addr := fmt.Sprintf("%s?access_token=%s", config.Schema.Detection.Face.List, token)
	data := fmt.Sprintf(`{"group_id":"%s", "user_id":"%s"}`, group, user)
	bts, er := httpPost(addr, data)
	if er != nil {
		return nil, er
	}
	result := gjson.ParseBytes(bts)
	code := result.Get("error_code").Int()
	if code != 0 {
		return nil, errors.New(result.Get("error_msg").String())
	}
	reply := new(FaceListResp)
	er = json.Unmarshal(bts, reply)
	return reply, er
}

func getUserMetas(group, user string) (*UserListResp, error) {
	token, er := getDetectAccessToken()
	if er != nil {
		return nil, er
	}
	addr := fmt.Sprintf("%s?access_token=%s", config.Schema.Detection.User.Get, token)
	data := fmt.Sprintf(`{"group_id":"%s", "user_id":"%s"}`, group, user)
	bts, er := httpPost(addr, data)
	if er != nil {
		return nil, er
	}
	result := gjson.ParseBytes(bts)
	code := result.Get("error_code").Int()
	if code != 0 {
		return nil, errors.New(result.Get("error_msg").String())
	}
	reply := new(UserListResp)
	er = json.Unmarshal(bts, reply)
	return reply, er
}

func removeUserFace(log uint64, face, user, group string) error {
	token, er := getDetectAccessToken()
	if er != nil {
		return er
	}
	addr := fmt.Sprintf("%s?access_token=%s", config.Schema.Detection.Face.Delete, token)
	data := fmt.Sprintf(`{"log_id":"%d", "group_id":"%s", "user_id":"%s", "face_token":"%s",}`, log, group, user, face)
	bts, er := httpPost(addr, data)
	if er != nil {
		return er
	}
	result := gjson.ParseBytes(bts)
	code := result.Get("error_code").Int()
	if code != 0 {
		return errors.New(result.Get("error_msg").String())
	}
	return nil
}
