package main

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
	index      int
}

type TaskHeap []*Task

func (h TaskHeap) Len() int           { return len(h) }
func (h TaskHeap) Less(i, j int) bool { return h[i].Priority > h[j].Priority } // max-heap
func (h TaskHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *TaskHeap) Push(x interface{}) {
	n := len(*h)
	task := x.(*Task)
	task.index = n
	*h = append(*h, task)
}

func (h *TaskHeap) Pop() interface{} {
	old := *h
	n := len(old)
	task := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually (comment from libs)
	task.index = -1 // for safety (comment from libs)
	*h = old[0 : n-1]
	return task
}

type Scheduler struct {
	taskHeap TaskHeap
	taskMap  map[int]*Task
}

func NewScheduler() Scheduler {
	return Scheduler{
		taskMap: make(map[int]*Task),
	}
}

// AddTask- запланировать задачу
func (s *Scheduler) AddTask(task Task) {
	t := Task{
		Identifier: task.Identifier,
		Priority:   task.Priority,
	}

	s.taskMap[task.Identifier] = &t
	heap.Push(&s.taskHeap, &t)
}

// ChangeTaskPriority - изменить приоритет задачи по идентификатору
func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	if task, ok := s.taskMap[taskID]; ok {
		task.Priority = newPriority
		heap.Fix(&s.taskHeap, task.index)
	}
}

// GetTask - получить задачу с наибольшим приоритетом для исполнения
func (s *Scheduler) GetTask() Task {
	if len(s.taskHeap) == 0 {
		return Task{}
	}
	task := heap.Pop(&s.taskHeap).(*Task)
	delete(s.taskMap, task.Identifier)
	return Task{
		Identifier: task.Identifier,
		Priority:   task.Priority,
	}
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	task1.Priority = 100
	assert.Equal(t, task1, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
