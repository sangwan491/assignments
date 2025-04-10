package main

import "math"

// --- hashGenerator Constants ---

// primeBase is the multiplier for the polynomial rolling hash. Using a prime number
// helps minimize collisions with patterned input. 257 is selected as it's a prime larger than the
// typical character set size (e.g., max ASCII 255), ensuring character positions
// contribute distinctly before the modulo operation occurs.
const primeBase uint64 = 257

// Implicit Modulus Explanation:
// This hash function utilizes the native overflow behavior of uint64 arithmetic.
// When the hash value exceeds math.MaxUint64, it wraps around. This is mathematically
// equivalent to performing the calculation modulo 2^64.
//
// Coprimality:
// An important property for polynomial rolling hashes is that the base (primeBase)
// and the modulus or the bucket size should ideally be coprime (having a greatest common divisor of 1).
// In this case, the implicit modulus is 2^64. Since primeBase (257) is a prime
// number other than 2, it is coprime to 2^64. This helps ensure better hash
// distribution and reduces the likelihood of certain types of collisions.
// If primeBase shared a common factor greater than 1 with the modulus, the effective bucket size
// would be reduced by that factor.
//
// Efficiency:
// Relying on the implicit modulo 2^64 via uint64 overflow is generally more
// computationally efficient than performing an explicit modulo operation (%) at each step.
const bucketSize uint64 = math.MaxUint64

// Constants for MurmurHash-like bit mixing (avalanche effect).
// These specific constants are chosen for their properties in spreading
// input bit changes across the output hash bits, improving distribution
// and reducing collisions.
const mixConstant1 uint64 = 0xff51afd7ed558ccd
const mixConstant2 uint64 = 0xc4ceb9fe1a85ec53
const mixShift uint64 = 33

// --- toBase62 Constants ---

// base62Chars represents the character set for Base62 encoding.
// It includes digits 0-9, uppercase letters A-Z, and lowercase letters a-z.
const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// base62 is the base used for the encoding (number of unique characters).
const base62 uint64 = 62

// targetHashLength specifies the desired fixed length of the final Base62 encoded string.
// The output will be padded with leading zeros or truncated to meet this length.
const targetHashLength = 10

// hashGenerator accepts an alphanumeric string of an arbitrary length and outputs a 10 character alphanumeric hash.
// It uses a polynomial rolling hash combined with bit mixing for better distribution.
func hashGenerator(str string) string {
	// Note: bucketSize was defined but not used. Removed for clarity.
	// const bucketSize uint64 = math.MaxUint64

	var hash uint64 = 0

	for i := 0; i < len(str); i++ {
		ascii := uint64(str[i])

		// Polynomial rolling hash step (implicit modulo 2^64)
		hash = hash*primeBase + ascii
	}

	// Bit mixing (avalanche effect) - inspired by MurmurHash
	hash ^= (hash >> mixShift)
	hash *= mixConstant1
	hash ^= (hash >> mixShift)
	hash *= mixConstant2
	hash ^= (hash >> mixShift)

	hashString := toBase62(hash)

	return hashString
}

// toBase62 converts a uint64 hash value to a fixed-length Base62 encoded string.
// The resulting string will have a length specified by targetHashLength.
func toBase62(hash uint64) string {
	if hash == 0 {
		// Pad with leading zeros to meet the target length
		result := make([]byte, targetHashLength)
		for i := range result {
			result[i] = '0'
		}
		// Special case: return "0" padded to the target length, e.g., "0000000000"
		// If you strictly want just "0" for input 0, use: return string(base62Chars[0])
		return string(result)
	}

	var result []byte

	for hash > 0 {
		remainder := hash % base62
		result = append([]byte{base62Chars[remainder]}, result...) // Prepend to reverse order
		hash /= base62
	}

	// Pad with leading zeros if the result is shorter than the target length
	paddingNeeded := targetHashLength - len(result)
	if paddingNeeded > 0 {
		padding := make([]byte, paddingNeeded)
		for i := range padding {
			padding[i] = base62Chars[0] // '0'
		}
		result = append(padding, result...)
	}

	// Truncate if the result is longer than the target length (shouldn't happen with uint64 and base62 often, but safety check)
	if len(result) > targetHashLength {
		result = result[:targetHashLength]
	}

	return string(result)
}
