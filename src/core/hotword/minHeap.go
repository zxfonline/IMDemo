package hotword

import "container/heap"

type TopKInfo struct {
	k       int
	MinHeap TimeTrieHeap
}

//前k数据 小顶堆
func NewTopKInfo(k int, nums []*TimeTrie) *TopKInfo {
	h := TimeTrieHeap(nums)
	heap.Init(&h)

	for len(h) > k {
		heap.Pop(&h)
	}

	return &TopKInfo{
		k:       k,
		MinHeap: h,
	}
}

func (t *TopKInfo) Add(val *TimeTrie) *TimeTrie {
	heap.Push(&t.MinHeap, val)
	if len(t.MinHeap) > t.k {
		heap.Pop(&t.MinHeap)
	}
	return t.MinHeap[0]
}

type TimeTrieHeap []*TimeTrie

func (h TimeTrieHeap) Len() int           { return len(h) }
func (h TimeTrieHeap) Less(i, j int) bool { return len(h[i].recordTss) < len(h[j].recordTss) }
func (h TimeTrieHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *TimeTrieHeap) Push(x interface{}) {
	*h = append(*h, x.(*TimeTrie))
}

func (h *TimeTrieHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
