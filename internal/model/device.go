package model

import (
	"encoding/json"
	"errors"
	"time"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/ncarlier/readflow/pkg/utils"
)

// DeviceNotification structure definition
type DeviceNotification struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
	Href  string `json:"href,omitempty"`
}

// Device structure definition
// Device key is a hash of the subscription payload and is used to prevent subscription duplication
type Device struct {
	ID           *uint                 `json:"id,omitempty"`
	UserID       *uint                 `json:"user_id,omitempty"`
	Key          string                `json:"key,omitempty"`
	Subscription *webpush.Subscription `json:"_"`
	LastSeenAt   *time.Time            `json:"last_seen_at,omitempty"`
	CreatedAt    *time.Time            `json:"created_at,omitempty"`
}

// GetSubscription get JSON string of the subscription
func (d *Device) GetSubscription() (string, error) {
	if d.Subscription != nil {
		b, err := json.Marshal(d.Subscription)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	return "", errors.New("subscription undefined")
}

// SetSubscription set subscription from JSON string
func (d *Device) SetSubscription(sub string) error {
	s := &webpush.Subscription{}
	if err := json.Unmarshal([]byte(sub), s); err != nil {
		return err
	}
	d.Subscription = s
	d.Key = utils.Hash(sub)
	return nil
}

// DeviceBuilder is a builder to create an Device
type DeviceBuilder struct {
	device *Device
}

// NewDeviceBuilder creates new Device builder instance
func NewDeviceBuilder() DeviceBuilder {
	device := &Device{}
	return DeviceBuilder{device}
}

// Build creates the device
func (ab *DeviceBuilder) Build() *Device {
	return ab.device
}

// UserID set device user ID
func (ab *DeviceBuilder) UserID(userID uint) *DeviceBuilder {
	ab.device.UserID = &userID
	return ab
}

// Subscription set device subscription
func (ab *DeviceBuilder) Subscription(sub string) *DeviceBuilder {
	ab.device.SetSubscription(sub)
	return ab
}
