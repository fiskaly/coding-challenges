package crypto

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Here I though of enumerating the implemented cryptographic algorithms.
// For future implementations of other algorithms, the enumeration can grow
type SigningAlgorithm uint

const (
	ECC SigningAlgorithm = iota + 1
	RSA
)

// custom SigningAlgorithm requires this Unmarshalling function to adhere
// to the unmarshalling interface
func (sa *SigningAlgorithm) UnmarshalJSON(bytes []byte) error {
	var str string
	var err error
	if err := json.Unmarshal(bytes, &str); err != nil {
		err := fmt.Errorf("could not unmarhal SigningAlgorithm: %s", err)
		log.Error(err.Error())
		return err
	}
	if *sa, err = stringToSA(str); err != nil {
		err := fmt.Errorf("could not convert string to SigningAlgorithm: %s", err)
		log.Error(err.Error())
		return err
	}
	return nil
}

// function returning a signing method matching the passed-in string
func stringToSA(str string) (SigningAlgorithm, error) {
	switch str {
	case "ECC":
		return ECC, nil
	case "RSA":
		return RSA, nil
	default:
		err := fmt.Errorf("passed string is not a valid SigningAlgorithm type")
		log.Error(err.Error())
		return SigningAlgorithm(0), err
	}
}

// custom SigningAlgorithm requires this Marshalling function to adhere
// to the marshalling interface
func (sa *SigningAlgorithm) MarshalJSON() ([]byte, error) {
	return json.Marshal(sa.String())
}

// function returning the string form of the enumerated signing method
func (sa SigningAlgorithm) String() string {
	switch sa {
	case ECC:
		return "ECC"
	case RSA:
		return "RSA"
	default:
		log.Errorf("the SigningAlgorithm '%d' is not a valid one", sa)
		return ""
	}
}
