package sssa

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
)

const (
	DefaultPrimeStr = "115792089237316195423570985008687907853269984665640564039457584007913129639747"
)

var (
	ErrCannotParsePrime        = errors.New("Cannot parse prime")
	ErrCannotRequireMoreShares = errors.New("Cannot require more shares then existing")
	ErrOneOfTheSharesIsInvalid = errors.New("one of the shares is invalid")
)

type SSSAGenerator interface {
	Create(minimum int, shares int, raw string) ([]string, error)
	Combine(shares []string) (string, error)
	IsValidShare(candidate string) bool
}

type DefaultSSSAGenerator struct {
	Prime *big.Int
}

func NewDefaultSSSAGenerator(prime string) (*DefaultSSSAGenerator, error) {
	if prime == "" {
		prime = DefaultPrimeStr
	}

	p, ok := big.NewInt(0).SetString(prime, 10)
	if !ok {
		return nil, ErrCannotParsePrime
	}

	return &DefaultSSSAGenerator{Prime: p}, nil
}

/**
 * Returns a new arary of secret shares (encoding x,y pairs as base64 strings)
 * created by Shamir's Secret Sharing Algorithm requring a minimum number of
 * share to recreate, of length shares, from the input secret raw as a string
**/
func (g *DefaultSSSAGenerator) Create(minimum int, shares int, raw string) ([]string, error) {
	// Verify minimum isn't greater than shares; there is no way to recreate
	// the original polynomial in our current setup, therefore it doesn't make
	// sense to generate fewer shares than are needed to reconstruct the secret.
	if minimum > shares {
		return nil, ErrCannotRequireMoreShares
	}

	// Convert the secret to its respective 256-bit big.Int representation
	var secret []*big.Int = splitByteToInt([]byte(raw))

	// List of currently used numbers in the polynomial
	var numbers []*big.Int = make([]*big.Int, 0)
	numbers = append(numbers, big.NewInt(0))

	// Create the polynomial of degree (minimum - 1); that is, the highest
	// order term is (minimum-1), though as there is a constant term with
	// order 0, there are (minimum) number of coefficients.
	//
	// However, the polynomial object is a 2d array, because we are constructing
	// a different polynomial for each part of the secret
	// polynomial[parts][minimum]
	var polynomial [][]*big.Int = make([][]*big.Int, len(secret))
	for i := range polynomial {
		polynomial[i] = make([]*big.Int, minimum)
		polynomial[i][0] = secret[i]

		for j := range polynomial[i][1:] {
			// Each coefficient should be unique
			number := random(g.Prime)
			for inNumbers(numbers, number) {
				number = random(g.Prime)
			}
			numbers = append(numbers, number)

			polynomial[i][j+1] = number
		}
	}

	// Create the secrets object; this holds the (x, y) points of each share.
	// Again, because secret is an array, each share could have multiple parts
	// over which we are computing Shamir's Algorithm. The last dimension is
	// always two, as it is storing an x, y pair of points.
	//
	// Note: this array is technically unnecessary due to creating result
	// in the inner loop. Can disappear later if desired. [TODO]
	//
	// secrets[shares][parts][2]
	var secrets [][][]*big.Int = make([][][]*big.Int, shares)
	var result []string = make([]string, shares)

	// For every share...
	for i := range secrets {
		secrets[i] = make([][]*big.Int, len(secret))
		// ...and every part of the secret...
		for j := range secrets[i] {
			secrets[i][j] = make([]*big.Int, 2)

			// ...generate a new x-coordinate...
			number := random(g.Prime)
			for inNumbers(numbers, number) {
				number = random(g.Prime)
			}
			numbers = append(numbers, number)

			// ...and evaluate the polynomial at that point...
			secrets[i][j][0] = number
			secrets[i][j][1] = evaluatePolynomial(
				polynomial[j],
				number,
				g.Prime,
			)

			// ...add it to results...
			result[i] += toBase64(secrets[i][j][0])
			result[i] += toBase64(secrets[i][j][1])
		}
	}

	// ...and return!
	return result, nil
}

