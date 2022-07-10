/*
 * Copyright (c) 2022. All rights reserved by XtraVisions.
 */

package xmgo

import (
	"context"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"xtravisions.com/xmgo/hooks"
	opts "xtravisions.com/xmgo/options"
)

type IQuery interface {
	Sort(fields ...string) IQuery
	Select(selector interface{}) IQuery
	Skip(n int64) IQuery
	Limit(n int64) IQuery
	One(result interface{}) error
	All(result interface{}) error
	Count() (n int64, err error)
	Exists() (b bool, err error)
	Distinct(key string, result interface{}) error
	Cursor() ICursor
	Apply(change Change, result interface{}) error
	Hint(hint interface{}) IQuery
}

type Change struct {
	Update    interface{} // update/replace document
	Replace   bool        // Whether to replace the document rather than updating
	Remove    bool        // Whether to remove the document found rather than updating
	Upsert    bool        // Whether to insert in case the document isn't found, take effect when Remove is false
	ReturnNew bool        // Should the modified document be returned rather than the old one, take effect when Remove is false
}

type Query struct {
	filter  interface{}
	sort    interface{}
	project interface{}
	hint    interface{}
	limit   *int64
	skip    *int64

	ctx        context.Context
	collection *mongo.Collection
	opts       []opts.FindOptions
	registry   *bsoncodec.Registry
}

func (q *Query) Sort(fields ...string) IQuery {
	if len(fields) == 0 {
		return q
	}

	var sorts bson.D
	for _, field := range fields {
		key, n := splitSortField(field)
		if key == "" {
			panic("Mongo Sort: 字段名不能为空")
		}

		sorts = append(sorts, bson.E{Key: key, Value: n})
	}

	newQ := q
	newQ.sort = sorts

	return newQ
}

func (q *Query) Select(projection interface{}) IQuery {
	newQ := q
	newQ.project = projection

	return newQ
}

func (q *Query) Skip(n int64) IQuery {
	newQ := q
	newQ.skip = &n

	return newQ
}

func (q *Query) Hint(hint interface{}) IQuery {
	newQ := q
	newQ.hint = hint

	return newQ
}

func (q *Query) Limit(n int64) IQuery {
	newQ := q
	newQ.limit = &n

	return newQ
}

func (q *Query) One(result interface{}) error {
	if len(q.opts) > 0 {
		if err := hooks.On(q.ctx, q.opts[0].QueryHook, hooks.BeforeQuery); err != nil {
			return err
		}
	}
	opt := options.FindOne()

	if q.sort != nil {
		opt.SetSort(q.sort)
	}
	if q.project != nil {
		opt.SetProjection(q.project)
	}
	if q.skip != nil {
		opt.SetSkip(*q.skip)
	}
	if q.hint != nil {
		opt.SetHint(q.hint)
	}

	err := q.collection.FindOne(q.ctx, q.filter, opt).Decode(result)

	if err != nil {
		return err
	}

	if len(q.opts) > 0 {
		if err := hooks.On(q.ctx, q.opts[0].QueryHook, hooks.AfterQuery); err != nil {
			return err
		}
	}

	return nil
}

func (q *Query) All(result interface{}) error {
	if len(q.opts) > 0 {
		if err := hooks.On(q.ctx, q.opts[0].QueryHook, hooks.BeforeQuery); err != nil {
			return err
		}
	}
	opt := options.Find()

	if q.sort != nil {
		opt.SetSort(q.sort)
	}
	if q.project != nil {
		opt.SetProjection(q.project)
	}
	if q.limit != nil {
		opt.SetLimit(*q.limit)
	}
	if q.skip != nil {
		opt.SetSkip(*q.skip)
	}
	if q.hint != nil {
		opt.SetHint(q.hint)
	}

	var err error
	var cursor *mongo.Cursor

	cursor, err = q.collection.Find(q.ctx, q.filter, opt)

	c := Cursor{
		ctx:    q.ctx,
		cursor: cursor,
		err:    err,
	}
	err = c.All(result)
	if err != nil {
		return err
	}
	if len(q.opts) > 0 {
		if err := hooks.On(q.ctx, q.opts[0].QueryHook, hooks.AfterQuery); err != nil {
			return err
		}
	}
	return nil
}

