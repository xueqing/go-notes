package round

import "testing"

func TestRound(t *testing.T) {
	input := []float64{0, 0.2, 0.5, 0.7, 1, 1.2, 1.5, 1.7}
	wanted := []float64{0, 0, 1, 1, 1, 1, 2, 2}
	const diff = 0.00000001
	for i, a := range input {
		res := Round(a)
		if wanted[i]-res > diff || res-wanted[i] > diff {
			t.Errorf("Round(%v):%v != wanted(%v)\n", a, res, wanted[i])
		}
	}
}
