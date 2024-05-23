package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Label struct {
	UID      primitive.ObjectID `bson:"_id"`
	ID       uint64             `json:"id" bson:"id"`
	Name     string             `json:"name" bson:"name"`
	Created  int64              `json:"created" bson:"created"`
	Updated  int64              `json:"updated" bson:"updated"`
	Deleted  int64              `json:"deleted" bson:"deleted"`
	Creator  string             `json:"creator" bson:"creator"`
	Operator string             `json:"operator" bson:"operator"`

	Scene  string `json:"scene" bson:"scene"`
	Remark string `json:"face" bson:"face"`
	Type   uint8  `json:"type" bson:"type"`
}

func CreateLabel(info *Label) error {
	_, err := insertOne(TableLabels, &info)
	return err
}

func GetLabelNextID() uint64 {
	num, _ := getSequenceNext(TableLabels)
	return num
}

func GetLabelCount() int64 {
	num, _ := getCount(TableLabels)
	return num
}

func RemoveLabel(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db thumb uid is empty ")
	}
	_, err := removeOne(TableLabels, uid, operator)
	return err
}

func GetLabelChildrenCount(parent string) int64 {
	filter := bson.M{"parent": parent, TimeDeleted: 0}
	num, _ := getCountByFilter(TableLabels, filter)
	return num
}

func GetLabel(uid string) (*Label, error) {
	if len(uid) < 2 {
		return nil, errors.New("db thumb uid is empty of GetLabel")
	}

	result, err := findOne(TableLabels, uid)
	if err != nil {
		return nil, err
	}
	model := new(Label)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetLabelsByScene(owner string) ([]*Label, error) {
	var items = make([]*Label, 0, 20)
	filter := bson.M{"scene": owner, TimeDeleted: 0}
	cursor, err1 := findMany(TableLabels, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	for cursor.Next(context.Background()) {
		var node = new(Label)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetLabelByName(name string) (*Label, error) {
	filter := bson.M{"name": name, TimeDeleted: 0}
	result, err := findOneBy(TableLabels, filter)
	if err != nil {
		return nil, err
	}
	model := new(Label)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func UpdateLabelBase(uid, name, remark, operator string) error {
	if len(uid) < 2 {
		return errors.New("db folder uid is empty")
	}

	msg := bson.M{"name": name, "remark": remark, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableLabels, uid, msg)
	return err
}
