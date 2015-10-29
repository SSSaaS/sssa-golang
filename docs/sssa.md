# SSSA - Ruby
## Constants
    prime = 99995644905598542077721161034987774965417302630805822064337798850767846245779
        Safe Prime; big.Int

## Functions
    create(minimum int, number int, raw string)
        minimum - number of shares required to recreate the secret
        number - total number of shares to generate
        raw - secret to protect, as a string

        returns shares as an array of base64 strings of variable length
            dependent on the size of the raw secret

    combine(shares []string)
        shares - array of base64 strings returned by create function

        returns a string of secret
            note: this string can be ill-formatted utf8 potentially, if the
            minimum number of shares was not met
