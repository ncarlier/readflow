package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/event"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/service"

	"github.com/rs/zerolog/log"
)

const (
	errorMessage                        = "unable to send notification"
	maxUserInactivityBeforeNotification = 6 * time.Hour
)

var status string = "inbox"

func init() {
	event.Subscribe(event.CreateArticle, func(payload ...interface{}) {
		if article, ok := payload[0].(model.Article); ok {
			uid := article.UserID
			ctx := context.WithValue(context.TODO(), constant.ContextUserID, uid)
			req := model.ArticlesPageRequest{Status: &status}

			user, err := service.Lookup().GetCurrentUser(ctx)
			if err != nil {
				log.Info().Err(err).Uint("id", uid).Msg(errorMessage)
				return
			}

			nb, err := service.Lookup().CountCurrentUserDevices(ctx)
			if err != nil {
				log.Info().Err(err).Uint("id", uid).Msg(errorMessage)
				return
			}
			if nb == 0 {
				// no devices, therefore no need to send a notification
				return
			}

			globalStrategy := true
			text := "You have a new article to read."
			if article.CategoryID != nil {
				category, err := service.Lookup().GetCategory(ctx, *article.CategoryID)
				if err != nil {
					log.Info().Err(err).Uint("id", uid).Msg(errorMessage)
					return
				}
				if category.NotificationStrategy == "none" {
					return
				}
				globalStrategy = category.NotificationStrategy == "global"
				text = fmt.Sprintf("You have a new article to read in %s", category.Title)
			}

			if globalStrategy {
				// Send notification only if user is inactive for a while
				lastLoginDelay := time.Now().Add(-maxUserInactivityBeforeNotification)
				if user.Enabled && user.LastLoginAt != nil && user.LastLoginAt.After(lastLoginDelay) {
					return
				}
				// Retrieve number of articles
				nb, err := service.Lookup().CountCurrentUserArticles(ctx, req)
				if err != nil {
					log.Info().Err(err).Uint("id", uid).Msg(errorMessage)
					return
				}
				// Send notification only every 10 articles
				if !(nb > 0 && math.Mod(float64(nb), 10) == 0) {
					return
				}
				// Format text message
				text = fmt.Sprintf("You have %d articles to read.", nb)
			}

			// Build notification
			notif := &model.DeviceNotification{
				Title: "New articles to read",
				Body:  text,
			}
			b, err := json.Marshal(notif)
			if err == nil {
				text = string(b)
			}

			// Notify all user devices
			service.Lookup().NotifyDevices(ctx, text)
		}
	})
}
