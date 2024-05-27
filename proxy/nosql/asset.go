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
	Created     int64              `json:"created" bson:"created"`
	Updated     int64              `json:"updated" bson:"updated"`
	Deleted     int64              `json:"deleted" bson:"deleted"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Status   uint8    `json:"status" bson:"status"`
	Type     uint8    `json:"type" bson:"type"`
	Scope    uint8    `json:"scope" bson:"scope"`
	Code     int      `json:"code" bson:"code"`
	Owner    string   `json:"owner" bson:"owner"`
	Size     uint64   `json:"size" bson:"size"`
	UUID     string   `json:"uuid" bson:"uuid"`
	Format   string   `json:"format" bson:"format"`
	MD5      string   `json:"md5" bson:"md5"`
	Version  string   `json:"version" bson:"version"`
	Language string   `json:"language" bson:"language"`
	Snapshot string   `json:"snapshot" bson:"snapshot"`
	Small    string   `json:"small" bson:"small"`
	Remark   string   `json:"remark" bson:"remark"`
	Meta     string   `json:"meta" bson:"meta"`
	Weight   uint32   `json:"weight" bson:"weight"`
	Width    uint32   `json:"width" bson:"width"`
	Height   uint32   `json:"height" bson:"height"`
	Quote    string   `json:"quote" bson:"quote"`
	Links    []string `json:"links" bson:"links"`
	Tags     []string `json:"tags" bson:"tags"`
}

func CreateAsset(info *Asset) error {
	_, err := insertOne(TableAssets, &info)
	return err
}

func GetAssetNextID() uint64 {
	num, _ := getSequenceNext(TableAssets)
	return num
}

func RemoveAsset(uid string) error {
	_, err := deleteOne(TableAssets, uid)
	return err
}

func GetAssetCount() int64 {
	num, _ := getCount(TableAssets)
	return num
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

func GetAssetByKey(uid string) (*Asset, error) {
	if len(uid) < 2 {
		return nil, errors.New("db Asset uid is empty of GetAsset")
	}
	filter := bson.M{"uuid": uid}
	result, err := findOneBy(TableAssets, filter)
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

func GetAllAssets() ([]*Asset, error) {
	var items = make([]*Asset, 0, 20)
	cursor, err1 := findAllEnable(TableAssets, 0)
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

func UpdateAssetSnapshot(uid, snapshot, operator string) error {
	if len(uid) < 2 {
		return errors.New("db asset uid is empty of UpdateAssetSnapshot")
	}

	msg := bson.M{"snapshot": snapshot, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableAssets, uid, msg)
	return err
}

func UpdateAssetSmall(uid, small, operator string) error {
	if len(uid) < 2 {
		return errors.New("db asset uid is empty of UpdateAssetSmall")
	}

	msg := bson.M{"small": small, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableAssets, uid, msg)
	return err
}

func UpdateAssetCode(uid string, code int) error {
	if len(uid) < 2 {
		return errors.New("db asset uid is empty of UpdateAssetCode")
	}

	msg := bson.M{"code": code, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableAssets, uid, msg)
	return err
}

func UpdateAssetBase(uid, name, remark, operator string) error {
	if len(uid) < 2 {
		return errors.New("db asset uid is empty of UpdateAssetBase")
	}

	msg := bson.M{"name": name, "remark": remark, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableAssets, uid, msg)
	return err
}

func UpdateAssetMeta(uid, meta, operator string) error {
	if len(uid) < 2 {
		return errors.New("db asset uid is empty of UpdateAssetMeta")
	}

	msg := bson.M{"meta": meta, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableAssets, uid, msg)
	return err
}

func UpdateAssetWeight(uid, operator string, weight uint32) error {
	if len(uid) < 2 {
		return errors.New("db asset uid is empty of UpdateAssetWeight")
	}

	msg := bson.M{"weight": weight, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableAssets, uid, msg)
	return err
}

func UpdateAssetStatus(uid, operator string, status uint8) error {
	if len(uid) < 2 {
		return errors.New("db asset uid is empty of UpdateAssetWeight")
	}

	msg := bson.M{"status": status, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableAssets, uid, msg)
	return err
}

func UpdateAssetScope(uid, operator string, scope uint8) error {
	if len(uid) < 2 {
		return errors.New("db asset uid is empty of UpdateAssetWeight")
	}

	msg := bson.M{"scope": scope, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableAssets, uid, msg)
	return err
}

func UpdateAssetLinks(uid, operator string, arr []string) error {
	if len(uid) < 2 {
		return errors.New("db asset uid is empty of UpdateAssetWeight")
	}

	msg := bson.M{"links": arr, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableAssets, uid, msg)
	return err
}

func UpdateAssetType(uid, operator string, tp uint8) error {
	if len(uid) < 2 {
		return errors.New("db asset uid is empty of UpdateAssetType")
	}

	msg := bson.M{"type": tp, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableAssets, uid, msg)
	return err
}

func UpdateAssetOwner(uid, owner, operator string) error {
	if len(uid) < 2 {
		return errors.New("db asset uid is empty of UpdateAssetType")
	}

	msg := bson.M{"owner": owner, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableAssets, uid, msg)
	return err
}

func UpdateAssetQuote(uid, quote, operator string) error {
	if len(uid) < 2 {
		return errors.New("db asset uid is empty of UpdateAssetType")
	}

	msg := bson.M{"quote": quote, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableAssets, uid, msg)
	return err
}

func UpdateAssetLanguage(uid, operator, lan string) error {
	if len(uid) < 2 {
		return errors.New("db asset uid is empty of UpdateAssetLanguage")
	}

	msg := bson.M{"language": lan, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableAssets, uid, msg)
	return err
}

func GetAssetsByOwner(owner string) ([]*Asset, error) {
	var items = make([]*Asset, 0, 20)
	filter := bson.M{"owner": owner, TimeDeleted: 0}
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

func GetAssetsByOwnerStatus(owner string, st uint8) ([]*Asset, error) {
	var items = make([]*Asset, 0, 20)
	filter := bson.M{"owner": owner, "status": st, TimeDeleted: 0}
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

func GetAssetsByQuote(quote string) ([]*Asset, error) {
	var items = make([]*Asset, 0, 20)
	filter := bson.M{"quote": quote, TimeDeleted: 0}
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

func GetAssetsCountByQuote(quote string) int64 {
	filter := bson.M{"quote": quote, TimeDeleted: 0}
	num, _ := getCountByFilter(TableAssets, filter)
	return num
}

func GetAssetsByOwnerQuote(owner, quote string) ([]*Asset, error) {
	var items = make([]*Asset, 0, 20)
	filter := bson.M{"owner": owner, "quote": quote, TimeDeleted: 0}
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

func GetAssetsByRegex(key string, from, to int64) ([]*Asset, error) {
	var items = make([]*Asset, 0, 20)
	filter := bson.M{"name": bson.M{"$regex": key}, TimeCreated: bson.M{"$gt": from, "$lt": to}, TimeDeleted: 0}
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

func GetAssetsByType(tp uint8) ([]*Asset, error) {
	var items = make([]*Asset, 0, 20)
	filter := bson.M{"type": tp, TimeDeleted: 0}
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

func GetAssetsByOwnerType(owner string, tp uint8) ([]*Asset, error) {
	var items = make([]*Asset, 0, 20)
	filter := bson.M{"owner": owner, "type": tp, TimeDeleted: 0}
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

func GetAssetsByCreator(user string) ([]*Asset, error) {
	var items = make([]*Asset, 0, 20)
	filter := bson.M{"creator": user, TimeDeleted: 0}
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

func GetAssetsByLink(link string) ([]*Asset, error) {
	var items = make([]*Asset, 0, 20)
	filter := bson.M{"links": link, TimeDeleted: 0}
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
