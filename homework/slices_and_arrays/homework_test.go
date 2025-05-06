package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type CircularQueue struct {
	values      []int
	front, rear int // изначально будут указывать на 0 индекс
	size        int
	count       int // счетчик количества элементов в очереди
}

func NewCircularQueue(size int) CircularQueue {
	return CircularQueue{
		values: make([]int, size),
		size:   size,
	}
}

func (q *CircularQueue) Push(value int) bool {
	if !q.Full() {
		q.values[q.rear] = value
		q.rear = (q.rear + 1) % q.size
		q.count++
		return true
	}

	return false
}

func (q *CircularQueue) Pop() bool {

	if !q.Empty() {
		q.values[q.front] = -1 // в целом, можем не заменять, так как указатель уже сместим и счетчик количества уменьшим
		q.front = (q.front + 1) % q.size
		q.count--
		return true
	}
	return false
}

func (q *CircularQueue) Front() int {
	if !q.Empty() {
		return q.values[q.front]
	}

	return -1
}

func (q *CircularQueue) Back() int {
	if !q.Empty() {
		lastIdx := (q.rear - 1 + q.size) % q.size
		return q.values[lastIdx]
	}

	return -1
}

func (q *CircularQueue) Empty() bool {
	return q.count == 0
}

func (q *CircularQueue) Full() bool {
	return q.count == q.size
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue(queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}
