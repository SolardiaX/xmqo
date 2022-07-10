/*
 * Copyright (c) 2022. All rights reserved by XtraVisions.
 */

package xmgo

type InsertOneResult struct {
	InsertedID interface{}
}

type InsertManyResult struct {
	InsertedIDs []interface{}
}

type UpdateResult struct {
	MatchedCount  int64
	ModifiedCount int64
	UpsertedCount int64
	UpsertedID    interface{}
}

type DeleteResult struct {
	DeletedCount int64
}
