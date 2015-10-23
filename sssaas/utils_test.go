package sssaas

import (
	"math/big"
	"testing"
)

func TestRandom(t *testing.T) {
	prime, _ = big.NewInt(0).SetString("99995644905598542077721161034987774965417302630805822064337798850767846245779", 10)
	for i := 0; i < 10000; i++ {
		if random().Cmp(prime) >= 0 {
			t.Fatal("Error! Random number out of bounds exception")
		}
	}
}

func TestBaseConversion(t *testing.T) {
	prime, _ = big.NewInt(0).SetString("99995644905598542077721161034987774965417302630805822064337798850767846245779", 10)
	for i := 0; i < 10000; i++ {
		point := random()
		if point.Cmp(fromBase64(toBase64(point))) != 0 {
			t.Fatal("Fatal: Base conversion failed")
		}
	}
}

func TestToBase64(t *testing.T) {
	prime, _ = big.NewInt(0).SetString("99995644905598542077721161034987774965417302630805822064337798850767846245779", 10)
	for i := 0; i < 10000; i++ {
		point := random()
		if len(toBase64(point)) != 44 {
			t.Fatal("Fatal: toBase64 returned wrong length")
		}
	}
}

func TestModInverse(t *testing.T) {
	prime, _ = big.NewInt(0).SetString("99995644905598542077721161034987774965417302630805822064337798850767846245779", 10)
	for i := 0; i < 10000; i++ {
		point := big.NewInt(0).Set(random())
		if point.Cmp(prime) >= 0 {
			t.Fatal("Error! Random point out of bounds exception")
		}
		inverse := big.NewInt(0).Set(modInverse(point))
		if inverse.Cmp(prime) >= 0 {
			t.Fatal("Error! Inverse out of bounds exception")
		}
		value := big.NewInt(0).Set(point)
		value = value.Mul(value, inverse)
		value = value.Mod(value, prime)
		expected := big.NewInt(1)
		if value.Cmp(expected) != 0 {
			t.Fatalf("Fatal: modInverse[%v] failed\nExpected: %v; Got: %v\nPoint %v\nInverse: %v\nPrime: %v", i, expected, value, point, inverse, prime)
		}
	}
}
