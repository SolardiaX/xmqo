/*
 * Copyright (c) 2022. All rights reserved by XtraVisions.
 */

package xmgo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type ICursor interface {
	Next(result interface{}) bool
	Close() error
	Err() error
	All(results interface{}) error
}

type Cursor struct {
	ctx    context.Context
	cursor *mongo.Cursor
	err    error
}

func (c *Cursor) Next(result interface{}) bool {
	if c.err != nil {
		return false
	}

	var err error
	if c.cursor.Next(c.ctx) {
		err = c.cursor.Decode(result)
		if err == nil {
			return true
		}
	}

	return false
}

func (c *Cursor) All(results interface{}) error {
	if c.err != nil {
		return c.err
	}

	return c.cursor.All(c.ctx, results)
}

func (c *Cursor) Close() error {
	if c.err != nil {
		return c.err
	}

	return c.cursor.Close(c.ctx)
}

func (c *Cursor) Err() error {
	if c.err != nil {
		return c.err
	}

	return c.cursor.Err()
}
