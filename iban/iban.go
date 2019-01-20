package iban

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

var countries = map[string]int{
	"UA": 29,
}

// Check make checks if IBAN code is valid
func Check(code string) error {
	pureCode := strings.ToUpper(strings.Replace(code, " ", "", -1))
	if len(pureCode) < 2 {
		return errors.New("You provided no code")
	}
	countryCode := pureCode[:2]

	//fmt.Printf("Recieved: %v\nCountry: %v\n", pureCode, countryCode)
	if _, ok := countries[countryCode]; !ok {
		err := errors.New("This country can't be validated or country code is invalid")
		return err
	}

	if len(pureCode) != countries[countryCode] {
		err := errors.New("IBAN length is invalid")
		return err
	}

	first, second := int(countryCode[0])-55, int(countryCode[1])-55
	// fmt.Printf("%s = %v;%s = %v\n", string(countryCode[0]), first, string(countryCode[1]), second)

	pureCode = fmt.Sprintf("%v%v%v%v", pureCode[4:], first, second, pureCode[2:4])
	// fmt.Println(pureCode)

	myInt, ok := new(big.Int).SetString(pureCode, 10)
	if !ok {
		return errors.New("Error converting string to int")
	}
	result := new(big.Int).Mod(myInt, big.NewInt(97))
	if result.Int64() != 1 {
		return errors.New("IBAN is incorrect. Please, check if you entered a correct IBAN code")
	}

	return nil
}

func main() {
	err := Check("ua 21 399622 0000026007233566001")
	if err != nil {
		fmt.Println(err)
	}

	nErr := Check("UA 21 399622 0000026007233566001")
	if nErr != nil {
		fmt.Println(nErr)
	}
}
