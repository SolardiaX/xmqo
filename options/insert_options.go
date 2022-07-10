/*
 * Copyright (c) 2022. All rights reserved by XtraVisions.
 */

package options

import "go.mongodb.org/mongo-driver/mongo/options"

type InsertOneOptions struct {
	InsertHook interface{}
	*options.InsertOneOptions
}
type InsertManyOptions struct {
	InsertHook interface{}
	*options.InsertManyOptions
}
