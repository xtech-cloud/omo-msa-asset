package cache

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"omo.msa.asset/config"
	"omo.msa.asset/proxy"
)

type FaceGroup struct {
	UID   string      `json:"id"`
	Users []*FaceUser `json:"users"`
}

type FaceUser struct {
	UID string `json:"id"`
}

func AddUserFace(group, user, meta string) (string, *proxy.LocationInfo, error) {
	token, er := getDetectAccessToken()
	if er != nil {
		return "", nil, er
	}
	addr := fmt.Sprintf("%s?access_token=%s", config.Schema.Detection.Face.Add, token)
	data := fmt.Sprintf(`{"group_id":"%s","user_id":"%s","image":"%s","image_type":"BASE64","quality_control":"NONE","liveness_control":"NONE"}`, group, user, meta)
	bts, er := httpPost(addr, data)
	if er != nil {
		return "", nil, er
	}
	result := gjson.ParseBytes(bts)
	code := result.Get("error_code").Int()
	if code != 0 {
		return "", nil, errors.New(result.Get("error_msg").String())
	}
	key := result.Get("result/face_token").String()
	re := result.Get("result/location")
	loc := new(proxy.LocationInfo)
	loc.Width = int(re.Get("width").Int())
	loc.Height = int(re.Get("height").Int())
	loc.Top = float32(re.Get("top").Float())
	loc.Left = float32(re.Get("left").Float())
	return key, loc, nil
}

func CreateGroup(uid string) error {
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

func GetGroups() ([]string, error) {
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

func GetUsers(group string) ([]string, error) {
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
