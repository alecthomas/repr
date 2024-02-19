package repr

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"
)

func equal(t *testing.T, want, have string) {
	t.Helper()
	if want != have {
		t.Errorf("\nWant: %q\nHave: %q", want, have)
	}
}

type anotherStruct struct {
	A []int
}

func (anotherStruct) String() string { return "anotherStruct" }

type testStruct struct {
	S string
	I *int
	A anotherStruct
}

type testStructWithInterfaceField struct {
	S string
	I fmt.Stringer
}

func TestHide(t *testing.T) {
	actual := testStruct{
		S: "str",
		A: anotherStruct{A: []int{1}},
	}
	equal(t, `repr.testStruct{S: "str"}`, String(actual, Hide[anotherStruct]()))
	equal(t, "repr.testStruct{\n  S: \"str\",\n}", String(actual, Indent("  "), Hide[anotherStruct]()))
	equal(t, "repr.testStructWithInterfaceField{S: \"str\"}",
		String(testStructWithInterfaceField{S: "str", I: anotherStruct{}}, Hide[fmt.Stringer]()))
}

func TestReprEmptyArray(t *testing.T) {
	equal(t, "[]string{}", String([]string{}, OmitEmpty(false)))
}

func TestReprEmptySliceMapFields(t *testing.T) {
	v := struct {
		S  []string
		M  map[string]string
		NZ []string
	}{[]string{}, map[string]string{}, []string{"a", "b"}}
	equal(t, `struct { S []string; M map[string]string; NZ []string }{NZ: []string{"a", "b"}}`, String(v, OmitEmpty(true)))
}

func TestReprStringArray(t *testing.T) {
	equal(t, "[]string{\"a\", \"b\"}", String([]string{"a", "b"}))
}

func TestReprIntArray(t *testing.T) {
	equal(t, "[]int{1, 2}", String([]int{1, 2}))
}

func TestReprPointerToInt(t *testing.T) {
	pi := new(int)
	*pi = 13
	equal(t, `&13`, String(pi))
}

func TestReprChannel(t *testing.T) {
	ch := make(<-chan map[string]*testStruct, 1)
	equal(t, `make(<-chan map[string]*repr.testStruct, 1)`, String(ch))
}

func TestReprEmptyMap(t *testing.T) {
	equal(t, "map[string]bool{}", String(map[string]bool{}))
}

func TestReprMap(t *testing.T) {
	m := map[string]int{"b": 3, "a": 1, "c": 5}
	for i := 0; i < 1000; i++ {
		equal(t, "map[string]int{\"a\": 1, \"b\": 3, \"c\": 5}", String(m))
	}
}

func TestReprIntMap(t *testing.T) {
	m := map[int]string{3: "b", 1: "a", 5: "c"}
	for i := 0; i < 1000; i++ {
		equal(t, "map[int]string{1: \"a\", 3: \"b\", 5: \"c\"}", String(m))
	}
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
	equal(t, `&repr.testStruct{
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
	equal(t, "[]byte(\"\\x01\\x02\\x03\")", String(b))
}

type privateTestStruct struct {
	a string
}

func TestReprPrivateField(t *testing.T) {
	s := privateTestStruct{"hello"}
	equal(t, `repr.privateTestStruct{a: "hello"}`, String(s))
}

type mixedTestStruct struct {
	A  string
	b  string
	C  string
	_D string
}

func TestReprPrivateMixed(t *testing.T) {
	s := mixedTestStruct{"hello", "world", "goodbye", "cruel world"}
	equal(t, `repr.mixedTestStruct{A: "hello", b: "world", C: "goodbye", _D: "cruel world"}`, String(s))
}

func TestReprPrivateMixedIgnorePrivate(t *testing.T) {
	s := mixedTestStruct{"hello", "world", "goodbye", "cruel world"}
	equal(t, `repr.mixedTestStruct{A: "hello", C: "goodbye"}`, String(s, IgnorePrivate()))
}

func TestReprNilAlone(t *testing.T) {
	var err error
	s := String(err)
	equal(t, "nil", s)
}

func TestExplicitTypes(t *testing.T) {
	arr := []*privateTestStruct{{"hello"}, nil}
	s := String(arr, ExplicitTypes(true))
	equal(t, "[]*repr.privateTestStruct{&repr.privateTestStruct{a: \"hello\"}, nil}", s)
}

func TestReprNilInsideArray(t *testing.T) {
	arr := []*privateTestStruct{{"hello"}, nil}
	s := String(arr)
	equal(t, "[]*repr.privateTestStruct{{a: \"hello\"}, nil}", s)
}

func TestReprEmptySlice(t *testing.T) {
	a := []int{}
	s := String(a)
	equal(t, "[]int{}", s)
}

func TestReprNilSlice(t *testing.T) {
	var a []int
	s := String(a)
	equal(t, "nil", s)
}

type intSliceStruct struct{ f []int }

func TestReprEmptySliceStruct(t *testing.T) {
	a := intSliceStruct{f: []int{}}
	s := String(a)
	equal(t, "repr.intSliceStruct{}", s)
}

func TestReprNilSliceStruct(t *testing.T) {
	var a intSliceStruct
	s := String(a)
	equal(t, "repr.intSliceStruct{}", s)
}

type Enum int

func (e Enum) String() string {
	return "Value"
}

func TestEnum(t *testing.T) {
	v := Enum(1)
	s := String(v)
	equal(t, "repr.Enum(Value)", s)
}

func TestShowType(t *testing.T) {
	a := map[string]privateTestStruct{"foo": {"bar"}}
	s := String(a, AlwaysIncludeType(), Indent("  "))
	equal(t, strings.TrimSpace(`
map[string]repr.privateTestStruct{
  string("foo"): repr.privateTestStruct{
    a: string("bar"),
  },
}
`), s)
}

func TestRecursiveIssue3(t *testing.T) {
	type data struct {
		parent   *data
		children []*data
	}
	child := &data{}
	root := &data{children: []*data{child}}
	child.parent = root
	want := "&repr.data{children: []*repr.data{{parent: &...}}}"
	have := String(root)
	equal(t, want, have)
}

type MyBuffer struct {
	buf *bytes.Buffer
}

func TestReprPrivateBytes(t *testing.T) {
	mb := MyBuffer{
		buf: bytes.NewBufferString("Hi th3re!"),
	}
	s := String(mb)

	switch v := runtime.Version(); {
	case strings.Contains(v, "go1.9"):
		equal(t, "repr.MyBuffer{buf: &bytes.Buffer{buf: []byte(\"Hi th3re!\"), bootstrap: [64]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}}", s)
	case strings.Contains(v, "go1.10"), strings.Contains(v, "go1.11"):
		equal(t, "repr.MyBuffer{buf: &bytes.Buffer{buf: []byte(\"Hi th3re!\"), bootstrap: [64]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, }}", s)
	default:
		equal(t, "repr.MyBuffer{buf: &bytes.Buffer{buf: []byte(\"Hi th3re!\")}}", s)
	}
}

func TestReprAnyNumeric(t *testing.T) {
	var value = []any{float64(123)}
	equal(t, "[]any{float64(123)}", String(value))
}

func TestReprFunc(t *testing.T) {
	in := func(any) {}
	equal(t, "func(any)", String(in))
	inout := func(interface{}) (any, error) { panic("not implemented") }
	equal(t, "func(any) (any, error)", String(inout))
}

func TestScalarLiterals(t *testing.T) {
	d := time.Second
	equal(t, "time.Duration(1000000000)", String(d, ScalarLiterals()))
}
