package crypto

import "github.com/fiskaly/coding-challenges/signing-service-challenge/domain"

var supportedGenerators = []domain.KeyPairGenerator{
	ECCGenerator{},
	RSAGenerator{},
}

func FindKeyPairGenerator(algorithmName string) (domain.KeyPairGenerator, bool) {
	for _, generator := range supportedGenerators {
		if generator.AlgorithmName() == algorithmName {
			return generator, true
		}
	}

	return nil, false
}
