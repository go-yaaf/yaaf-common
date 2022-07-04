/*
 * Priority queue implementation
 *
 * Based on https://github.com/ReneKroon/ttlcache
 */
package cache

import (
	"container/heap"
)

func newPriorityQueue() *priorityQueue {
	queue := &priorityQueue{}
	heap.Init(queue)
	return queue
}

type priorityQueue struct {
	items []*cachedItem
}

func (pq *priorityQueue) update(item *cachedItem) {
	heap.Fix(pq, item.queueIndex)
}

func (pq *priorityQueue) push(item *cachedItem) {
	heap.Push(pq, item)
}

func (pq *priorityQueue) pop() *cachedItem {
	if pq.Len() == 0 {
		return nil
	}
	return heap.Pop(pq).(*cachedItem)
}

func (pq *priorityQueue) remove(item *cachedItem) {
	heap.Remove(pq, item.queueIndex)
}

func (pq priorityQueue) Len() int {
	length := len(pq.items)
	return length
}

// Less will consider items with time.Time default value (epoch start) as more than set items.
func (pq priorityQueue) Less(i, j int) bool {
	if pq.items[i].expireAt.IsZero() {
		return false
	}
	if pq.items[j].expireAt.IsZero() {
		return true
	}
	return pq.items[i].expireAt.Before(pq.items[j].expireAt)
}

func (pq priorityQueue) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].queueIndex = i
	pq.items[j].queueIndex = j
}

func (pq *priorityQueue) Push(x any) {
	item := x.(*cachedItem)
	item.queueIndex = len(pq.items)
	pq.items = append(pq.items, item)
}

func (pq *priorityQueue) Pop() any {
	old := pq.items
	n := len(old)
	item := old[n-1]
	item.queueIndex = -1
	pq.items = old[0 : n-1]
	return item
}
