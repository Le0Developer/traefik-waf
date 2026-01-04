package internal

import (
	"crypto/rand"
	"strings"

	"github.com/corazawaf/coraza/v3"
)

type Instance struct {
	cfg      *Config
	engine   coraza.WAF
	response *responseTemplater
	secret   []byte
}

func New(cfg *Config) (*Instance, error) {
	var engine coraza.WAF

	if cfg.RuleSetEnabled {
		ccfg := coraza.NewWAFConfig()
		for path := range strings.SplitSeq(cfg.RuleSetPath, ",") {
			ccfg = ccfg.WithDirectivesFromFile(strings.TrimSpace(path))
		}
		var err error
		engine, err = coraza.NewWAF(ccfg)
		if err != nil {
			return nil, err
		}
	}

	secret := make([]byte, 32)
	_, _ = rand.Read(secret)

	response := newResponseTemplater()

	return &Instance{cfg: cfg, engine: engine, response: response, secret: secret}, nil
}
