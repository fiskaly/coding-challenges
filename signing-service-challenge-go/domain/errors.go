package domain

import "fmt"

func ErrorDeviceNotFound(deviceId string) error {
	return fmt.Errorf("device with id %s does not exist", deviceId)
}
