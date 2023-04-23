package hot

import (
	"container/heap"
	"errors"
	"github.com/patrickmn/go-cache"
)

type ColumnOrder struct {
	columns map[string]int
	heap    StringHeap
}

var c *cache.Cache

func (co *ColumnOrder) Init() {
	co.columns = make(map[string]int, 10)
	co.heap = make(StringHeap, 0)
	heap.Init(&co.heap)
	//c = cache.New(5*time.Minute, 10*time.Minute)
	c = cache.New(1, 1)
}

func (co *ColumnOrder) RecordClick(column string) {
	co.columns[column]++
}

func (co *ColumnOrder) RecordClickChecked(column string) error {
	if _, ok := co.columns[column]; ok {
		co.columns[column]++
	} else {
		return errors.New("not found")
	}
	return nil
}

func (co *ColumnOrder) GetOrder() []string {
	// Check the cache first
	order, found := c.Get("column_order")
	if found {
		return order.([]string)
	}

	// Update the heap based on the click counts
	for column, count := range co.columns {
		heap.Push(&co.heap, &ColumnCount{column, count})
	}

	// Retrieve the order from the heap
	torder := make([]string, 0, len(co.heap))
	for len(co.heap) > 0 {
		column := heap.Pop(&co.heap).(*ColumnCount).Column
		torder = append(torder, column)
	}

	// Store the order in the cache
	c.Set("column_order", torder, cache.DefaultExpiration)

	return torder
}

type ColumnCount struct {
	Column string
	Count  int
}

type StringHeap []*ColumnCount

func (h StringHeap) Len() int           { return len(h) }
func (h StringHeap) Less(i, j int) bool { return h[i].Count < h[j].Count }
func (h StringHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *StringHeap) Push(x interface{}) {
	*h = append(*h, x.(*ColumnCount))
}

func (h *StringHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}
