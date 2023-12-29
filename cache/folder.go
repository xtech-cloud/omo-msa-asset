package cache

import (
	"errors"
	pb "github.com/xtech-cloud/omo-msp-asset/proto/asset"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.asset/proxy"
	"omo.msa.asset/proxy/nosql"
	"time"
)

type FolderInfo struct {
	ID       uint64 `json:"-"`
	Created  int64
	Updated  int64
	UID      string `json:"uid"`
	Creator  string
	Operator string

	Access uint8

	Name   string
	Remark string
	Parent string
	Scene  string
	Cover  string

	Tags     []string
	Users    []string
	Contents []*proxy.PairInfo
}

func (mine *cacheContext) CreateFolder(in *pb.ReqFolderAdd) (*FolderInfo, error) {
	db := new(nosql.Folder)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetFolderNextID()
	db.Created = time.Now().Unix()
	db.Creator = in.Operator
	db.Name = in.Name
	db.Remark = in.Remark
	db.Scene = in.Owner
	if db.Scene == "" {
		db.Scene = "system"
	}
	db.Parent = in.Parent
	db.Access = 0
	db.Cover = in.Cover
	db.Tags = in.Tags
	db.Users = in.Users
	db.Contents = make([]*proxy.PairInfo, 0, len(in.Contents))
	for _, content := range in.Contents {
		db.Contents = append(db.Contents, &proxy.PairInfo{
			Key:   content.Key,
			Value: content.Value,
			Count: content.Count,
		})
	}
	err := nosql.CreateFolder(db)
	if err == nil {
		info := new(FolderInfo)
		info.initInfo(db)
		return info, nil
	}
	return nil, err
}

func (mine *cacheContext) RemoveFolder(uid, operator string) error {
	num := nosql.GetFolderChildrenCount(uid)
	if num > 0 {
		return errors.New("the folder not empty")
	}
	return nosql.RemoveFolder(uid, operator)
}

func (mine *cacheContext) GetFolder(uid string) (*FolderInfo, error) {
	db, err := nosql.GetFolder(uid)
	if err != nil {
		return nil, err
	}
	info := new(FolderInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetFoldersByScene(uid string) ([]*FolderInfo, error) {
	dbs, err := nosql.GetFoldersByScene(uid)
	if err != nil {
		return nil, err
	}
	list := make([]*FolderInfo, 0, len(dbs))
	for _, db := range dbs {
		if db.Parent == "" {
			info := new(FolderInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}

	return list, nil
}

func (mine *cacheContext) GetFoldersByParent(uid string) ([]*FolderInfo, error) {
	dbs, err := nosql.GetFoldersByParent(uid)
	if err != nil {
		return nil, err
	}
	list := make([]*FolderInfo, 0, len(dbs))
	for _, db := range dbs {
		info := new(FolderInfo)
		info.initInfo(db)
		list = append(list, info)
	}

	return list, nil
}

func (mine *FolderInfo) initInfo(db *nosql.Folder) {
	mine.ID = db.ID
	mine.UID = db.UID.Hex()
	mine.Created = db.Created
	mine.Updated = db.Updated
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Scene = db.Scene
	mine.Name = db.Name
	mine.Remark = db.Remark
	mine.Parent = db.Parent
	mine.Cover = db.Cover
	mine.Access = db.Access
	mine.Tags = db.Tags
	mine.Users = db.Users
	mine.Contents = db.Contents
}

func (mine *FolderInfo) UpdateBase(name, remark, operator string) error {
	err := nosql.UpdateFolderBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
		mine.Updated = time.Now().Unix()
	}
	return err
}

func (mine *FolderInfo) UpdateAccess(operator string, acc uint32) error {
	if mine.Access == uint8(acc) {
		return nil
	}
	err := nosql.UpdateFolderAccess(mine.UID, operator, uint8(acc))
	if err == nil {
		mine.Access = uint8(acc)
		mine.Operator = operator
	}
	return err
}

func (mine *FolderInfo) UpdateParent(operator, parent string) error {
	if mine.Parent == parent {
		return nil
	}
	err := nosql.UpdateFolderParent(mine.UID, parent, operator)
	if err == nil {
		mine.Parent = parent
		mine.Operator = operator
	}
	return err
}

func (mine *FolderInfo) UpdateCover(operator, cover string) error {
	if mine.Cover == cover {
		return nil
	}
	err := nosql.UpdateFolderCover(mine.UID, cover, operator)
	if err == nil {
		mine.Cover = cover
		mine.Operator = operator
	}
	return err
}

func (mine *FolderInfo) UpdateContents(operator string, list []*pb.PairInt) error {
	arr := make([]*proxy.PairInfo, 0, len(list))
	for _, pair := range list {
		arr = append(arr, &proxy.PairInfo{
			Key:   pair.Key,
			Value: pair.Value,
			Count: pair.Count,
		})
	}
	err := nosql.UpdateFolderContents(mine.UID, operator, arr)
	if err == nil {
		mine.Contents = arr
		mine.Operator = operator
	}
	return err
}
