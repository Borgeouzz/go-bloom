package bloomfilter

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
