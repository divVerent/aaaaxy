// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package engine

import (
	"errors"
	"log"
)

type listIndex int

const (
	allList listIndex = iota
	opaqueList
	zList
	numLists
)

type entityList struct {
	index       listIndex
	items       []*Entity
	shadowItems []*Entity
}

func makeList(index listIndex) entityList {
	return entityList{index: index, items: nil}
}

func (l *entityList) insert(e *Entity) {
	if e.indexInListPlusOne[l.index] != 0 {
		log.Panicf("inserting into the same intrusive items twice: entity %v, items %v", e, l.index)
	}
	l.items = append(l.items, e)
	e.indexInListPlusOne[l.index] = len(l.items)
}

func (l *entityList) remove(e *Entity) {
	last := len(l.items) - 1
	idxPlusOne := e.indexInListPlusOne[l.index]
	if idxPlusOne == 0 {
		log.Panicf("removing from an intrusive items the entity isn't in: entity %v, list %v", e, l.index)
	}
	idx := idxPlusOne - 1
	if idx != last {
		other := l.items[last]
		l.items[idx] = other
		other.indexInListPlusOne[l.index] = idxPlusOne
	}
	l.items = l.items[0:last]
	e.indexInListPlusOne[l.index] = 0
}

var breakError = errors.New("break")

func (l *entityList) forEach(f func(e *Entity) error) error {
	// TODO: Can we make this work without copying?
	// The hard part is not breaking the iteration when existing elements are relinked.
	n := len(l.items)
	for len(l.shadowItems) < n {
		l.shadowItems = append(l.shadowItems, nil)
	}
	copy(l.shadowItems, l.items)
	for _, e := range l.shadowItems[:n] {
		err := f(e)
		if err != nil {
			return err
		}
	}
	return nil
}
