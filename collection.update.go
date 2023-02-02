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

func (c *Collection) UpdateById(id interface{}, update interface{}, opts ...opts.UpdateOptions) error {
	return c.UpdateByIdWithCtx(context.TODO(), id, update, opts...)
}

func (c *Collection) UpdateByIdWithCtx(ctx context.Context, id interface{}, update interface{}, opts ...opts.UpdateOptions) (err error) {
	updateOpts := options.Update()

	if len(opts) > 0 {
		if opts[0].UpdateOptions != nil {
			updateOpts = opts[0].UpdateOptions
		}
		if opts[0].UpdateHook != nil {
			if err = hooks.On(ctx, opts[0].UpdateHook, hooks.BeforeUpdate); err != nil {
				return
			}
		}
	}

	var res *mongo.UpdateResult

	res, err = c.collection.UpdateOne(ctx, bson.M{"_id": id}, update, updateOpts)
	if res != nil && res.MatchedCount == 0 {
		err = ErrNoSuchDocuments
	}

	if err != nil {
		return
	}

	if len(opts) > 0 && opts[0].UpdateHook != nil {
		if err = hooks.On(ctx, opts[0].UpdateHook, hooks.AfterUpdate); err != nil {
			return
		}
	}

	return
}

func (c *Collection) UpdateOne(filter interface{}, update interface{}, opts ...opts.UpdateOptions) error {
	return c.UpdateOneWithCtx(context.TODO(), filter, update, opts...)
}

func (c *Collection) UpdateOneWithCtx(ctx context.Context, filter interface{}, update interface{}, opts ...opts.UpdateOptions) (err error) {
	updateOpts := options.Update()

	if len(opts) > 0 {
		if opts[0].UpdateOptions != nil {
			updateOpts = opts[0].UpdateOptions
		}
		if opts[0].UpdateHook != nil {
			if err = hooks.On(ctx, opts[0].UpdateHook, hooks.BeforeUpdate); err != nil {
				return
			}
		}
	}

	var res *mongo.UpdateResult

	res, err = c.collection.UpdateOne(ctx, filter, update, updateOpts)
	if res != nil && res.MatchedCount == 0 {
		if updateOpts.Upsert == nil || !*updateOpts.Upsert {
			err = ErrNoSuchDocuments
		}
	}

	if err != nil {
		return err
	}

	if len(opts) > 0 && opts[0].UpdateHook != nil {
		if err = hooks.On(ctx, opts[0].UpdateHook, hooks.AfterUpdate); err != nil {
			return
		}
	}

	return
}

func (c *Collection) UpdateAll(filter interface{}, update interface{}, opts ...opts.UpdateOptions) (*UpdateResult, error) {
	return c.UpdateAllWithCtx(context.TODO(), filter, update, opts...)
}

func (c *Collection) UpdateAllWithCtx(ctx context.Context, filter interface{}, update interface{}, opts ...opts.UpdateOptions) (result *UpdateResult, err error) {
	updateOpts := options.Update()
	if len(opts) > 0 {
		if opts[0].UpdateOptions != nil {
			updateOpts = opts[0].UpdateOptions
		}
		if opts[0].UpdateHook != nil {
			if err = hooks.On(ctx, opts[0].UpdateHook, hooks.BeforeUpdate); err != nil {
				return
			}
		}
	}

	var res *mongo.UpdateResult

	res, err = c.collection.UpdateMany(ctx, filter, update, updateOpts)
	if res != nil {
		result = translateUpdateResult(res)
	}

	if err != nil {
		return
	}

	if len(opts) > 0 && opts[0].UpdateHook != nil {
		if err = hooks.On(ctx, opts[0].UpdateHook, hooks.AfterUpdate); err != nil {
			return
		}
	}

	return
}

func (c *Collection) Upsert(filter interface{}, replacement interface{}, opts ...opts.UpsertOptions) (*UpdateResult, error) {
	return c.UpsertWithCtx(context.TODO(), filter, replacement, opts...)
}

