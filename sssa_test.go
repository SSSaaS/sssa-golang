package sssa

import (
	"bytes"
	"math/big"
	"testing"
)

func TestCreateCombine(t *testing.T) {
	g, err := NewDefaultSSSAGenerator("")
	if err != nil {
		t.Fatal("Cannot NewDefaultSSSAGenerator: %s", err)
	}

	// Short, medium, and long tests
	strings := []string{
		"N17FigASkL6p1EOgJhRaIquQLGvYV0",
		"0y10VAfmyH7GLQY6QccCSLKJi8iFgpcSBTLyYOGbiYPqOpStAf1OYuzEBzZR",
		"KjRHO1nHmIDidf6fKvsiXWcTqNYo2U9U8juO94EHXVqgearRISTQe0zAjkeUYYBvtcB8VWzZHYm6ktMlhOXXCfRFhbJzBUsXaHb5UDQAvs2GKy6yq0mnp8gCj98ksDlUultqygybYyHvjqR7D7EAWIKPKUVz4of8OzSjZlYg7YtCUMYhwQDryESiYabFID1PKBfKn5WSGgJBIsDw5g2HB2AqC1r3K8GboDN616Swo6qjvSFbseeETCYDB3ikS7uiK67ErIULNqVjf7IKoOaooEhQACmZ5HdWpr34tstg18rO",
	}

	minimum := []int{4, 6, 20}
	shares := []int{5, 100, 100}

	for i := range strings {
		created, err := g.Create(minimum[i], shares[i], strings[i])
		if err != nil {
			t.Fatal("Fatal: creating: ", err)
		}
		combined, err := g.Combine(created)
		if err != nil {
			t.Fatal("Fatal: combining: ", err)
		}
		if combined != strings[i] {
			t.Fatal("Fatal: combining returned invalid data")
		}
	}
}

func TestLibraryCombine(t *testing.T) {
	g, err := NewDefaultSSSAGenerator("")
	if err != nil {
		t.Fatal("Cannot NewDefaultSSSAGenerator: %s", err)
	}

	shares := []string{
		"U1k9koNN67-og3ZY3Mmikeyj4gEFwK4HXDSglM8i_xc=yA3eU4_XYcJP0ijD63Tvqu1gklhBV32tu8cHPZXP-bk=",
		"O7c_iMBaGmQQE_uU0XRCPQwhfLBdlc6jseTzK_qN-1s=ICDGdloemG50X5GxteWWVZD3EGuxXST4UfZcek_teng=",
		"8qzYpjk7lmB7cRkOl6-7srVTKNYHuqUO2WO31Y0j1Tw=-g6srNoWkZTBqrKA2cMCA-6jxZiZv25rvbrCUWVHb5g=",
		"wGXxa_7FPFSVqdo26VKdgFxqVVWXNfwSDQyFmCh2e5w=8bTrIEs0e5FeiaXcIBaGwtGFxeyNtCG4R883tS3MsZ0=",
		"j8-Y4_7CJvL8aHxc8WMMhP_K2TEsOkxIHb7hBcwIBOo=T5-EOvAlzGMogdPawv3oK88rrygYFza3KSki2q8WEgs=",
	}

	combined, err := g.Combine(shares)
	if err != nil {
		t.Fatal("Fatal: combining: ", err)
	}
	if combined != "test-pass" {
		t.Fatal("Failed library cross-language check")
	}
}

func TestIsValidShare(t *testing.T) {
	g, err := NewDefaultSSSAGenerator("")
	if err != nil {
		t.Fatal("Cannot NewDefaultSSSAGenerator: %s", err)
	}

	shares := []string{
		"U1k9koNN67-og3ZY3Mmikeyj4gEFwK4HXDSglM8i_xc=yA3eU4_XYcJP0ijD63Tvqu1gklhBV32tu8cHPZXP-bk=",
		"O7c_iMBaGmQQE_uU0XRCPQwhfLBdlc6jseTzK_qN-1s=ICDGdloemG50X5GxteWWVZD3EGuxXST4UfZcek_teng=",
		"8qzYpjk7lmB7cRkOl6-7srVTKNYHuqUO2WO31Y0j1Tw=-g6srNoWkZTBqrKA2cMCA-6jxZiZv25rvbrCUWVHb5g=",
		"wGXxa_7FPFSVqdo26VKdgFxqVVWXNfwSDQyFmCh2e5w=8bTrIEs0e5FeiaXcIBaGwtGFxeyNtCG4R883tS3MsZ0=",
		"j8-Y4_7CJvL8aHxc8WMMhP_K2TEsOkxIHb7hBcwIBOo=T5-EOvAlzGMogdPawv3oK88rrygYFza3KSki2q8WEgs=",
		"Hello world",
	}

	results := []bool{
		true,
		true,
		true,
		true,
		true,
		false,
	}

	for i := range shares {
		if g.IsValidShare(shares[i]) != results[i] {
			t.Fatal("Checking for valid shares failed:", i)
		}
	}
}

