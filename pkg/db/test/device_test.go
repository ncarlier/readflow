package dbtest

import (
	"encoding/json"
	"testing"

	webpush "github.com/SherClockHolmes/webpush-go"

	"github.com/ncarlier/readflow/pkg/assert"
	"github.com/ncarlier/readflow/pkg/model"
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

func assertDeviceExists(t *testing.T, device model.Device) *model.Device {
	result, err := testDB.GetDeviceByUserIDAndKey(*device.UserID, device.Key)
	assert.Nil(t, err, "error on getting device by user and key should be nil")
	if result != nil {
		return result
	}

	result, err = testDB.CreateDevice(device)
	assert.Nil(t, err, "error on create device should be nil")
	assert.NotNil(t, result, "device shouldn't be nil")
	assert.NotNil(t, result.ID, "device ID shouldn't be nil")
	assert.Equal(t, *device.UserID, *result.UserID, "")
	assert.NotEqual(t, "", result.Key, "")
	assert.Equal(t, device.Key, result.Key, "")
	return result
}

func TestCreateDevice(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	sub := newTestSubscription()
	b, err := json.Marshal(sub)
	assert.Nil(t, err, "error on marshaling test subscription to JSON")

	builder := model.NewDeviceBuilder()
	device := builder.UserID(*testUser.ID).Subscription(string(b)).Build()

	// Create device
	newDevice := assertDeviceExists(t, *device)

	devices, err := testDB.GetDevicesByUserID(*testUser.ID)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, len(devices) > 0, "devices should not be empty")

	// Cleanup
	err = testDB.DeleteDevice(*newDevice)
	assert.Nil(t, err, "error on cleanup should be nil")

	device, err = testDB.GetDeviceByUserIDAndKey(*testUser.ID, newDevice.Key)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, device == nil, "device should be nil")
}
