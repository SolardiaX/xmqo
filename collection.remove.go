/*
 * Copyright (c) 2022. All rights reserved by XtraVisions.
 */

package xmgo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"xtravisions.com/xmgo/hooks"
	opts "xtravisions.com/xmgo/options"
)

func (c *Collection) Remove(filter interface{}, opts ...opts.RemoveOptions) error {
	return c.RemoveWithCtx(context.TODO(), filter, opts...)
}

func (c *Collection) RemoveWithCtx(ctx context.Context, filter interface{}, opts ...opts.RemoveOptions) (err error) {
	deleteOptions := options.Delete()
	if len(opts) > 0 {
		if opts[0].DeleteOptions != nil {
			deleteOptions = opts[0].DeleteOptions
		}
		if opts[0].RemoveHook != nil {
			if err = hooks.On(ctx, opts[0].RemoveHook, hooks.BeforeRemove); err != nil {
				return err
			}
		}
	}

	var res *mongo.DeleteResult

	res, err = c.collection.DeleteOne(ctx, filter, deleteOptions)
	if res != nil && res.DeletedCount == 0 {
		err = ErrNoSuchDocuments
	}

	if err != nil {
		return err
	}

	if len(opts) > 0 && opts[0].RemoveHook != nil {
		if err = hooks.On(ctx, opts[0].RemoveHook, hooks.AfterRemove); err != nil {
			return err
		}
	}

	return
}

func (c *Collection) RemoveById(id interface{}, opts ...opts.RemoveOptions) error {
	return c.RemoveByIdWithCtx(context.TODO(), id, opts...)
}

func (c *Collection) RemoveByIdWithCtx(ctx context.Context, id interface{}, opts ...opts.RemoveOptions) (err error) {
	deleteOptions := options.Delete()
	if len(opts) > 0 {
		if opts[0].DeleteOptions != nil {
			deleteOptions = opts[0].DeleteOptions
		}
		if opts[0].RemoveHook != nil {
			if err = hooks.On(ctx, opts[0].RemoveHook, hooks.BeforeRemove); err != nil {
				return err
			}
		}
	}

	var res *mongo.DeleteResult

	res, err = c.collection.DeleteOne(ctx, bson.M{"_id": id}, deleteOptions)
	if res != nil && res.DeletedCount == 0 {
		err = ErrNoSuchDocuments
	}

	if err != nil {
		return err
	}

	if len(opts) > 0 && opts[0].RemoveHook != nil {
		if err = hooks.On(ctx, opts[0].RemoveHook, hooks.AfterRemove); err != nil {
			return err
		}
	}

	return
}
