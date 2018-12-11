// Copyright 2018 Jordon Bedwell. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package umap

import (
	"fmt"
	"reflect"
)

// KeyMap is the key -> key map between
// the two structs we are mapping out from
// upstream into the downstream.
type KeyMap map[string]string

// structOK makes sure the value is a struct
// we can only map struct to struct, mostly because
// we only want to map struct to struct
func structOK(v reflect.Value) error {
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("input must be a struct")
	}

	return nil
}

// nPtrOf converts a Ptr to the actual
func nPtrOf(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		return v.Elem()
	}

	return v
}

// valueOf gives back the reflect.Value of the interface
func valueOf(i interface{}) (v reflect.Value, err error) {
	v, ok := i.(reflect.Value)
	if !ok {
		v = reflect.ValueOf(i)
	}

	v = nPtrOf(v)
	err = structOK(v)
	return
}

// Map maps an upstream struct from a downstream struct
// containting an upstream field (as a wrapper) into the
// named downstream that it's mapping upstream.
//
// 	type MyStruct struct { upstream Upstream, key string }
//	type Upstream struct { key string, val string }
//	type UpstreamMap map[string]string {
//		"key" => "key"
//	}
//
//	Map(&MyStruct{
// 		upstream: &Upstream{
//			key: "hello",
//			val: "world",
//		}
//	}, UpstreamMap)
func Map(in interface{}, keymap KeyMap) (err error) {
	var upstream, downstream reflect.Value
	downstream, err = valueOf(in)
	if err == nil {
		val := downstream.FieldByName("upstream")
		upstream, err = valueOf(val)
		if err != nil {
			return
		}
	}

	for fromKey, toKey := range keymap {
		from, to := upstream.FieldByName(fromKey), downstream.FieldByName(toKey)
		if to.Kind() != from.Kind() {
			return fmt.Errorf("%s(%s) mismatches %s(%s)",
				fromKey, from.Kind(), toKey,
				to.Kind())
		}

		switch from.Kind() {

		/**
		 * True, False
		 * Duh
		 */
		case reflect.Bool:
			to.SetBool(from.Bool())
			continue

		/**
		 * Works for all types of int..
		 * It's magic 🤷‍♂️
		 */
		case reflect.Int8, reflect.Int16, reflect.Int32,
			reflect.Int64, reflect.Int:
			to.SetInt(from.Int())
			continue

		/**
		 * Strings
		 * Duh
		 */
		case reflect.String:
			to.SetString(from.String())
			continue

		/**
		 * This should be rare
		 * but if there is a type
		 * send a bug report
		 */
		default:
			// I don't support every type bruv.
			return fmt.Errorf("unsupported %s(%s), %s(%s)",
				fromKey, from.Kind(), toKey,
				to.Kind())
		}
	}

	return
}
