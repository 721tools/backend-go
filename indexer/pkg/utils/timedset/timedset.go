package timedset

import (
	"container/list"
	"sync"
	"time"
)

type TimedSet struct {
	m    map[interface{}]*list.Element
	lst  *list.List
	dur  time.Duration
	lock sync.Mutex
}

type timedSetEntry struct {
	value interface{}
	timer *time.Timer
}

func NewTimedSet(dur time.Duration) *TimedSet {
	return &TimedSet{
		m:   make(map[interface{}]*list.Element),
		lst: list.New(),
		dur: dur,
	}
}

func (s *TimedSet) Add(value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.m[value]; ok {
		return
	}

	timer := time.AfterFunc(s.dur, func() {
		s.Remove(value)
	})

	entry := &timedSetEntry{
		value: value,
		timer: timer,
	}

	node := s.lst.PushFront(entry)
	s.m[value] = node
}

func (s *TimedSet) Remove(value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	node, ok := s.m[value]
	if !ok {
		return
	}

	delete(s.m, value)
	s.lst.Remove(node)

	entry := node.Value.(*timedSetEntry)
	entry.timer.Stop()
}

func (s *TimedSet) Len() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return len(s.m)
}

func (s *TimedSet) GetAll() []interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()

	result := make([]interface{}, 0, len(s.m))
	for e := s.lst.Front(); e != nil; e = e.Next() {
		entry := e.Value.(*timedSetEntry)
		result = append(result, entry.value)
	}
	return result
}
