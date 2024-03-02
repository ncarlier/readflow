package service

import (
	"context"
	"errors"
	"time"

	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/internal/global"
	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/pkg/event"
)

func getCurrentUserIDFromContext(ctx context.Context) uint {
	return ctx.Value(global.ContextUserID).(uint)
}

func isAdmin(ctx context.Context) bool {
	if value := ctx.Value(global.ContextIsAdmin); value != nil {
		return value.(bool)
	}
	return false
}

// GetCurrentUser get current user
func (reg *Registry) GetCurrentUser(ctx context.Context) (*model.User, error) {
	if value := ctx.Value(global.ContextUser); value != nil {
		user := value.(model.User)
		return &user, nil
	}
	uid := getCurrentUserIDFromContext(ctx)
	return reg.GetUserByID(ctx, uid)
}

// GetOrRegisterUser get an existing user or creates new one
func (reg *Registry) GetOrRegisterUser(ctx context.Context, username string) (*model.User, error) {
	logger := reg.logger.With().Str("username", username).Logger()

	logger.Debug().Msg("user login...")
	// Try to fetch existing user...
	user, err := reg.db.GetUserByUsername(username)
	if err != nil {
		logger.Info().Err(err).Msg("unable to login")
		return nil, err
	}
	// If user already exists...
	if user != nil {
		// Checks that the user is not disabled
		if !user.Enabled {
			err = errors.New("user disabled")
			logger.Info().Err(err).Msg("unable to login")
			return nil, err
		}
		// Update user login date
		lastLoginDate := time.Now()
		user.LastLoginAt = &lastLoginDate
		user, err = reg.db.CreateOrUpdateUser(*user)
		if err != nil {
			logger.Info().Err(err).Msg("unable to login")
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
		logger.Info().Err(err).Msg("unable to register user")
		return nil, err
	}
	logger.Info().Uint("uid", *user.ID).Msg("user registered")
	reg.events.Publish(event.NewEvent(EventCreateUser, *user))
	return user, nil
}

// GetCurrentUserPlan get current user plan
func (reg *Registry) GetCurrentUserPlan(ctx context.Context) (*config.UserPlan, error) {
	user, err := reg.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	return reg.conf.GetUserPlan(user.Plan), nil
}

// DeleteAccount delete current user account
func (reg *Registry) DeleteAccount(ctx context.Context) (bool, error) {
	user, err := reg.GetCurrentUser(ctx)
	if err != nil {
		return false, err
	}
	if err := reg.db.DeleteUser(*user); err != nil {
		return false, err
	}
	reg.events.Publish(event.NewEvent(EventDeleteUser, *user))
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
			"uid", uid,
		).Msg("unable to find user")
		return nil, err
	}

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

	return user, nil
}

// GetUserByHashID get user by hash ID
func (reg *Registry) GetUserByHashID(ctx context.Context, hashid string) (*model.User, error) {
	ids, err := reg.hashid.Decode(hashid)
	if err != nil {
		return nil, err
	}
	uid := uint(ids[0])
	return reg.GetUserByID(ctx, uid)
}

// UpdateUser update user account (required admin access)
func (reg *Registry) UpdateUser(ctx context.Context, form model.UserForm) (*model.User, error) {
	start := time.Now()
	uid := getCurrentUserIDFromContext(ctx)
	if !isAdmin(ctx) {
		err := errors.New("forbidden")
		reg.logger.Info().Err(err).Uint(
			"uid", form.ID,
		).Msg("unable to update user")
		return nil, err
	}
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

	user, err = reg.db.CreateOrUpdateUser(*user)
	if err != nil {
		reg.logger.Info().Err(err).Str(
			"username", user.Username,
		).Msg("unable to update user")
		return nil, err
	}
	reg.logger.Info().Uint(
		"uid", *user.ID,
	).Str(
		"username", user.Username,
	).Dur("took", time.Since(start)).Msg("user updated")

	reg.events.Publish(event.NewEvent(EventUpdateUser, *user))

	return user, nil
}

// GetUserHashID returns user hashid
func (reg *Registry) GetUserHashID(uid uint) string {
	return reg.hashid.Encode([]int{int(uid)})
}

// GetUserPlans returns user plans
func (reg *Registry) GetUserPlans() []config.UserPlan {
	return reg.conf.UserPlans
}
