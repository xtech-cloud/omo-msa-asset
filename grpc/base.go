package grpc

import (
	"encoding/json"
	"github.com/micro/go-micro/v2/logger"
	pb "github.com/xtech-cloud/omo-msp-asset/proto/asset"
)

func outError(name, msg string, code pb.ResultStatus) *pb.ReplyStatus {
	logger.Warnf("[error.%s]:code = %d, msg = %s", name, code, msg)
	tmp := &pb.ReplyStatus{
		Code: uint32(code),
		Error: msg,
	}
	return tmp
}

func inLog(name, data interface{})  {
	bytes, _ := json.Marshal(data)
	msg := byteString(bytes)
	logger.Infof("[in.%s]:data = %s", name, msg)
}

func outLog(name, data interface{}) *pb.ReplyStatus {
	//bytes, _ := json.Marshal(data)
	//msg := byteString(bytes)
	//logger.Infof("[out.%s]:data = %s", name, msg)
	tmp := &pb.ReplyStatus{
		Code: 0,
		Error: "",
	}
	return tmp
}

func byteString(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}