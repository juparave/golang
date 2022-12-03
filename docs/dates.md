# Time and Dates in Golang

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    var t time.Time
    fmt.Println(t)
    fmt.Println(t.IsZero())
}
```

output

    0001-01-01 00:00:00 +0000 UTC
    true

## Date format

The layout string used by the Parse function and Format method
shows by example how the reference time should be represented.
We stress that one must show how the reference time is formatted,
not a time of the user's choosing. Thus each layout string is a
representation of the time stamp,

    `Jan 2 15:04:05 2006 MST`

An easy way to remember this value is that it holds, when presented
in this order, the values (lined up with the elements above):

    `1 2  3  4  5    6  -7`
