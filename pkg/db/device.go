package db

import "github.com/ncarlier/readflow/pkg/model"

// DeviceRepository is the repository interface to manage Devices
type DeviceRepository interface {
	GetDeviceByID(id uint) (*model.Device, error)
	GetDeviceByUserIDAndKey(uid uint, key string) (*model.Device, error)
	GetDevicesByUserID(uid uint) ([]model.Device, error)
	CreateDevice(device model.Device) (*model.Device, error)
	DeleteDevice(device model.Device) error
	DeleteDevices(uid uint, ids []uint) (int64, error)
}
