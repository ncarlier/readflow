package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/ncarlier/readflow/internal/global"
	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/pkg/event"
	"github.com/ncarlier/readflow/pkg/logger"
)

const (
	notificationErrorMessage            = "unable to send notification"
	maxUserInactivityBeforeNotification = 8 * time.Hour
)

func newNotificationEventHandler(srv *Registry) event.EventHandler {
	var status string = "inbox"
	return func(evt event.Event) {
		if article, ok := evt.Payload.(model.Article); ok {
			if evt.Option.Has(NoNotificationEventOption) {
				return
			}
			uid := article.UserID
			ctx := context.WithValue(context.TODO(), global.ContextUserID, uid)
			req := model.ArticlesPageRequest{Status: &status}

			logger := logger.With().Uint("uid", uid).Logger()

			user, err := srv.GetCurrentUser(ctx)
			if err != nil {
				logger.Info().Err(err).Msg(notificationErrorMessage)
				return
			}

			nb, err := srv.CountCurrentUserDevices(ctx)
			if err != nil {
				logger.Info().Err(err).Msg(notificationErrorMessage)
				return
			}
			if nb == 0 {
				// no devices, therefore no need to send a notification
				return
			}

			// Send notification only if user is inactive for a while
			lastLoginDelay := time.Now().Add(-maxUserInactivityBeforeNotification)
			if user.Enabled && user.LastLoginAt != nil && user.LastLoginAt.After(lastLoginDelay) {
				return
			}
			// Retrieve number of articles
			nb, err = srv.CountCurrentUserArticles(ctx, req)
			if err != nil {
				logger.Info().Err(err).Msg(notificationErrorMessage)
				return
			}
			// Send notification only every 10 articles
			if !(nb > 0 && math.Mod(float64(nb), 10) == 0) {
				return
			}

			text := fmt.Sprintf("You have %d articles to read.", nb)

			// Build notification
			notif := &model.DeviceNotification{
				Title: "New articles to read",
				Body:  text,
				Href:  "/",
			}
			// Notify all user devices
			srv.NotifyDevices(ctx, notif)
		}
	}
}
