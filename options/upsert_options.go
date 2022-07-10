/*
 * Copyright (c) 2022. All rights reserved by XtraVisions.
 */

package options

import "go.mongodb.org/mongo-driver/mongo/options"

type UpsertOptions struct {
	UpsertHook interface{}
	*options.ReplaceOptions
}
