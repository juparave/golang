# Slices

### Unpack slices on assignment

python: > a, b, c = "foo;bar;baz".split(";"

As done in python is not supported, you can solve it with your own
ad-hoc function using multiple returns:

```go
func splitCoord(s, sep string) (string, string, string) {
    x := strings.Split(s, sep)
    return x[0], x[1], x[1]
}
```

And yot can then write:

```go
x, y, z := splitLink("21,34,6", ",")
```
