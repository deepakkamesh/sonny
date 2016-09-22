package main

import (
	"container/heap"
	"fmt"
)

type item struct {
	pkt   []byte
	pri   int
	index int
}

type cmdQueue []*item

func (cq cmdQueue) Len() int {
	return len(cq)
}

func (cq cmdQueue) Less(i, j int) bool {
	return cq[i].pri > cq[j].pri
}

func (cq cmdQueue) Swap(i, j int) {
	cq[i], cq[j] = cq[j], cq[i]
	cq[i].index = i
	cq[j].index = j
}

func (cq *cmdQueue) Push(x interface{}) {
	n := len(*cq)
	i := x.(*item)
	i.index = n
	*cq = append(*cq, i)
}

func (cq *cmdQueue) Pop() interface{} {
	old := *cq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*cq = old[0 : n-1]
	return item
}

func main() {

	cmd := cmdQueue{}
	heap.Init(&cmd)
	heap.Push(&cmd, &item{
		pkt: []byte{'a'},
		pri: 10,
	})
	heap.Push(&cmd, &item{
		pkt: []byte{'b'},
		pri: 12,
	})

	for cmd.Len() > 0 {
		item := heap.Pop(&cmd).(*item)
		fmt.Printf("%v   --", item.pkt)
	}

}
