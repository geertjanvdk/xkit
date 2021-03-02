// Copyright (c) 2021, Geert JM Vanderkelen

package xutil

import (
	"testing"

	"github.com/geertjanvdk/xkit/xt"
)

func TestOrderedMap(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		om := OrderedMap{}
		xt.Eq(t, 0, om.Count())
		xt.Eq(t, 0, len(om.Keys()))
		xt.Eq(t, 0, len(om.Values()))
		keys, values := om.KeysValues()
		xt.Eq(t, 0, len(keys))
		xt.Eq(t, 0, len(values))
	})

	t.Run("retrieve keys and values", func(t *testing.T) {
		om := OrderedMap{}
		om.Set("key3", "value3")
		om.Set("key1", 1.1)
		om.Set("key2", 2)

		expKeys := []string{"key3", "key1", "key2"}
		expValues := []interface{}{"value3", 1.1, 2}

		xt.Eq(t, expKeys, om.Keys())
		xt.Eq(t, expValues, om.Values())
		keys, values := om.KeysValues()
		xt.Eq(t, expKeys, keys)
		xt.Eq(t, expValues, values)
	})

	t.Run("set already set does not change order", func(t *testing.T) {
		om := OrderedMap{}
		om.Set("key3", "value3")
		om.Set("key1", 1.1)
		om.Set("key2", 2)

		om.Set("key3", "value number 3")

		expKeys := []string{"key3", "key1", "key2"}
		expValues := []interface{}{"value number 3", 1.1, 2}

		keys, values := om.KeysValues()
		xt.Eq(t, expKeys, keys)
		xt.Eq(t, expValues, values)
		xt.Eq(t, len(expKeys), om.Count())
	})
}
