package bigint

import (
	"testing"

)

type argMod struct {
	x, y, m, want *Int
}
var modVec = []argMod{
	{NewInt(0), NewInt(1), NewInt(1152921504606489097), NewInt(0)},
	{NewInt(1), NewInt(1), NewInt(1152921504606489097), NewInt(1)},
	{NewInt(112233445566778899), NewInt(998877665544332211), NewInt(1152921504606489097), NewInt(1103875254192881828)},
}

func TestMod(t *testing.T) {
	var z Int
	for i, testPair := range modVec {
		z.Mul(testPair.x,testPair.y)
		if !z.Mod(&z, testPair.m).EqualTo(testPair.want) {
			t.Errorf("Error Mod test pair %v", i)
		}
	}
}

func BenchmarkMod(b *testing.B) {
	var z Int
	x := NewInt(112233445566778899)
	y := NewInt(998877665544332211)
	q := NewInt(1152921504606489097)
	for i := 0; i < b.N; i++ {
		z.Mul(x, y)
		z.Mod(&z, q)
	}
}