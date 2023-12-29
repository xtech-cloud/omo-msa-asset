package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.asset/proxy"
	"time"
)

type Folder struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Created     int64              `json:"created" bson:"created"`
	Updated     int64              `json:"updated" bson:"updated"`
	Deleted     int64              `json:"deleted" bson:"deleted"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Scene    string            `json:"scene" bson:"scene"`
	Remark   string            `json:"face" bson:"face"`
	Parent   string            `json:"parent" bson:"parent"`
	Access   uint8             `json:"access" bson:"access"`
	Cover    string            `json:"cover" bson:"cover"`
	Tags     []string          `json:"tags" bson:"tags"`
	Users    []string          `json:"users" bson:"users"`
	Contents []*proxy.PairInfo `json:"contents" bson:"contents"`
}

func CreateFolder(info *Folder) error {
	_, err := insertOne(TableFolders, &info)
	return err
}

func GetFolderNextID() uint64 {
	num, _ := getSequenceNext(TableFolders)
	return num
}

func GetFolderCount() int64 {
	num, _ := getCount(TableFolders)
	return num
}

func RemoveFolder(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db thumb uid is empty ")
	}
	_, err := removeOne(TableFolders, uid, operator)
	return err
}

func GetFolderChildrenCount(parent string) int64 {
	filter := bson.M{"parent": parent, TimeDeleted: 0}
	num, _ := getCountByFilter(TableFolders, filter)
	return num
}

func GetFolder(uid string) (*Folder, error) {
	if len(uid) < 2 {
		return nil, errors.New("db thumb uid is empty of GetFolder")
	}

	result, err := findOne(TableFolders, uid)
	if err != nil {
		return nil, err
	}
	model := new(Folder)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetFoldersByScene(owner string) ([]*Folder, error) {
	var items = make([]*Folder, 0, 20)
	filter := bson.M{"scene": owner, TimeDeleted: 0}
	cursor, err1 := findMany(TableFolders, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	for cursor.Next(context.Background()) {
		var node = new(Folder)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetFoldersByParent(parent string) ([]*Folder, error) {
	var items = make([]*Folder, 0, 20)
	filter := bson.M{"parent": parent, TimeDeleted: 0}
	cursor, err1 := findMany(TableFolders, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	for cursor.Next(context.Background()) {
		var node = new(Folder)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateFolderBase(uid, name, remark, operator string) error {
	if len(uid) < 2 {
		return errors.New("db folder uid is empty")
	}

	msg := bson.M{"name": name, "remark": remark, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableAssets, uid, msg)
	return err
}

func UpdateFolderAccess(uid, operator string, acc uint8) error {
	if len(uid) < 2 {
		return errors.New("db folder uid is empty")
	}

	msg := bson.M{"access": acc, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableAssets, uid, msg)
	return err
}

func UpdateFolderContents(uid, operator string, list []*proxy.PairInfo) error {
	if len(uid) < 2 {
		return errors.New("db folder uid is empty")
	}

	msg := bson.M{"contents": list, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableAssets, uid, msg)
	return err
}

func UpdateFolderParent(uid, parent, operator string) error {
	if len(uid) < 2 {
		return errors.New("db folder uid is empty")
	}

	msg := bson.M{"parent": parent, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableAssets, uid, msg)
	return err
}

func UpdateFolderCover(uid, cover, operator string) error {
	if len(uid) < 2 {
		return errors.New("db folder uid is empty")
	}

	msg := bson.M{"cover": cover, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableAssets, uid, msg)
	return err
}
