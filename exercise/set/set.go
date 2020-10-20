package set

import (
	"bytes"
	"fmt"
)

// Set ...
type Set struct {
	m map[interface{}]bool
}

// NewSet return a pointer to set
func NewSet() *Set {
	return &Set{
		m: make(map[interface{}]bool),
	}
}

// Add add e to set
func (s *Set) Add(e interface{}) bool {
	if s.m[e] {
		return false
	}
	s.m[e] = true
	return true
}

// Del remove e from set
func (s *Set) Del(e interface{}) {
	delete(s.m, e)
}

// Clear delete all elements
func (s *Set) Clear() {
	s.m = make(map[interface{}]bool)
}

// Contain e or not
func (s *Set) Contain(e interface{}) bool {
	return s.m[e]
}

// Len return number of elements
func (s *Set) Len() int {
	return len(s.m)
}

// Equal to another set or not: only true if contain the same elements
func (s *Set) Equal(o *Set) bool {
	if o == nil {
		return false
	}

	if s.Len() != o.Len() {
		return false
	}

	for e := range s.m {
		if !o.Contain(e) {
			return false
		}
	}

	return true
}

// Elements return all elements in set now: set may be modified during iteration
func (s *Set) Elements() []interface{} {
	nowLen := s.Len()
	eles := make([]interface{}, nowLen)
	num := 0
	for e := range s.m {
		if num < nowLen {
			eles[num] = e
		} else {
			// set may be larger during iteration
			eles = append(eles, e)
		}
		num++
	}
	// set may be smaller during iteration
	if num < nowLen {
		eles = eles[:num]
	}
	return eles
}

// String return a string with all elements
func (s *Set) String() string {
	var buf bytes.Buffer
	var sep string
	buf.WriteString(fmt.Sprintf("Set{ "))
	for k := range s.m {
		buf.WriteString(fmt.Sprintf("%v%v", sep, k))
		sep = " "
	}
	buf.WriteString(fmt.Sprintf(" }"))
	return buf.String()
}

// IsSuperSet return if s is super set of o
func (s *Set) IsSuperSet(o *Set) bool {
	if o == nil {
		return true
	}

	sLen, oLen := s.Len(), o.Len()
	if sLen < oLen {
		return false
	}
	if sLen > 0 && oLen == 0 {
		return true
	}

	for _, e := range o.Elements() {
		if !s.Contain(e) {
			return false
		}
	}

	return true
}

// Union return elements including in s or o
func (s *Set) Union(o *Set) *Set {
	ret := NewSet()
	for _, e := range s.Elements() {
		ret.Add(e)
	}
	for _, e := range o.Elements() {
		ret.Add(e)
	}
	return ret
}

// Intersect return elements including in s and o
func (s *Set) Intersect(o *Set) *Set {
	ret := NewSet()
	for _, e := range s.Elements() {
		if o.Contain(e) {
			ret.Add(e)
		}
	}
	return ret
}

// Difference return elements including in s but not in o
func (s *Set) Difference(o *Set) *Set {
	ret := NewSet()
	for _, e := range s.Elements() {
		if !o.Contain(e) {
			ret.Add(e)
		}
	}
	return ret
}

// SymmetricDifference return elements including in s or o, but not including by s and o
func (s *Set) SymmetricDifference(o *Set) *Set {
	ret := s.Difference(o)
	for _, e := range o.Elements() {
		if !s.Contain(e) {
			ret.Add(e)
		}
	}
	return ret
}
