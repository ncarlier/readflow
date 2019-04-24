package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/ncarlier/readflow/pkg/constant"

	"github.com/ncarlier/readflow/pkg/event"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/service"
)

func init() {
	status := "unread"
	event.Subscribe(event.CreateArticle, func(payload ...interface{}) {
		if article, ok := payload[0].(model.Article); ok {
			uid := article.UserID
			ctx := context.WithValue(context.TODO(), constant.UserID, uid)
			req := model.ArticlesPageRequest{Status: &status}

			nb, err := service.Lookup().CountArticles(ctx, req)
			if err != nil {
				return
			}
			if !(nb > 0 && math.Mod(float64(nb), 10) == 0) {
				// No need to notify the user
				return
			}

			text := fmt.Sprintf("You have %d articles to read.", nb)
			notif := &model.DeviceNotification{
				Title: "New articles to read",
				Body:  text,
			}
			b, err := json.Marshal(notif)
			if err == nil {
				text = string(b)
			}

			// User notification
			service.Lookup().NotifyDevices(ctx, text)
		}
	})
}
