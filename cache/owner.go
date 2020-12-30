package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.asset/proxy/nosql"
	"time"
)

const (
	OwnerTypePerson = 1
	OwnerTypeUnit = 0
)

type OwnerInfo struct {
	UID    string
	assets []*AssetInfo
}

func (mine *cacheContext)AllOwners() []*OwnerInfo {
	return cacheCtx.owners
}

func (mine *cacheContext)GetOwner(uid string) *OwnerInfo {
	for i := 0;i < len(cacheCtx.owners);i += 1{
		if cacheCtx.owners[i].UID == uid {
			return cacheCtx.owners[i]
		}
	}
	owner := new(OwnerInfo)
	owner.initInfo(uid)
	cacheCtx.owners = append(cacheCtx.owners, owner)
	return owner
}

func (mine *OwnerInfo)initInfo(owner string)  {
	mine.UID = owner
	array,err := nosql.GetAssetsByOwner(owner)
	if err == nil{
		mine.assets = make([]*AssetInfo, 0, len(array))
		for _, value := range array {
			t := new(AssetInfo)
			t.initInfo(value)
			mine.assets = append(mine.assets, t)
		}
	}else{
		mine.assets = make([]*AssetInfo, 0, 1)
	}
}

func (mine *OwnerInfo)Assets() []*AssetInfo {
	return mine.assets
}

func (mine *OwnerInfo)CreateAsset(info *AssetInfo) error {
	db := new(nosql.Asset)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetAssetNextID()
	db.CreatedTime = time.Now()
	db.Creator = info.Creator
	db.Operator = info.Operator
	db.Name = info.Name

	db.Owner = mine.UID
	db.Type = info.Type
	db.Size = info.Size
	db.UUID = info.UUID
	db.Format = info.Format
	db.MD5 = info.MD5
	db.Version = info.Version
	db.Language = info.Language
	db.Snapshot = info.Snapshot

	err := nosql.CreateAsset(db)
	if err == nil {
		info.initInfo(db)
		mine.assets = append(mine.assets, info)
	}
	return err
}

func (mine *OwnerInfo)AddAsset(info *AssetInfo) bool {
	if info == nil {
		return false
	}
	if mine.HadAsset(info.UID){
		return true
	}
	mine.assets = append(mine.assets, info)
	return true
}

func (mine *OwnerInfo)HadAsset(uid string) bool {
	for _, asset := range mine.assets {
		if asset.UID == uid {
			return true
		}
	}
	return false
}

func (mine *OwnerInfo)GetAsset(uid string) *AssetInfo {
	for i := 0;i < len(mine.assets);i += 1{
		if mine.assets[i].UID == uid {
			return mine.assets[i]
		}
	}
	return nil
}

func (mine *OwnerInfo)deleteAsset(uid string) {
	for i := 0;i < len(mine.assets);i += 1{
		if mine.assets[i].UID == uid {
			mine.assets = append(mine.assets[:i], mine.assets[i+1:]...)
			break
		}
	}
}

func (mine *OwnerInfo)RemoveAsset(uid, operator string) error {
	err := nosql.RemoveAsset(uid, operator)
	if err == nil {
		mine.deleteAsset(uid)
	}
	return err
}