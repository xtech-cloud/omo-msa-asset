package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"omo.msa.asset/config"
	"omo.msa.asset/proxy"
	"omo.msa.asset/tool"
)

const (
	ErrorCodeNotMatch   = 222207 //人脸用户不匹配
	ErrorCodeFaceExist  = 223105 //人脸已经存在
	ErrorCodeFaceNone   = 222202 //没有人脸
	ErrorCodeFaceFailed = 222203 //人脸解析失败
)

const (
	QualityNone   = "NONE"
	QualityLow    = "LOW"
	QualityNormal = "NORMAL"
	QualityHigh   = "HIGH"
)

const (
	LiveNone = "NONE"
)

const (
	ImageTypeBase64 = "BASE64"
	ImageTypeFace   = "FACE_TOKEN"
)

const (
	FaceGroupDefault = "default_users"
)

type FaceSearchReq struct {
	Image     string `json:"image"`
	Type      string `json:"image_type"`
	Groups    string `json:"group_id_list"` //用户组列表，多个就用逗号分割
	Quality   string `json:"quality_control"`
	User      string `json:"user_id,omitempty"` //如果指定该字段，则是人脸认证
	MaxUser   int    `json:"max_user_num"`      //返回的用户数量[1,50]，默认1
	Threshold int    `json:"match_threshold"`   //匹配阈值[0, 100]，默认80
}

type FaceSearchResult struct {
	Token string        `json:"face_token"`
	Users []*UserResult `json:"user_list"`
}

type FaceMultiSearchResult struct {
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

type FaceAddReq struct {
	Image   string `json:"image"`
	Type    string `json:"image_type"`
	Group   string `json:"group_id"`
	User    string `json:"user_id"`
	Meta    string `json:"user_info"`
	Quality string `json:"quality_control"`
}

func searchFaceByOne(info *FaceSearchReq) (*FaceSearchResult, error) {
	token, er := getDetectAccessToken()
	if er != nil {
		return nil, er
	}
	addr := fmt.Sprintf("%s?access_token=%s", config.Schema.Detection.OneSearch, token)
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
	reply := new(FaceSearchResult)
	if code == ErrorCodeNotMatch {
		reply.Token = ""
		return reply, nil
	}
	if code != 0 {
		return nil, errors.New(result.Get("error_msg").String())
	}

	re := result.Get("result").String()
	er = json.Unmarshal([]byte(re), reply)
	return reply, er
}

func searchFaceByMulti(info *FaceSearchReq) (*FaceMultiSearchResult, error) {
	token, er := getDetectAccessToken()
	if er != nil {
		return nil, er
	}
	addr := fmt.Sprintf("%s?access_token=%s", config.Schema.Detection.MultiSearch, token)
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
	reply := new(FaceMultiSearchResult)
	if code == ErrorCodeNotMatch {
		return reply, nil
	}
	if code != 0 {
		return nil, errors.New(result.Get("error_msg").String())
	}

	re := result.Get("result").String()
	er = json.Unmarshal([]byte(re), reply)
	return reply, er
}

func registerUserFace(info *FaceAddReq) (*ImageFaceResult, int, error) {
	token, er := getDetectAccessToken()
	if er != nil {
		return nil, -1, er
	}
	addr := fmt.Sprintf("%s?access_token=%s", config.Schema.Detection.Face.Add, token)
	data, err := json.Marshal(info)
	if err != nil {
		return nil, -1, err
	}
	bts, er := httpPost(addr, string(data))
	if er != nil {
		return nil, -1, er
	}
	result := gjson.ParseBytes(bts)
	code := result.Get("error_code").Int()
	if code != 0 {
		return nil, int(code), errors.New(result.Get("error_msg").String())
	}
	reply := new(ImageFaceResult)
	er = json.Unmarshal(bts, reply)
	return reply, 0, er
}

func updateUserFace(info *FaceAddReq) (*ImageFaceResult, error) {
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

func createFaceGroup(uid string) error {
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

func removeFaceGroup(uid string) error {
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

func getFaceGroups() ([]string, error) {
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
	arr := result.Get("result.group_id_list").Array()
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

func getUserCountByGroup(group string) int {
	list, err := getUsersByGroup(group)
	if err != nil {
		return -1
	}
	return len(list)
}

func getFacesByGroup(group, user string) (*FaceListResp, error) {
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

func getFaceMetas(group, user string) (*UserListResp, error) {
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

func removeFace(log uint64, face, user, group string) error {
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

func CheckFaceGroup(group string) error {
	list, er := getFaceGroups()
	if er != nil {
		return er
	}
	if tool.HasItem(list, group) {
		return nil
	}
	return createFaceGroup(group)
}
