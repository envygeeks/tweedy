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
	actual := valueOf(val).Kind()
	assert.Equal(t, reflect.Struct,
		actual)
}

func TestMapValues(t *testing.T) {
	a := &struct {
		A int
		B string
		C bool
	}{
		A: 1,
		B: "Hello",
		C: true,
	}

	b := &struct {
		A int
		B string
		C bool
	}{}

	c := Map{
		"A": "A",
		"B": "B",
		"C": "C",
	}

	MapValues(a, b, c)
	assert.Equal(t, a.A, b.A)
	assert.Equal(t, a.B, b.B)
	assert.Equal(t, a.C, b.C)
}
