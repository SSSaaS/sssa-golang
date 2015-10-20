package sssaas

import (
	"math/big"
)

func create(minimum int, shares int, secret []big.Int) [][]big.Int {
	var result [][]big.Int = make([][]big.Int, len(secret))
	prime, _ = big.NewInt(0).SetString("96911788199763998566185843021439103838446442331102965305766889944557597472419", 10)

	for s := range secret {
		var polynomial []big.Int = make([]big.Int, minimum)

		polynomial[0] = secret[s]

		for i := range polynomial[1:] {
			polynomial[i] = random()
		}
		result[s] = polynomial
	}

	return result 
}

func Create(minimum int, shares int, raw string) []string {
	var secret []big.Int = split([]byte(raw))

	return pretty(create(minimum, shares, secret))
}
