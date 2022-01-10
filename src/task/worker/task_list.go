package worker

import "sync"

type taskList struct {
	mu    sync.RWMutex
	tasks map[string]*Task
}

func NewTaskList() *taskList {
	return &taskList{
		tasks: make(map[string]*Task),
	}
}

// Put 新增
func (t *taskList) Put(k string, task *Task) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.tasks[k] = task
}

// Delete 删除
func (t *taskList) Delete(k string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.tasks, k)
}

// CopyGet 获取单个数据
func (t *taskList) CopyGet(k string) Task {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return *t.tasks[k]
}

// Get 获取单个数据
func (t *taskList) Get(k string) *Task {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.tasks[k]
}

// List 获取全部数据
func (t *taskList) List() map[string]*Task {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.tasks
}

// Len 长度
func (t *taskList) Len() int {
	return len(t.tasks)
}

// Exists Task是否存在
func (t *taskList) Exists(k string) bool {
	_, ok := t.tasks[k]
	return ok
}
