package shortuuid

import (
	"math/rand"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"

func ShortUUID(length int) string {
	// get int value of uuid
	uuid := rand.Int63()
	numerals := alphabet

	// encode it to base62 until we get the length we want
	shortuuid := encode(uuid, len(numerals), numerals)
	for len(shortuuid) < length {
		uuid = rand.Int63()
		shortuuid = encode(uuid, len(numerals), numerals)
		// if length is too long, append a second encoded uuid
		if length > 11 { // 12 is the max length of base62
			shortuuid += encode(rand.Int63(), len(numerals), numerals)
		}
	}
	return shortuuid[0:length]
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
