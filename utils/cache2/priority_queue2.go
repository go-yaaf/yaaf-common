/*
 * Priority queue implementation
 *
 * Based on https://github.com/ReneKroon/ttlcache
 */

package cache2

import (
	"container/heap"
)

func newPriorityQueue() *priorityQueue2 {
	queue := &priorityQueue2{}
	heap.Init(queue)
	return queue
}

type priorityQueue2 struct {
	items []*cachedItem2
}

func (pq *priorityQueue2) update(item *cachedItem2) {
	heap.Fix(pq, item.queueIndex)
}

func (pq *priorityQueue2) push(item *cachedItem2) {
	heap.Push(pq, item)
}

func (pq *priorityQueue2) pop() *cachedItem2 {
	if pq.Len() == 0 {
		return nil
	}
	return heap.Pop(pq).(*cachedItem2)
}

func (pq *priorityQueue2) remove(item *cachedItem2) {
	heap.Remove(pq, item.queueIndex)
}

func (pq priorityQueue2) Len() int {
	length := len(pq.items)
	return length
}

// Less will consider items with time.Time default value (epoch start) as more than set items.
func (pq priorityQueue2) Less(i, j int) bool {
	if pq.items[i].expireAt.IsZero() {
		return false
	}
	if pq.items[j].expireAt.IsZero() {
		return true
	}
	return pq.items[i].expireAt.Before(pq.items[j].expireAt)
}

func (pq priorityQueue2) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].queueIndex = i
	pq.items[j].queueIndex = j
}

func (pq *priorityQueue2) Push(x any) {
	item := x.(*cachedItem2)
	item.queueIndex = len(pq.items)
	pq.items = append(pq.items, item)
}

func (pq *priorityQueue2) Pop() any {
	old := pq.items
	n := len(old)
	item := old[n-1]
	item.queueIndex = -1
	pq.items = old[0 : n-1]
	return item
}
