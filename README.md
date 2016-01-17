# Python's repr() for Go [![](https://godoc.org/github.com/alecthomas/repr?status.svg)](http://godoc.org/github.com/alecthomas/repr) [![Build Status](https://travis-ci.org/alecthomas/repr.png)](https://travis-ci.org/alecthomas/repr)

This package attempts to represent Go values in a form that can be copy-and-pasted into source code
directly.

Unfortunately some values (such as pointers to basic types) can not be represented directly in Go.
These values will be represented as `&<value>`. eg. `&23`

## Example

```go
type test struct {
  S string
  I int
  A []int
}

func main() {
  repr.Print(&test{
    S: "String",
    I: 123,
    A: []int{1, 2, 3},
  })
}
```

Outputs

```
&main.test{
  "S": "String",
  "I": 123,
  A: []int{
    1,
    2,
    3,
  },
}
```
