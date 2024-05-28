package cache

import (
	"bytes"
	"context"
	"errors"
	"github.com/qiniu/api.v7/v7/storage"
)

type MyPutRet struct {
	UUID   string
	Key    string
	Hash   string
	Size   int
	Bucket string
	Name   string
	Type   string
}

func uploadToQiNiu(key string, data []byte) (string, int64, error) {
	if data == nil {
		return "", 0, errors.New("the data is nil")
	}
	if key == "" {
		return "", 0, errors.New("the key is empty")
	}
	token := cacheCtx.GetUpToken(key)
	upToken := token
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Region = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	uploader := storage.NewFormUploader(&cfg)
	ret := MyPutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"key": key,
		},
	}

	dataLen := int64(len(data))
	err := uploader.Put(context.Background(), &ret, upToken, key, bytes.NewReader(data), dataLen, &putExtra)
	if err != nil {
		return "", 0, err
	}
	return ret.Hash, dataLen, nil
}
