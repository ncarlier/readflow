package service

import (
	"bytes"
	"context"
	"errors"

	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/pkg/event"
	"github.com/ncarlier/readflow/pkg/logger"
	"github.com/ncarlier/readflow/pkg/thumbhash"
	"github.com/ncarlier/readflow/pkg/utils"
)

const thumbhashErrorMessage = "unable to create thumbhash"

func newThumbhashEventHandler(srv *Registry) event.EventHandler {
	return func(evt event.Event) {
		article, ok := evt.Payload.(model.Article)
		if !ok || article.Status == "read" || utils.IsNilOrEmpty(article.Image) || !utils.IsNilOrEmpty(article.ThumbHash) {
			// Ignore if not a article event
			// OR if the article is marked as read
			// OR if the article have no image
			// OR if the article have already a thumbhash
			return
		}
		logger := logger.With().Uint("id", article.ID).Logger()

		// download article image
		src := *article.Image
		if srv.imageProxy.URL() != "" {
			src = srv.imageProxy.Encode(src, "")
		}
		logger = logger.With().Str("src", src).Logger()
		asset, res, err := srv.dl.Get(context.Background(), src, nil)
		if err != nil {
			logger.Info().Err(err).Msg(thumbhashErrorMessage)
			return
		}

		if res != nil && res.StatusCode != 200 {
			err := errors.New("bad status code")
			logger.Info().Err(err).Int("status", res.StatusCode).Msg(thumbhashErrorMessage)
			return
		}

		if asset == nil {
			return
		}

		// generate thumbhash
		r := bytes.NewReader(asset.Data)
		hash, err := thumbhash.GetThumbhash(r)
		if err != nil {
			logger.Info().Err(err).Msg(thumbhashErrorMessage)
			return
		}

		// save article thumbhash
		if _, err := srv.db.SetArticleThumbHash(article.ID, hash); err != nil {
			logger.Info().Err(err).Msg(thumbhashErrorMessage)
			return
		}

		logger.Debug().Uint("id", article.ID).Str("hash", hash).Msg("acrticle thumbash created")
	}
}
