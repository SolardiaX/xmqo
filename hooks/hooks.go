/*
 * Copyright (c) 2022. All rights reserved by XtraVisions.
 */

package hooks

import (
	"context"
	"reflect"
)

var hookHandler = map[OpType]func(ctx context.Context, hook interface{}) error{
	BeforeInsert:  beforeInsert,
	AfterInsert:   afterInsert,
	BeforeUpdate:  beforeUpdate,
	AfterUpdate:   afterUpdate,
	BeforeQuery:   beforeQuery,
	AfterQuery:    afterQuery,
	BeforeRemove:  beforeRemove,
	AfterRemove:   afterRemove,
	BeforeUpsert:  beforeUpsert,
	AfterUpsert:   afterUpsert,
	BeforeReplace: beforeUpdate,
	AfterReplace:  afterUpdate,
}

func On(ctx context.Context, hook interface{}, opType OpType, opts ...interface{}) error {
	if len(opts) > 0 {
		hook = opts[0]
	}

	to := reflect.TypeOf(hook)
	if to == nil {
		return nil
	}
	switch to.Kind() {
	case reflect.Slice:
		return sliceHandle(ctx, hook, opType)
	case reflect.Ptr:
		v := reflect.ValueOf(hook).Elem()
		switch v.Kind() {
		case reflect.Slice:
			return sliceHandle(ctx, v.Interface(), opType)
		default:
			return do(ctx, hook, opType)
		}
	default:
		return do(ctx, hook, opType)
	}
}

func sliceHandle(ctx context.Context, hook interface{}, opType OpType) error {
	// []interface{}{UserType{}...}
	if h, ok := hook.([]interface{}); ok {
		for _, v := range h {
			if err := do(ctx, v, opType); err != nil {
				return err
			}
		}
		return nil
	}
	// []UserType{}
	s := reflect.ValueOf(hook)
	for i := 0; i < s.Len(); i++ {
		if err := do(ctx, s.Index(i).Interface(), opType); err != nil {
			return err
		}
	}
	return nil
}

type BeforeInsertHook interface {
	BeforeInsert(ctx context.Context) error
}

type AfterInsertHook interface {
	AfterInsert(ctx context.Context) error
}

func beforeInsert(ctx context.Context, hook interface{}) error {
	if ih, ok := hook.(BeforeInsertHook); ok {
		return ih.BeforeInsert(ctx)
	}
	return nil
}

func afterInsert(ctx context.Context, hook interface{}) error {
	if ih, ok := hook.(AfterInsertHook); ok {
		return ih.AfterInsert(ctx)
	}
	return nil
}

type BeforeUpdateHook interface {
	BeforeUpdate(ctx context.Context) error
}

type AfterUpdateHook interface {
	AfterUpdate(ctx context.Context) error
}

func beforeUpdate(ctx context.Context, hook interface{}) error {
	if ih, ok := hook.(BeforeUpdateHook); ok {
		return ih.BeforeUpdate(ctx)
	}
	return nil
}

func afterUpdate(ctx context.Context, hook interface{}) error {
	if ih, ok := hook.(AfterUpdateHook); ok {
		return ih.AfterUpdate(ctx)
	}
	return nil
}

type BeforeQueryHook interface {
	BeforeQuery(ctx context.Context) error
}

type AfterQueryHook interface {
	AfterQuery(ctx context.Context) error
}

func beforeQuery(ctx context.Context, hook interface{}) error {
	if ih, ok := hook.(BeforeQueryHook); ok {
		return ih.BeforeQuery(ctx)
	}
	return nil
}

func afterQuery(ctx context.Context, hook interface{}) error {
	if ih, ok := hook.(AfterQueryHook); ok {
		return ih.AfterQuery(ctx)
	}
	return nil
}

type BeforeRemoveHook interface {
	BeforeRemove(ctx context.Context) error
}

type AfterRemoveHook interface {
	AfterRemove(ctx context.Context) error
}

func beforeRemove(ctx context.Context, hook interface{}) error {
	if ih, ok := hook.(BeforeRemoveHook); ok {
		return ih.BeforeRemove(ctx)
	}
	return nil
}

func afterRemove(ctx context.Context, hook interface{}) error {
	if ih, ok := hook.(AfterRemoveHook); ok {
		return ih.AfterRemove(ctx)
	}
	return nil
}

type BeforeUpsertHook interface {
	BeforeUpsert(ctx context.Context) error
}

type AfterUpsertHook interface {
	AfterUpsert(ctx context.Context) error
}

func beforeUpsert(ctx context.Context, hook interface{}) error {
	if ih, ok := hook.(BeforeUpsertHook); ok {
		return ih.BeforeUpsert(ctx)
	}
	return nil
}

func afterUpsert(ctx context.Context, hook interface{}) error {
	if ih, ok := hook.(AfterUpsertHook); ok {
		return ih.AfterUpsert(ctx)
	}
	return nil
}

func do(ctx context.Context, hook interface{}, opType OpType) error {
	if f, ok := hookHandler[opType]; !ok {
		return nil
	} else {
		return f(ctx, hook)
	}
}
