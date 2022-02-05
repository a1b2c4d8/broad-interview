package util

import (
	. "broad-interview/internal/service"
	"sort"
)

type RouteIdSet struct {
	Set map[RouteId]struct{}
}

func NewRouteIdSet() *RouteIdSet {
	return &RouteIdSet{make(map[RouteId]struct{})}
}

func NewRouteIdSetFrom(init map[RouteId]struct{}) *RouteIdSet {
	s := NewRouteIdSet()
	for k, v := range init {
		s.Set[k] = v
	}
	return s
}

func (s *RouteIdSet) Add(rid RouteId) {
	s.Set[rid] = struct{}{}
}

func (s *RouteIdSet) AddAll(from *RouteIdSet) {
	s.addAll(from.Set)
}

func (s *RouteIdSet) AddAllFromMap(from map[RouteId]struct{}) {
	s.addAll(from)
}

func (s *RouteIdSet) addAll(from map[RouteId]struct{}) {
	for k, v := range from {
		s.Set[k] = v
	}
}

func (s *RouteIdSet) ContainsAll(other *RouteIdSet) bool {
	if len(s.Set) < len(other.Set) {
		return false
	}

	for k := range other.Set {
		_, found := s.Set[k]
		if !found {
			return false
		}
	}

	return true
}

func (s *RouteIdSet) Delete(rid RouteId) {
	delete(s.Set, rid)
}

func (s *RouteIdSet) DeleteAll(from *RouteIdSet) {
	for k := range from.Set {
		delete(s.Set, k)
	}
}

func (s *RouteIdSet) Sorted() []RouteId {
	sorted := make([]RouteId, len(s.Set))
	index := 0

	for k := range s.Set {
		sorted[index] = k
		index++
	}

	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	return sorted
}
