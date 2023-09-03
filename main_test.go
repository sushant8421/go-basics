package main

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func populateProducts(sliceCount, sliceSize, maxScore int) [][]Product {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	productSlices := make([][]Product, sliceCount)
	for i := 0; i < sliceCount; i++ {
		slice := make([]Product, sliceSize)
		for j := 0; j < sliceSize; j++ {
			sku := "slice" + strconv.Itoa(i+1) + "-sku" + strconv.Itoa(j+1)
			score := int64(r.Intn(maxScore) + 1)
			slice[j] = Product{SKU: sku, Score: score}
		}
		productSlices[i] = slice
	}

	return productSlices
}
func BenchmarkLazyApproach(b *testing.B) {
	// Initialize product slices here...
	productSlices := populateProducts(20, 150, 500)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mergeAndSortProducts(productSlices, 10)
	}
}
func BenchmarkTechieApproach(b *testing.B) {
	// Initialize sets here...
	productSlices := populateProducts(20, 150, 500)
	sets := convertToSortedSets(productSlices)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		its := make([]*Iterator, len(sets))
		for j, set := range sets {
			its[j] = NewIterator(set)
		}
		getTopX(its, 10)
	}
}
