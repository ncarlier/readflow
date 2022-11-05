package scripting

import (
	"context"
	"fmt"

	"github.com/golang/groupcache/lru"
	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ScriptEngine is the scripting engine
type ScriptEngine struct {
	cache  *lru.Cache
	logger zerolog.Logger
}

// NewScriptEngine create new script engine
func NewScriptEngine(cacheSize int) *ScriptEngine {
	logger := log.With().Str("component", "scripting-engine").Logger()
	return &ScriptEngine{
		cache:  lru.New(cacheSize),
		logger: logger,
	}
}

// Exec a script by the script engine
func (s *ScriptEngine) Exec(ctx context.Context, script string, input ScriptInput) (OperationStack, error) {
	key := helper.Hash(script)
	logger := s.logger.With().Str("key", key).Str("title", input.Title).Logger()
	cacheItem, ok := s.cache.Get(key)
	var interpreter *Interpreter
	if ok {
		var valid bool
		if interpreter, valid = cacheItem.(*Interpreter); !valid {
			return nil, fmt.Errorf("script engine cache is poisoned")
		}
		logger.Debug().Msg("interpreter loaded from cache")
	} else {
		var err error
		logger.Debug().Msg("creating interpreter...")
		if interpreter, err = NewInterpreter(script, s.logger); err != nil {
			return nil, err
		}
		s.cache.Add(key, interpreter)
	}
	logger.Debug().Msg("executing...")
	return interpreter.Exec(ctx, input)
}
