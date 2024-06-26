/*
 * Priority queue implementation
 *
 * Based on https://github.com/ReneKroon/ttlcache
 */

package cache

import (
	"container/heap"
)

func newPriorityQueue[K comparable, T any]() *priorityQueue[K, T] {
	queue := &priorityQueue[K, T]{}
	heap.Init(queue)
	return queue
}

type priorityQueue[K comparable, T any] struct {
	items []*cachedItem[K, T]
}

func (pq *priorityQueue[K, T]) update(item *cachedItem[K, T]) {
	heap.Fix(pq, item.queueIndex)
}

func (pq *priorityQueue[K, T]) push(item *cachedItem[K, T]) {
	heap.Push(pq, item)
}

func (pq *priorityQueue[K, T]) pop() *cachedItem[K, T] {
	if pq.Len() == 0 {
		return nil
	}
	return heap.Pop(pq).(*cachedItem[K, T])
}

func (pq *priorityQueue[K, T]) remove(item *cachedItem[K, T]) {
	heap.Remove(pq, item.queueIndex)
}

func (pq priorityQueue[K, T]) Len() int {
	length := len(pq.items)
	return length
}

// Less will consider items with time.Time default value (epoch start) as more than set items.
func (pq priorityQueue[K, T]) Less(i, j int) bool {
	if pq.items[i].expireAt.IsZero() {
		return false
	}
	if pq.items[j].expireAt.IsZero() {
		return true
	}
	return pq.items[i].expireAt.Before(pq.items[j].expireAt)
}

func (pq priorityQueue[K, T]) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].queueIndex = i
	pq.items[j].queueIndex = j
}

func (pq *priorityQueue[K, T]) Push(x any) {
	item := x.(*cachedItem[K, T])
	item.queueIndex = len(pq.items)
	pq.items = append(pq.items, item)
}

func (pq *priorityQueue[K, T]) Pop() any {
	old := pq.items
	n := len(old)
	item := old[n-1]
	item.queueIndex = -1
	pq.items = old[0 : n-1]
	return item
}
