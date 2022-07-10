/*
 * Copyright (c) 2022. All rights reserved by XtraVisions.
 */

package xmgo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// alias mongo drive bson primitives
// thus user don't need to import go.mongodb.org/mongo-driver/mongo, it's all in qmgo
type (
	// M is an alias of bson.M
	M = bson.M
	// A is an alias of bson.A
	A = bson.A
	// D is an alias of bson.D
	D = bson.D
	// E is an alias of bson.E
	E = bson.E
	// ObjectId is an alias of primitive.ObjectID
	ObjectId = primitive.ObjectID
)

// NewObjectId 生成 ObjectID
func NewObjectId() primitive.ObjectID {
	return primitive.NewObjectID()
}

func GetObjectId(s string) (id ObjectId) {
	if primitive.IsValidObjectID(s) {
		id, _ = primitive.ObjectIDFromHex(s)
	} else {
		id = primitive.NewObjectID()
	}

	return
}

func IsNilObjectId(id ObjectId) bool {
	return id == primitive.NilObjectID
}
