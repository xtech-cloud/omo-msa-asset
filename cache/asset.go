package cache

import (
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
	CreateTime time.Time
	UpdateTime time.Time
}

func GetAsset(uid string) *AssetInfo {
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
		owner := GetOwner(info.Owner)
		owner.AddAsset(info)
		return info
	}
	return nil
}

func RemoveAsset(uid, operator string) error {
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
}

func (mine *AssetInfo)Remove(operator string) error {
	return nosql.RemoveAsset(mine.UID, operator)
}

