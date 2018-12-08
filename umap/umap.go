// Copyright 2018 Jordon Bedwell. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package umap

import (
	"fmt"
	"reflect"
)

// Map allows you to k,k map from an
// upstream struct into a downstream cleaned
// struct, without much labor
type Map map[string]string

// valueOf gives back the reflect.Value, but
// if it happens to be a ptr, we'll run extra
// step and get the reflect.Value of the ptr
func valueOf(i interface{}) (v reflect.Value) {
	v = reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return
}

// MapValues allows you to map upstream
func MapValues(_from, _to interface{}, m Map) (err error) {
	from, to := valueOf(_from), valueOf(_to)
	if from.Kind() != reflect.Struct || to.Kind() != reflect.Struct {
		return fmt.Errorf("%q, and %q must be structs",
			"from", "to")
	}

	for fromKey, toKey := range m {
		fromField, toField := from.FieldByName(fromKey), to.FieldByName(toKey)
		if toField.Kind() != fromField.Kind() {
			return fmt.Errorf("%s(%s) mismatches %s(%s)",
				fromKey, fromField.Kind(), toKey,
				toField.Kind())
		}

		switch fromField.Kind() {

		/**
		 * True, False
		 * Duh
		 */
		case reflect.Bool:
			toField.SetBool(fromField.Bool())
			continue

		/**
		 * Works for all types of int..
		 * It's magic 🤷‍♂️
		 */
		case reflect.Int8, reflect.Int16, reflect.Int32,
			reflect.Int64, reflect.Int:

			toField.SetInt(fromField.Int())
			continue

		/**
		 * Strings
		 * Duh
		 */
		case reflect.String:
			toField.SetString(fromField.String())
			continue

		/**
		 * This should be rare
		 * but if there is a type
		 * send a bug report
		 */
		default:
			// I don't support every type bruv.
			return fmt.Errorf("unsupported %s(%s), %s(%s)",
				fromKey, fromField.Kind(), toKey,
				toField.Kind())
		}
	}

	return
}
