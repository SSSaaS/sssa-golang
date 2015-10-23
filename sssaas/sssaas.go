package sssaas

import (
	"fmt"
	"math/big"
)

func Create(minimum int, shares int, raw string) []string {
	if minimum > shares {
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

	fmt.Println("secrets:", secrets)

	fmt.Println("Result: ", result)

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
		count := len(share) / 88
		secrets[i] = make([][]*big.Int, count)

		for j := range secrets[i] {
			cshare := share[j*88 : (j+1)*88]
			secrets[i][j] = make([]*big.Int, 2)
			secrets[i][j][0] = fromBase64(cshare[0:44])
			secrets[i][j][1] = fromBase64(cshare[44:])
		}
	}

	fmt.Println("secrets", secrets)

	var secret []*big.Int = make([]*big.Int, len(secrets[0]))
	for j := range secret {
		secret[j] = big.NewInt(0)
		for i := range secrets {
			origin := secrets[i][j][0]
			originy := secrets[i][j][1]
			numerator := big.NewInt(1)
			denominator := big.NewInt(1)
			for k := range secrets {
				if k != i {
					current := secrets[k][j][0]
					fmt.Println("current: ", current)
					negative := big.NewInt(0)
					negative = negative.Mul(current, big.NewInt(-1))
					fmt.Println("negative: ", negative)
					added := big.NewInt(0)
					added = added.Sub(origin, current)
					fmt.Println("subtracted: ", added)

					fmt.Println("Numerator0: ", numerator)
					numerator = numerator.Mul(numerator, negative)
					fmt.Println("Numerator1: ", numerator)
					numerator = numerator.Mod(numerator, prime)
					fmt.Println("Numerator2: ", numerator)

					fmt.Println("denominator0: ", denominator)
					denominator = denominator.Mul(denominator, added)
					fmt.Println("denominator1: ", denominator)
					denominator = denominator.Mod(denominator, prime)
					fmt.Println("denominator2: ", denominator)
					fmt.Println("current: ", current)
					fmt.Println("prime: ", prime)
					fmt.Println("negative: ", negative)
					fmt.Println("added: ", added)
					fmt.Println("")
				}
			}

			fmt.Println("originy: ", originy)
			working := big.NewInt(0).Set(originy)
			fmt.Println("working: ", working)
			fmt.Println("Numerator: ", numerator)
			working = working.Mul(working, numerator)
			fmt.Println("working: ", working)
			fmt.Println("Denominator: ", denominator)
			fmt.Println("Inverse: ", modInverse(denominator))
			test := big.NewInt(0).Set(denominator)
			test = test.Mul(test, modInverse(denominator))
			test = test.Mod(test, prime)
			fmt.Println("test: ", test)
			working = working.Mul(working, modInverse(denominator))
			fmt.Println("working: ", working)
			fmt.Println("prime: ", prime)
			working = working.Add(working, prime)
			fmt.Println("working: ", working)

			fmt.Println("secret: ", secret[j])
			secret[j] = secret[j].Add(secret[j], working)
			fmt.Println("secret: ", secret[j])
			secret[j] = secret[j].Mod(secret[j], prime)
			fmt.Println("secret: ", secret[j])
		}
	}

	fmt.Println("Secret: ", secret)

	return "Hi!"
}
