/*
 * Priority queue implementation for managing cache item expiration.
 *
 * This priority queue is built on Go's `container/heap` package and is used
 * by the cache to efficiently manage item expiration. Items are ordered by their
 * expiration time, allowing the cache to quickly find the next item to expire.
 *
 * Based on https://github.com/ReneKroon/ttlcache
 */

package cache

import (
	"container/heap"
)

// newPriorityQueue creates and initializes a new priorityQueue.
func newPriorityQueue[K comparable, T any]() *priorityQueue[K, T] {
	queue := &priorityQueue[K, T]{
		items: make([]*cachedItem[K, T], 0),
	}
	heap.Init(queue)
	return queue
}

// priorityQueue implements heap.Interface and holds cachedItems.
// The priority is determined by the item's expiration time.
type priorityQueue[K comparable, T any] struct {
	items []*cachedItem[K, T]
}

// update modifies the priority and position of an item in the queue.
// This is called when an item's expiration time is updated.
func (pq *priorityQueue[K, T]) update(item *cachedItem[K, T]) {
	heap.Fix(pq, item.queueIndex)
}

// push adds an item to the priority queue.
func (pq *priorityQueue[K, T]) push(item *cachedItem[K, T]) {
	heap.Push(pq, item)
}

// pop removes and returns the item with the highest priority (the one that expires soonest).
func (pq *priorityQueue[K, T]) pop() *cachedItem[K, T] {
	if pq.Len() == 0 {
		return nil
	}
	return heap.Pop(pq).(*cachedItem[K, T])
}

// remove removes an item from the priority queue.
func (pq *priorityQueue[K, T]) remove(item *cachedItem[K, T]) {
	if item.queueIndex < 0 || item.queueIndex >= len(pq.items) {
		return
	}
	heap.Remove(pq, item.queueIndex)
}

// Len returns the number of items in the priority queue.
func (pq priorityQueue[K, T]) Len() int {
	return len(pq.items)
}

// Less compares two items based on their expiration time.
// Items with a zero expiration time are considered to have a lower priority.
func (pq priorityQueue[K, T]) Less(i, j int) bool {
	if pq.items[i].expireAt.IsZero() {
		return false
	}
	if pq.items[j].expireAt.IsZero() {
		return true
	}
	return pq.items[i].expireAt.Before(pq.items[j].expireAt)
}

// Swap swaps the positions of two items in the priority queue.
func (pq priorityQueue[K, T]) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].queueIndex = i
	pq.items[j].queueIndex = j
}

// Push implements the heap.Interface Push method.
func (pq *priorityQueue[K, T]) Push(x any) {
	item := x.(*cachedItem[K, T])
	item.queueIndex = len(pq.items)
	pq.items = append(pq.items, item)
}

// Pop implements the heap.Interface Pop method.
func (pq *priorityQueue[K, T]) Pop() any {
	old := pq.items
	n := len(old)
	item := old[n-1]
	item.queueIndex = -1 // for safety
	pq.items = old[0 : n-1]
	return item
}
