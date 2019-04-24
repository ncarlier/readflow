package listener

import (
	"fmt"

	"github.com/ncarlier/readflow/pkg/event"
	eventbroker "github.com/ncarlier/readflow/pkg/event-broker"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/rs/zerolog/log"
)

const errorMsg = "Unable to send new user external event broker"

func init() {
	event.Subscribe(event.CreateUser, func(payload ...interface{}) {
		if user, ok := payload[0].(model.User); ok {
			broker := eventbroker.Lookup()
			if broker == nil {
				log.Debug().Err(fmt.Errorf("event broker not configured")).Uint("uid", *user.ID).Msg(errorMsg)
				return
			}
			evt := eventbroker.NewUserEvent(user)
			if err := broker.Send(evt.Buffer()); err != nil {
				log.Error().Err(err).Uint("uid", *user.ID).Msg(errorMsg)
			}
		}
	})
}
