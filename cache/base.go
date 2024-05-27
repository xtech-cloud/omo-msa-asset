package cache

import (
	"errors"
	"github.com/micro/go-micro/v2/logger"
	"github.com/qiniu/api.v7/v7/auth"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/cdn"
	"github.com/qiniu/api.v7/v7/storage"
	"go.uber.org/zap"
	"omo.msa.asset/config"
	"omo.msa.asset/proxy/nosql"
)

const DefaultScene = "system"

type cacheContext struct {
}

var cacheCtx *cacheContext

func InitData() error {
	cacheCtx = &cacheContext{}

	err := nosql.InitDB(config.Schema.Database.IP, config.Schema.Database.Port, config.Schema.Database.Name, config.Schema.Database.Type)
	if err == nil {
		num := nosql.GetAssetCount()
		count := nosql.GetThumbCount()
		logger.Infof("the asset count = %d and the thumb count = %d", num, count)
		//nosql.CheckTimes()
	}
	go CheckFaceGroup()
	return err
}

func Context() *cacheContext {
	return cacheCtx
}

func PublishSystemAssets() {
	dbs, _ := nosql.GetAssetsByOwner("system")
	for _, db := range dbs {
		if db.Status != StatusPublish {
			_ = nosql.UpdateAssetStatus(db.UID.Hex(), db.Operator, StatusPublish)
		}
	}
}

func TestDetectFaces() {
	//url := "https://rdpdown.suii.cn/000c0f54-3dd7-40c6-aa2b-f67378947978"
	//url := "https://rdpdown.suii.cn/00278e27e030ac05"
	asset := cacheCtx.GetAsset("6654571ce182e1766ab4b012")
	validateAsset(asset)
}

func checkPage[T any](page, number uint32, all []T) (uint32, uint32, []T) {
	if len(all) < 1 {
		return 0, 0, make([]T, 0, 1)
	}
	if number < 1 {
		number = 10
	}
	total := uint32(len(all))
	if len(all) <= int(number) {
		return total, 1, all
	}
	//array := reflect.ValueOf(all)
	maxPage := total/number + 1
	if page < 1 {
		return total, maxPage, all
	}

	var start = (page - 1) * number
	var end = start + number
	if end > total {
		end = total
	}
	list := make([]T, 0, number)
	list = append(all[start:end])
	return total, maxPage, list
}

func (mine *cacheContext) GetThumb(uid string) *ThumbInfo {
	db, _ := nosql.GetThumb(uid)
	if db != nil {
		info := new(ThumbInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *cacheContext) GetThumbsByOwner(uid string) []*ThumbInfo {
	array, err := nosql.GetThumbsByOwner(uid)
	if err == nil {
		list := make([]*ThumbInfo, 0, len(array))
		for _, item := range array {
			info := new(ThumbInfo)
			info.initInfo(item)
			list = append(list, info)
		}
		return list
	} else {
		list := make([]*ThumbInfo, 0, 1)
		return list
	}
}

func (mine *cacheContext) GetUpToken(key string) string {
	cof := config.Schema.Storage
	mac := auth.New(cof.AccessKey, cof.SecretKey)
	// 设置上传凭证有效期
	putPolicy := storage.PutPolicy{
		Scope: config.Schema.Storage.Bucket,
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","size":$(fsize),"type":"$(mimeType)", 
		"img":$(imageInfo), "uuid":"$(uuid)", "bucket":"$(bucket)","name":"$(fname)"}`,
	}
	if len(key) > 2 {
		putPolicy.Scope = config.Schema.Storage.Bucket + ":" + key
	}
	putPolicy.Expires = uint64(config.Schema.Storage.Expire) //有效期
	auth.New(cof.AccessKey, cof.SecretKey)
	return putPolicy.UploadToken(mac)
}

func RefreshCDN(url string) bool {
	mac := qbox.NewMac(config.Schema.Storage.AccessKey, config.Schema.Storage.SecretKey)
	cdnManager := cdn.NewCdnManager(mac)

	urlsToRefresh := []string{
		url,
	}
	_, err := cdnManager.RefreshUrls(urlsToRefresh)
	if err != nil {
		logger.Warn("cache: refresh cdn failed from qiniu cache!!!", zap.String("url", url))
		return false
	}
	return true
}

func deleteContentFromCloud(key string) error {
	if len(key) < 1 {
		return errors.New("cache: the key is empty")
	}
	mac := qbox.NewMac(config.Schema.Storage.AccessKey, config.Schema.Storage.SecretKey)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	//cfg.Zone=&storage.ZoneHuabei
	bucketManager := storage.NewBucketManager(mac, &cfg)
	err := bucketManager.Delete(config.Schema.Storage.Bucket, key)
	if err != nil {
		return err
	}
	return nil
}

func GetContentFromCloud(key string) *storage.FileInfo {
	mac := qbox.NewMac(config.Schema.Storage.AccessKey, config.Schema.Storage.SecretKey)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	//cfg.Zone=&storage.ZoneHuabei
	bucketManager := storage.NewBucketManager(mac, &cfg)
	fileInfo, err := bucketManager.Stat(config.Schema.Storage.Bucket, key)
	if err == nil {
		return &fileInfo
	}
	logger.Warn("cache: check file info failed from qiniu cache!!!", zap.String("key", key))
	return nil
}
