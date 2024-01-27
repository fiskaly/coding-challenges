package domain

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func SignTransaction(
	device SignatureDevice,
	deviceRepository SignatureDeviceRepository,
	dataToBeSigned string,
) (
	base64EncodedSignature string,
	signedData string,
	err error,
) {
	securedDataToBeSigned := SecureDataToBeSigned(device, dataToBeSigned)

	signature, err := device.SignTransaction(securedDataToBeSigned)
	if err != nil {
		return "", "", errors.New(fmt.Sprintf("failed to sign transaction: %s", err))
	}
	encodedSignature := base64.StdEncoding.EncodeToString(signature)

	device.Base64EncodedLastSignature = encodedSignature
	device.SignatureCounter++
	err = deviceRepository.Update(device)
	if err != nil {
		return "", "", errors.New(fmt.Sprintf("failed to update signature device: %s", err))
	}

	return encodedSignature, securedDataToBeSigned, nil
}

func SecureDataToBeSigned(device SignatureDevice, data string) string {
	components := []string{
		strconv.Itoa(int(device.SignatureCounter)),
		data,
	}

	if device.SignatureCounter == 0 {
		// when the device has not yet been used, the `lastSignature` is blank,
		// so use the device ID instead
		encodedID := base64.StdEncoding.EncodeToString([]byte(device.ID.String()))
		components = append(components, encodedID)
	} else {
		encodedLastSignature := base64.StdEncoding.EncodeToString([]byte(device.Base64EncodedLastSignature))
		components = append(components, encodedLastSignature)
	}

	return strings.Join(components, "_")
}
