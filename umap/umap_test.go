// Copyright 2018 Jordon Bedwell. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package umap

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValueOf(t *testing.T) {
	val := &struct{}{}
	actual, err := valueOf(val)
	assert.Equal(t, reflect.Struct, actual.Kind())
	assert.Nil(t, err)
}

func TestMap(t *testing.T) {
	type a struct {
		A int
		B string
		C bool
	}

	type b struct {
		upstream *a

		A int
		B string
		C bool
	}

	upstream := &a{A: 1, B: "Hello", C: true}
	obj := &b{upstream: upstream}
	err := Map(obj, KeyMap{
		"A": "A",
		"B": "B",
		"C": "C",
	})

	assert.Nil(t, err)
	assert.Equal(t, upstream.A, obj.A)
	assert.Equal(t, upstream.B, obj.B)
	assert.Equal(t, upstream.C, obj.C)
}
