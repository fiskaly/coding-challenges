package mocks

import "signing-service-challenge/domain"

type MockDeviceRepository struct {
	UpdateDeviceCallsCount  int
	GetDeviceByIdCallsCount int
	ListDevicesCallsCount   int
	DeviceToUpdate          *domain.Device
	DeviceToReturn          *domain.Device
	DevicesListToReturn     []*domain.Device
	GetDeviceByIdArg        string
	GetDeviceByIdFound      bool
}

func NewMockDeviceRepository() *MockDeviceRepository {
	return &MockDeviceRepository{}
}

func (m *MockDeviceRepository) UpdateDevice(device *domain.Device) {
	m.UpdateDeviceCallsCount++
	m.DeviceToUpdate = device
}

func (m *MockDeviceRepository) GetDeviceById(deviceId string) (*domain.Device, bool) {
	m.GetDeviceByIdCallsCount++
	m.GetDeviceByIdArg = deviceId
	return m.DeviceToReturn, m.GetDeviceByIdFound
}

func (*MockDeviceRepository) ListDevices() ([]*domain.Device, error) {
	return nil, nil
}
