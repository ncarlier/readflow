package event

import (
	"bytes"
	"encoding/json"
	"fmt"

	eventbroker "github.com/ncarlier/reader/pkg/event-broker"
	"github.com/ncarlier/reader/pkg/model"
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
			b := new(bytes.Buffer)
			json.NewEncoder(b).Encode(user)
			if err := broker.Send(b); err != nil {
				log.Error().Err(err).Uint("uid", *user.ID).Msg(errorMsg)
			}
		}
	})
}
