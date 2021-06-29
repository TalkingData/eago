package worker

import "sync"

type taskList struct {
	mu   sync.RWMutex
	data map[string]*Task
}

// New 新增
func (t *taskList) New(k string, task *Task) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.data[k] = task
}

// Del 删除
func (t *taskList) Del(k string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.data, k)
}

// Get 获取单个数据
func (t *taskList) CopyGet(k string) Task {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return *t.data[k]
}

// Get 获取单个数据
func (t *taskList) Get(k string) *Task {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.data[k]
}

// List 获取全部数据
func (t *taskList) List() map[string]*Task {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.data
}

// Len 长度
func (t *taskList) Len() int {
	return len(t.data)
}

// Exists Task是否存在
func (t *taskList) Exists(k string) bool {
	_, ok := t.data[k]
	return ok
}
