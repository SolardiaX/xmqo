/*
 * Copyright (c) 2022. All rights reserved by XtraVisions.
 */

package options

import "go.mongodb.org/mongo-driver/mongo/options"

// IndexOptions 索引配置
type IndexOptions struct {
	Key []string // Index key fields; prefix name with dash (-) for descending order
	*options.IndexOptions
}
