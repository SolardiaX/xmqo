/*
 * Copyright (c) 2022. All rights reserved by XtraVisions.
 */

package xmgo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver"

	opts "xtravisions.com/xmgo/options"
)

type Session struct {
	session mongo.Session
}

func (s *Session) StartTransaction(ctx context.Context, cb func(sessCtx context.Context) (interface{}, error), opts ...*opts.TransactionOptions) (interface{}, error) {
	transactionOpts := options.Transaction()
	if len(opts) > 0 && opts[0].TransactionOptions != nil {
		transactionOpts = opts[0].TransactionOptions
	}
	result, err := s.session.WithTransaction(ctx, wrapperCustomCb(cb), transactionOpts)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Session) EndSession(ctx context.Context) {
	s.session.EndSession(ctx)
}

func (s *Session) AbortTransaction(ctx context.Context) error {
	return s.session.AbortTransaction(ctx)
}

func wrapperCustomCb(cb func(ctx context.Context) (interface{}, error)) func(sessCtx mongo.SessionContext) (interface{}, error) {
	return func(sessCtx mongo.SessionContext) (interface{}, error) {
		result, err := cb(sessCtx)
		if err == ErrTransactionRetry {
			return nil, mongo.CommandError{Labels: []string{driver.TransientTransactionError}}
		}
		return result, err
	}
}
