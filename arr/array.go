package arr

import "sync"

type ArrayType interface {
	Add(value any)
	Check(value any) bool
	Del(value any)
	Len() int
	Values() []any
}

// CheckType type for check array data
type CheckArrayType struct {
	sync.Mutex
	values map[any]bool
	len    int
}

// NewCheckArrayType
func NewCheckArrayType(len int) *CheckArrayType {
	return &CheckArrayType{values: make(map[any]bool, len)}
}

// Add
func (ct *CheckArrayType) Add(value any) {
	defer ct.Unlock()
	ct.Lock()
	ct.values[value] = true
	ct.len++
}

// AddMutil
func (ct *CheckArrayType) AddMutil(values ...any) {
	for _, v := range values {
		v := v
		ct.Add(v)
	}
}

// Check
func (ct *CheckArrayType) Check(value any) bool {
	defer ct.Unlock()
	ct.Lock()
	if b, ok := ct.values[value]; ok && b {
		return true
	}
	return false
}

// Len
func (ct *CheckArrayType) Len() int {
	defer ct.Unlock()
	ct.Lock()
	return ct.len
}

// Del
func (ct *CheckArrayType) Del(value any) {
	defer ct.Unlock()
	ct.Lock()
	delete(ct.values, value)
	ct.len--
}

// Values
func (ct *CheckArrayType) Values() []any {
	defer ct.Unlock()
	ct.Lock()
	values := make([]any, 0, len(ct.values))
	for k := range ct.values {
		values = append(values, k)
	}
	return values
}