func (c *Collection) UpsertWithCtx(ctx context.Context, filter interface{}, replacement interface{}, opts ...opts.UpsertOptions) (result *UpdateResult, err error) {
	h := replacement
	officialOpts := options.Replace().SetUpsert(true)

	if len(opts) > 0 {
		if opts[0].ReplaceOptions != nil {
			opts[0].ReplaceOptions.SetUpsert(true)
			officialOpts = opts[0].ReplaceOptions
		}
		if opts[0].UpsertHook != nil {
			h = opts[0].UpsertHook
		}
	}

	if err = hooks.On(ctx, replacement, hooks.BeforeUpsert, h); err != nil {
		return
	}

	var res *mongo.UpdateResult

	res, err = c.collection.ReplaceOne(ctx, filter, replacement, officialOpts)
	if res != nil {
		result = translateUpdateResult(res)
	}

	if err != nil {
		return
	}

	if err = hooks.On(ctx, replacement, hooks.AfterUpsert, h); err != nil {
		return
	}

	return
}

func (c *Collection) UpsertById(id interface{}, replacement interface{}, opts ...opts.UpsertOptions) (*UpdateResult, error) {
	return c.UpsertByIdWithCtx(context.TODO(), id, replacement, opts...)
}

func (c *Collection) UpsertByIdWithCtx(ctx context.Context, id interface{}, replacement interface{}, opts ...opts.UpsertOptions) (result *UpdateResult, err error) {
	h := replacement
	officialOpts := options.Replace().SetUpsert(true)

	if len(opts) > 0 {
		if opts[0].ReplaceOptions != nil {
			opts[0].ReplaceOptions.SetUpsert(true)
			officialOpts = opts[0].ReplaceOptions
		}
		if opts[0].UpsertHook != nil {
			h = opts[0].UpsertHook
		}
	}

	if err = hooks.On(ctx, replacement, hooks.BeforeUpsert, h); err != nil {
		return
	}

	var res *mongo.UpdateResult

	res, err = c.collection.ReplaceOne(ctx, bson.M{"_id": id}, replacement, officialOpts)
	if res != nil {
		result = translateUpdateResult(res)
	}

	if err != nil {
		return
	}

	if err = hooks.On(ctx, replacement, hooks.AfterUpsert, h); err != nil {
		return
	}

	return
}

func (c *Collection) ReplaceOne(filter interface{}, doc interface{}, opts ...opts.ReplaceOptions) error {
	return c.ReplaceOneWithCtx(context.TODO(), filter, doc, opts...)

}

func (c *Collection) ReplaceOneWithCtx(ctx context.Context, filter interface{}, doc interface{}, opts ...opts.ReplaceOptions) (err error) {
	h := doc
	replaceOpts := options.Replace()

	if len(opts) > 0 {
		if opts[0].ReplaceOptions != nil {
			replaceOpts = opts[0].ReplaceOptions
		}
		if opts[0].UpdateHook != nil {
			h = opts[0].UpdateHook
		}
	}

	if err = hooks.On(ctx, doc, hooks.BeforeReplace, h); err != nil {
		return
	}

	var res *mongo.UpdateResult

	res, err = c.collection.ReplaceOne(ctx, filter, doc, replaceOpts)
	if res != nil && res.MatchedCount == 0 {
		err = ErrNoSuchDocuments
	}

	if err != nil {
		return
	}

	if err = hooks.On(ctx, doc, hooks.AfterReplace, h); err != nil {
		return
	}

	return
}

func translateUpdateResult(res *mongo.UpdateResult) (result *UpdateResult) {
	result = &UpdateResult{
		MatchedCount:  res.MatchedCount,
		ModifiedCount: res.ModifiedCount,
		UpsertedCount: res.UpsertedCount,
		UpsertedID:    res.UpsertedID,
	}
	return
}
