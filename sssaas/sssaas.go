package sssaas

import (
	"math/big"
)

func Create(minimum int, shares int, raw string) []string {
	var secret []*big.Int = splitByteToInt([]byte(raw))
	prime, _ = big.NewInt(0).SetString("96911788199763998566185843021439103838446442331102965305766889944557597472419", 10)

	var polynomial [][]*big.Int = make([][]*big.Int, len(secret))
	for i := range polynomial {
		polynomial[i] = make([]*big.Int, minimum)
	}

	var secrets [][][]*big.Int = make([][][]*big.Int, shares)
	for i := range secrets {
		secrets[i] = make([][]*big.Int, len(secret))
		for j := range secrets[i] {
			secrets[i][j] = make([]*big.Int, 2)
		}
	}

	var numbers []*big.Int = make([]*big.Int, 0)

	for s := range secret {
		polynomial[s][0] = secret[s]

		for i := range polynomial[s][1:] {
			number := random()
			for inNumbers(numbers, number) {
				number = random()
			}
			numbers = append(numbers, number)

			polynomial[s][i+1] = number
		}
	}

	var result []string = make([]string, shares)

	for i := range secrets {
		for j := range secret {
			number := random()
			for inNumbers(numbers, number) {
				number = random()
			}
			numbers = append(numbers, number)

			secrets[i][j][0] = number
			secrets[i][j][1] = evaluatePolynomial(polynomial[j], number)

			result[i] += toBase64(secrets[i][j][0])
			result[i] += toBase64(secrets[i][j][1])
		}
	}

	return result
}
