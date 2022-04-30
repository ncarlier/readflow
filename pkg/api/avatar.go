package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ncarlier/readflow/pkg/avatar"
	"github.com/ncarlier/readflow/pkg/config"
	"github.com/rs/zerolog/log"
)

const (
	MaxAge = 864000
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
	generator, err := newAvatarGenerator(conf.Integration.AvatarProvider)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create avatar generator")
	}
	log.Info().Str("component", "api").Str("provider", conf.Integration.AvatarProvider).Msg("using Avatar provider")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seed := strings.TrimPrefix(r.URL.Path, "/avatar/")
		if seed == "" {
			http.Error(w, "URL param 'seed' is missing", http.StatusBadRequest)
			return
		}
		if generator == nil {
			redirect := strings.ReplaceAll(conf.Integration.AvatarProvider, "{seed}", seed)
			http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
			return
		}
		set := r.URL.Query().Get("set")
		log.Debug().Str("seed", seed).Str("set", set).Msg("generating avatar image")
		img, err := generator.Generate(seed, set)
		if err != nil {
			log.Error().Err(err).Str("seed", seed).Msg("unable to generate avatar image")
			http.Error(w, "unable to generate avatar image", http.StatusInternalServerError)
			return
		}
		expires := time.Now().Add(MaxAge * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Pragma", "public")
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", MaxAge))
		w.Header().Set("Expires", expires.Local().String())
		img.WriteTo(w)
	})
}
