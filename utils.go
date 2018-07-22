package shamir

import (
	"crypto/rand"
	"math/big"
)

/**
 * Returns a random number from the range (0, P-1) inclusive
**/
func random() *big.Int {
	result := big.NewInt(0).Set(P)
	result = result.Sub(result, big.NewInt(1))
	result, _ = rand.Int(rand.Reader, result)
	return result
}

/**
 * Evauluates a polynomial with coefficients specified in reverse order:
 * evaluatePolynomial([a, b, c, d], x):
 * 		returns a + bx + cx^2 + dx^3
**/
func evaluatePolynomial(polynomial []*big.Int, value *big.Int) *big.Int {
	last := len(polynomial) - 1
	var result *big.Int = big.NewInt(0).Set(polynomial[last])

	for s := last - 1; s >= 0; s-- {
		result = result.Mul(result, value)
		result = result.Add(result, polynomial[s])
		result = result.Mod(result, P)
	}

	return result
}

/**
 * inNumbers(array, value) returns boolean whether or not value is in array
**/
func inNumbers(numbers []*big.Int, value *big.Int) bool {
	for n := range numbers {
		if numbers[n].Cmp(value) == 0 {
			return true
		}
	}

	return false
}

/**
 * Computes the multiplicative inverse of the number on the field P; more
 * specifically, number * inverse == 1; Note: number should never be zero
**/
func modInverse(number *big.Int) *big.Int {
	copy := big.NewInt(0).Set(number)
	copy = copy.Mod(copy, P)
	pcopy := big.NewInt(0).Set(P)
	x := big.NewInt(0)
	y := big.NewInt(0)

	copy.GCD(x, y, pcopy, copy)

	result := big.NewInt(0).Set(P)

	result = result.Add(result, y)
	result = result.Mod(result, P)
	return result
}
