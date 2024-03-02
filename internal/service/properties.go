package service

import (
	webpush "github.com/SherClockHolmes/webpush-go"

	"github.com/ncarlier/readflow/internal/model"
)

func (reg *Registry) initProperties() error {
	properties, err := reg.db.GetProperties()
	if err != nil {
		return err
	}
	if properties == nil {
		// Initialize properties
		privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
		if err != nil {
			return err
		}
		properties = &model.Properties{
			VAPIDPrivateKey: privateKey,
			VAPIDPublicKey:  publicKey,
		}
		properties, err = reg.db.CreateProperties(*properties)
		if err != nil {
			return err
		}
		reg.logger.Debug().Str(
			"pub", properties.VAPIDPublicKey,
		).Msg("new VAPID key created")
	}
	reg.logger.Info().Uint(
		"rev", *properties.Rev,
	).Msg("properties loaded")
	reg.properties = properties
	return nil
}

// GetProperties retrieve service properties
func (reg *Registry) GetProperties() model.Properties {
	return *reg.properties
}
