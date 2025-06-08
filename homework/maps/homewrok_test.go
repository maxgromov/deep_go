package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type OrderedMap struct {
	root *node
	size int
}

type node struct {
	key         int
	value       int
	left, right *node
}

// NewOrderedMap - создать упорядоченный словарь
func NewOrderedMap() OrderedMap {
	return OrderedMap{}
}

// Insert - добавить элемент в словарь
func (m *OrderedMap) Insert(key, value int) {
	var inserted bool
	m.root, inserted = insert(m.root, key, value)
	if inserted {
		m.size++
	}
}

func insert(n *node, key, value int) (*node, bool) {
	if n == nil {
		return &node{
			key:   key,
			value: value,
		}, true
	}

	switch {
	case key < n.key:
		var inserted bool
		n.left, inserted = insert(n.left, key, value)
		return n, inserted

	case key > n.key:
		var inserted bool
		n.right, inserted = insert(n.right, key, value)
		return n, inserted

	default:
		fmt.Println("key already exist")
		return n, false
	}
}

// Erase - удалить элемент из словаря
func (m *OrderedMap) Erase(key int) {
	var deleted bool
	m.root, deleted = deleteNode(m.root, key)
	if deleted {
		m.size--
	}
}

func deleteNode(n *node, key int) (*node, bool) {
	if n == nil {
		return nil, false
	}

	var isDeleted bool

	switch {
	case key < n.key:
		n.left, isDeleted = deleteNode(n.left, key)

	case key > n.key:
		n.right, isDeleted = deleteNode(n.right, key)

	default:
		isDeleted = true
		if n.left == nil && n.right == nil {
			return nil, isDeleted
		}

		if n.left == nil {
			return n.right, isDeleted
		}
		if n.right == nil {
			return n.left, isDeleted
		}

		// преемник
		s := minNode(n.right)
		n.key = s.key
		n.value = s.value
		n.right, _ = deleteNode(n.right, s.key)
	}

	return n, isDeleted
}

// поиск преемника
func minNode(n *node) *node {
	for n.left != nil {
		n = n.left
	}

	return n
}

// Contains - проверить существование элемента в словаре
func (m *OrderedMap) Contains(key int) bool {
	return contain(m.root, key)
}

func contain(n *node, key int) bool {
	if n == nil {
		return false
	}
	switch {
	case key < n.key:
		return contain(n.left, key)
	case key > n.key:
		return contain(n.right, key)
	default:
		return true
	}
}

// Size - получить количество элементов в словаре
func (m *OrderedMap) Size() int {
	return m.size
}

// ForEach - применить функцию к каждому элементу словаря от меньшего к большему
func (m *OrderedMap) ForEach(action func(int, int)) {
	inOrder(m.root, action)
}

// обход: левое поддерево - узел - правое поддерево
func inOrder(n *node, action func(int, int)) {
	if n == nil {
		return
	}
	inOrder(n.left, action)
	action(n.key, n.value)
	inOrder(n.right, action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
