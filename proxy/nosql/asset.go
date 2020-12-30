package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Asset struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Owner    string `json:"owner" bson:"owner"`
	Type     uint8  `json:"type" bson:"type"`
	Size     uint64 `json:"size" bson:"size"`
	UUID     string `json:"uuid" bson:"uuid"`
	Format   string `json:"format" bson:"format"`
	MD5      string `json:"md5" bson:"md5"`
	Version  string `json:"version" bson:"version"`
	Language string `json:"language" bson:"language"`
	Snapshot string `json:"snapshot" bson:"snapshot"`
	Small string `json:"small" bson:"small"`
}

func CreateAsset(info *Asset) error {
	_, err := insertOne(TableAssets, &info)
	return err
}

func GetAssetNextID() uint64 {
	num, _ := getSequenceNext(TableAssets)
	return num
}

func RemoveAsset(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db Asset uid is empty ")
	}
	_, err := removeOne(TableAssets, uid, operator)
	return err
}

func GetAsset(uid string) (*Asset, error) {
	if len(uid) < 2 {
		return nil, errors.New("db Asset uid is empty of GetAsset")
	}

	result, err := findOne(TableAssets, uid)
	if err != nil {
		return nil, err
	}
	model := new(Asset)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func UpdateAssetSnapshot(uid, snapshot,operator string) error {
	if len(uid) < 2 {
		return errors.New("db asset uid is empty of GetAsset")
	}

	msg := bson.M{"snapshot": snapshot,"operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableAssets, uid, msg)
	return err
}

func UpdateAssetSmall(uid, small,operator string) error {
	if len(uid) < 2 {
		return errors.New("db asset uid is empty of GetAsset")
	}

	msg := bson.M{"small": small,"operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableAssets, uid, msg)
	return err
}

func GetAssetsByOwner(owner string) ([]*Asset, error) {
	var items = make([]*Asset, 0, 20)
	def := new(time.Time)
	filter := bson.M{"owner": owner, "deleteAt": def}
	cursor, err1 := findMany(TableAssets, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	for cursor.Next(context.Background()) {
		var node = new(Asset)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}
