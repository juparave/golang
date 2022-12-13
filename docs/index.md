# Go

[go-intersect lib](https://github.com/juliangruber/go-intersect/blob/master/intersect.go)

## Vim-Go

- [Vim-go](vimgo.md)
- [Vim tips](vimtips.md)
- [Vim Slices](slices.md)

### Web servers

- [fiber](fiber.md)
- [http](http.md)

### Tests

[Tests](tests.md)

### Templates

[Practical Go Lessons. Chap32 Templates](https://www.practical-go-lessons.com/chap-32-templates)

### Databases

[gorm](gorm.md)

### Strings

[Go strings handling overview](https://yourbasic.org/golang/string-functions-reference-cheat-sheet/)

### Dates

[Dates](dates.md)

### Constants

Use CamelCase to name constants, use capital letter to export outside package

```go
const (
	AuthURL      = "https://api.someservice.com/api/v1/auth"
	DialURL      = "https://api.someservice.com/api/v1/calls/dial"
    /// These constants should be stored in environment variables
	clientID     = "CLIENTSTRING"
	clientSecret = "CLIENTSECRET"
)
```
