package persistence

import (
	"slices"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

type signatureDeviceRepository struct {
	devices []domain.SignatureDevice
	m       sync.Mutex
}

func NewInMemoryDeviceRepository() domain.SignatureDeviceRepository {
	return &signatureDeviceRepository{}
}

func (repo *signatureDeviceRepository) StoreSignatureDevice(device *domain.SignatureDevice) error {
	// TODO: deep clone, checkout "github.com/barkimedes/go-deepcopy"
	repo.m.Lock()
	defer repo.m.Unlock()

	// Update if exists
	i := slices.IndexFunc(repo.devices, func(d domain.SignatureDevice) bool {
		return d.Id == device.Id
	})

	if i == -1 {
		repo.devices = append(repo.devices, *device)
	} else {
		repo.devices[i] = *device
	}

	return nil
}

func (repo *signatureDeviceRepository) ListSignatureDevices() ([]domain.SignatureDevice, error) {
	repo.m.Lock()
	defer repo.m.Unlock()

	// TODO: deep clone
	return repo.devices, nil
}

func (repo *signatureDeviceRepository) GetSignatureDeviceById(id string) *domain.SignatureDevice {
	repo.m.Lock()
	defer repo.m.Unlock()

	i := slices.IndexFunc(repo.devices, func(d domain.SignatureDevice) bool {
		return d.Id == id
	})

	if i == -1 {
		return nil
	}

	// TODO: deep clone
	device := repo.devices[i]
	return &device
}
