package shortuuid

import (
	"math/rand"
	"time"
)

const numerals = "0123456789ABCDEFGHIJKLMNPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"

func ShortUUID(length int) string {
	// seed a new source of randon numbers
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	// get int value of uuid
	uuid := r.Int63()
	numerals := alphabet

	// encode it until we get the length we want
	shortuuid := encode(uuid, len(numerals), numerals)
	for len(shortuuid) < length {
		uuid = r.Int63()
		shortuuid = encode(uuid, len(numerals), numerals)
		// if length is too long, append a second encoded uuid
		if length > 11 { // 11 is the max length using alphabet, probably =)
			shortuuid += encode(r.Int63(), len(numerals), numerals)
		}
	}
	return shortuuid[0:length]
}

func SortableShortUUID() string {
	// use time.Now() to get a unix timestamp
	// the length of the timestamp is 11 using all numerals
	// this UUID is sortable and can be used as a key in a map
	return encode(time.Now().UnixNano(), len(alphabet), alphabet)
}

func encode(num int64, base int, numerals string) string {
	var result string
	if numerals == "" {
		numerals = alphabet
	}
	for num > 0 {
		result = string(numerals[num%int64(base)]) + result
		num /= int64(base)
	}
	return result
}

// func decode(s string, base int, numerals string) int {
// 	var result int
// 	var fbase float64 = float64(base)
// 	if numerals == "" {
// 		numerals = alphabet
// 	}
// 	// reverse string
// 	runes := []rune(s)
// 	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
// 		runes[i], runes[j] = runes[j], runes[i]
// 	}
// 	s = string(runes)
// 	for i, c := range s {
// 		var fi float64 = float64(i)
// 		var pow float64 = math.Pow(fbase, fi)
// 		result += strings.Index(numerals, string(c)) * int(pow)
// 		fmt.Println(string(c), strings.Index(numerals, string(c)), base, pow)
// 	}
// 	return result
// }
