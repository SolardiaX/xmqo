/*
 * Copyright (c) 2022. All rights reserved by XtraVisions.
 */

package xmgo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	opts "xtravisions.com/xmgo/options"
)

type IAggregate interface {
	All(results interface{}) error
	One(result interface{}) error
	Iter() ICursor
}

type Aggregate struct {
	ctx        context.Context
	pipeline   interface{}
	collection *mongo.Collection
	options    []opts.AggregateOptions
}

func (a *Aggregate) All(results interface{}) error {
	aOpts := options.Aggregate()
	if len(a.options) > 0 {
		aOpts = a.options[0].AggregateOptions
	}

	c, err := a.collection.Aggregate(a.ctx, a.pipeline, aOpts)
	if err != nil {
		return err
	}

	return c.All(a.ctx, results)
}

func (a *Aggregate) One(result interface{}) error {
	aOpts := options.Aggregate()
	if len(a.options) > 0 {
		aOpts = a.options[0].AggregateOptions
	}

	c, err := a.collection.Aggregate(a.ctx, a.pipeline, aOpts)
	if err != nil {
		return err
	}

	cr := Cursor{
		ctx:    a.ctx,
		cursor: c,
		err:    err,
	}
	defer func(cr *Cursor) {
		_ = cr.Close()
	}(&cr)

	if !cr.Next(result) {
		return ErrNoSuchDocuments
	}

	return err
}

func (a *Aggregate) Iter() ICursor {
	aOpts := options.Aggregate()
	if len(a.options) > 0 {
		aOpts = a.options[0].AggregateOptions
	}
	c, err := a.collection.Aggregate(a.ctx, a.pipeline, aOpts)
	return &Cursor{
		ctx:    a.ctx,
		cursor: c,
		err:    err,
	}
}
