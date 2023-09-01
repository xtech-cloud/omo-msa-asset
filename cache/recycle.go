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

type RecycleInfo struct {
	Type      uint8
	Status    uint8
	Width     uint32
	Height    uint32
	Weight    uint32
	Size      uint64
	ID        uint64 `json:"-"`
	UID       string `json:"uid"`
	Name      string `json:"name"`
	Remark    string
	Meta      string
	Creator   string
	Operator  string
	Scavenger string

	Owner    string
	UUID     string
	Version  string
	Format   string
	MD5      string
	Language string

	// 快照，中图
	Snapshot string
	// 封面小图
	Small string

	CreateTime time.Time
	UpdateTime time.Time
	Links      []string
}

func (mine *cacheContext) CreateRecycle(info *RecycleInfo) error {
	db := new(nosql.Recycle)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetRecycleNextID()
	db.CreatedTime = time.Now()
	db.Creator = info.Creator
	db.Operator = info.Operator
	db.Name = info.Name
	db.Remark = info.Remark
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
	db.Links = info.Links
	err := nosql.CreateRecycle(db)
	if err == nil {
		info.UID = db.UID.Hex()
		info.ID = db.ID
		info.CreateTime = db.CreatedTime
	}
	return err
}

func (mine *cacheContext) GetRecycle(uid string) *RecycleInfo {
	db, err := nosql.GetRecycle(uid)
	if err == nil {
		info := new(RecycleInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *cacheContext) GetRecyclesByOwner(uid string) []*RecycleInfo {
	array, err := nosql.GetRecyclesByOwner(uid)
	if err != nil {
		return make([]*RecycleInfo, 0, 1)
	}
	list := make([]*RecycleInfo, 0, len(array))
	for _, asset := range array {
		info := new(RecycleInfo)
		info.initInfo(asset)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetRecyclesByCreator(uid string) []*RecycleInfo {
	array, err := nosql.GetRecyclesByCreator(uid)
	if err != nil {
		return make([]*RecycleInfo, 0, 1)
	}
	list := make([]*RecycleInfo, 0, len(array))
	for _, asset := range array {
		info := new(RecycleInfo)
		info.initInfo(asset)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetRecyclesByException(page, number uint32) (uint32, uint32, []*RecycleInfo) {
	array, err := nosql.GetAllRecycles()
	if err != nil {
		return 0, 0, make([]*RecycleInfo, 0, 1)
	}
	list := make([]*RecycleInfo, 0, len(array))
	for _, asset := range array {
		if asset.Creator == asset.Owner {
			info := new(RecycleInfo)
			info.initInfo(asset)
			list = append(list, info)
		}
	}
	return checkPage(page, number, list)
}

func (mine *RecycleInfo) initInfo(db *nosql.Recycle) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.CreateTime = db.CreatedTime
	mine.Scavenger = db.Scavenger
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
	mine.Links = db.Links
}

func (mine *RecycleInfo) Remove() error {
	err := nosql.RemoveRecycle(mine.UID)
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

func (mine *RecycleInfo) getURL(key string, cdn bool) string {
	if len(key) < 2 {
		return ""
	}
	if strings.Contains(key, "http") {
		return key
	}
	domain := config.Schema.Storage.Domain
	if !cdn {
		domain = config.Schema.Storage.Source
	}
	if config.Schema.Storage.Type == UP_QINIU {
		if config.Schema.Storage.ACM > 0 {
			mac := qbox.NewMac(config.Schema.Storage.AccessKey, config.Schema.Storage.SecretKey)
			return storage.MakePrivateURL(mac, domain, key, config.Schema.Storage.Period)
		} else {
			return storage.MakePublicURL(domain, key)
		}
	} else {
		return mine.UID
	}
}

func (mine *RecycleInfo) URL() string {
	return mine.getURL(mine.UUID, true)
}

func (mine *RecycleInfo) SourceURL() string {
	if mine.Snapshot != "" {
		return mine.getURL(mine.Snapshot, false)
	}
	return mine.getURL(mine.UUID, false)
}

func (mine *RecycleInfo) SnapshotURL() string {
	return mine.getURL(mine.Snapshot, true)
}

func (mine *RecycleInfo) SmallImageURL() string {
	return mine.getURL(mine.Small, false)
}

func (mine *RecycleInfo) HadThumbByFace(face string) bool {
	info := mine.GetThumbByFace(face)
	if info == nil {
		return false
	}
	return true
}

func (mine *RecycleInfo) GetThumbByFace(face string) *ThumbInfo {
	db, err := nosql.GetThumbByFace(mine.UID, face)
	if err != nil {
		return nil
	}
	tmp := new(ThumbInfo)
	tmp.initInfo(db)
	return tmp
}

func (mine *RecycleInfo) GetThumb(uid string) *ThumbInfo {
	db, err := nosql.GetThumb(uid)
	if err != nil {
		return nil
	}
	tmp := new(ThumbInfo)
	tmp.initInfo(db)
	return tmp
}

func (mine *RecycleInfo) hadThumb(uid string) bool {
	info := mine.GetThumb(uid)
	if info == nil {
		return false
	}
	return true
}

func (mine *RecycleInfo) RemoveThumb(uid, operator string) error {
	if !mine.hadThumb(uid) {
		return nil
	}
	return nosql.RemoveThumb(uid, operator)
}
