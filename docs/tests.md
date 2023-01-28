# Go Tests

* [Best Practices](https://medium.com/@benbjohnson/structuring-tests-in-go-46ddee7a25c)
* [Testing functions for Go](https://github.com/benbjohnson/testing)

```go
var testsWithInterface = []struct {
	inte interface{}
	find interface{}
	want bool
}{
	{
		[]int{2, 3, 4, 5, 6, 1990},
		7,
		false,
	},
	{
		[]string{"a", "t", "l", "E", "R", "q"},
		"q",
		true,
	},
}

func TestContains(t *testing.T) {

	for _, tt := range tests {
		testname := fmt.Sprintf("%v, %v, %v", tt.str, tt.find, tt.want)
		t.Run(testname, func(t *testing.T) {
			arr := strings.Split(tt.str, "")
			res := Contains(arr, tt.find)
			if res != tt.want {
				t.Errorf("got %v, want %v", res, tt.want)
			}
		})
	}

	for _, tt := range testsWithInterface {
		testname := fmt.Sprintf("interface %v, %v, %v", tt.inte, tt.find, tt.want)
		t.Run(testname, func(t *testing.T) {
			res := Contains(tt.inte, tt.find)
			if res != tt.want {
				t.Errorf("got %v, want %v", res, tt.want)
			}
		})
	}
}
```

Run tests with `entr`

    $ <<< scanner_test.go entr -cc go test -v
