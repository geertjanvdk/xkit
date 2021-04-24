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
		xt.Eq(t, false, om.Has("somekey"))
	})

	t.Run("retrieve keys and values", func(t *testing.T) {
		om := OrderedMap{}
		om.Set("key3", "value3")
		om.Set("key1", 1.1)
		om.Set("key2", 2)
		xt.Eq(t, true, om.Has("key2"))

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

	t.Run("retrieve key", func(t *testing.T) {
		om := OrderedMap{}
		om.Set("key3", "value3")
		om.Set("key1", 1.1)
		om.Set("key2", 2)
		om.Set("key4", nil)

		cases := map[string]struct {
			key     string
			exp     interface{}
			expHave bool
		}{
			"nil":      {key: "key4", exp: nil, expHave: true},
			"notkey":   {key: "notkey", exp: nil, expHave: false},
			"key3":     {key: "key3", exp: "value3", expHave: true},
			"nilnokey": {key: "nilnokey", exp: nil, expHave: false},
		}

		for name, cs := range cases {
			t.Run(name, func(t *testing.T) {
				v, have := om.Value(cs.key)
				xt.Eq(t, cs.expHave, have)
				xt.Eq(t, cs.exp, v)
			})
		}
	})

	t.Run("delete element", func(t *testing.T) {
		om := OrderedMap{}
		om.Set("key3", "value3")
		om.Set("key1", 1.1)
		om.Set("key2", 2)
		om.Set("key4", nil)

		xt.Eq(t, 4, om.Count())
		om.Delete("noSuchKey")

		xt.Eq(t, 4, om.Count())

		om.Delete("key2")
		xt.Eq(t, 3, om.Count())
		_, have := om.Value("key2")
		xt.Eq(t, false, have)

		om.Delete("key3")
		xt.Eq(t, 2, om.Count())
		_, have = om.Value("key3")
		xt.Eq(t, false, have)

		om.Delete("key3")
		xt.Eq(t, 2, om.Count())
		_, have = om.Value("key3")
		xt.Eq(t, false, have)
	})
}
