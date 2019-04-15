package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ncarlier/readflow/pkg/config"
	eventbroker "github.com/ncarlier/readflow/pkg/event-broker"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/rs/zerolog/log"
)

const errorMsg = "Unable to send new user external event broker"

func init() {
	bus.Subscribe(CreateUser, func(payload ...interface{}) {
		if user, ok := payload[0].(model.User); ok {
			broker := eventbroker.Lookup()
			if broker == nil {
				log.Debug().Err(fmt.Errorf("event broker not configured")).Uint("uid", *user.ID).Msg(errorMsg)
				return
			}
			event := UserEvent{
				Payload: user,
			}
			event.Action = CreateUser
			event.Issue = Issue{
				Date: time.Now(),
				URL:  config.Get().PublicURL,
			}

			b := new(bytes.Buffer)
			json.NewEncoder(b).Encode(event)
			if err := broker.Send(b); err != nil {
				log.Error().Err(err).Uint("uid", *user.ID).Msg(errorMsg)
			}
		}
	})
}
