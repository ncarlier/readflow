package service

import (
	"context"
	"errors"
	"time"

	"github.com/ncarlier/reader/pkg/model"
)

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
	reg.logger.Debug().Uint32(
		"uid", *user.ID,
	).Str("username", username).Msg("user registered")
	return user, nil
}

// GetUserByID get an user by its ID
func (reg *Registry) GetUserByID(ctx context.Context, id uint32) (*model.User, error) {
	reg.logger.Debug().Uint32(
		"id", id,
	).Msg("getting user...")

	// Try to fetch existing user...
	user, err := reg.db.GetUserByID(id)
	if err != nil {
		reg.logger.Info().Err(err).Uint32(
			"id", id,
		).Msg("unable to find user")
		return nil, err
	}
	return user, nil
}
