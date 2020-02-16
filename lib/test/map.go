package test

import (
	"fmt"
)

type CounterMap struct {
	inner map[string]int32
}

func NewCounterMap() *CounterMap {
	return &CounterMap{
		inner: make(map[string]int32),
	}
}

func FromSlice(elements []string) *CounterMap {
	cm := NewCounterMap()
	for _, element := range elements {
		cm.Inc(element)
	}
	return cm
}

func (cm *CounterMap) ToMap() map[string]int32 {
	return cm.inner
}

func (cm *CounterMap) ToSlice() []string {
	slice := make([]string, 0)
	for key := range cm.inner {
		slice = append(slice, key)
	}
	return slice
}

func (cm *CounterMap) get(key string) int32 {
	return cm.inner[key]
}

func (cm *CounterMap) set(key string, value int32) {
	cm.inner[key] = value
}

func (cm *CounterMap) Inc(key string) {
	cm.inner[key]++
}

func (cm *CounterMap) Diff(other *CounterMap) *CounterMap {
	result := NewCounterMap()

	for key := range cm.inner {
		val := cm.get(key) - other.get(key)
		if val > 0 {
			result.set(key, val)
		}
	}

	return result
}

func (cm *CounterMap) Contains(other *CounterMap) (bool, string) {
	for key := range other.inner {
		if other.get(key) > cm.get(key) {
			return false, fmt.Sprintf(
				"key: %v, other's value: %v > cm's value: %v\n",
				key,
				other.get(key),
				cm.get(key),
			)
		}
	}
	return true, ""
}

func (cm *CounterMap) Equal(other *CounterMap) (bool, string) {
	if ok, desc := cm.Contains(other); !ok {
		return false, desc
	}

	if ok, desc := other.Contains(cm); !ok {
		return false, desc
	}

	return true, ""
}