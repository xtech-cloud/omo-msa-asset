package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.asset/proxy"
	"time"
)

type Thumb struct {
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

	Owner    string             `json:"owner" bson:"owner"`
	Probably float32            `json:"probably" bson:"probably"`
	Blur     float32            `json:"blur" bson:"blur"`
	File     string             `json:"file" bson:"file"` //远程文件名
	Asset    string             `json:"asset" bson:"asset"`
	Similar  float32            `json:"similar" bson:"similar"`
	Meta     string             `json:"meta" bson:"meta"`
	User     string             `json:"user" bson:"user"`
	Quote    string             `json:"quote" bson:"quote"`
	Location proxy.LocationInfo `json:"location" bson:"location"`
}

func CreateThumb(info *Thumb) error {
	_, err := insertOne(TableThumbs, &info)
	return err
}

func GetThumbNextID() uint64 {
	num, _ := getSequenceNext(TableThumbs)
	return num
}

func GetThumbCount() int64 {
	num, _ := getCount(TableThumbs)
	return num
}

func RemoveThumb(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db thumb uid is empty ")
	}
	_, err := removeOne(TableThumbs, uid, operator)
	return err
}

func GetThumb(uid string) (*Thumb, error) {
	if len(uid) < 2 {
		return nil, errors.New("db thumb uid is empty of GetThumb")
	}

	result, err := findOne(TableThumbs, uid)
	if err != nil {
		return nil, err
	}
	model := new(Thumb)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetThumbByFace(asset, face string) (*Thumb, error) {
	if len(face) < 2 {
		return nil, errors.New("db thumb face is empty of GetThumbByFace")
	}
	filter := bson.M{"asset": asset, "face": face, TimeDeleted: 0}
	result, err := findOneBy(TableThumbs, filter)
	if err != nil {
		return nil, err
	}
	model := new(Thumb)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetThumbsByOwner(owner string) ([]*Thumb, error) {
	var items = make([]*Thumb, 0, 20)
	filter := bson.M{"owner": owner, TimeDeleted: 0}
	cursor, err1 := findMany(TableThumbs, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	for cursor.Next(context.Background()) {
		var node = new(Thumb)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetThumbsByQuote(quote string) ([]*Thumb, error) {
	var items = make([]*Thumb, 0, 20)
	filter := bson.M{"quote": quote, TimeDeleted: 0}
	cursor, err1 := findMany(TableThumbs, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	for cursor.Next(context.Background()) {
		var node = new(Thumb)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetThumbsByAsset(asset string) ([]*Thumb, error) {
	var items = make([]*Thumb, 0, 20)
	filter := bson.M{"asset": asset, TimeDeleted: 0}
	cursor, err1 := findMany(TableThumbs, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	for cursor.Next(context.Background()) {
		var node = new(Thumb)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateThumbBase(uid, owner string, similar float32) error {
	if len(uid) < 2 {
		return errors.New("db thumb uid is empty of UpdateThumbBase")
	}

	msg := bson.M{"owner": owner, "similar": similar, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableThumbs, uid, msg)
	return err
}

func UpdateThumbMeta(uid, meta, operator string) error {
	if len(uid) < 2 {
		return errors.New("db thumb uid is empty of GetAsset")
	}

	msg := bson.M{"meta": meta, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableThumbs, uid, msg)
	return err
}

func UpdateThumbUser(uid, user, operator string) error {
	if len(uid) < 2 {
		return errors.New("db thumb uid is empty of GetAsset")
	}

	msg := bson.M{"user": user, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableThumbs, uid, msg)
	return err
}
