package arr

import "sync"

type ArrayType interface {
	Add(value interface{})
	Check(value interface{}) bool
	Del(value interface{})
	Len() int
	Values() map[interface{}]bool
}

// CheckType type for check array data
type CheckArrayType struct {
	values map[interface{}]bool
	sm     sync.Mutex
	len    int
}

// NewCheckArrayType
func NewCheckArrayType(len int) *CheckArrayType {
	return &CheckArrayType{values: make(map[interface{}]bool, len)}
}

// Add
func (ct *CheckArrayType) Add(value interface{}) {
	defer ct.sm.Unlock()
	ct.sm.Lock()
	ct.values[value] = true
	ct.len++
}

// AddMutil
func (ct *CheckArrayType) AddMutil(values ...interface{}) {
	for _, v := range values {
		v := v
		ct.Add(v)
	}
}

// Check
func (ct *CheckArrayType) Check(value interface{}) bool {
	defer ct.sm.Unlock()
	ct.sm.Lock()
	if b, ok := ct.values[value]; ok && b {
		return true
	}
	return false
}

// Len
func (ct *CheckArrayType) Len() int {
	defer ct.sm.Unlock()
	ct.sm.Lock()
	return ct.len
}

// Del
func (ct *CheckArrayType) Del(value interface{}) {
	defer ct.sm.Unlock()
	ct.sm.Lock()
	delete(ct.values, value)
	ct.len--
}

// Values
func (ct *CheckArrayType) Values() map[interface{}]bool {
	defer ct.sm.Unlock()
	ct.sm.Lock()
	return ct.values
}
