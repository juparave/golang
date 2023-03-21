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

## Handle zero dates

Is better to persist a null value than the zero date

	if lead.AssignmentDate.IsZero() {
		lead.AssignmentDate = nil
	}

## Date parser

```go
// ParseDate converts a string 'yyyy-mm-dd' to a time.Time
func ParseDate(s string) time.Time {
	format := "2006-01-02"
	parsed, err := time.Parse(format, s)
	if err != nil {
		// defaults to current date
		return time.Now()
	}
	return parsed
}

// ParseDate converts a string with defined format to a time.Time
func ParseDateF(s, format string) time.Time {
	parsed, err := time.Parse(format, s)
	if err != nil {
		// defaults to current date
		return time.Now()
	}
	return parsed
}
```
