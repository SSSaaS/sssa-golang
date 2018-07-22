package shamir

import (
	"math/big"
	"testing"
)

func TestRandom(t *testing.T) {
	for i := 0; i < 100; i++ {
		if random().Cmp(P) >= 0 {
			t.Fatal("Error! Random number out of bounds exception")
		}
	}
}

func TestModInverse(t *testing.T) {
	for i := 0; i < 100; i++ {
		point := big.NewInt(0).Set(random())
		if point.Cmp(P) >= 0 {
			t.Fatal("Error! Random point out of bounds exception")
		}
		inverse := big.NewInt(0).Set(modInverse(point))
		if inverse.Cmp(P) >= 0 {
			t.Fatal("Error! Inverse out of bounds exception")
		}
		value := big.NewInt(0).Set(point)
		value = value.Mul(value, inverse)
		value = value.Mod(value, P)
		expected := big.NewInt(1)
		if value.Cmp(expected) != 0 {
			t.Fatalf("Fatal: modInverse[%v] failed\nExpected: %v; Got: %v\nPoint %v\nInverse: %v\nP: %v", i, expected, value, point, inverse, P)
		}
	}
}

func TestEvaluatePolynomial(t *testing.T) {
	values := [][][]*big.Int{
		[][]*big.Int{
			[]*big.Int{big.NewInt(20), big.NewInt(21), big.NewInt(42)},
			[]*big.Int{big.NewInt(0)},
		},
		[][]*big.Int{
			[]*big.Int{big.NewInt(0), big.NewInt(0), big.NewInt(0)},
			[]*big.Int{big.NewInt(4)},
		},
		[][]*big.Int{
			[]*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3), big.NewInt(4), big.NewInt(5)},
			[]*big.Int{big.NewInt(10)},
		},
	}

	actual := []*big.Int{big.NewInt(20), big.NewInt(0), big.NewInt(54321)}

	for i := range values {
		result := evaluatePolynomial(values[i][0], values[i][1][0])
		if result.Cmp(actual[i]) != 0 {
			t.Fatalf("Fatal: EvaluatePolynomial[%v] failed\nExpected: %v; Got: %v\n", i, actual[i], result)
		}
	}
}
