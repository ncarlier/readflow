package dbtest

import (
	"encoding/json"
	"testing"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/stretchr/testify/assert"

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
	return result
}

func TestCreateDevice(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	uid := *testUser.ID
	sub := newTestSubscription()
	b, err := json.Marshal(sub)
	assert.Nil(t, err)

	builder := model.NewDeviceBuilder()
	device := builder.UserID(uid).Subscription(string(b)).Build()

	// Create device
	newDevice := assertDeviceExists(t, *device)

	devices, err := testDB.GetDevicesByUser(uid)
	assert.Nil(t, err)
	assert.True(t, len(devices) > 0, "devices should not be empty")

	// Cleanup
	err = testDB.DeleteDevice(*newDevice.ID)
	assert.Nil(t, err)

	device, err = testDB.GetDeviceByUserAndKey(uid, newDevice.Key)
	assert.Nil(t, err)
	assert.Nil(t, device)
}
