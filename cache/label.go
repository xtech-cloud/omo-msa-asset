package cache

import (
	"errors"
	pb "github.com/xtech-cloud/omo-msp-asset/proto/asset"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.asset/proxy/nosql"
	"omo.msa.asset/tool"
	"strings"
	"time"
)

type LabelInfo struct {
	ID       uint64 `json:"-"`
	Created  int64
	Updated  int64
	UID      string `json:"uid"`
	Creator  string
	Operator string

	Type uint8

	Name   string
	Remark string
	Scene  string
}

func (mine *cacheContext) CreateLabel(in *pb.ReqLabelAdd) (*LabelInfo, error) {
	db := new(nosql.Label)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetLabelNextID()
	db.Created = time.Now().Unix()
	db.Creator = in.Operator
	db.Name = in.Name
	db.Remark = in.Remark
	db.Scene = in.Scene
	if len(db.Scene) < 1 {
		db.Scene = DefaultScene
	}
	db.Type = uint8(in.Type)
	err := nosql.CreateLabel(db)
	if err == nil {
		info := new(LabelInfo)
		info.initInfo(db)
		return info, nil
	}
	return nil, err
}

func (mine *cacheContext) RemoveLabel(uid, operator string) error {
	num := nosql.GetLabelChildrenCount(uid)
	if num > 0 {
		return errors.New("the folder not empty")
	}
	return nosql.RemoveLabel(uid, operator)
}

func (mine *cacheContext) GetLabel(uid string) (*LabelInfo, error) {
	db, err := nosql.GetLabel(uid)
	if err != nil {
		return nil, err
	}
	info := new(LabelInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) HadLabel(name string) (bool, error) {
	db, err := nosql.GetLabelByName(name)
	if err != nil && !strings.Contains(err.Error(), "no documents") {
		return false, err
	}
	if db != nil {
		return true, nil
	}
	return false, nil
}

func (mine *cacheContext) GetLabelsByScene(uid string) ([]*LabelInfo, error) {
	if len(uid) < 1 {
		uid = DefaultScene
	}
	dbs, err := nosql.GetLabelsByScene(uid)
	if err != nil {
		return nil, err
	}
	list := make([]*LabelInfo, 0, len(dbs))
	for _, db := range dbs {
		info := new(LabelInfo)
		info.initInfo(db)
		list = append(list, info)
	}

	return list, nil
}

func (mine *cacheContext) GetLabelsByQuote(uid string) ([]*LabelInfo, error) {
	all, err := nosql.GetLabelsByScene(DefaultScene)
	if err != nil {
		return nil, err
	}
	assets, er := nosql.GetAssetsByQuote(uid, 0, 0)
	if er != nil {
		return nil, er
	}
	arr := make([]string, 0, len(assets))
	for _, asset := range assets {
		for _, tag := range asset.Tags {
			if !tool.HasItem(arr, tag) {
				arr = append(arr, tag)
			}
		}
	}
	list := make([]*LabelInfo, 0, len(all))
	for _, label := range all {
		if tool.HasItem(arr, label.Name) {
			tmp := new(LabelInfo)
			tmp.initInfo(label)
			list = append(list, tmp)
		}
	}

	return list, nil
}

func (mine *LabelInfo) initInfo(db *nosql.Label) {
	mine.ID = db.ID
	mine.UID = db.UID.Hex()
	mine.Created = db.Created
	mine.Updated = db.Updated
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Scene = db.Scene
	mine.Name = db.Name
	mine.Type = db.Type
	mine.Remark = db.Remark
}

func (mine *LabelInfo) GetChildCount() uint32 {
	num := nosql.GetLabelChildrenCount(mine.UID)
	return uint32(num)
}

func (mine *LabelInfo) UpdateBase(name, remark, operator string) error {
	had, er := cacheCtx.HadLabel(name)
	if er != nil {
		return er
	}
	if had {
		return errors.New("the name had exited")
	}

	err := nosql.UpdateLabelBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
		mine.Updated = time.Now().Unix()
	}
	return err
}
