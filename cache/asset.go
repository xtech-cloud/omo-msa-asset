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

const (
	OwnerTypePerson = 1
	OwnerTypeUnit = 0
)

const UP_QINIU = "qiniu"

const (
	StatusIdle uint8 = 0
	StatusHide uint8 = 1
)

type AssetInfo struct {
	Type uint8
	Status uint8
	Size uint64
	Width uint32
	Height uint32
	ID         uint64 `json:"-"`
	UID        string `json:"uid"`
	Name       string `json:"name"`
	Remark    string
	Meta  string
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
	Weight uint32
	CreateTime time.Time
	UpdateTime time.Time
}

func (mine *cacheContext)CreateAsset(info *AssetInfo) error {
	db := new(nosql.Asset)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetAssetNextID()
	db.CreatedTime = time.Now()
	db.Creator = info.Creator
	db.Operator = info.Operator
	db.Name = info.Name

	db.Owner = info.Owner
	db.Type = info.Type
	db.Size = info.Size
	db.UUID = info.UUID
	db.Format = info.Format
	db.MD5 = info.MD5
	db.Version = info.Version
	db.Language = info.Language
	db.Snapshot = info.Snapshot
	db.Small = info.Small
	db.Width = info.Width
	db.Height = info.Height
	db.Meta = info.Meta
	db.Weight = 0
	db.Status = StatusIdle
	err := nosql.CreateAsset(db)
	if err == nil {
		info.UID = db.UID.Hex()
		info.ID = db.ID
		info.CreateTime = db.CreatedTime
	}
	return err
}

func (mine *cacheContext)GetAsset(uid string) *AssetInfo {
	db,err := nosql.GetAsset(uid)
	if err == nil {
		info := new(AssetInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *cacheContext)GetAssetsByOwner(uid string) []*AssetInfo {
	array,err := nosql.GetAssetsByOwner(uid)
	if err != nil {
		return make([]*AssetInfo, 0, 1)
	}
	list := make([]*AssetInfo, 0, len(array))
	for _, asset := range array {
		info := new(AssetInfo)
		info.initInfo(asset)
		list = append(list, info)
	}
	return list
}

func (mine *AssetInfo)initInfo(db *nosql.Asset)  {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.Remark = db.Remark
	mine.Meta = db.Meta
	mine.Weight = db.Weight

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
	mine.Weight = db.Weight
	mine.Status = db.Status
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
	err := nosql.RemoveAsset(mine.UID)
	if err == nil {
		_ = deleteContentFromCloud(mine.UUID)
		if len(mine.Snapshot) > 2 {
			_ = deleteContentFromCloud(mine.Snapshot)
		}
		if len(mine.Small) > 2 {
			_ = deleteContentFromCloud(mine.Small)
		}
	}
	return err
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
		mine.Small = small
		mine.Operator = operator
	}
	return err
}

func (mine *AssetInfo)UpdateBase(operator, name, remark string) error {
	err := nosql.UpdateAssetBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *AssetInfo)UpdateMeta(operator, meta string) error {
	err := nosql.UpdateAssetMeta(mine.UID, meta, operator)
	if err == nil {
		mine.Meta = meta
		mine.Operator = operator
	}
	return err
}

func (mine *AssetInfo)UpdateWeight(weight uint32, operator string) error {
	err := nosql.UpdateAssetWeight(mine.UID, operator, weight)
	if err == nil {
		mine.Weight = weight
		mine.Operator = operator
	}
	return err
}

func (mine *AssetInfo)UpdateStatus(st uint8, operator string) error {
	err := nosql.UpdateAssetStatus(mine.UID, operator, st)
	if err == nil {
		mine.Status = st
		mine.Operator = operator
	}
	return err
}

func (mine *AssetInfo)UpdateType(st uint8, operator string) error {
	err := nosql.UpdateAssetType(mine.UID, operator, st)
	if err == nil {
		mine.Type = st
		mine.Operator = operator
	}
	return err
}

func (mine *AssetInfo)UpdateLanguage(lan, operator string) error {
	err := nosql.UpdateAssetLanguage(mine.UID, operator, lan)
	if err == nil {
		mine.Language = lan
		mine.Operator = operator
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

