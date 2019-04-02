package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/ncarlier/reader/pkg/constant"
	"github.com/ncarlier/reader/pkg/event"
	"github.com/ncarlier/reader/pkg/model"
	"github.com/ncarlier/reader/pkg/tooling"
)

func getCurrentUserFromContext(ctx context.Context) uint {
	return ctx.Value(constant.UserID).(uint)
}

// GetOrRegisterUser get an existing user or creates new one
func (reg *Registry) GetOrRegisterUser(ctx context.Context, username string) (*model.User, error) {
	reg.logger.Debug().Str(
		"username", username,
	).Msg("user login...")

	// Try to fetch existing user...
	user, err := reg.db.GetUserByUsername(username)
	if err != nil {
		reg.logger.Info().Err(err).Str(
			"username", username,
		).Msg("unable to login")
		return nil, err
	}
	// If user already exists...
	if user != nil {
		// Checks that the user is not disabled
		if !user.Enabled {
			err = errors.New("user diabled")
			reg.logger.Info().Err(err).Str(
				"username", username,
			).Msg("unable to login")
			return nil, err
		}
		// Update user login date
		lastLoginDate := time.Now()
		user.LastLoginAt = &lastLoginDate
		user, err = reg.db.CreateOrUpdateUser(*user)
		if err != nil {
			reg.logger.Info().Err(err).Str(
				"username", username,
			).Msg("unable to login")
			return nil, err
		}
		// Returns existing user
		return user, nil
	}
	// ... else create a new user...
	user = &model.User{
		Username: username,
		Enabled:  true,
	}
	user, err = reg.db.CreateOrUpdateUser(*user)
	if err != nil {
		reg.logger.Info().Err(err).Str(
			"username", username,
		).Msg("unable to register user")
		return nil, err
	}
	reg.logger.Debug().Uint(
		"uid", *user.ID,
	).Str("username", username).Msg("user registered")
	event.Emit(event.CreateUser, *user)
	return user, nil
}

// GetCurrentUser get current user
func (reg *Registry) GetCurrentUser(ctx context.Context) (*model.User, error) {
	uid := getCurrentUserFromContext(ctx)
	reg.logger.Debug().Uint(
		"id", uid,
	).Msg("getting user...")

	// Try to fetch existing user...
	user, err := reg.db.GetUserByID(uid)
	if err != nil || user == nil {
		if user == nil {
			err = errors.New("user not found")
		}
		reg.logger.Info().Err(err).Uint(
			"id", uid,
		).Msg("unable to find user")
		return nil, err
	}
	// Compute user hash
	user.Hash = tooling.Hash(strings.ToLower(user.Username))

	return user, nil
}
