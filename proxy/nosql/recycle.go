package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Recycle struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Created     int64              `json:"created" bson:"created"`
	Deleted     int64              `json:"deleted" bson:"deleted"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Scavenger string   `json:"scavenger" bson:"scavenger"`
	Status    uint8    `json:"status" bson:"status"`
	Type      uint8    `json:"type" bson:"type"`
	Owner     string   `json:"owner" bson:"owner"`
	Size      uint64   `json:"size" bson:"size"`
	UUID      string   `json:"uuid" bson:"uuid"`
	Format    string   `json:"format" bson:"format"`
	MD5       string   `json:"md5" bson:"md5"`
	Version   string   `json:"version" bson:"version"`
	Language  string   `json:"language" bson:"language"`
	Snapshot  string   `json:"snapshot" bson:"snapshot"`
	Small     string   `json:"small" bson:"small"`
	Remark    string   `json:"remark" bson:"remark"`
	Meta      string   `json:"meta" bson:"meta"`
	Weight    uint32   `json:"weight" bson:"weight"`
	Width     uint32   `json:"width" bson:"width"`
	Height    uint32   `json:"height" bson:"height"`
	Links     []string `json:"links" bson:"links"`
}

func CreateRecycle(info *Recycle) error {
	_, err := insertOne(TableRecycles, &info)
	return err
}

func GetRecycleNextID() uint64 {
	num, _ := getSequenceNext(TableRecycles)
	return num
}

func RemoveRecycle(uid string) error {
	_, err := deleteOne(TableRecycles, uid)
	return err
}

func GetRecycleCount() int64 {
	num, _ := getCount(TableRecycles)
	return num
}

func GetRecycle(uid string) (*Recycle, error) {
	if len(uid) < 2 {
		return nil, errors.New("db Recycle uid is empty of GetRecycle")
	}

	result, err := findOne(TableRecycles, uid)
	if err != nil {
		return nil, err
	}
	model := new(Recycle)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAllRecycles() ([]*Recycle, error) {
	var items = make([]*Recycle, 0, 20)
	def := new(time.Time)
	filter := bson.M{"deleteAt": def}
	cursor, err1 := findMany(TableRecycles, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	for cursor.Next(context.Background()) {
		var node = new(Recycle)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetRecyclesByOwner(owner string) ([]*Recycle, error) {
	var items = make([]*Recycle, 0, 20)
	def := new(time.Time)
	filter := bson.M{"owner": owner, "deleteAt": def}
	cursor, err1 := findMany(TableRecycles, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	for cursor.Next(context.Background()) {
		var node = new(Recycle)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetRecyclesByCreator(user string) ([]*Recycle, error) {
	var items = make([]*Recycle, 0, 20)
	def := new(time.Time)
	filter := bson.M{"creator": user, "deleteAt": def}
	cursor, err1 := findMany(TableRecycles, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	for cursor.Next(context.Background()) {
		var node = new(Recycle)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetRecyclesByLink(link string) ([]*Recycle, error) {
	var items = make([]*Recycle, 0, 20)
	def := new(time.Time)
	filter := bson.M{"links": link, "deleteAt": def}
	cursor, err1 := findMany(TableRecycles, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	for cursor.Next(context.Background()) {
		var node = new(Recycle)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}
