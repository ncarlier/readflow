package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	webpush "github.com/SherClockHolmes/webpush-go"

	"github.com/ncarlier/readflow/pkg/model"
)

const errNotification = "unable to notify user's devices"

// GetDevices get devices from current user
func (reg *Registry) GetDevices(ctx context.Context) (*[]model.Device, error) {
	uid := getCurrentUserIDFromContext(ctx)

	devices, err := reg.db.GetDevicesByUser(uid)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to get devices")
		return nil, err
	}

	return &devices, err
}

// CountCurrentUserDevices get total categories of current user
func (reg *Registry) CountCurrentUserDevices(ctx context.Context) (uint, error) {
	uid := getCurrentUserIDFromContext(ctx)
	return reg.db.CountDevicesByUser(uid)
}

// GetDevice get a device of the current user
func (reg *Registry) GetDevice(ctx context.Context, id uint) (*model.Device, error) {
	uid := getCurrentUserIDFromContext(ctx)

	device, err := reg.db.GetDeviceByID(id)
	if err != nil || device == nil || *device.UserID != uid {
		if err == nil {
			err = ErrDeviceNotFound
		}
		return nil, err
	}
	return device, nil
}

// CreateDevice create or update a device for current user
func (reg *Registry) CreateDevice(ctx context.Context, sub string) (*model.Device, error) {
	uid := getCurrentUserIDFromContext(ctx)

	builder := model.NewDeviceBuilder()
	device := builder.UserID(uid).Subscription(sub).Build()

	if device.Subscription == nil {
		err := errors.New("invalid subscription")
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to configure device")
		return nil, err
	}

	result, err := reg.db.CreateDevice(*device)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("key", device.Key).Msg("unable to create device")
		return nil, err
	}
	return result, err
}

// DeleteDevice delete a device of the current user
func (reg *Registry) DeleteDevice(ctx context.Context, id uint) (*model.Device, error) {
	uid := getCurrentUserIDFromContext(ctx)

	device, err := reg.GetDevice(ctx, id)
	if err != nil {
		return nil, err
	}

	err = reg.db.DeleteDevice(id)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Uint("id", id).Msg("unable to delete device")
		return nil, err
	}
	return device, nil
}

// DeleteDevices delete devices of the current user
func (reg *Registry) DeleteDevices(ctx context.Context, ids []uint) (int64, error) {
	uid := getCurrentUserIDFromContext(ctx)
	idsStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")

	nb, err := reg.db.DeleteDevicesByUser(uid, ids)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("ids", idsStr).Msg(errNotification)
		return 0, err
	}
	reg.logger.Debug().Err(err).Uint(
		"uid", uid,
	).Str("ids", idsStr).Int64("nb", nb).Msg("devices deleted")
	return nb, nil
}

// NotifyDevices send a notification to all user devices
func (reg *Registry) NotifyDevices(ctx context.Context, msg string) (int, error) {
	user, err := reg.GetCurrentUser(ctx)
	if err != nil {
		reg.logger.Info().Err(err).Msg(errNotification)
		return 0, err
	}
	uid := *user.ID

	devices, err := reg.GetDevices(ctx)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg(errNotification)
		return 0, err
	}
	counter := 0
	for _, device := range *devices {
		// Rate limiting
		if _, _, _, ok, err := reg.notificationRateLimiter.Take(ctx, user.Username); err != nil || !ok {
			if !ok {
				err = errors.New("rate limiting activated")
			}
			reg.logger.Info().Err(err).Uint(
				"uid", uid,
			).Uint("device", *device.ID).Msg(errNotification)
			continue
		}
		// Send notification
		res, err := webpush.SendNotification([]byte(msg), device.Subscription, &webpush.Options{
			Subscriber:      user.Username,
			VAPIDPublicKey:  reg.properties.VAPIDPublicKey,
			VAPIDPrivateKey: reg.properties.VAPIDPrivateKey,
			TTL:             30,
		})
		if err != nil {
			reg.logger.Info().Err(err).Uint(
				"uid", uid,
			).Uint("device", *device.ID).Msg(errNotification)
			continue
		}
		if res.StatusCode == 410 {
			// Registration is gone... we should remove the device
			err = reg.db.DeleteDevice(*device.ID)
			reg.logger.Info().Err(err).Uint(
				"uid", uid,
			).Uint("device", *device.ID).Msg("registration gone: device deleted")
			continue
		}
		if res.StatusCode >= 400 {
			reg.logger.Info().Err(errors.New(res.Status)).Uint(
				"uid", uid,
			).Uint(
				"device", *device.ID,
			).Int("status", res.StatusCode).Msg(errNotification)
			continue
		}
		counter++
		reg.logger.Info().Uint(
			"uid", uid,
		).Uint(
			"device", *device.ID,
		).Int("status", res.StatusCode).Msg("notification sent to user device")
	}
	return counter, nil
}
