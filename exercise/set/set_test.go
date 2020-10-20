package set

import (
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	s := NewSet()
	s.Add(1)
	s.Add("2")
	s.Add(3.0)
	s.Add(struct {
		name string
		age  int
	}{"kiki", 28})
	fmt.Println("s:", s)

	s2 := NewSet()
	s2.Add(3.0)
	s2.Add("2")
	s2.Add(1)
	s2.Add(struct {
		name string
		age  int
	}{"kiki", 28})
	fmt.Println("s2:", s2)
	fmt.Println("s.Equal(s2):", s.Equal(s2))

	s2.Del(struct {
		name string
		age  int
	}{"kiki", 28})
	fmt.Println("s2:", s2)
	fmt.Println("s.IsSuperSet(s2):", s.IsSuperSet(s2))

	s2.Add(4)
	fmt.Println("s2:", s2)
	fmt.Println("s.IsSuperSet(s2):", s.IsSuperSet(s2))
	fmt.Println("s.Union(s2):", s.Union(s2))
	fmt.Println("s.Intersect(s2):", s.Intersect(s2))
	fmt.Println("s.Difference(s2):", s.Difference(s2))
	fmt.Println("s.SymmetricDifference(s2):", s.SymmetricDifference(s2))
}
