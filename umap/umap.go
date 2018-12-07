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
func MapValues(from interface{}, to interface{}, m Map) (err error) {
	f, t := valueOf(from), valueOf(to)
	if f.Kind() != reflect.Struct || t.Kind() != reflect.Struct {
		return fmt.Errorf("%q, and %q must be structs",
			"from", "to")
	}

	for fk, tk := range m {
		ff, tf := f.FieldByName(fk), t.FieldByName(tk)
		if tf.Kind() != ff.Kind() {
			return fmt.Errorf("%s(%s) mismatches %s(%s)",
				fk, ff.Kind(), tk, tf.Kind())
		}

		switch ff.Kind() {

		// True, False
		case reflect.Bool:
			v := ff.Bool()
			tf.SetBool(v)
			continue

		// Works for all types of int.. magic 🤷‍♂️
		case reflect.Int8, reflect.Int16, reflect.Int32,
			reflect.Int64, reflect.Int:

			v := ff.Int()
			tf.SetInt(v)
			continue

		// Strings
		case reflect.String:
			v := ff.String()
			tf.SetString(v)
			continue

		/**
		 * This should be rare
		 * but if there is a type
		 * send a bug report
		 */
		default:
			// I don't support every type bruv.
			return fmt.Errorf("unsupported %s(%s), %s(%s)",
				fk, ff.Kind(), tk, tf.Kind())
		}
	}

	return
}
