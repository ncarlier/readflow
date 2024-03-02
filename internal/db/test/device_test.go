package dbtest

import (
	"encoding/json"
	"testing"
	"time"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/internal/model"
)

func newTestSubscription() webpush.Subscription {
	return webpush.Subscription{
		Endpoint: "https://android.googleapis.com/gcm/send/a-subscription-id",
		Keys: webpush.Keys{
			Auth:   "AEl35...7fG",
			P256dh: "Fg5t8...2rC",
		},
	}
}

func newTestDevice(t *testing.T, uid uint) *model.Device {
	sub := newTestSubscription()
	b, err := json.Marshal(sub)
	assert.Nil(t, err)

	builder := model.NewDeviceBuilder()
	return builder.UserID(uid).Subscription(string(b)).Build()
}

func assertDeviceExists(t *testing.T, device model.Device) *model.Device {
	result, err := testDB.GetDeviceByUserAndKey(*device.UserID, device.Key)
	assert.Nil(t, err)
	if result != nil {
		return result
	}

	result, err = testDB.CreateDevice(device)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.ID)
	assert.Equal(t, *device.UserID, *result.UserID)
	assert.NotEqual(t, "", result.Key)
	assert.Equal(t, device.Key, result.Key)
	assert.True(t, time.Now().After(*result.LastSeenAt))
	return result
}

func TestCreateDevice(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	uid := *testUser.ID
	device := newTestDevice(t, uid)

	// Create device
	newDevice := assertDeviceExists(t, *device)

	// Get device
	newDevice, err := testDB.GetDeviceByID(*newDevice.ID)
	assert.Nil(t, err)
	assert.Equal(t, device.Key, newDevice.Key)

	// Delete the device
	err = testDB.DeleteDevice(*newDevice.ID)
	assert.Nil(t, err)

	// Try to get the device again
	device, err = testDB.GetDeviceByID(*newDevice.ID)
	assert.Nil(t, err)
	assert.Nil(t, device)
}

func TestListDevice(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	uid := *testUser.ID
	device := newTestDevice(t, uid)

	// Create device
	newDevice := assertDeviceExists(t, *device)

	// List devices
	devices, err := testDB.GetDevicesByUser(uid)
	assert.Nil(t, err)
	assert.Positive(t, len(devices), "devices should not be empty")

	// Delete the device
	deleted, err := testDB.DeleteDevicesByUser(uid, []uint{*newDevice.ID})
	assert.Nil(t, err)
	assert.Positive(t, deleted)
}

func TestDeviceCleanup(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	uid := *testUser.ID
	device := newTestDevice(t, uid)

	// Create device
	assertDeviceExists(t, *device)

	// Count devices
	nb, err := testDB.CountDevicesByUser(uid)
	assert.Nil(t, err)
	assert.Positive(t, nb)

	// Cleanup
	deleted, err := testDB.DeleteInactiveDevicesOlderThan(0)
	assert.Nil(t, err)
	assert.Positive(t, deleted)
}
