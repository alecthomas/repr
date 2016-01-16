// Package repr attempts to represent Go values in a form that can be copy-and-pasted into source
// code directly.
//
// Some values can not be represented directly, specifically pointers to basic types, so they are
// represented as `&<value>`. eg. &23
package repr

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

// Repr returns a string representing v.
func Repr(v interface{}) string {
	w := bytes.NewBuffer(nil)
	Write(w, v)
	return w.String()
}

// Write writes a representation of v to w.
func Write(w io.Writer, v interface{}) {
	reprValue(w, reflect.ValueOf(v))
}

func reprValue(w io.Writer, v reflect.Value) {
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		fmt.Fprintf(w, "%s", v.Type())
		fmt.Fprintf(w, "{")
		for i := 0; i < v.Len(); i++ {
			e := v.Index(i)
			reprValue(w, e)
			if i != v.Len()-1 {
				fmt.Fprintf(w, ", ")
			}
		}
		fmt.Fprintf(w, "}")
	case reflect.Chan:
		fmt.Fprintf(w, "make(")
		fmt.Fprintf(w, "%s", v.Type())
		fmt.Fprintf(w, ", %d)", v.Cap())
	case reflect.Map:
		fmt.Fprintf(w, "%s", v.Type())
		fmt.Fprintf(w, "{")
		for i, k := range v.MapKeys() {
			kv := v.MapIndex(k)
			reprValue(w, k)
			fmt.Fprintf(w, ": ")
			reprValue(w, kv)
			if i != v.Len()-1 {
				fmt.Fprintf(w, ", ")
			}
		}
		fmt.Fprintf(w, "}")
	case reflect.Struct:
		fmt.Fprintf(w, "%s{", v.Type())
		for i := 0; i < v.NumField(); i++ {
			t := v.Type().Field(i)
			f := v.Field(i)
			fmt.Fprintf(w, "%q: ", t.Name)
			reprValue(w, f)
			if i != v.NumField()-1 {
				fmt.Fprintf(w, ", ")
			}
		}
		fmt.Fprintf(w, "}")
	case reflect.Ptr:
		fmt.Fprintf(w, "&")
		reprValue(w, v.Elem())
	case reflect.String:
		fmt.Fprintf(w, "%q", v.Interface())
	default:
		fmt.Fprintf(w, "%v", v)
	}
}
