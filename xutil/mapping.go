// Copyright (c) 2021, Geert JM Vanderkelen

package xutil

import "sync"

// OrderedMap wraps around a Go map keeping the order with which
// elements have been added. Keys must be strings, but values
// can be anything (interface{}).
// Unlike map, index assigment is not possible. Use the `Set`
// method to set a key with a particular value.
// Use the Keys method to retrieves keys, Values to get the
// values. To get both, which probably what you want, use
// the KeysValues method.
type OrderedMap struct {
	mapMU sync.RWMutex
	map_  map[string]interface{}
	order []string
}

// Count returns the number of elements in the map.
func (om *OrderedMap) Count() int {
	return len(om.order)
}

// Set key in OrderedMap to value. Previously stored values
// are overwritten, but the order does not change.
func (om *OrderedMap) Set(key string, value interface{}) {
	om.mapMU.Lock()
	defer om.mapMU.Unlock()

	if om.map_ == nil {
		om.map_ = map[string]interface{}{}
	}

	om.map_[key] = value
	if !HasString(om.order, key) {
		om.order = append(om.order, key)
	}
}

// Keys returns keys as slice of string.
func (om *OrderedMap) Keys() []string {
	om.mapMU.RLock()
	defer om.mapMU.RUnlock()

	return om.order
}

func (om *OrderedMap) values() []interface{} {
	res := make([]interface{}, len(om.order))
	for i, k := range om.order {
		res[i] = om.map_[k]
	}
	return res
}

// Values returns the values as slice of interfaces.
func (om *OrderedMap) Values() []interface{} {
	om.mapMU.RLock()
	defer om.mapMU.RUnlock()

	return om.values()
}

// KeysValues returns the keys as slice of strings, and values as slice of interfaces.
func (om *OrderedMap) KeysValues() ([]string, []interface{}) {
	om.mapMU.RLock()
	defer om.mapMU.RUnlock()

	return om.order, om.values()
}
