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
