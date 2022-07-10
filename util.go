/*
 * Copyright (c) 2022. All rights reserved by XtraVisions.
 */

package xmgo

import (
	"math"
	"strconv"
	"strings"
)

func splitSortField(field string) (key string, sort int32) {
	sort = 1
	key = field

	if len(field) != 0 {
		switch field[0] {
		case '+':
			key = strings.TrimPrefix(field, "+")
			sort = 1
		case '-':
			key = strings.TrimPrefix(field, "-")
			sort = -1
		}
	}

	return key, sort
}

func compareVersions(v1 string, v2 string) (int, error) {
	n1 := strings.Split(v1, ".")
	n2 := strings.Split(v2, ".")

	for i := 0; i < int(math.Min(float64(len(n1)), float64(len(n2)))); i++ {
		i1, err := strconv.Atoi(n1[i])
		if err != nil {
			return 0, err
		}
		i2, err := strconv.Atoi(n2[i])
		if err != nil {
			return 0, err
		}
		difference := i1 - i2
		if difference != 0 {
			return difference, nil
		}
	}

	return 0, nil
}
