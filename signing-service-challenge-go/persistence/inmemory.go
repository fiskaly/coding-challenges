package persistence

import (
	"bytes"
	"fmt"
	"slices"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

var lock sync.Mutex
var buffer bytes.Buffer
var signatureDevices = []domain.SignatureDevice{}

func AddSignatureDevice(device *domain.SignatureDevice) error {
	if slices.ContainsFunc(signatureDevices, func(d domain.SignatureDevice) bool {
		return d.Id == device.Id
	}) {
		return fmt.Errorf("device with id %s already exists", device.Id)
	}

	return nil
}
