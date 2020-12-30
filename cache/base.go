package cache

import (
	"omo.msa.asset/config"
	"omo.msa.asset/proxy/nosql"
)

type cacheContext struct {
	owners []*OwnerInfo
}

var cacheCtx *cacheContext

func InitData() error {
	cacheCtx = &cacheContext{}
	cacheCtx.owners = make([]*OwnerInfo, 0, 100)

	err := nosql.InitDB(config.Schema.Database.IP, config.Schema.Database.Port, config.Schema.Database.Name, config.Schema.Database.Type)
	return err
}

func Context() *cacheContext {
	return cacheCtx
}

func (mine *cacheContext)GetThumb(uid string) *ThumbInfo {
	for _, owner := range mine.owners {
		for _, asset := range owner.assets {
			for _, thumb := range asset.Thumbs {
				if thumb.UID == uid {
					return thumb
				}
			}
		}
	}
	return nil
}

func (mine *cacheContext) GetThumbsByOwner(uid string) []*ThumbInfo {
	array,err := nosql.GetThumbsByOwner(uid)
	if err == nil {
		list := make([]*ThumbInfo, 0, len(array))
		for _, item := range array {
			info := new(ThumbInfo)
			info.initInfo(item)
			list = append(list, info)
		}
		return list
	}else{
		list := make([]*ThumbInfo, 0, 1)
		return list
	}
}
