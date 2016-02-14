package repr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type anotherStruct struct {
	A []int
}

type testStruct struct {
	S string
	I *int
	A anotherStruct
}

func TestReprEmptyArray(t *testing.T) {
	assert.Equal(t, "[]string{}", Repr([]string{}))
}

func TestReprStringArray(t *testing.T) {
	assert.Equal(t, "[]string{\"a\", \"b\"}", Repr([]string{"a", "b"}))
}

func TestReprIntArray(t *testing.T) {
	assert.Equal(t, "[]int{1, 2}", Repr([]int{1, 2}))
}

func TestReprPointerToInt(t *testing.T) {
	pi := new(int)
	*pi = 13
	assert.Equal(t, `&13`, Repr(pi))
}

func TestReprChannel(t *testing.T) {
	ch := make(<-chan map[string]*testStruct, 1)
	assert.Equal(t, `make(<-chan map[string]*repr.testStruct, 1)`, Repr(ch))
}

func TestReprEmptyMap(t *testing.T) {
	assert.Equal(t, "map[string]bool{}", Repr(map[string]bool{}))
}

func TestReprMap(t *testing.T) {
	m := map[string]int{"a": 1}
	assert.Equal(t, "map[string]int{\"a\": 1}", Repr(m))
}

func TestReprStructWithIndent(t *testing.T) {
	pi := new(int)
	*pi = 13
	s := &testStruct{
		S: "String",
		I: pi,
		A: anotherStruct{
			A: []int{1, 2, 3},
		},
	}
	assert.Equal(t, `&repr.testStruct{
  S: "String",
  I: &13,
  A: repr.anotherStruct{
    A: []int{
      1,
      2,
      3,
    },
  },
}`, Repr(s, Indent("  ")))

}

func TestReprByteArray(t *testing.T) {
	b := []byte{1, 2, 3}
	assert.Equal(t, `[]uint8{1, 2, 3}`, Repr(b))
}

type privateTestStruct struct {
	a string
}

func TestReprPrivateField(t *testing.T) {
	s := privateTestStruct{"hello"}
	assert.Equal(t, `repr.privateTestStruct{a: "hello"}`, Repr(s))
}
