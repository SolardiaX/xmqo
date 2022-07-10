/*
 * Copyright (c) 2022. All rights reserved by XtraVisions.
 */

package xmgo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	opts "xtravisions.com/xmgo/options"
)

// Collection mongodb collection 操作封装
type Collection struct {
	collection *mongo.Collection
	registry   *bsoncodec.Registry
}

// Name 获取 collection 名称
func (c *Collection) Name() string {
	return c.collection.Name()
}

// Drop 删除 collection
func (c *Collection) Drop() error {
	return c.collection.Drop(context.TODO())
}

func (c *Collection) Watch(pipeline interface{}, opts ...*opts.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	return c.WatchWithCtx(context.TODO(), pipeline, opts...)
}

func (c *Collection) WatchWithCtx(ctx context.Context, pipeline interface{}, opts ...*opts.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	changeStreamOption := options.ChangeStream()
	if len(opts) > 0 && opts[0].ChangeStreamOptions != nil {
		changeStreamOption = opts[0].ChangeStreamOptions
	}

	return c.collection.Watch(ctx, pipeline, changeStreamOption)
}

func (c *Collection) Aggregate(pipeline interface{}, opts ...opts.AggregateOptions) IAggregate {
	return c.AggregateWithCtx(context.TODO(), pipeline, opts...)
}

func (c *Collection) AggregateWithCtx(ctx context.Context, pipeline interface{}, opts ...opts.AggregateOptions) IAggregate {
	return &Aggregate{
		ctx:        ctx,
		collection: c.collection,
		pipeline:   pipeline,
		options:    opts,
	}
}

func (c *Collection) Find(filter interface{}, opts ...opts.FindOptions) IQuery {
	return c.FindWithCtx(context.TODO(), filter, opts...)
}

func (c *Collection) FindWithCtx(ctx context.Context, filter interface{}, opts ...opts.FindOptions) IQuery {
	return &Query{
		ctx:        ctx,
		collection: c.collection,
		filter:     filter,
		opts:       opts,
		registry:   c.registry,
	}
}
