package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/pkg/avatar"
	"github.com/ncarlier/readflow/pkg/defaults"
	"github.com/ncarlier/readflow/pkg/logger"
)

func newAvatarGenerator(provider string) (*avatar.Generator, error) {
	u, err := url.Parse(provider)
	if err != nil {
		return nil, err
	}
	defaultSet := u.Query().Get("default")
	switch u.Scheme {
	case "file":
		return avatar.NewGenerator(u.Host+u.Path, defaultSet)
	case "https":
		return nil, nil
	}
	return nil, fmt.Errorf("invalid avatar provider: %s", provider)
}

func avatarHandler(conf *config.Config) http.Handler {
	generator, err := newAvatarGenerator(conf.Avatar.ServiceProvider)
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to create avatar generator")
	}
	logger.Info().Str("component", "api").Str("provider", conf.Avatar.ServiceProvider).Msg("using Avatar provider")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seed := strings.TrimPrefix(r.URL.Path, "/avatar/")
		if seed == "" {
			http.Error(w, "URL param 'seed' is missing", http.StatusBadRequest)
			return
		}
		if generator == nil {
			redirect := strings.ReplaceAll(conf.Avatar.ServiceProvider, "{seed}", seed)
			http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
			return
		}
		set := r.URL.Query().Get("set")
		logger.Debug().Str("seed", seed).Str("set", set).Msg("generating avatar image")
		img, err := generator.Generate(seed, set)
		if err != nil {
			logger.Error().Err(err).Str("seed", seed).Msg("unable to generate avatar image")
			http.Error(w, "unable to generate avatar image", http.StatusInternalServerError)
			return
		}
		expires := time.Now().Add(defaults.CacheMaxAge * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Pragma", "public")
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", defaults.CacheMaxAge))
		w.Header().Set("Expires", expires.Local().String())
		img.WriteTo(w)
	})
}
