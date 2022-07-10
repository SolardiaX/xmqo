/*
 * Copyright (c) 2022. All rights reserved by XtraVisions.
 */

package xmgo

import (
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"xtravisions.com/xmgo/hooks"
	opts "xtravisions.com/xmgo/options"
)

func (c *Collection) InsertOne(doc interface{}, opts ...opts.InsertOneOptions) (*InsertOneResult, error) {
	return c.InsertOneWithCtx(context.TODO(), doc, opts...)
}

func (c *Collection) InsertOneWithCtx(ctx context.Context, doc interface{}, opts ...opts.InsertOneOptions) (result *InsertOneResult, err error) {
	h := doc
	insertOneOpts := options.InsertOne()
	if len(opts) > 0 {
		if opts[0].InsertOneOptions != nil {
			insertOneOpts = opts[0].InsertOneOptions
		}
		if opts[0].InsertHook != nil {
			h = opts[0].InsertHook
		}
	}

	if err = hooks.On(ctx, doc, hooks.BeforeInsert, h); err != nil {
		return
	}

	var res *mongo.InsertOneResult

	if res, err = c.collection.InsertOne(ctx, doc, insertOneOpts); err != nil {
		return
	}

	result = &InsertOneResult{InsertedID: res.InsertedID}

	if err = hooks.On(ctx, doc, hooks.AfterInsert, h); err != nil {
		return
	}

	return
}

func (c *Collection) InsertMany(docs interface{}, opts ...opts.InsertManyOptions) (result *InsertManyResult, err error) {
	return c.InsertManyWithCtx(context.TODO(), docs, opts...)
}

func (c *Collection) InsertManyWithCtx(ctx context.Context, docs interface{}, opts ...opts.InsertManyOptions) (result *InsertManyResult, err error) {
	h := docs
	insertManyOpts := options.InsertMany()
	if len(opts) > 0 {
		if opts[0].InsertManyOptions != nil {
			insertManyOpts = opts[0].InsertManyOptions
		}
		if opts[0].InsertHook != nil {
			h = opts[0].InsertHook
		}
	}

	if err = hooks.On(ctx, docs, hooks.BeforeInsert, h); err != nil {
		return
	}

	sDocs := interfaceToSliceInterface(docs)
	if sDocs == nil {
		return nil, ErrNotValidSliceToInsert
	}

	var res *mongo.InsertManyResult
	if res, err = c.collection.InsertMany(ctx, sDocs, insertManyOpts); err != nil {
		return
	}

	result = &InsertManyResult{InsertedIDs: res.InsertedIDs}

	if err = hooks.On(ctx, docs, hooks.AfterInsert, h); err != nil {
		return
	}

	return
}

func interfaceToSliceInterface(docs interface{}) []interface{} {
	if reflect.Slice != reflect.TypeOf(docs).Kind() {
		return nil
	}

	s := reflect.ValueOf(docs)
	if s.Len() == 0 {
		return nil
	}

	var sDocs []interface{}
	for i := 0; i < s.Len(); i++ {
		sDocs = append(sDocs, s.Index(i).Interface())
	}

	return sDocs
}
