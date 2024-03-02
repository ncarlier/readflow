package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	webpush "github.com/SherClockHolmes/webpush-go"

	"github.com/ncarlier/readflow/internal/model"
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

	logger := reg.logger.With().Uint("uid", uid).Logger()

	builder := model.NewDeviceBuilder()
	device := builder.UserID(uid).Subscription(sub).Build()

	if device.Subscription == nil {
		err := errors.New("invalid subscription")
		logger.Info().Err(err).Msg("unable to configure device")
		return nil, err
	}
	logger = logger.With().Str("key", device.Key).Logger()

	logger.Debug().Msg("creating device...")
	result, err := reg.db.CreateDevice(*device)
	if err != nil {
		logger.Info().Err(err).Msg("unable to create device")
		return nil, err
	}
	logger.Info().Uint("id", *result.ID).Msg("device created")
	return result, err
}

// DeleteDevice delete a device of the current user
func (reg *Registry) DeleteDevice(ctx context.Context, id uint) (*model.Device, error) {
	uid := getCurrentUserIDFromContext(ctx)

	logger := reg.logger.With().Uint("uid", uid).Uint("id", id).Logger()

	device, err := reg.GetDevice(ctx, id)
	if err != nil {
		return nil, err
	}

	logger.Debug().Msg("deleting device...")
	err = reg.db.DeleteDevice(id)
	if err != nil {
		logger.Info().Err(err).Msg("unable to delete device")
		return nil, err
	}
	logger.Info().Msg("device deleted")
	return device, nil
}

// DeleteDevices delete devices of the current user
func (reg *Registry) DeleteDevices(ctx context.Context, ids []uint) (int64, error) {
	uid := getCurrentUserIDFromContext(ctx)
	idsStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")

	logger := reg.logger.With().Uint("uid", uid).Str("ids", idsStr).Logger()

	logger.Debug().Msg("deleting devices...")
	nb, err := reg.db.DeleteDevicesByUser(uid, ids)
	if err != nil {
		logger.Info().Err(err).Msg("unable to delete devices")
		return 0, err
	}
	logger.Info().Int64("nb", nb).Msg("devices deleted")
	return nb, nil
}

// NotifyDevices send a notification to all user devices
func (reg *Registry) NotifyDevices(ctx context.Context, payload *model.DeviceNotification) (int, error) {
	user, err := reg.GetCurrentUser(ctx)
	if err != nil {
		reg.logger.Info().Err(err).Msg(errNotification)
		return 0, err
	}
	uid := *user.ID

	logger := reg.logger.With().Uint("uid", uid).Logger()

	devices, err := reg.GetDevices(ctx)
	if err != nil {
		logger.Info().Err(err).Msg(errNotification)
		return 0, err
	}
	msg, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}
	counter := 0
	for _, device := range *devices {
		start := time.Now()
		logger = reg.logger.With().Uint("uid", uid).Uint("device", *device.ID).Logger()
		// Rate limiting
		if _, _, _, ok, err := reg.notificationRateLimiter.Take(ctx, user.Username); err != nil || !ok {
			if !ok {
				err = errors.New("rate limiting activated")
			}
			logger.Info().Err(err).Msg(errNotification)
			continue
		}
		// Send notification
		res, err := webpush.SendNotification(msg, device.Subscription, &webpush.Options{
			Subscriber:      user.Username,
			VAPIDPublicKey:  reg.properties.VAPIDPublicKey,
			VAPIDPrivateKey: reg.properties.VAPIDPrivateKey,
			TTL:             30,
		})
		if err != nil {
			logger.Info().Err(err).Msg(errNotification)
			continue
		}
		if res.StatusCode == 410 {
			// Registration is gone... we should remove the device
			err = reg.db.DeleteDevice(*device.ID)
			logger.Info().Err(err).Msg("registration gone: device deleted")
			continue
		}
		if res.StatusCode >= 400 {
			logger.Info().Err(errors.New(res.Status)).Int("status", res.StatusCode).Dur("took", time.Since(start)).Msg(errNotification)
			continue
		}
		counter++
		logger.Info().Str("title", payload.Title).Dur("took", time.Since(start)).Msg("notification sent to user device")
	}
	return counter, nil
}
