package sssaas

import (
	"math/big"
	"fmt"
)

func Create(minimum int, shares int, raw string) []string {
	if (minimum > shares) {
		return []string{""}
	}

	var secret []*big.Int = splitByteToInt([]byte(raw))
	prime, _ = big.NewInt(0).SetString("99995644905598542077721161034987774965417302630805822064337798850767846245779", 10)
	var numbers []*big.Int = make([]*big.Int, 0)

	fmt.Println(secret)

	var polynomial [][]*big.Int = make([][]*big.Int, len(secret))
	for i := range polynomial {
		polynomial[i] = make([]*big.Int, minimum)
		polynomial[i][0] = secret[i]

		for j := range polynomial[i][1:] {
			number := random()
			for inNumbers(numbers, number) {
				number = random()
			}
			numbers = append(numbers, number)

			polynomial[i][j+1] = number
		}
	}

	var secrets [][][]*big.Int = make([][][]*big.Int, shares)
	var result []string = make([]string, shares)
	for i := range secrets {
		secrets[i] = make([][]*big.Int, len(secret))
		for j := range secrets[i] {
			secrets[i][j] = make([]*big.Int, 2)

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

	fmt.Println("Secrets: ", result)

	return result
}

/**
 * Note: the polynomial will converge if the specified minimum number of shares
 *       or more are passed to this function. Passing thus does not affect it
 *       Passing fewer however, simply means that the returned secret is wrong.
**/
func Combine(shares []string) string {
	var secrets [][][]*big.Int = make([][][]*big.Int, len(shares))
	prime, _ = big.NewInt(0).SetString("99995644905598542077721161034987774965417302630805822064337798850767846245779", 10)

	for i := range shares {
		if (len(shares[i]) % 44) != 0 {
			return ""
		}

		share := shares[i]
		count := len(share)/88
		secrets[i] = make([][]*big.Int, count)

		for j := range secrets[i] {
			cshare := share[j*88:(j+1)*88]
			secrets[i][j] = make([]*big.Int, 2)
			secrets[i][j][0] = fromBase64(cshare[0:44])
			secrets[i][j][1] = fromBase64(cshare[44:])
		}
	}

	var secret []*big.Int = make([]*big.Int, len(secrets[0]))
	for j := range secret {
		secret[j] = big.NewInt(0)
		for i := range secrets {
			x := secrets[i][j][0]
			y := secrets[i][j][1]
			lagrange := big.NewInt(0)
			first := true
			for k := range secrets {
				if k != i {
					value := big.NewInt(0)
					numerator := big.NewInt(0)
					denominator := big.NewInt(0).Set(x)
					numerator = numerator.Sub(value, secrets[k][j][0])
					denominator = denominator.Sub(denominator, secrets[k][j][0])
					value, _ = value.DivMod(numerator, denominator, prime)
					if first {
						lagrange = big.NewInt(0).Set(value)
						lagrange = lagrange.Mod(lagrange, prime)
					} else {
						lagrange = lagrange.Mul(lagrange, value)
						lagrange = lagrange.Mod(lagrange, prime)
					}
				}
			}
			lagrange = lagrange.Mul(lagrange, y)
			lagrange = lagrange.Mod(lagrange, prime)
			fmt.Println("La: ", lagrange.String())
			secret[j] = secret[j].Add(secret[j], lagrange)
			secret[j] = secret[j].Mod(secret[j], prime)
		}
	}

	fmt.Println("Secret: ", secret)

	return "Hi!"
}
