// Package repr attempts to represent Go values in a form as close to real Go code as possible.
//
// Some values can not be represented directly, specifically pointers to basic types as this
// requires two steps; new() then assign. These values are represented as `&<value>`. eg. &23
package repr

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
)

// Repr returns a string representing v.
func Repr(v interface{}) string {
	w := bytes.NewBuffer(nil)
	Write(w, v)
	return w.String()
}

// Print v to os.Stdout on one line.
func Print(v interface{}) {
	Write(os.Stdout, v)
	fmt.Fprintln(os.Stdout)
}

// Write writes a representation of v to w.
func Write(w io.Writer, v interface{}) {
	reprValue(w, reflect.ValueOf(v), "")
}

func reprValue(w io.Writer, v reflect.Value, indent string) {
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		fmt.Fprintf(w, "%s{\n", v.Type())
		for i := 0; i < v.Len(); i++ {
			e := v.Index(i)
			fmt.Fprintf(w, "%s  ", indent)
			reprValue(w, e, indent+"  ")
			if i != v.Len()-1 {
				fmt.Fprintf(w, ",\n")
			}
		}
		fmt.Fprintf(w, ",\n%s}", indent)
	case reflect.Chan:
		fmt.Fprintf(w, "make(")
		fmt.Fprintf(w, "%s", v.Type())
		fmt.Fprintf(w, ", %d)", v.Cap())
	case reflect.Map:
		fmt.Fprintf(w, "%s{\n", v.Type())
		for i, k := range v.MapKeys() {
			kv := v.MapIndex(k)
			fmt.Fprintf(w, "%s  ", indent)
			reprValue(w, k, indent+"  ")
			fmt.Fprintf(w, ": ")
			reprValue(w, kv, indent)
			if i != v.Len()-1 {
				fmt.Fprintf(w, ",\n")
			}
		}
		fmt.Fprintf(w, ",\n%s}", indent)
	case reflect.Struct:
		fmt.Fprintf(w, "%s{\n", v.Type())
		for i := 0; i < v.NumField(); i++ {
			t := v.Type().Field(i)
			f := v.Field(i)
			fmt.Fprintf(w, "%s%s: ", indent+"  ", t.Name)
			reprValue(w, f, indent+"  ")
			if i != v.NumField()-1 {
				fmt.Fprintf(w, ",\n")
			}
		}
		fmt.Fprintf(w, ",\n%s}", indent)
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Fprintf(w, "nil")
			return
		}
		fmt.Fprintf(w, "&")
		reprValue(w, v.Elem(), indent)
	case reflect.String:
		fmt.Fprintf(w, "%q", v.Interface())
	default:
		fmt.Fprintf(w, "%v", v)
	}
}
