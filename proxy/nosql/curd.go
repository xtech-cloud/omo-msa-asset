package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const timeOut = 10 * time.Second

func insertOne(collection string, info interface{}) (interface{}, error) {
	if len(collection) < 1 {
		return "",	errors.New("the collection is empty")
	}
	c := noSql.Collection(collection)
	if c == nil {
		return "", errors.New("can not found the collection of" + collection)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	result, err := c.InsertOne(ctx, info)
	if err != nil {
		return "", err
	}
	return result.InsertedID, nil
}

func getCount(collection string) (int64, error) {
	if len(collection) < 1 {
		return 0,	errors.New("the collection is empty")
	}
	c := noSql.Collection(collection)
	if c == nil {
		return 0, errors.New("can not found the collection of" + collection)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	result, err := c.EstimatedDocumentCount(ctx, nil)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func deleteOne(collection string, uid string) (int64, error) {
	if len(collection) < 1 {
		return 0,	errors.New("the collection is empty")
	}
	objID, e := primitive.ObjectIDFromHex(uid)
	if e != nil {
		return 0, e
	}
	c := noSql.Collection(collection)
	if c == nil {
		return 0, errors.New("can not found the collection of" + collection)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	node := bson.M{"_id": objID}
	result, err := c.DeleteOne(ctx, node)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func removeOne(collection string, uid, operator string) (int64, error) {
	if len(collection) < 1 {
		return 0,	errors.New("the collection is empty")
	}
	objID, e := primitive.ObjectIDFromHex(uid)
	if e != nil {
		return 0, e
	}
	c := noSql.Collection(collection)
	if c == nil {
		return 0, errors.New("can not found the collection of" + collection)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	filter := bson.M{"_id": objID}
	node := bson.M{"$set": bson.M{"operator":operator, "deleteAt": time.Now()}}
	result, err := c.UpdateOne(ctx, filter, node)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

func hadOne(collection string, filter bson.M) (bool, error) {
	if len(collection) < 1 {
		return false, errors.New("the collection is empty")
	}
	c := noSql.Collection(collection)
	if c == nil {
		return false, errors.New("can not found the collection of" + collection)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	result := c.FindOne(ctx, filter)
	if result.Err() != nil {
		return false, result.Err()
	}
	return true, nil
}

func updateOne(collection string, uid string, data bson.M) (int64, error) {
	if len(collection) < 1 {
		return 0,	errors.New("the collection is empty")
	}
	objID, e := primitive.ObjectIDFromHex(uid)
	if e != nil {
		return 0, e
	}
	c := noSql.Collection(collection)
	if c == nil {
		return 0, errors.New("can not found the collection of" + collection)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	filter := bson.M{"_id": objID}
	node := bson.M{"$set": data}
	result, err := c.UpdateOne(ctx, filter, node)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

/**
往数组里面追加一个元素
*/
func appendElement(collection string, uid string, data bson.M) (int64, error) {
	if len(collection) < 1 {
		return 0,	errors.New("the collection is empty")
	}
	objID, e := primitive.ObjectIDFromHex(uid)
	if e != nil {
		return 0, e
	}
	c := noSql.Collection(collection)
	if c == nil {
		return 0, errors.New("can not found the collection of" + collection)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	filter := bson.M{"_id": objID}
	node := bson.M{"$push": data, "$set": bson.M{"updatedAt": time.Now()}}
	result, err := c.UpdateOne(ctx, filter, node)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

/**
从数组里面移除一个元素
*/
func removeElement(collection string, uid string, data bson.M) (int64, error) {
	if len(collection) < 1 {
		return 0,	errors.New("the collection is empty")
	}
	objID, e := primitive.ObjectIDFromHex(uid)
	if e != nil {
		return 0, e
	}
	c := noSql.Collection(collection)
	if c == nil {
		return 0, errors.New("can not found the collection of" + collection)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	filter := bson.M{"_id": objID}
	node := bson.M{"$pull": data, "$set": bson.M{"updatedAt": time.Now()}}
	result, err := c.UpdateOne(ctx, filter, node)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

func updateOneBy(collection string, filter bson.M, update bson.M) (int64, error) {
	if len(collection) < 1 {
		return 0,	errors.New("the collection is empty")
	}
	c := noSql.Collection(collection)
	if c == nil {
		return 0, errors.New("can not found the collection of" + collection)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	result, err := c.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

func findOne(collection string, uid string) (*mongo.SingleResult, error) {
	if len(collection) < 1 {
		return nil,	errors.New("the collection is empty")
	}
	objID, e := primitive.ObjectIDFromHex(uid)
	if e != nil {
		return nil, e
	}
	c := noSql.Collection(collection)
	if c == nil {
		return nil, errors.New("can not found the collection of" + collection)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	filter := bson.M{"_id": objID}
	result := c.FindOne(ctx, filter)
	if result.Err() != nil {
		return nil, result.Err()
	}
	return result, nil
}

func findOneBy(collection string, filter bson.M) (*mongo.SingleResult, error) {
	if len(collection) < 1 {
		return nil,	errors.New("the collection is empty")
	}
	c := noSql.Collection(collection)
	if c == nil {
		return nil, errors.New("can not found the collection of" + collection)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	result := c.FindOne(ctx, filter)
	if result.Err() != nil {
		return nil, result.Err()
	}
	return result, nil
}

func findOneOfField(collection string, uid string, selector bson.M) (*mongo.SingleResult, error) {
	if len(collection) < 1 {
		return nil,	errors.New("the collection is empty")
	}
	objID, e := primitive.ObjectIDFromHex(uid)
	if e != nil {
		return nil, e
	}
	c := noSql.Collection(collection)
	if c == nil {
		return nil, errors.New("can not found the collection of" + collection)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	filter := bson.M{"_id": objID}
	result := c.FindOne(ctx, filter, options.FindOne().SetProjection(selector))
	if result.Err() != nil {
		return nil, result.Err()
	}
	return result, nil
}

func findOneByOpt(collection string, filter bson.M, selector bson.M) (*mongo.SingleResult, error) {
	if len(collection) < 1 {
		return nil,	errors.New("the collection is empty")
	}
	c := noSql.Collection(collection)
	if c == nil {
		return nil, errors.New("can not found the collection of" + collection)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	result := c.FindOne(ctx, filter, options.FindOne().SetProjection(selector))
	if result.Err() != nil {
		return nil, result.Err()
	}
	return result, nil
}

func findMany(collection string, filter bson.M, limit int64) (*mongo.Cursor, error) {
	if len(collection) < 1 {
		return nil,	errors.New("the collection is empty")
	}
	c := noSql.Collection(collection)
	if c == nil {
		return nil, errors.New("can not found the collection of" + collection)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	var cursor *mongo.Cursor
	var err error
	if limit > 0 {
		cursor, err = c.Find(ctx, filter, options.Find().SetLimit(limit))
	} else {
		cursor, err = c.Find(ctx, filter)
	}
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

func findManyByOpts(collection string, filter bson.M, opts *options.FindOptions) (*mongo.Cursor, error) {
	if len(collection) < 1 {
		return nil,	errors.New("the collection is empty")
	}
	c := noSql.Collection(collection)
	if c == nil {
		return nil, errors.New("can not found the collection of" + collection)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	cursor, err := c.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

func findAllByOpts(collection string, opts *options.FindOptions) (*mongo.Cursor, error) {
	if len(collection) < 1 {
		return nil,	errors.New("the collection is empty")
	}
	c := noSql.Collection(collection)
	if c == nil {
		return nil, errors.New("can not found the collection of" + collection)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	def := new(time.Time)
	filter := bson.M{"deleteAt": def}
	cursor, err := c.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

func findAll(collection string, limit int64) (*mongo.Cursor, error) {
	if len(collection) < 1 {
		return nil,	errors.New("the collection is empty")
	}
	c := noSql.Collection(collection)
	if c == nil {
		return nil, errors.New("can not found the collection of" + collection)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	def := new(time.Time)
	filter := bson.M{"deleteAt": def}
	var cursor *mongo.Cursor
	var err error
	if limit > 0 {
		cursor, err = c.Find(ctx, filter, options.Find().SetLimit(limit))
	} else {
		cursor, err = c.Find(ctx, filter)
	}
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

func dropOne(collection string) error {
	if len(collection) < 1 {
		return errors.New("the collection is empty")
	}
	c := noSql.Collection(collection)
	if c == nil {
		return errors.New("can not found the collection of" + collection)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	err := c.Drop(ctx)
	if err != nil {
		return err
	}

	return nil
}

func copyOne(collection string) (*mongo.Collection, error) {
	if len(collection) < 1 {
		return nil,	errors.New("the collection is empty")
	}
	c := noSql.Collection(collection)
	if c == nil {
		return nil, errors.New("can not found the collection of" + collection)
	}
	tmp, err := c.Clone()
	if err != nil {
		return nil, err
	}
	return tmp, nil
}
