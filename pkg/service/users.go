package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/event"
	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/model"
	userplan "github.com/ncarlier/readflow/pkg/user-plan"
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
			err = errors.New("user disabled")
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
		Plan:     "default",
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
	// TODO retrieve current user object from context
	uid := getCurrentUserFromContext(ctx)
	return reg.GetUserByID(ctx, uid)
}

// GetCurrentUserPlan get current user plan
func (reg *Registry) GetCurrentUserPlan(ctx context.Context) (*userplan.UserPlan, error) {
	user, err := reg.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	return reg.UserPlans.GetPlan(user.Plan), nil
}

// DeleteAccount delete current user account
func (reg *Registry) DeleteAccount(ctx context.Context) (bool, error) {
	user, err := reg.GetCurrentUser(ctx)
	if err != nil {
		return false, err
	}
	if err = reg.db.DeleteUser(*user); err != nil {
		return false, err
	}
	event.Emit(event.DeleteUser, *user)
	return true, nil
}

// GetUserByID get user by id
func (reg *Registry) GetUserByID(ctx context.Context, uid uint) (*model.User, error) {
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
	user.Hash = helper.Hash(strings.ToLower(user.Username))

	return user, nil
}

// GetUserByUsername get user by username
func (reg *Registry) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	reg.logger.Debug().Str(
		"username", username,
	).Msg("getting user...")

	// Try to fetch existing user...
	user, err := reg.db.GetUserByUsername(username)
	if err != nil || user == nil {
		if user == nil {
			err = errors.New("user not found")
		}
		reg.logger.Info().Err(err).Str(
			"username", username,
		).Msg("unable to find user")
		return nil, err
	}
	// Compute user hash
	user.Hash = helper.Hash(strings.ToLower(user.Username))

	return user, nil
}

// UpdateUser update user account
func (reg *Registry) UpdateUser(ctx context.Context, form model.UserForm) (*model.User, error) {
	uid := getCurrentUserFromContext(ctx)
	user, err := reg.db.GetUserByID(form.ID)
	if err != nil || user == nil {
		if user == nil {
			err = errors.New("user not found")
		}
		reg.logger.Info().Err(err).Uint(
			"uid", form.ID,
		).Msg("unable to update user")
		return nil, err
	}

	if form.Enabled != nil {
		user.Enabled = *form.Enabled
	}
	if form.Plan != nil {
		user.Plan = *form.Plan
	}
	if form.CustomerID != nil {
		user.CustomerID = *form.CustomerID
	}

	// Self protection
	if !user.Enabled && *user.ID == uid {
		err = errors.New("disabling himself is forbidden")
		reg.logger.Info().Err(err).Uint(
			"uid", form.ID,
		).Msg("unable to update user")
		return nil, err
	}

	event.Emit(event.UpdateUser, *user)

	return reg.db.CreateOrUpdateUser(*user)
}
