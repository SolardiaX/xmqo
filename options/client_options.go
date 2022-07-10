/*
 * Copyright (c) 2022. All rights reserved by XtraVisions.
 */

package options

import "go.mongodb.org/mongo-driver/mongo/options"

// ClientOptions 源生连接参数引用
//	参见 go.mongodb.org/mongo-driver/mongo/options/ClientOptions
type ClientOptions struct {
	*options.ClientOptions
}
