package bloomfilter

import (
	"math"

	"github.com/spaolacci/murmur3"
)

type BloomFilter struct {
	vec        []byte                                // byte array to store the bloom filter
	bits       int                                   // number of bits in bloom filter array
	hashNumber int                                   // number of hash functions
	hashFunc   func(data []byte, seed uint32) uint64 // hash function
}

func New(n int, fpr float64) *BloomFilter {
	bytes := calculateBytes(n, fpr)
	bits := bytes * 8
	hashNumber := calculateHashNumber(n, bits)

	return &BloomFilter{
		vec:        make([]byte, bytes),
		bits:       bits,
		hashNumber: hashNumber,
		hashFunc:   murmur3.Sum64WithSeed,
	}
}

func calculateBytes(n int, fpr float64) int {
	m := int(-(float64(n) * math.Log(fpr)) / (math.Ln2 * math.Ln2))
	return int(math.Ceil(float64(m) / 8))
}

func calculateHashNumber(n int, m int) int {
	return int(math.Ceil(float64(m) / float64(n) * math.Ln2))
}

/*
 * The bloom filter uses a byte array to store information.
 * To minimize memory allocation, the bloom filter must work with bits, not bytes.
 * Every add operation will calculate the index of the byte to work with.
 * Then calculates the index of the bit and set the byte selected bit to 1.
 * Example: index is 20, byte array is [0, 0, 0, 0, 0, 0, 0, 0]
 * bit array 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000
 * blockIndex: 20 / 8 = 2
 * bitIndex: 20 % 8 = 4
 * b.vec[2] = b.vec[2] | (1 << 4)
 * b.vec[2] = 00000000 | (1 << 4)
 * b.vec[2] = 00000000 | 00001000
 * b.vec[2] = 00001000
 * After add operation, byte array is [0, 0, 8, 0, 0, 0, 0, 0]
 */

// getBit returns the value of the bit at the given index.
func (b *BloomFilter) getBit(index uint64) bool {
	blockIndex := index / 8
	bitIndex := index % 8
	return b.vec[blockIndex]&(1<<bitIndex) != 0
}

// setBit sets the value of the bit at the given index to 1.
func (b *BloomFilter) setBit(index uint64) {
	blockIndex := index / 8
	bitIndex := index % 8
	b.vec[blockIndex] = b.vec[blockIndex] | (1 << bitIndex)
}

// Add adds an element to the bloom filter.
func (b *BloomFilter) Add(elem string) error {
	// hash the element that needs to be added to the bloom filter
	// for the number of hash functions.
	for i := 0; i < b.hashNumber; i++ {
		hash := b.hashFunc([]byte(elem), uint32(i))
		index := hash % uint64(b.bits)
		b.setBit(index)
	}
	return nil
}

// Contains checks if an element is in the bloom filter.
func (b *BloomFilter) Contains(elem string) bool {
	for i := 0; i < b.hashNumber; i++ {
		hash := b.hashFunc([]byte(elem), uint32(i))
		index := hash % uint64(b.bits)
		if !b.getBit(index) {
			return false
		}
	}
	return true
}
