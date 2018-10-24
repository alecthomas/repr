package repr

import (
	"strings"
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
	assert.Equal(t, "[]string{}", String([]string{}, OmitEmpty(false)))
}

func TestReprStringArray(t *testing.T) {
	assert.Equal(t, "[]string{\"a\", \"b\"}", String([]string{"a", "b"}))
}

func TestReprIntArray(t *testing.T) {
	assert.Equal(t, "[]int{1, 2}", String([]int{1, 2}))
}

func TestReprPointerToInt(t *testing.T) {
	pi := new(int)
	*pi = 13
	assert.Equal(t, `&13`, String(pi))
}

func TestReprChannel(t *testing.T) {
	ch := make(<-chan map[string]*testStruct, 1)
	assert.Equal(t, `make(<-chan map[string]*repr.testStruct, 1)`, String(ch))
}

func TestReprEmptyMap(t *testing.T) {
	assert.Equal(t, "map[string]bool{}", String(map[string]bool{}))
}

func TestReprMap(t *testing.T) {
	m := map[string]int{"a": 1}
	assert.Equal(t, "map[string]int{\"a\": 1}", String(m))
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
}`, String(s, Indent("  ")))

}

func TestReprByteArray(t *testing.T) {
	b := []byte{1, 2, 3}
	assert.Equal(t, "[]byte(\"\\x01\\x02\\x03\")", String(b))
}

type privateTestStruct struct {
	a string
}

func TestReprPrivateField(t *testing.T) {
	s := privateTestStruct{"hello"}
	assert.Equal(t, `repr.privateTestStruct{a: "hello"}`, String(s))
}

func TestReprNilAlone(t *testing.T) {
	var err error
	s := String(err)
	assert.Equal(t, "nil", s)
}

func TestReprNilInsideArray(t *testing.T) {
	arr := []*privateTestStruct{{"hello"}, nil}
	s := String(arr)
	assert.Equal(t, "[]*repr.privateTestStruct{&repr.privateTestStruct{a: \"hello\"}, nil}", s)
}

type Enum int

func (e Enum) String() string {
	return "Value"
}

func TestEnum(t *testing.T) {
	v := Enum(1)
	s := String(v)
	assert.Equal(t, "repr.Enum(Value)", s)
}

func TestShowType(t *testing.T) {
	a := map[string]privateTestStruct{"foo": {"bar"}}
	s := String(a, AlwaysIncludeType(), Indent("  "))
	t.Log(s)
	assert.Equal(t, strings.TrimSpace(`
map[string]repr.privateTestStruct{
  string("foo"): repr.privateTestStruct{
    a: string("bar"),
  },
}
`), s)
}
