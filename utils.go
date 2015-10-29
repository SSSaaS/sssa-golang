package sssa

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
)

var prime *big.Int

/**
 * Returns a random number from the range (0, prime-1) inclusive
**/
func random() *big.Int {
	result := big.NewInt(0).Set(prime)
	result = result.Sub(result, big.NewInt(1))
	result, _ = rand.Int(rand.Reader, result)
	return result
}

/**
 * Converts a byte array into an a 256-bit big.Int, arraied based upon size of
 * the input byte; all values are right-padded to length 256, even if the most
 * significant bit is zero.
**/
func splitByteToInt(secret []byte) []*big.Int {
	count := int(math.Ceil(float64(len(secret)) / 32))

	var result []*big.Int = make([]*big.Int, count)
	for i := range result {
		data := ""
		if (i+1)*32 < len(secret) {
			data = hex.EncodeToString(secret[i*32 : (i+1)*32])
		} else {
			tmp := make([]byte, 32)
			copy(tmp, secret[i*32:])
			data = hex.EncodeToString(tmp)
		}

		result[i], _ = big.NewInt(0).SetString(data, 16)
	}

	return result
}

/**
 * Converts an array of big.Ints to the original byte array, removing any
 * least significant nulls
**/
func mergeIntToByte(secret []*big.Int) []byte {
	var result []byte
	for i := range secret {
		tmp, _ := hex.DecodeString(fmt.Sprintf("%x", secret[i]))
		result = append(result, tmp...)
	}

	return result
}

/**
 * Evauluates a polynomial with coefficients specified in reverse order:
 * evaluatePolynomial([a, b, c, d], x):
 * 		returns a + bx + cx^2 + dx^3
**/
func evaluatePolynomial(polynomial []*big.Int, value *big.Int) *big.Int {
	var result *big.Int = big.NewInt(0).Set(polynomial[0])

	for s := range polynomial[1:] {
		tmp := big.NewInt(0).Set(value)
		tmp = tmp.Exp(tmp, big.NewInt(int64(s)+1), prime)
		tmp = tmp.Mul(tmp, polynomial[s+1])
		result = result.Add(result, tmp)
		result = result.Mod(result, prime)
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
 * Returns the big.Int number base10 in base64 representation; note: this is
 * not a string representation; the base64 output is exactly 256 bits long
**/
func toBase64(number *big.Int) string {
	hexdata := fmt.Sprintf("%x", number)
	for i := 0; len(hexdata) < 64; i++ {
		hexdata = "0" + hexdata
	}
	bytedata, success := hex.DecodeString(hexdata)
	if success != nil {
		fmt.Println("Error!")
		fmt.Println("hexdata: ", hexdata)
		fmt.Println("bytedata: ", bytedata)
		fmt.Println(success)
	}
	return base64.URLEncoding.EncodeToString(bytedata)
}

/**
 * Returns the number base64 in base 10 big.Int representation; note: this is
 * not coming from a string representation; the base64 input is exactly 256
 * bits long, and the output is an arbitrary size base 10 integer
**/
func fromBase64(number string) *big.Int {
	bytedata, _ := base64.URLEncoding.DecodeString(number)
	hexdata := hex.EncodeToString(bytedata)
	result, _ := big.NewInt(0).SetString(hexdata, 16)
	return result
}

/**
 * Computes the multiplicitive inverse of the number on the field prime; more
 * specifically, number * inverse == 1; Note: number should neve rbe zero
**/
func modInverse(number *big.Int) *big.Int {
	copy := big.NewInt(0).Set(number)
	copy = copy.Mod(copy, prime)
	pcopy := big.NewInt(0).Set(prime)
	x := big.NewInt(0)
	y := big.NewInt(0)

	copy.GCD(x, y, pcopy, copy)

	result := big.NewInt(0).Set(prime)

	result = result.Add(result, y)
	result = result.Mod(result, prime)
	return result
}