func TestRandom(t *testing.T) {
	prime, _ := big.NewInt(0).SetString(DefaultPrimeStr, 10)
	for i := 0; i < 10000; i++ {
		if random(prime).Cmp(prime) >= 0 {
			t.Fatal("Error! Random number out of bounds exception")
		}
	}
}

func TestBaseConversion(t *testing.T) {
	prime, _ := big.NewInt(0).SetString(DefaultPrimeStr, 10)
	for i := 0; i < 10000; i++ {
		point := random(prime)
		if point.Cmp(fromBase64(toBase64(point))) != 0 {
			t.Fatal("Fatal: Base conversion failed")
		}
	}
}

func TestToBase64(t *testing.T) {
	prime, _ := big.NewInt(0).SetString(DefaultPrimeStr, 10)
	for i := 0; i < 10000; i++ {
		point := random(prime)
		if len(toBase64(point)) != 44 {
			t.Fatal("Fatal: toBase64 returned wrong length")
		}
	}
}

func TestSplitMerge(t *testing.T) {
	// Short, medium, and long tests
	tests := [][]byte{
		[]byte("N17FigASkL6p1EOgJhRaIquQLGvYV0"),
		[]byte("0y10VAfmyH7GLQY6QccCSLKJi8iFgpcSBTLyYOGbiYPqOpStAf1OYuzEBzZR"),
		[]byte("KjRHO1nHmIDidf6fKvsiXWcTqNYo2U9U8juO94EHXVqgearRISTQe0zAjkeUYYBvtcB8VWzZHYm6ktMlhOXXCfRFhbJzBUsXaHb5UDQAvs2GKy6yq0mnp8gCj98ksDlUultqygybYyHvjqR7D7EAWIKPKUVz4of8OzSjZlYg7YtCUMYhwQDryESiYabFID1PKBfKn5WSGgJBIsDw5g2HB2AqC1r3K8GboDN616Swo6qjvSFbseeETCYDB3ikS7uiK67ErIULNqVjf7IKoOaooEhQACmZ5HdWpr34tstg18rO"),
	}

	for i := range tests {
		result := mergeIntToByte(splitByteToInt(tests[i]))
		if !bytes.Equal(result, tests[i]) {
			t.Fatal("Fatal: splitting and merging returned invalid data: ", result, tests[i])
		}
	}
}

func TestSplitMergeOdds(t *testing.T) {
	// Short, medium, and long tests
	tests := [][]byte{
		[]byte("a\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00a"),
		[]byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa哈囉世界"),
		[]byte("こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界こんにちは、世界"),
	}

	for i := range tests {
		result := mergeIntToByte(splitByteToInt(tests[i]))
		if !bytes.Equal(result, tests[i]) {
			t.Fatal("Fatal: splitting and merging returned invalid data on test", i)
		}
	}
}

func TestModInverse(t *testing.T) {
	prime, _ := big.NewInt(0).SetString(DefaultPrimeStr, 10)
	for i := 0; i < 10000; i++ {
		point := big.NewInt(0).Set(random(prime))
		if point.Cmp(prime) >= 0 {
			t.Fatal("Error! Random point out of bounds exception")
		}
		inverse := big.NewInt(0).Set(modInverse(point, prime))
		if inverse.Cmp(prime) >= 0 {
			t.Fatal("Error! Inverse out of bounds exception")
		}
		value := big.NewInt(0).Set(point)
		value = value.Mul(value, inverse)
		value = value.Mod(value, prime)
		expected := big.NewInt(1)
		if value.Cmp(expected) != 0 {
			t.Fatalf("Fatal: modInverse[%v] failed\nExpected: %v; Got: %v\nPoint %v\nInverse: %v\nPrime: %v", i, expected, value, point, inverse, prime)
		}
	}
}

func TestEvaluatePolynomial(t *testing.T) {
	prime, _ := big.NewInt(0).SetString(DefaultPrimeStr, 10)

	values := [][][]*big.Int{
		[][]*big.Int{
			[]*big.Int{big.NewInt(20), big.NewInt(21), big.NewInt(42)},
			[]*big.Int{big.NewInt(0)},
		},
		[][]*big.Int{
			[]*big.Int{big.NewInt(0), big.NewInt(0), big.NewInt(0)},
			[]*big.Int{big.NewInt(4)},
		},
		[][]*big.Int{
			[]*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3), big.NewInt(4), big.NewInt(5)},
			[]*big.Int{big.NewInt(10)},
		},
	}

	actual := []*big.Int{big.NewInt(20), big.NewInt(0), big.NewInt(54321)}

	for i := range values {
		result := evaluatePolynomial(values[i][0], values[i][1][0], prime)
		if result.Cmp(actual[i]) != 0 {
			t.Fatalf("Fatal: EvaluatePolynomial[%v] failed\nExpected: %v; Got: %v\n", i, actual[i], result)
		}
	}
}
