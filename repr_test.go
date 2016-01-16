package repr

import (
	"testing"

	"github.com/alecthomas/assert"
)

type anotherStruct struct {
	A []int
}

type testStruct struct {
	S string
	I *int
	A anotherStruct
}

func TestRepr(t *testing.T) {
	assert.Equal(t, "[]string{\n  \"a\",\n  \"b\",\n}", Repr([]string{"a", "b"}))
	assert.Equal(t, "[]int{\n  1,\n  2,\n}", Repr([]int{1, 2}))
	pi := new(int)
	*pi = 13
	assert.Equal(t, `&13`, Repr(pi))

	ch := make(<-chan map[string]*testStruct, 1)
	assert.Equal(t, `make(<-chan map[string]*repr.testStruct, 1)`, Repr(ch))

	m := map[string]int{"a": 1}
	assert.Equal(t, "map[string]int{\n  \"a\": 1,\n}", Repr(m))

	s := &testStruct{
		S: "String",
		I: pi,
		A: anotherStruct{
			A: []int{1, 2, 3},
		},
	}
	assert.Equal(t, `&repr.testStruct{
  "S": "String",
  "I": &13,
  "A": repr.anotherStruct{
    "A": []int{
      1,
      2,
      3,
    },
  },
}`, Repr(s))

	b := []uint8{1, 2, 3}
	assert.Equal(t, `[]uint8{
  1,
  2,
  3,
}`, Repr(b))
}
