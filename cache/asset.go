package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.asset/proxy/nosql"
	"time"
)

type AssetInfo struct {
	Type uint8
	Size uint64

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
	Thumbs []*ThumbInfo
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

	list,err := nosql.GetThumbsByAsset(mine.UID)
	if err == nil {
		mine.Thumbs = make([]*ThumbInfo, 0, len(list)+5)
		for _, thumb := range list {
			info := new(ThumbInfo)
			info.initInfo(thumb)
			mine.Thumbs = append(mine.Thumbs, info)
		}
	}else{
		mine.Thumbs = make([]*ThumbInfo, 0, 10)
	}
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

func (mine *AssetInfo)HadThumbByFace(face string) bool {
	for _, thumb := range mine.Thumbs {
		if thumb.Face == face {
			return true
		}
	}
	return false
}

func (mine *AssetInfo)GetThumbByFace(face string) *ThumbInfo {
	for _, thumb := range mine.Thumbs {
		if thumb.Face == face {
			return thumb
		}
	}
	return nil
}

func (mine *AssetInfo)GetThumb(uid string) *ThumbInfo {
	for _, thumb := range mine.Thumbs {
		if thumb.UID == uid {
			return thumb
		}
	}
	return nil
}

func (mine *AssetInfo)hadThumb(uid string) bool {
	for _, thumb := range mine.Thumbs {
		if thumb.UID == uid {
			return true
		}
	}
	return false
}

func (mine *AssetInfo)RemoveThumb(uid, operator string) error {
	if !mine.hadThumb(uid) {
		return nil
	}
	err := nosql.RemoveThumb(uid, operator)
	if err == nil {
		for i := 0;i < len(mine.Thumbs);i += 1 {
			if mine.Thumbs[i].UID == uid {
				mine.Thumbs = append(mine.Thumbs[:i], mine.Thumbs[i+1:]...)
				break
			}
		}
	}
	return err
}

func (mine *AssetInfo)CreateThumb(face, url, operator string) (*ThumbInfo,error) {
	t := mine.GetThumbByFace(face)
	if t != nil {
		return t, nil
	}
	db := new(nosql.Thumb)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetAssetNextID()
	db.CreatedTime = time.Now()
	db.Creator = operator
	db.Operator = operator
	db.FaceID = face
	db.URL = url
	db.Asset = mine.UID
	err := nosql.CreateThumb(db)
	if err == nil {
		info := new(ThumbInfo)
		info.initInfo(db)
		mine.Thumbs = append(mine.Thumbs, info)
		return info, nil
	}
	return nil, err
}

