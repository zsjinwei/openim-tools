// Copyright © 2024 OpenIM open source community. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cacheutil

import "sync"

// Cache is a Generic sync.Map structure.
type Cache[K comparable, V any] struct {
	m sync.Map
}

func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{}
}

// Load returns the value stored in the map for a key, or nil if no value is present.
func (c *Cache[K, V]) Load(key K) (value V, ok bool) {
	rawValue, ok := c.m.Load(key)
	if !ok {
		return
	}
	return rawValue.(V), ok
}

// Store sets the value for a key.
func (c *Cache[K, V]) Store(key K, value V) {
	c.m.Store(key, value)
}

// StoreAll sets all value by f's key.
func (c *Cache[K, V]) StoreAll(f func(value V) K, values []V) {
	for _, v := range values {
		c.m.Store(f(v), v)
	}
}

// LoadOrStore returns the existing value for the key if present.
func (c *Cache[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	rawValue, loaded := c.m.LoadOrStore(key, value)
	return rawValue.(V), loaded
}

// Delete deletes the value for a key.
func (c *Cache[K, V]) Delete(key K) {
	c.m.Delete(key)
}

// DeleteAll deletes all values.
func (c *Cache[K, V]) DeleteAll() {
	c.m.Range(func(key, value interface{}) bool {
		c.m.Delete(key)
		return true
	})
}

// RangeAll returns all values in the map.
func (c *Cache[K, V]) RangeAll() (values []V) {
	c.m.Range(func(rawKey, rawValue interface{}) bool {
		values = append(values, rawValue.(V))
		return true
	})
	return values
}

// RangeCon returns values in the map that satisfy condition f.
func (c *Cache[K, V]) RangeCon(f func(key K, value V) bool) (values []V) {
	c.m.Range(func(rawKey, rawValue interface{}) bool {
		if f(rawKey.(K), rawValue.(V)) {
			values = append(values, rawValue.(V))
		}
		return true
	})
	return values
}
