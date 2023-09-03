package main

import (
	"fmt"
	"sort"

	"github.com/wangjia184/sortedset"
)

type Product struct {
	SKU   string
	Score int64
}

type Iterator struct {
	set   *sortedset.SortedSet
	index int
}

func NewIterator(set *sortedset.SortedSet) *Iterator {
	return &Iterator{set, set.GetCount()} // starting from the end
}

func (it *Iterator) Next() *sortedset.SortedSetNode {
	if it.index >= 1 {
		node := it.set.GetByRank(it.index, false)
		it.index--
		return node
	}
	return nil
}

func getTopX(its []*Iterator, x int) []*sortedset.SortedSetNode {
	result := make([]*sortedset.SortedSetNode, 0, x)
	nodes := make([]*sortedset.SortedSetNode, len(its))

	for i := range its {
		nodes[i] = its[i].Next()
	}

	for len(result) < x {
		maxNode := (*sortedset.SortedSetNode)(nil)
		maxItIndex := -1
		for i, node := range nodes {
			if node != nil && (maxNode == nil || node.Score() > maxNode.Score()) {
				maxNode = node
				maxItIndex = i
			}
		}
		if maxNode == nil {
			break
		}
		result = append(result, maxNode)
		nodes[maxItIndex] = its[maxItIndex].Next()
	}

	return result
}

func mergeAndSortProducts(productSlices [][]Product, x int) []Product {
	mergedProducts := make([]Product, 0)
	for _, products := range productSlices {
		mergedProducts = append(mergedProducts, products...)
	}

	sort.Slice(mergedProducts, func(i, j int) bool {
		return mergedProducts[i].Score > mergedProducts[j].Score
	})

	if x > len(mergedProducts) {
		x = len(mergedProducts)
	}
	return mergedProducts[:x]
}

func convertToSortedSets(productSlices [][]Product) []*sortedset.SortedSet {
	sets := make([]*sortedset.SortedSet, len(productSlices))

	for i, products := range productSlices {
		set := sortedset.New()
		for _, product := range products {
			set.AddOrUpdate(product.SKU, sortedset.SCORE(product.Score), nil)
		}
		sets[i] = set
	}

	return sets
}

func main() {
	productSlices := [][]Product{
		{
			{"sku1", 5},
			{"sku2", 3},
			{"sku3", 1},
		},
		{
			{"sku4", 8},
			{"sku5", 6},
			{"sku6", 7},
		},
		{
			{"sku7", 9},
			{"sku8", 11},
			{"sku9", 10},
		},
	}

	fmt.Printf("For lazy approach\n")
	topProducts := mergeAndSortProducts(productSlices, 3)
	for _, product := range topProducts {
		fmt.Printf("SKU: %s, Score: %d\n", product.SKU, product.Score)
	}

	sets := convertToSortedSets(productSlices)

	its := make([]*Iterator, len(sets))
	for i, set := range sets {
		its[i] = NewIterator(set)
	}

	fmt.Println("--------------------------------------------------")
	fmt.Printf("For techie approach\n")
	topX := getTopX(its, 3)
	for _, node := range topX {
		fmt.Printf("Value: %s, Score: %d\n", node.Key(), node.Score())
	}
}
