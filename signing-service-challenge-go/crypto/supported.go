package crypto

import "github.com/fiskaly/coding-challenges/signing-service-challenge/domain"

var SupportedAlgorithms = []domain.SignatureAlgorithm{
	ECCAlgorithm{},
	RSAAlgorithm{},
}

func FindSupportedAlgorithm(name string) (domain.SignatureAlgorithm, bool) {
	for _, algorithm := range SupportedAlgorithms {
		if algorithm.Name() == name {
			return algorithm, true
		}
	}

	return nil, false
}
