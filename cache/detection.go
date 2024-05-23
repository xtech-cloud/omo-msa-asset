package cache

import (
	"encoding/json"
	"errors"
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
	Token       string  `json:"face_token"`
	Probability float32 `json:"face_probability"`
	Age         int     `json:"age"`
	//Data        string             `json:"corp_image_base64"`
	Quality  ImageQuality       `json:"quality"`
	Location proxy.LocationInfo `json:"location"`
	Angle    *AngleInfo         `json:"angle"`
	Gender   *ProbabilityInfo   `json:"gender"`     //性别，male:男性 female:女性
	Glasses  *ProbabilityInfo   `json:"glasses"`    //是否带眼镜，none:无眼镜，common:普通眼镜，sun:墨镜
	Shape    *ProbabilityInfo   `json:"face_shape"` //情绪，angry:愤怒 disgust:厌恶 fear:恐惧 happy:高兴 sad:伤心 surprise:惊讶 neutral:无表情 pouty: 撅嘴 grimace:鬼脸
	Type     *ProbabilityInfo   `json:"face_type"`  //真实人脸、卡通人脸；human: 真实人脸 cartoon: 卡通人脸
	Mask     *ProbabilityInfo   `json:"mask"`       //口罩识别，取值0或1； 0代表没戴口罩 1 代表戴口罩
	Emotion  *ProbabilityInfo   `json:"emotion"`    //表情
}

type ImageQuality struct {
	Blur float32 `json:"blur"`
}

type AngleInfo struct {
	Yaw   float32 `json:"yaw"`
	Pitch float32 `json:"pitch"`
	Roll  float32 `json:"roll"`
}

type ProbabilityInfo struct {
	Type        string  `json:"type"`
	Probability float32 `json:"probability"`
}

type ImageMatchReq struct {
	Image           string `json:"image"`
	Type            string `json:"image_type"`
	FaceType        string `json:"face_type"`
	QualityControl  string `json:"quality_control"`
	LivenessControl string `json:"liveness_control"`
}

type ImageMatchResp struct {
	Score    float32      `json:"score"`
	FaceList []*FaceBrief `json:"face_list"`
}

type FaceBrief struct {
	Token string `json:"face_token"`
	Stamp string `json:"ctime"`
}

/**
 * 使用 AK，SK 生成鉴权签名（Access Token）
 * @return string 鉴权签名信息（Access Token）
 */
func getDetectAccessToken() (string, error) {
	url := config.Schema.Detection.Token
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
func detectFaces(img string) (*FaceResponse, error) {
	token, er := getDetectAccessToken()
	if er != nil {
		return nil, er
	}
	addr := fmt.Sprintf("%s?access_token=%s", config.Schema.Detection.Address, token)
	data := fmt.Sprintf(`{"image":"%s","image_type":"URL","face_type":"LIVE", "max_face_num":50,"liveness_control":"NONE", "face_field":"age,gender,glasses,face_shape,quality"}`, img)
	bts, er := httpPost(addr, data)
	if er != nil {
		return nil, er
	}
	fmt.Println(string(bts))
	reply := new(FaceResponse)
	er = json.Unmarshal(bts, reply)
	if reply.Code > 0 {
		return nil, errors.New(reply.Message)
	}
	return reply, er
}

//对比人脸
func compareImages(one *ImageMatchReq, two *ImageMatchReq) (*ImageMatchResp, error) {
	token, er := getDetectAccessToken()
	if er != nil {
		return nil, er
	}
	addr := fmt.Sprintf("%s?access_token=%s", config.Schema.Detection.Match, token)
	data := fmt.Sprintf(`{"image":"%s","image_type":"URL","face_type":"LIVE", "max_face_num":50,"liveness_control":"NONE", "face_field":"age,gender,glasses,face_shape"}`, one.Image)
	msg, er := httpPost(addr, data)
	if er != nil {
		return nil, er
	}
	reply := new(ImageMatchResp)
	er = json.Unmarshal(msg, reply)
	return reply, er
}