func (q *Query) Count() (n int64, err error) {
	opt := options.Count()

	if q.limit != nil {
		opt.SetLimit(*q.limit)
	}
	if q.skip != nil {
		opt.SetSkip(*q.skip)
	}

	return q.collection.CountDocuments(q.ctx, q.filter, opt)
}

func (q *Query) Exists() (b bool, err error) {
	var n int64 = 0
	if n, err = q.Count(); err == nil {
		if n > 0 {
			b = true
		}
	}

	return
}

func (q *Query) Distinct(key string, result interface{}) error {
	resultVal := reflect.ValueOf(result)

	if resultVal.Kind() != reflect.Ptr {
		return ErrQueryNotSlicePointer
	}

	resultElmVal := resultVal.Elem()
	if resultElmVal.Kind() != reflect.Interface && resultElmVal.Kind() != reflect.Slice {
		return ErrQueryNotSliceType
	}

	opt := options.Distinct()
	res, err := q.collection.Distinct(q.ctx, key, q.filter, opt)
	if err != nil {
		return err
	}
	registry := q.registry
	if registry == nil {
		registry = bson.DefaultRegistry
	}
	valueType, valueBytes, err_ := bson.MarshalValueWithRegistry(registry, res)
	if err_ != nil {
		fmt.Printf("bson.MarshalValue err: %+v\n", err_)
		return err_
	}

	rawValue := bson.RawValue{Type: valueType, Value: valueBytes}
	err = rawValue.Unmarshal(result)
	if err != nil {
		fmt.Printf("rawValue.Unmarshal err: %+v\n", err)
		return ErrQueryResultTypeInconsistent
	}

	return nil
}

func (q *Query) Cursor() ICursor {
	opt := options.Find()

	if q.sort != nil {
		opt.SetSort(q.sort)
	}
	if q.project != nil {
		opt.SetProjection(q.project)
	}
	if q.limit != nil {
		opt.SetLimit(*q.limit)
	}
	if q.skip != nil {
		opt.SetSkip(*q.skip)
	}

	var err error
	var cur *mongo.Cursor
	cur, err = q.collection.Find(q.ctx, q.filter, opt)

	return &Cursor{
		ctx:    q.ctx,
		cursor: cur,
		err:    err,
	}
}

func (q *Query) Apply(change Change, result interface{}) error {
	var err error

	if change.Remove {
		err = q.findOneAndDelete(change, result)
	} else if change.Replace {
		err = q.findOneAndReplace(change, result)
	} else {
		err = q.findOneAndUpdate(change, result)
	}

	return err
}

func (q *Query) findOneAndDelete(_ Change, result interface{}) error {
	opt := options.FindOneAndDelete()
	if q.sort != nil {
		opt.SetSort(q.sort)
	}
	if q.project != nil {
		opt.SetProjection(q.project)
	}

	return q.collection.FindOneAndDelete(q.ctx, q.filter, opt).Decode(result)
}

func (q *Query) findOneAndReplace(change Change, result interface{}) error {
	opt := options.FindOneAndReplace()
	if q.sort != nil {
		opt.SetSort(q.sort)
	}
	if q.project != nil {
		opt.SetProjection(q.project)
	}
	if change.Upsert {
		opt.SetUpsert(change.Upsert)
	}
	if change.ReturnNew {
		opt.SetReturnDocument(options.After)
	}

	err := q.collection.FindOneAndReplace(q.ctx, q.filter, change.Update, opt).Decode(result)
	if change.Upsert && !change.ReturnNew && err == mongo.ErrNoDocuments {
		return nil
	}

	return err
}

func (q *Query) findOneAndUpdate(change Change, result interface{}) error {
	opt := options.FindOneAndUpdate()
	if q.sort != nil {
		opt.SetSort(q.sort)
	}
	if q.project != nil {
		opt.SetProjection(q.project)
	}
	if change.Upsert {
		opt.SetUpsert(change.Upsert)
	}
	if change.ReturnNew {
		opt.SetReturnDocument(options.After)
	}

	err := q.collection.FindOneAndUpdate(q.ctx, q.filter, change.Update, opt).Decode(result)
	if change.Upsert && !change.ReturnNew && err == mongo.ErrNoDocuments {
		return nil
	}

	return err
}
