# SSSA - Golang
## Package
    sssa

## Constants
    prime = 115792089237316195423570985008687907853269984665640564039457584007913129639747
        Safe Prime; big.Int; not exported

## Functions
    Create(minimum int, number int, raw string)
        minimum - number of shares required to recreate the secret
        number - total number of shares to generate
        raw - secret to protect, as a string

        returns shares as an array of base64 strings of variable length
            dependent on the size of the raw secret

    Combine(shares []string)
        shares - array of base64 strings returned by create function

        returns a string of secret
            note: this string can be ill-formatted utf8 potentially, if the
            minimum number of shares was not met

    IsValidShare(candidate string)
        candidate - candidate to check whether or not a valid share

        returns a boolean of validity
