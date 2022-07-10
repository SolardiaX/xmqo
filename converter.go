/*
 * Copyright (c) 2022. All rights reserved by XtraVisions.
 */

package xmgo

import (
	"github.com/jinzhu/copier"
)

type TypeConverter = copier.TypeConverter

var DefaultConverters = []TypeConverter{
	{
		SrcType: "",
		DstType: ObjectId{},
		Fn: func(src interface{}) (dst interface{}, err error) {
			id := src.(string)
			dst = GetObjectId(id)

			return
		},
	},
	{
		SrcType: ObjectId{},
		DstType: "",
		Fn: func(src interface{}) (dst interface{}, err error) {
			id := src.(ObjectId)
			dst = id.String()
			return
		},
	},
}

var modelConverters = make(map[string][]TypeConverter, 0)

func RegisterConverterWithDefault(typ string, converters []TypeConverter) {
	ct := append(converters, DefaultConverters...)
	RegisterConverter(typ, ct)
}

func RegisterConverter(typ string, converters []TypeConverter) {
	var ct []copier.TypeConverter
	var ok bool

	if ct, ok = modelConverters[typ]; !ok {
		ct = make([]copier.TypeConverter, 0)
	}

	ct = append(ct, converters...)
	modelConverters[typ] = ct
}

func Convert(typ string, from interface{}, to interface{}) error {
	converters := modelConverters[typ]
	if converters == nil {
		converters = DefaultConverters
	}

	return copier.CopyWithOption(to, from, copier.Option{Converters: converters})
}
