package shortener

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"

	"github.com/itchyny/base58-go"
)

func sha2560f(input string) []byte{
	algorthim := sha256.New()
	algorthim.Write([]byte(input))
	return algorthim.Sum(nil)
}

func base58Encoding(bytes []byte) string{
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		log.Fatalf("error happend in encoding: %v", err)
	}
	return string(encoded)
}


/* The final algorithm will be super straightforward now as we now have our 2
main building blocks already setup, it will go as follow :

- Hashing  initialUrl + userId url with sha256.  Here userId is added to 
	prevent providing similar shortened urls to separate users in case they 
	want to shorten exact same link, it's a design decision, so some implementations 
	do this differently.
- Derive a big integer number from the hash bytes generated during the hasing.
- Finally apply base58  on the derived big integer value and pick the first 8 characters
*/
func GenertateShortLink(initalLink string, userId string) string {
	urlHashBytes := sha2560f(initalLink+userId)
	genreteNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString := base58Encoding([]byte(fmt.Sprintf("%d", genreteNumber)))
	return finalString[:8]
}


