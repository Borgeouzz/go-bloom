package bloom

import (
	"math"

	"github.com/spaolacci/murmur3"
)

type BloomFilter struct {
	vec        []byte
	hashNumber int
	hashFunc   func(data []byte) uint64
}

func New(elemNumber int, falsePositiveRate float64) *BloomFilter {
	byteVecSize := calculateByteVecSize(elemNumber, falsePositiveRate)
	hashFunctions := calculateHashFunctions(elemNumber, byteVecSize)
	return &BloomFilter{
		vec:        make([]byte, byteVecSize),
		hashNumber: hashFunctions,
		hashFunc:   murmur3.Sum64,
	}
}

func (b *BloomFilter) Add(elem string) {
	// hash the element that needs to be added to the bloom filter
	// for the number of hash functions.
	for i := 0; i < b.hashNumber; i++ {
		hash := murmur3.Sum64([]byte(elem))
		index := hash % uint64(len(b.vec))
		b.vec[index] = 1
	}
}

func (b *BloomFilter) Contains(elem string) bool {
	for i := 0; i < b.hashNumber; i++ {
		hash := b.hashFunc([]byte(elem))
		index := hash % uint64(len(b.vec))
		if b.vec[index] == 0 {
			return false
		}
	}
	return true
}

func calculateByteVecSize(n int, p float64) int {
	// formula is m = -(n x ln(p)) / (ln(2))^2
	// where m is the number of bits
	m := int(-(float64(n) * math.Log(p)) / (math.Ln2 * math.Ln2))
	numBytes := int(math.Ceil(float64(m) / 8))
	return numBytes
}

func calculateHashFunctions(n int, m int) int {
	return int(math.Ceil(float64(m) / float64(n) * math.Ln2))
}
