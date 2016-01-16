# Python's repr() for Go

This package attempts to represent Go values in a form that can be copy-and-pasted into source code
directly.

Unfortunately, some values such as pointers to basic types, can not be represented directly in Go,
so they are represented as `&<value>`. eg. `&23`

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
