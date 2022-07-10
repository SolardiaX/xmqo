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

// Database mongodb 数据库
type Database struct {
	database *mongo.Database
	registry *bsoncodec.Registry
}

// Name 获取当前数据库名
func (d *Database) Name() string {
	return d.database.Name()
}

// Collection 获取指定名称的 mongodb collection
//	@param name mongodb collection 名称
func (d *Database) Collection(name string) *Collection {
	var cp *mongo.Collection
	cp = d.database.Collection(name)

	return &Collection{
		collection: cp,
		registry:   d.registry,
	}
}

// ModelCollection 获取 IModel 实体对应的 mongodb collection
//	@param model IModel 实体
func (d *Database) ModelCollection(model IModel) *Collection {
	return d.Collection(model.CollectionName())
}

// Drop 使用默认上下文删除当前数据库
func (d *Database) Drop() error {
	return d.DropWithCtx(context.TODO())
}

// DropWithCtx 删除当前数据库
func (d *Database) DropWithCtx(ctx context.Context) error {
	return d.database.Drop(ctx)
}

// RunCommand 使用默认上下文在当前数据库直接执行
//	参见 https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Database.RunCommand
func (d *Database) RunCommand(runCommand interface{}, opts ...opts.RunCommandOptions) *mongo.SingleResult {
	option := options.RunCmd()
	if len(opts) > 0 && opts[0].RunCmdOptions != nil {
		option = opts[0].RunCmdOptions
	}

	return d.database.RunCommand(context.TODO(), runCommand, option)
}

// RunCommandWithCtx 在当前数据库直接执行
//	参见 https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Database.RunCommand
func (d *Database) RunCommandWithCtx(ctx context.Context, runCommand interface{}, opts ...opts.RunCommandOptions) *mongo.SingleResult {
	option := options.RunCmd()
	if len(opts) > 0 && opts[0].RunCmdOptions != nil {
		option = opts[0].RunCmdOptions
	}

	return d.database.RunCommand(ctx, runCommand, option)
}
