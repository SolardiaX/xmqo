/*
 * Copyright (c) 2022. All rights reserved by XtraVisions.
 */

package options

// IndexOptions 索引配置
type IndexOptions struct {
	Key                []string // Index key fields; prefix name with dash (-) for descending order
	Unique             bool     // Prevent two documents from having the same index key
	Background         bool     // Build index in background and return immediately
	Sparse             bool     // Only index documents containing the Key fields
	ExpireAfterSeconds *int32   // Periodically delete docs with indexed time.Time older than that.
}
