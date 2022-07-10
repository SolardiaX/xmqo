/*
 * Copyright (c) 2022. All rights reserved by XtraVisions.
 */

package xmgo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BulkResult struct {
	InsertedCount int64
	MatchedCount  int64
	ModifiedCount int64
	DeletedCount  int64
	UpsertedCount int64
	UpsertedIDs   map[int64]interface{}
}

type Bulk struct {
	coll *Collection

	queue   []mongo.WriteModel
	ordered *bool
}

func (c *Collection) Bulk() *Bulk {
	return &Bulk{
		coll:    c,
		queue:   nil,
		ordered: nil,
	}
}

func (b *Bulk) SetOrdered(ordered bool) *Bulk {
	b.ordered = &ordered
	return b
}

func (b *Bulk) InsertOne(doc interface{}) *Bulk {
	wm := mongo.NewInsertOneModel().SetDocument(doc)
	b.queue = append(b.queue, wm)
	return b
}

func (b *Bulk) Remove(filter interface{}) *Bulk {
	wm := mongo.NewDeleteOneModel().SetFilter(filter)
	b.queue = append(b.queue, wm)
	return b
}

func (b *Bulk) RemoveId(id interface{}) *Bulk {
	b.Remove(bson.M{"_id": id})
	return b
}

func (b *Bulk) RemoveAll(filter interface{}) *Bulk {
	wm := mongo.NewDeleteManyModel().SetFilter(filter)
	b.queue = append(b.queue, wm)
	return b
}

func (b *Bulk) Upsert(filter interface{}, replacement interface{}) *Bulk {
	wm := mongo.NewReplaceOneModel().SetFilter(filter).SetReplacement(replacement).SetUpsert(true)
	b.queue = append(b.queue, wm)
	return b
}

func (b *Bulk) UpsertId(id interface{}, replacement interface{}) *Bulk {
	b.Upsert(bson.M{"_id": id}, replacement)
	return b
}

func (b *Bulk) UpdateOne(filter interface{}, update interface{}) *Bulk {
	wm := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update)
	b.queue = append(b.queue, wm)
	return b
}

func (b *Bulk) UpdateId(id interface{}, update interface{}) *Bulk {
	b.UpdateOne(bson.M{"_id": id}, update)
	return b
}

func (b *Bulk) UpdateAll(filter interface{}, update interface{}) *Bulk {
	wm := mongo.NewUpdateManyModel().SetFilter(filter).SetUpdate(update)
	b.queue = append(b.queue, wm)
	return b
}

func (b *Bulk) Run() (*BulkResult, error) {
	return b.RunWithCtx(context.TODO())
}

func (b *Bulk) RunWithCtx(ctx context.Context) (*BulkResult, error) {
	opts := options.BulkWriteOptions{
		Ordered: b.ordered,
	}
	result, err := b.coll.collection.BulkWrite(ctx, b.queue, &opts)
	if err != nil {
		// In original mgo, queue is not reset in case of error.
		return nil, err
	}

	// Empty the queue for possible reuse, as per mgo's behavior.
	b.queue = nil

	return &BulkResult{
		InsertedCount: result.InsertedCount,
		MatchedCount:  result.MatchedCount,
		ModifiedCount: result.ModifiedCount,
		DeletedCount:  result.DeletedCount,
		UpsertedCount: result.UpsertedCount,
		UpsertedIDs:   result.UpsertedIDs,
	}, nil
}
