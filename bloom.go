package bloom

/* This file implements a classic bloom filter.
 * It provides a New function to create a new bloom filter and the following methods:
 * - Add: adds an element to the bloom filter
 * - Contains: checks if an element is in the bloom filter
 * - InsertionsCount: returns the number of elements inserted in the bloom filter
 * - MemoryUsage: returns the memory usage of the bloom filter
 * - SaveToFile: saves the bloom filter to a file
 * - LoadFromFile: loads the bloom filter from a file
 */

import (
	"encoding/binary"
	"math"
	"os"

	"github.com/spaolacci/murmur3"
)

var hashFunc = murmur3.Sum64WithSeed

// BloomFilter is the main struct that represents the classic bloom filter.
type BloomFilter struct {
	vec        []byte                                // byte array to store the bloom filter
	bits       int                                   // number of bits in bloom filter array
	hashNumber int                                   // number of hash functions
	hashFunc   func(data []byte, seed uint32) uint64 // hash function
	insertions int                                   // number of elements inserted in the bloom filter
}

// New creates a new bloom filter.
func New(n int, fpr float64) *BloomFilter {
	bytes := calculateBytes(n, fpr)
	bits := bytes * 8
	hashNumber := calculateHashNumber(n, bits)

	return &BloomFilter{
		vec:        make([]byte, bytes),
		bits:       bits,
		hashNumber: hashNumber,
		hashFunc:   hashFunc,
		insertions: 0,
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
	b.insertions++
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

// InsertionsCount returns the number of elements inserted in the bloom filter.
func (b *BloomFilter) InsertionsCount() int {
	return b.insertions
}

// MemoryUsage returns the memory usage of the bloom filter.
// It returns the number of bytes used by the bloom filter.
func (b *BloomFilter) MemoryUsage() int {
	return b.bits / 8
}

// SaveToFile saves the bloom filter to a file.
// It save the bloom filter metadata (bits, hash number, insertions) and byte array to the file.
func (b *BloomFilter) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// write bytes to save bloom filter metadata and byte array
	// use int32 to ensure consistent size on all platforms
	binary.Write(file, binary.LittleEndian, int32(b.bits))
	binary.Write(file, binary.LittleEndian, int32(b.hashNumber))
	binary.Write(file, binary.LittleEndian, int32(b.insertions))

	// write byte array to file
	_, err = file.Write(b.vec)
	return err
}

// LoadFromFile loads the bloom filter from a file.
// It loads the bloom filter metadata (bits, hash number, insertions) and byte array from the file.
func LoadFromFile(filename string) (*BloomFilter, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// read bytes from file using int32 to match SaveToFile
	b := &BloomFilter{}
	var bits, hashNumber, insertions int32
	binary.Read(file, binary.LittleEndian, &bits)
	binary.Read(file, binary.LittleEndian, &hashNumber)
	binary.Read(file, binary.LittleEndian, &insertions)

	b.bits = int(bits)
	b.hashNumber = int(hashNumber)
	b.insertions = int(insertions)
	b.hashFunc = hashFunc

	// read byte array from file
	b.vec = make([]byte, (b.bits+7)/8)
	_, err = file.Read(b.vec)
	return b, err
}
