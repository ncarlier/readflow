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

func init() {
	status := "unread"
	event.Subscribe(event.CreateArticle, func(payload ...interface{}) {
		if article, ok := payload[0].(model.Article); ok {
			uid := article.UserID
			ctx := context.WithValue(context.TODO(), constant.UserID, uid)
			req := model.ArticlesPageRequest{Status: &status}

			user, err := service.Lookup().GetCurrentUser(ctx)
			if err != nil {
				log.Info().Err(err).Uint("id", uid).Msg("unable to send notification")
				return
			}

			// Send notification only if user logged in more than 5 minutes ago
			lastLoginDelay := time.Now().Add(-5 * time.Minute)
			if user.Enabled && user.LastLoginAt != nil && user.LastLoginAt.After(lastLoginDelay) {
				return
			}

			nb, err := service.Lookup().CountArticles(ctx, req)
			if err != nil {
				log.Info().Err(err).Uint("id", uid).Msg("unable to send notification")
				return
			}

			// Send notification only every 10 articles
			if !(nb > 0 && math.Mod(float64(nb), 10) == 0) {
				return
			}

			// Build notification
			text := fmt.Sprintf("You have %d articles to read.", nb)
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
