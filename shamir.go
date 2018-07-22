package shamir

import (
	"errors"
	"math/big"
)

var (
	ErrCannotRequireMoreShares = errors.New("cannot require more shares then existing")
	ErrOneOfTheSharesIsInvalid = errors.New("one of the shares is invalid")
)

const (
	DefaultP = "82FD4FA0549ABF545984EB566FDE65F015FF1FF0FB48EDDB419C5906D54CD8591AF8CB0292BE1F43E212690808EC0F76ABA7E07895CB7855316DB5037F51F746B52D09F9B77921E06623615926DE8ED3D9B6AB81B3F468551505E9FAFAA802328A6C5785E341440CDA0F9BC1D27DB2D5A613DC56BD27C36531589CF3300F4CA2DC259F7AABDF0C29AE98B95ACA4CA0153B840A9D2B658920E7BBAF7C96109C93D41FEBF0AE024FC63BB6DA4FD032D760B294FAFD3B2862520BB117A56CE6650D3C45151D0C9EF8B64E54224238AF1911756ABD8F28ABE140FE2E5BFFF94C9D9F22BC0E22841A79948C3CA5D4BFD09315A41049068FC1CFEB4A662D2DEE19DE7B"
)

var P, _ = new(big.Int).SetString(DefaultP, 16)


type Share struct {
	X *big.Int
	Y *big.Int
}


func CreateShares(minimum int, shares int, raw []byte) ([]*Share, error) {
	if minimum > shares {
		return nil, ErrCannotRequireMoreShares
	}

	// Convert the secret to its big.Int representation
	secret := new(big.Int).SetBytes(raw)

	if secret.Cmp(P) != -1 {
		return nil, errors.New("secret too large for encoding")
	}

	// List of currently used numbers in the polynomial
	var numbers []*big.Int = make([]*big.Int, 0)
	numbers = append(numbers, big.NewInt(0))

	// Create the polynomial of degree (minimum - 1); that is, the highest
	// order term is (minimum-1), though as there is a constant term with
	// order 0, there are (minimum) number of coefficients.

	polynomial := make([]*big.Int, minimum)
	polynomial[0] = secret
	for j := range polynomial[1:] {
		// Each coefficient should be unique
		number := random()
		for inNumbers(numbers, number) {
			number = random()
		}
		numbers = append(numbers, number)

		polynomial[j+1] = number
	}

	// Create the secrets object; this holds the (x, y) points of each share.
	// Again, because secret is an array, each share could have multiple parts
	// over which we are computing Shamir's Algorithm.

	result := make([]*Share, shares)

	// For every share...
	for i := range result {
		// new x-coordinate...
		number := big.NewInt(int64(i+1))

		// evaluate the polynomial at that point...
		result[i] = &Share{
			X: number,
			Y: evaluatePolynomial(polynomial, number),
		}

	}

	// ...and return!
	return result, nil
}


func CombineShares(shares []*Share, minimum int) (*big.Int, error) {

	if len(shares) < minimum {
		return nil, errors.New("not enough shares for interoplation")
	}
	// Recreate the original object of x, y points, based upon number of shares

	// Use Lagrange Polynomial Interpolation (LPI) to reconstruct the secret.
	secret := big.NewInt(0)
	for i, s := range shares { 
		origin := s.X
		originy := s.Y
		numerator := big.NewInt(1)  
		denominator := big.NewInt(1) 
		for k := range shares {
			if k != i {
				current := shares[k].X
				negative := big.NewInt(0)
				negative = negative.Mul(current, big.NewInt(-1))
				added := big.NewInt(0)
				added = added.Sub(origin, current)

				numerator = numerator.Mul(numerator, negative)
				numerator = numerator.Mod(numerator, P)

				denominator = denominator.Mul(denominator, added)
				denominator = denominator.Mod(denominator, P)
			}
		}

		working := big.NewInt(0).Set(originy)
		working = working.Mul(working, numerator)
		working = working.Mul(working, modInverse(denominator))

		secret.Add(secret, working)
		secret.Mod(secret, P)
	}

	return secret, nil
}
