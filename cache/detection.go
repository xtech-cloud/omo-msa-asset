package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"omo.msa.asset/config"
	"omo.msa.asset/proxy"
	"strings"
)

type FaceResponse struct {
	Code      int          `json:"error_code"`
	Message   string       `json:"error_msg"`
	LogID     int          `json:"log_id"`
	Timestamp int          `json:"timestamp"`
	Cached    int          `json:"cached"`
	Result    *FacesResult `json:"result"`
}

type FacesResult struct {
	Number int          `json:"face_num"`
	List   []*ImageFace `json:"face_list"`
}

type ImageFace struct {
	Token       string             `json:"face_token"`
	Probability int                `json:"face_probability"`
	Age         int                `json:"age"`
	Location    proxy.LocationInfo `json:"location"`
	Angle       *AngleInfo         `json:"angle"`
	Gender      *ProbabilityInfo   `json:"gender"`
	Glasses     *ProbabilityInfo   `json:"glasses"`
	Shape       *ProbabilityInfo   `json:"face_shape"`
	Type        *ProbabilityInfo   `json:"face_type"`
}

type AngleInfo struct {
	Yaw   float32 `json:"yaw"`
	Pitch float32 `json:"pitch"`
	Roll  float32 `json:"roll"`
}

type ProbabilityInfo struct {
	Type        string `json:"type"`
	Probability string `json:"probability"`
}

/**
 * 使用 AK，SK 生成鉴权签名（Access Token）
 * @return string 鉴权签名信息（Access Token）
 */
func getDetectAccessToken() (string, error) {
	url := "https://aip.baidubce.com/oauth/2.0/token"
	postData := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", config.Schema.Detection.AccessKey, config.Schema.Detection.SecretKey)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(postData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	accessTokenObj := map[string]string{}
	json.Unmarshal([]byte(body), &accessTokenObj)
	return accessTokenObj["access_token"], nil
}

func httpPost(address, data string) ([]byte, error) {
	payload := strings.NewReader(data)
	client := &http.Client{}
	req, err := http.NewRequest("POST", address, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

//检测图片中的人脸
func DetectFaces(img string) (*FaceResponse, error) {
	token, er := getDetectAccessToken()
	if er != nil {
		return nil, er
	}
	addr := fmt.Sprintf("%s?access_token=%s", config.Schema.Detection.Address, token)
	data := fmt.Sprintf(`{"image":"%s","image_type":"URL","face_type":"LIVE", "max_face_num":50,"liveness_control":"NONE", "face_field":"age,gender,glasses,face_shape"}`, img)
	msg, er := httpPost(addr, data)
	if er != nil {
		return nil, er
	}
	reply := new(FaceResponse)
	er = json.Unmarshal(msg, reply)
	return reply, er
}
