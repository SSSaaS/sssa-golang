package sssaas

import (
	"math"
	"math/big"
	"crypto/rand"
	"encoding/base64"
)

var prime *big.Int

func random() big.Int {
	result, _ := rand.Int(rand.Reader, prime.Sub(prime, big.NewInt(1))) 
	return *result
}

func split(secret []byte) []big.Int {
	count := int(math.Ceil(float64(len(secret))/32))

	var result []big.Int = make([]big.Int, count)
	for i := range result {
		tmp := make([]byte, 32)
		if (i+1)*32 < len(secret) {
			copy(tmp, secret[i*32:(i+1)*32])
		} else {
			copy(tmp, secret[i*32:])
		}
		result[i] = big.NewInt(0).SetString(string(tmp))
	}

	return result
}

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

	for i := range shares {
		result[i] = base64.URLEncoding.EncodeToString(shares[i].String())
	}

	return result
}
