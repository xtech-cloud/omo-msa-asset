package cache

import (
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.asset/config"
	"omo.msa.asset/proxy/nosql"
	"strings"
	"time"
)

const UP_QINIU = "qiniu"

type AssetInfo struct {
	Type uint8
	Size uint64
	Width uint32
	Height uint32
	ID         uint64 `json:"-"`
	UID        string `json:"uid"`
	Name       string `json:"name"`
	Creator string
	Operator string

	Owner string
	UUID string
	Version string
	Format string
	MD5 string
	Language string
	// 快照，中图
	Snapshot string
	// 封面小图
	Small string
	CreateTime time.Time
	UpdateTime time.Time
}

func (mine *cacheContext)GetAsset(uid string) *AssetInfo {
	for i := 0;i < len(cacheCtx.owners);i += 1{
		info := cacheCtx.owners[i].GetAsset(uid)
		if info != nil {
			return info
		}
	}
	db,err := nosql.GetAsset(uid)
	if err == nil {
		info := new(AssetInfo)
		info.initInfo(db)
		owner := cacheCtx.GetOwner(info.Owner)
		owner.AddAsset(info)
		return info
	}
	return nil
}

func (mine *cacheContext)RemoveAsset(uid, operator string) error {
	err := nosql.RemoveAsset(uid, operator)
	if err == nil {
		for i := 0;i < len(cacheCtx.owners);i += 1{
			if cacheCtx.owners[i].HadAsset(uid) {
				cacheCtx.owners[i].deleteAsset(uid)
				break
			}
		}
	}
	return err
}

func (mine *AssetInfo)initInfo(db *nosql.Asset)  {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name

	mine.Size = db.Size
	mine.UUID = db.UUID
	mine.Type = db.Type
	mine.Owner = db.Owner
	mine.Version = db.Version
	mine.MD5 = db.MD5
	mine.Format = db.Format
	mine.Language = db.Language
	mine.Snapshot = db.Snapshot
	mine.Small = db.Small

	mine.Width = db.Width
	mine.Height = db.Height
}

func (mine *AssetInfo)GetThumbs() ([]*ThumbInfo,error) {
	array,err := nosql.GetThumbsByAsset(mine.UID)
	if err != nil {
		return nil,err
	}
	list := make([]*ThumbInfo, 0, len(array))
	for _, thumb := range array {
		tmp := new(ThumbInfo)
		tmp.initInfo(thumb)
		list = append(list, tmp)
	}
	return list,nil
}

func (mine *AssetInfo)Remove(operator string) error {
	return nosql.RemoveAsset(mine.UID, operator)
}

func (mine *AssetInfo)UpdateSnapshot(operator, snapshot string) error {
	err := nosql.UpdateAssetSnapshot(mine.UID, snapshot,operator)
	if err == nil {
		mine.Snapshot = snapshot
	}
	return err
}

func (mine *AssetInfo)UpdateSmall(operator, small string) error {
	err := nosql.UpdateAssetSmall(mine.UID, small, operator)
	if err == nil {
		mine.Small = operator
	}
	return err
}

func (mine *AssetInfo)getURL(key string) string {
	if len(key) < 2 {
		return ""
	}
	if strings.Contains(key, "http") {
		return key
	}
	if config.Schema.Storage.Type == UP_QINIU {
		if config.Schema.Storage.ACM > 0 {
			mac := qbox.NewMac(config.Schema.Storage.AccessKey, config.Schema.Storage.SecretKey)
			return storage.MakePrivateURL(mac, config.Schema.Storage.Domain, key, config.Schema.Storage.Period)
		}else{
			return storage.MakePublicURL(config.Schema.Storage.Domain, key)
		}
	} else {
		return mine.UID
	}
}

func (mine *AssetInfo) URL() string {
	return mine.getURL(mine.UUID)
}

func (mine *AssetInfo) SnapshotURL() string {
	return mine.getURL(mine.Snapshot)
}

func (mine *AssetInfo) SmallImageURL() string {
	return mine.getURL(mine.Small)
}

func (mine *AssetInfo)HadThumbByFace(face string) bool {
	info := mine.GetThumbByFace(face)
	if info == nil {
		return false
	}
	return true
}

func (mine *AssetInfo)GetThumbByFace(face string) *ThumbInfo {
	db,err := nosql.GetThumbByFace(mine.UID, face)
	if err != nil {
		return nil
	}
	tmp := new(ThumbInfo)
	tmp.initInfo(db)
	return tmp
}

func (mine *AssetInfo)GetThumb(uid string) *ThumbInfo {
	db,err := nosql.GetThumb(uid)
	if err != nil {
		return nil
	}
	tmp := new(ThumbInfo)
	tmp.initInfo(db)
	return tmp
}

func (mine *AssetInfo)hadThumb(uid string) bool {
	info := mine.GetThumb(uid)
	if info == nil {
		return false
	}
	return true
}

func (mine *AssetInfo)RemoveThumb(uid, operator string) error {
	if !mine.hadThumb(uid) {
		return nil
	}
	return nosql.RemoveThumb(uid, operator)
}

func (mine *AssetInfo)CreateThumb(face, url, operator, owner string, score,similar,blur float32) (*ThumbInfo,error) {
	t := mine.GetThumbByFace(face)
	if t != nil {
		return t, nil
	}
	db := new(nosql.Thumb)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetThumbNextID()
	db.CreatedTime = time.Now()
	db.Creator = operator
	db.Operator = operator
	db.FaceID = face
	db.URL = url
	db.Asset = mine.UID
	db.Blur = blur
	db.Owner = owner
	db.Probably = score
	db.Similar = similar
	err := nosql.CreateThumb(db)
	if err == nil {
		info := new(ThumbInfo)
		info.initInfo(db)
		return info, nil
	}
	return nil, err
}

