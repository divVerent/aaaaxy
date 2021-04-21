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
	index listIndex
	items []*Entity
}

func makeList(index listIndex) entityList {
	return entityList{index: index, items: nil}
}

func (l *entityList) verify(step string) {
	// Uncomment this if you suspect lists may get corrupted.
	/*
		if l == nil {
			log.Panicf("verifying nil list")
		}
		for i, e := range l.items {
			if e == nil {
				continue
			}
			idx := e.indexInListPlusOne[l.index] - 1
			if idx != i {
				log.Panicf("%v: corrupted entity list %v: entity %v: got index %v, want %v", step, l.index, e, e.indexInListPlusOne[l.index]-1, i)
			}
		}
	*/
}

func (l *entityList) insert(e *Entity) {
	l.verify("insert pre")
	if e.indexInListPlusOne[l.index] != 0 {
		log.Panicf("inserting into the same entity list twice: entity %v, items %v", e, l.index)
	}
	l.items = append(l.items, e)
	e.indexInListPlusOne[l.index] = len(l.items)
	l.verify("insert post")
}

func (l *entityList) remove(e *Entity) {
	l.verify("remove pre")
	idxPlusOne := e.indexInListPlusOne[l.index]
	if idxPlusOne == 0 {
		log.Panicf("removing from an entity list the entity isn't in: entity %v, list %v", e, l.index)
	}
	idx := idxPlusOne - 1
	if l.items[idx] != e {
		log.Panicf("removing from a corrupted entity list %v: entity %v isn't actually at index %v", l.index, e, idx)
	}
	l.items[idx] = nil
	e.indexInListPlusOne[l.index] = 0
	l.verify("remove post")
}

func (l *entityList) compact() {
	l.verify("compact pre")
	n := 0
	for _, e := range l.items {
		if e == nil {
			continue
		}
		l.items[n] = e
		n++
		e.indexInListPlusOne[l.index] = n
	}
	l.items = l.items[:n]
	l.verify("compact post")
}

var breakError = errors.New("break")

func (l *entityList) forEach(f func(e *Entity) error) error {
	for _, e := range l.items {
		if e == nil {
			continue
		}
		err := f(e)
		if err != nil {
			return err
		}
	}
	return nil
}
