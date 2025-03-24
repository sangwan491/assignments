package main

import (
	"fmt"
	"math"
)

func hashGenerator(str string) uint64 {
	var base, hash uint64

	base = 257
	bucketSize = uint64(math.Pow(2, 64))

	hash = 0
	for i := 0; i < len(str); i++ {
		ascii := uint64(str[i])
		
		hash = (hash * base + ascii) % bucketSize

		// Bit mixing (avalanche effect)
		hash ^= (hash >> 33)
		hash *= 0xff51afd7ed558ccd
		hash ^= (hash >> 33)
		hash *= 0xc4ceb9fe1a85ec53
		hash ^= (hash >> 33)
	}

	return hash
}

// Base62 characters
const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// Converts uint64 hash to a Base62 encoded string
func toBase62(hash uint64) string {
	if hash == 0 {
		return "0"
	}
	
	var result []byte
	base := uint64(62) // Base62 encoding

	for hash > 0 {
		remainder := hash % base
		result = append([]byte{base62Chars[remainder]}, result...)
		hash /= base
	}

    for i := 0; i < (10 - len(result)); i++ {
        result = append([]byte{'0'}, result...)
    }
    
    if len(result) > 10 {
        result = result[:10]
    }
    
	return string(result)
}

func main() {
	hash := hashGenerator("harshit")
	alphanumericHash := toBase62(hash)
	fmt.Println(alphanumericHash)
}