/**
 * Takes a string array of shares encoded in base64 created via Shamir's
 * Algorithm; each string must be of equal length of a multiple of 88 characters
 * as a single 88 character share is a pair of 256-bit numbers (x, y).
 *
 * Note: the polynomial will converge if the specified minimum number of shares
 *       or more are passed to this function. Passing thus does not affect it
 *       Passing fewer however, simply means that the returned secret is wrong.
**/
func (g *DefaultSSSAGenerator) Combine(shares []string) (string, error) {
	// Recreate the original object of x, y points, based upon number of shares
	// and size of each share (number of parts in the secret).
	var secrets [][][]*big.Int = make([][][]*big.Int, len(shares))

	// For each share...
	for i := range shares {
		// ...ensure that it is valid...
		if g.IsValidShare(shares[i]) == false {
			return "", ErrOneOfTheSharesIsInvalid
		}

		// ...find the number of parts it represents...
		share := shares[i]
		count := len(share) / 88
		secrets[i] = make([][]*big.Int, count)

		// ...and for each part, find the x,y pair...
		for j := range secrets[i] {
			cshare := share[j*88 : (j+1)*88]
			secrets[i][j] = make([]*big.Int, 2)
			// ...decoding from base 64.
			secrets[i][j][0] = fromBase64(cshare[0:44])
			secrets[i][j][1] = fromBase64(cshare[44:])
		}
	}

	// Use Lagrange Polynomial Interpolation (LPI) to reconstruct the secret.
	// For each part of the secert (clearest to iterate over)...
	var secret []*big.Int = make([]*big.Int, len(secrets[0]))
	for j := range secret {
		secret[j] = big.NewInt(0)
		// ...and every share...
		for i := range secrets { // LPI sum loop
			// ...remember the current x and y values...
			origin := secrets[i][j][0]
			originy := secrets[i][j][1]
			numerator := big.NewInt(1)   // LPI numerator
			denominator := big.NewInt(1) // LPI denominator
			// ...and for every other point...
			for k := range secrets { // LPI product loop
				if k != i {
					// ...combine them via half products...
					current := secrets[k][j][0]
					negative := big.NewInt(0)
					negative = negative.Mul(current, big.NewInt(-1))
					added := big.NewInt(0)
					added = added.Sub(origin, current)

					numerator = numerator.Mul(numerator, negative)
					numerator = numerator.Mod(numerator, g.Prime)

					denominator = denominator.Mul(denominator, added)
					denominator = denominator.Mod(denominator, g.Prime)
				}
			}

			// LPI product
			// ...multiply together the points (y)(numerator)(denominator)^-1...
			working := big.NewInt(0).Set(originy)
			working = working.Mul(working, numerator)
			working = working.Mul(working, modInverse(denominator, g.Prime))

			// LPI sum
			secret[j] = secret[j].Add(secret[j], working)
			secret[j] = secret[j].Mod(secret[j], g.Prime)
		}
	}

	// ...and return the result!
	return string(mergeIntToByte(secret)), nil
}

/**
 * Takes in a given string to check if it is a valid secret
 *
 * Requirements:
 * 	Length multiple of 88
 *	Can decode each 44 character block as base64
 *
 * Returns only success/failure (bool)
**/
func (g *DefaultSSSAGenerator) IsValidShare(candidate string) bool {
	if len(candidate)%88 != 0 {
		return false
	}

	count := len(candidate) / 44
	for j := 0; j < count; j++ {
		part := candidate[j*44 : (j+1)*44]
		decode := fromBase64(part)
		if decode.Cmp(big.NewInt(0)) == -1 || decode.Cmp(g.Prime) == 1 {
			return false
		}
	}

	return true
}

/**
 * Returns a random number from the range (0, prime-1) inclusive
**/
func random(prime *big.Int) *big.Int {
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
	hex_data := hex.EncodeToString(secret)
	count := int(math.Ceil(float64(len(hex_data)) / 64.0))

	result := make([]*big.Int, count)

	for i := 0; i < count; i++ {
		if (i+1)*64 < len(hex_data) {
			result[i], _ = big.NewInt(0).SetString(hex_data[i*64:(i+1)*64], 16)
		} else {
			data := strings.Join([]string{hex_data[i*64:], strings.Repeat("0", 64-(len(hex_data)-i*64))}, "")
			result[i], _ = big.NewInt(0).SetString(data, 16)
		}
	}

	return result
}

/**
 * Converts an array of big.Ints to the original byte array, removing any
 * least significant nulls
**/
func mergeIntToByte(secret []*big.Int) []byte {
	var hex_data = ""
	for i := range secret {
		tmp := fmt.Sprintf("%x", secret[i])
		hex_data += strings.Join([]string{strings.Repeat("0", (64 - len(tmp))), tmp}, "")
	}

	result, _ := hex.DecodeString(hex_data)
	result = bytes.TrimRight(result, "\x00")

	return result
}

/**
 * Evauluates a polynomial with coefficients specified in reverse order:
 * evaluatePolynomial([a, b, c, d], x, prime):
 * 		returns a + bx + cx^2 + dx^3
**/
func evaluatePolynomial(polynomial []*big.Int, value *big.Int, prime *big.Int) *big.Int {
	last := len(polynomial) - 1
	var result *big.Int = big.NewInt(0).Set(polynomial[last])

	for s := last - 1; s >= 0; s-- {
		result = result.Mul(result, value)
		result = result.Add(result, polynomial[s])
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
 * bits long, and the output is an arbitrary size base 10 integer.
 *
 * Returns -1 on failure
**/
func fromBase64(number string) *big.Int {
	bytedata, err := base64.URLEncoding.DecodeString(number)
	if err != nil {
		return big.NewInt(-1)
	}

	hexdata := hex.EncodeToString(bytedata)
	result, ok := big.NewInt(0).SetString(hexdata, 16)
	if ok == false {
		return big.NewInt(-1)
	}

	return result
}

/**
 * Computes the multiplicative inverse of the number on the field prime; more
 * specifically, number * inverse == 1; Note: number should never be zero
**/
func modInverse(number *big.Int, prime *big.Int) *big.Int {
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
