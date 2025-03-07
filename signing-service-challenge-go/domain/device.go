package domain

import (
	"encoding/base64"
	"fmt"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

func MakeSignatureDomain(dvp persistence.DevicePersister) SignatureDomainImpl {
	return SignatureDomainImpl{
		dvp: dvp,
	}
}

var PERSISTENCE_RETRY_LIMIT = 3

type SignatureDomainImpl struct {
	dvp persistence.DevicePersister
}

func (s *SignatureDomainImpl) CreateSignatureDevice(d Device) error {
	_, sec, err := crypto.GenerateAndEncode(d.Algorithm)
	if err == crypto.ErrUnsupportedAlgorithm {
		return crypto.ErrUnsupportedAlgorithm
	}
	if err != nil {
		return fmt.Errorf("failed to gererate keys: %w", err)
	}

	err = s.dvp.AddDevice(persistence.Device{
		ID:             persistence.ID(d.ID),
		Algorithm:      persistence.Algorithm(d.Algorithm),
		Label:          &d.Label,
		LastSignature:  persistence.Signature(base64.StdEncoding.EncodeToString([]byte(d.ID))),
		SignatureCount: 0,
		PrivateKey:     sec,
	})

	if err == persistence.ErrAlreadyExists {
		return err
	}
	if err != nil {
		return fmt.Errorf("failed to create signature device: %w", err)
	}

	return nil
}

func (s *SignatureDomainImpl) ListSignatureDevices() ([]ID, error) {
	ids, err := s.dvp.ListDevices()
	if err != nil {
		return nil, fmt.Errorf("failed to list signature devices: %w", err)
	}

	domainIDs := make([]ID, len(ids))
	for i, id := range ids {
		domainIDs[i] = ID(id)
	}

	return domainIDs, nil
}

func (s *SignatureDomainImpl) SignTransaction(id ID, data Data) (*CreatedSignature, error) {
	for range PERSISTENCE_RETRY_LIMIT {
		d, err := s.dvp.GetDevice(persistence.ID(id))
		if err != nil {
			return nil, fmt.Errorf("failed to get signature device details: %w", err)
		}
		r := &CreatedSignature{
			SignedData: fmt.Sprint(d.SignatureCount) + "_" + string(data) + "_" + string(d.LastSignature),
		}
		signer, err := crypto.MakeSigner(string(d.Algorithm), d.PrivateKey)
		if err != nil {
			return r, fmt.Errorf("failed initialize signer: %w", err)
		}

		b, err := signer.Sign([]byte(r.SignedData))
		if err != nil {
			return r, fmt.Errorf("failed to sign transaction: %w", err)
		}

		r.Signature = base64.StdEncoding.EncodeToString(b)

		err = s.dvp.CompareAndSwapSignature(d.ID, d.LastSignature, persistence.Signature(r.Signature))

		if err == persistence.ErrRace {
			//retry
			continue
		}
		if err != nil {
			return r, fmt.Errorf("failed to persist signature update: %w", err)
		}

		return r, nil
	}
	return nil, fmt.Errorf("persistence retry limit reached")
}

func (s *SignatureDomainImpl) GetSignatureDeviceDetails(id ID) (Device, error) {
	device, err := s.dvp.GetDevice(persistence.ID(id))
	if err != nil {
		return Device{}, fmt.Errorf("failed to get signature device details: %w", err)
	}

	return Device{
		ID:        ID(device.ID),
		Algorithm: string(device.Algorithm),
		Label:     *device.Label,
	}, nil
}
