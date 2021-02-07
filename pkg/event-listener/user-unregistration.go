package listener

import (
	"fmt"

	"github.com/ncarlier/readflow/pkg/event"
	eventbroker "github.com/ncarlier/readflow/pkg/event-broker"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/rs/zerolog/log"
)

func init() {
	errorMsg := "unable to send user unregistration"
	event.Subscribe(event.DeleteUser, func(payload ...interface{}) {
		if user, ok := payload[0].(model.User); ok {
			broker := eventbroker.Lookup()
			if broker == nil {
				log.Debug().Err(fmt.Errorf("event broker not configured")).Uint("uid", *user.ID).Msg(errorMsg)
				return
			}
			evt := eventbroker.NewUserEvent(user, event.DeleteUser)
			if err := broker.Send(evt.Buffer()); err != nil {
				log.Error().Err(err).Uint("uid", *user.ID).Msg(errorMsg)
			}
		}
	})
}
