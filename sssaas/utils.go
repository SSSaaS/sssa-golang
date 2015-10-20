package sssaas

import (
	"math"
	"math/big"
	"crypto/rand"
	"encoding/hex"
)

var prime *big.Int

func random() *big.Int {
	var result *big.Int
	result, _ = rand.Int(rand.Reader, prime.Sub(prime, big.NewInt(1))) 
	return result
}

func splitByteToInt(secret []byte) []*big.Int {
	count := int(math.Ceil(float64(len(secret))/32))

	var result []*big.Int = make([]*big.Int, count)
	for i := range result {
		data := ""
		if (i+1)*32 < len(secret) {
			data = hex.EncodeToString(secret[i*32:(i+1)*32])
		} else {
			tmp := make([]byte, 32)
			copy(tmp, secret[i*32:])
			data = hex.EncodeToString(tmp)
		}
		
		result[i], _ = big.NewInt(0).SetString(data, 16)
	}

	return result
}

func evaluatePolynomial(polynomial []*big.Int, value *big.Int) *big.Int {
	var result *big.Int = polynomial[0] 
	
	for s := range polynomial[1:] {
		result = result.Add(result, result.Exp(polynomial[s+1], big.NewInt(int64(s)+1), prime))
		result = result.Mod(result, prime)
	}
	
	return result
}

func inNumbers(numbers []*big.Int, value *big.Int) bool {
	for n := range numbers {
		if (numbers[n].Cmp(value) == 0) {
			return true
		}
	}
	
	return false
}

/*
func join(secret []bigInt) []byte {
	var result []byte
	for i := range secret {
		result = append(result, secret[i]...)
	}

	return result
}

func chomp(array []byte) []byte {
        last := len(array)
        for array[last-1] == byte(0) {
                last -= 1
        }
        
        return array[0:last]
}

func pretty(shares [][]big.Int) []string {
	var result []string = make([]string, len(shares))

	for i := range shares[0] {
		result[i] = base64.URLEncoding.EncodeToString(shares[i].String())
	}

	return result
}
*/