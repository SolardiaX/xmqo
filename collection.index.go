/*
 * Copyright (c) 2022. All rights reserved by XtraVisions.
 */

package xmgo

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"

	opts "xtravisions.com/xmgo/options"
)

// DropIndex 使用默认上下文删除索引
//	@param indexes 待删除索引名
func (c *Collection) DropIndex(indexes []string) error {
	return c.DropIndexWithCtx(context.TODO(), indexes)
}

// DropIndexWithCtx 删除索引
//	@param indexes 待删除索引名
func (c *Collection) DropIndexWithCtx(ctx context.Context, indexes []string) (err error) {
	var res string
	for _, e := range indexes {
		key, sort := splitSortField(e)
		n := key + "_" + fmt.Sprint(sort)
		if len(res) == 0 {
			res = n
		} else {
			res += "_" + n
		}
	}

	_, err = c.collection.Indexes().DropOne(ctx, res)
	return
}

// DropAllIndex 使用默认上下文删除全部索引
func (c *Collection) DropAllIndex() error {
	return c.DropAllIndexWithCtx(context.TODO())
}

// DropAllIndexWithCtx 删除全部索引
func (c *Collection) DropAllIndexWithCtx(ctx context.Context) (err error) {
	_, err = c.collection.Indexes().DropAll(ctx)
	return
}

// CreateIndexes 使用默认上下文创建索引
func (c *Collection) CreateIndexes(indexes []opts.IndexOptions) error {
	return c.CreateIndexesWithCtx(context.TODO(), indexes)
}

// CreateIndexesWithCtx 创建索引
//	注意：不支持在 `local` 模式读策略下的操作
func (c *Collection) CreateIndexesWithCtx(ctx context.Context, indexes []opts.IndexOptions) error {
	return c.ensureIndex(ctx, indexes)
}

// EnsureIndexes 使用默认上下文确保使用索引
//	@param uniques 唯一索引
//	@param indexes 普通索引
func (c *Collection) EnsureIndexes(uniques []string, indexes []string) error {
	return c.EnsureIndexesWithCtx(context.TODO(), uniques, indexes)
}

// EnsureIndexesWithCtx 确保使用索引
//	注意：不支持在 `local` 模式读策略下的操作
//	@param ctx 上下文
//	@param uniques 唯一索引
//	@param indexes 普通索引
func (c *Collection) EnsureIndexesWithCtx(ctx context.Context, uniques []string, indexes []string) (err error) {
	var uniqueModel, indexesModel []opts.IndexOptions

	for _, v := range uniques {
		vv := strings.Split(v, ",")
		model := opts.IndexOptions{Key: vv, Unique: true}
		uniqueModel = append(uniqueModel, model)
	}

	if err = c.CreateIndexesWithCtx(ctx, uniqueModel); err != nil {
		return
	}

	for _, v := range indexes {
		vv := strings.Split(v, ",")
		model := opts.IndexOptions{Key: vv}
		indexesModel = append(indexesModel, model)
	}

	if err = c.CreateIndexesWithCtx(ctx, indexesModel); err != nil {
		return
	}
	return
}

func (c *Collection) ensureIndex(ctx context.Context, indexes []opts.IndexOptions) error {
	var indexModels []mongo.IndexModel
	for _, idx := range indexes {
		var model mongo.IndexModel
		var keysDoc bsonx.Doc

		for _, field := range idx.Key {
			key, n := splitSortField(field)

			keysDoc = keysDoc.Append(key, bsonx.Int32(n))
		}
		iOptions := options.Index().SetUnique(idx.Unique).SetBackground(idx.Background).SetSparse(idx.Sparse)
		if idx.ExpireAfterSeconds != nil {
			iOptions.SetExpireAfterSeconds(*idx.ExpireAfterSeconds)
		}
		model = mongo.IndexModel{
			Keys:    keysDoc,
			Options: iOptions,
		}

		indexModels = append(indexModels, model)
	}

	if len(indexModels) == 0 {
		return nil
	}

	res, err := c.collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil || len(res) == 0 {
		fmt.Println("<MongoDB.C>: ", c.collection.Name(), " Index: ", indexes, " error: ", err, "res: ", res)
		return err
	}

	return nil
}
