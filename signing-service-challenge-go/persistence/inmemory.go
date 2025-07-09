package persistence

import (
	"encoding/base64"
	"fmt"
	"sync"

	crypto "github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// Variables defined at package level. This is smt that I consider bad practice,
// but for the purpose of this test it's ok. For more complex and real-case scenarios
// I would use a db and I would prefer to handle concurrency at DB level rather
// than here
var (
	// mutex used to manage concurrency. Defining it at package level
	mutex sync.Mutex
	// in-memory storage of the devices. It leverages the uniqueness
	// of the UUID to have a fast-to-access map structure
	devices map[string]*domain.Device = make(map[string]*domain.Device, 0)
)

// return all the devices. In a real case this function should recover
// them from a db
func GetDevices() ([]domain.Device, error) {
	mutex.Lock()
	defer mutex.Unlock()

	// retrieve the devices from the in-memory map
	var devs []domain.Device = make([]domain.Device, 0)
	for _, dev := range devices {
		devs = append(devs, *dev)
	}

	if len(devs) != 1 {
		log.Debugf("GetDevices found %d devices", len(devs))
	} else {
		log.Debugf("GetDevices found %d device", len(devs))
	}

	return devs, nil
}

// return device by UUID. In a real case this function should recover
// it from a db
func GetDevice(uuid string) (domain.Device, error) {
	mutex.Lock()
	defer mutex.Unlock()

	// retrieve the devices from the in-memory map
	if dev, found := devices[uuid]; !found {
		return domain.Device{}, fmt.Errorf("device with UUID '%s' not found", uuid)
	} else {
		log.Debugf("GetDevice found device with UUID %s", dev.UUID)
		return *dev, nil
	}
}

// add a device to the map. In a real case this function should add
// the device to the db
func AddDevice(signingAlgorithm crypto.SigningAlgorithm, label *string) (domain.Device, error) {
	mutex.Lock()
	defer mutex.Unlock()

	// generate UUID
	uuid := uuid.New().String()

	// make sure uuid is unique (maybe unecessary here?)
	for {
		if _, found := devices[uuid]; !found {
			break
		}
	}

	var signer crypto.SignerI
	var err error

	// switch between implemented algorithm
	switch signingAlgorithm {
	case crypto.ECC:
		signer, err = crypto.NewECCSigner()
		if err != nil {
			return domain.Device{}, fmt.Errorf("could not create a new ECC Signer: %s", err)
		}
	case crypto.RSA:
		signer, err = crypto.NewRSASigner()
		if err != nil {
			return domain.Device{}, fmt.Errorf("could not create a new RSA Signer: %s", err)
		}
	default:
		// verification on the payload of the REST API should prevent to reach here
		return domain.Device{}, fmt.Errorf("selected algorithm '%s' not yet implemented", signingAlgorithm.String())
	}

	// create device
	dev := domain.Device{
		UUID:                             string(uuid),
		LastSignatureBase64EncodedString: base64.StdEncoding.EncodeToString([]byte(string(uuid))),
		SignatureCounter:                 0,
		Signer:                           signer,
	}

	// use UUID as label if it has not been provided
	if label == nil {
		dev.Label = dev.UUID
	} else {
		dev.Label = *label
	}

	// add it to the map
	devices[uuid] = &dev

	log.Debugf("AddDevice succesfully added the device with UUID %s", dev.UUID)

	return dev, nil
}

// sign the data and returned encoded signature with signed data
func SignTransaction(uuid string, dataIn string) (string, string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	// retrieve device
	dev, found := devices[uuid]
	if !found {
		return "", "", fmt.Errorf("device with UUID '%s' not found", uuid)
	}

	// compose the data to be signed
	dataToBeSigned := fmt.Sprint(dev.SignatureCounter) + "_" + dataIn + "_" + dev.LastSignatureBase64EncodedString

	// this is where it is key point to provide easy support in implementing
	// other cryptographic algorithms.
	signature, err := dev.Signer.Sign([]byte(dataToBeSigned))
	if err != nil {
		return "", "", fmt.Errorf("error while signing the data: %s", err)
	}

	// encode to string the signature
	encodedSignature := base64.StdEncoding.EncodeToString(signature)

	// store the string for of encoded signature in the dev for next signing process
	dev.LastSignatureBase64EncodedString = encodedSignature

	// if everything is successful, increment the counter
	dev.SignatureCounter++

	log.Debugf("SignTransaction succesfully signed the transaction with the device UUID %s", dev.UUID)

	return encodedSignature, dataToBeSigned, nil
}

// This is why I don't like to have stuff defined at pkg level. I need
// this function to clean the db between each test, otherwise the test
// results influence each other
func CleanMemory() {
	devices = make(map[string]*domain.Device)
}

type DeviceNotFoundError error
