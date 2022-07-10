/*
 * Copyright (c) 2022. All rights reserved by XtraVisions.
 */

package xmgo

import (
	"context"
	"time"
)

type IModel interface {
	CollectionName() string
}

type BaseModel struct {
	Id        ObjectId  `json:"id,string" bson:"_id"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

func (m *BaseModel) BeforeInsert(_ context.Context) error {
	m.Id = NewObjectId()
	m.CreatedAt = time.Now().Local()
	m.UpdatedAt = time.Now().Local()

	return nil
}

func (m *BaseModel) BeforeUpdate(_ context.Context) error {
	m.UpdatedAt = time.Now().Local()

	return nil
}

func (m *BaseModel) BeforeUpsert(_ context.Context) error {
	if m.Id.IsZero() {
		m.Id = NewObjectId()
	}

	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now().Local()
	}

	m.UpdatedAt = time.Now().Local()

	return nil
}
