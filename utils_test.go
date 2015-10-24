package sssa

import (
	"math/big"
	"bytes"
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

func TestSplitMerge(t *testing.T) {
	// Short, medium, and long tests
	tests := [][]byte {
		[]byte("N17FigASkL6p1EOgJhRaIquQLGvYV0"),
		[]byte("0y10VAfmyH7GLQY6QccCSLKJi8iFgpcSBTLyYOGbiYPqOpStAf1OYuzEBzZR"),
		[]byte("KjRHO1nHmIDidf6fKvsiXWcTqNYo2U9U8juO94EHXVqgearRISTQe0zAjkeUYYBvtcB8VWzZHYm6ktMlhOXXCfRFhbJzBUsXaHb5UDQAvs2GKy6yq0mnp8gCj98ksDlUultqygybYyHvjqR7D7EAWIKPKUVz4of8OzSjZlYg7YtCUMYhwQDryESiYabFID1PKBfKn5WSGgJBIsDw5g2HB2AqC1r3K8GboDN616Swo6qjvSFbseeETCYDB3ikS7uiK67ErIULNqVjf7IKoOaooEhQACmZ5HdWpr34tstg18rO"),
	}

	for i := range(tests) {
		if (bytes.Equal(mergeIntToByte(splitByteToInt(tests[i])), tests[i])) {
			t.Fatal("Fatal: splitting and merging returned invalid data")
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
