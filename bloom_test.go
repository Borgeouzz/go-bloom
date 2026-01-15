package bloom

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBloomFilter(t *testing.T) {
	bf := New(1000, 0.01)
	bf.Add("hello")
	bf.Add("world")
	if !bf.Contains("hello") {
		t.Errorf("hello should be in the bloom filter")
	}
	if !bf.Contains("world") {
		t.Errorf("world should be in the bloom filter")
	}
	if bf.Contains("foo") {
		t.Errorf("foo should not be in the bloom filter")
	}
}

func TestBloomFilterFalsePositiveRate(t *testing.T) {
	bf := New(100, 0.01)
	for i := range 100 {
		bf.Add(fmt.Sprintf("hello%d", i))
	}

	var falsePositiveCount int
	for i := range 100000 {
		if bf.Contains(fmt.Sprintf("rand%d", i)) {
			falsePositiveCount++
		}
	}

	falsePositiveRate := float64(falsePositiveCount) / 100000
	t.Logf("false positive rate: %f", falsePositiveRate)
	require.Less(t, falsePositiveRate, 0.01)
}

func TestBloomFilterSaveAndLoadFromFile(t *testing.T) {
	bf := New(1000, 0.01)
	bf.Add("hello")
	bf.Add("world")
	bf.SaveToFile("test.bf")

	bf2, err := LoadFromFile("test.bf")
	if err != nil {
		t.Errorf("failed to load bloom filter from file: %v", err)
	}

	// check if bloom metadata are the same
	require.Equal(t, bf.hashNumber, bf2.hashNumber)
	require.Equal(t, bf.bits, bf2.bits)
	require.Equal(t, bf.vec, bf2.vec)

	// check if bloom filter is the same
	require.Equal(t, bf.MemoryUsage(), bf2.MemoryUsage())
	require.Equal(t, bf.InsertionsCount(), bf2.InsertionsCount())

	// check if bloom filter contains the same elements
	require.True(t, bf2.Contains("hello"))
	require.True(t, bf2.Contains("world"))
	require.False(t, bf2.Contains("foo"))
}
