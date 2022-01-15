package db

import "github.com/ncarlier/readflow/pkg/model"

// DeviceRepository is the repository interface to manage Devices
type DeviceRepository interface {
	GetDeviceByID(id uint) (*model.Device, error)
	GetDeviceByUserAndKey(uid uint, key string) (*model.Device, error)
	GetDevicesByUser(uid uint) ([]model.Device, error)
	CountDevicesByUser(uid uint) (uint, error)
	CreateDevice(device model.Device) (*model.Device, error)
	DeleteDevice(id uint) error
	DeleteDevicesByUser(uid uint, ids []uint) (int64, error)
}
