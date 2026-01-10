package internal

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func (i *Instance) Mux() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/healthz":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
			return
		case "/assets/w.wasm":
			w.Header().Set("Content-Type", "application/wasm")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(wasmData)
			return
		}

		if r.Header.Get("X-Forwarded-Host") != "" {
			host := r.Header.Get("X-Forwarded-Host")
			r.Host = host
			r.URL.Host = host
		} else if r.Host != "" {
			r.URL.Host = r.Host
		}
		if r.Header.Get("X-Forwarded-Proto") != "" {
			r.URL.Scheme = r.Header.Get("X-Forwarded-Proto")
		} else if r.TLS != nil {
			r.URL.Scheme = "https"
		} else {
			r.URL.Scheme = "http"
		}
		if r.Header.Get("X-Forwarded-Uri") != "" {
			uri := r.Header.Get("X-Forwarded-Uri")
			parsed, err := url.Parse(uri)
			if err != nil {
				fmt.Printf("failed to parse X-Forwarded-Uri: %v\n", err)
			} else {
				r.URL.Path = parsed.Path
				r.URL.RawQuery = parsed.RawQuery
			}
		}

		ref := i.reference(r)
		if ref == "" {
			ref = "WAF-" + strconv.FormatInt(time.Now().UnixNano(), 36)
		}

		if interruption := i.evaluateRules(r); interruption != nil {
			if i.cfg.Verbosity >= 1 {
				fmt.Printf("request blocked: ref=%s rule=%d action=%s\n", ref, interruption.RuleID, interruption.Action)
			}
			i.replyWithBlocked(w, r, ref)
			return
		}

		if err := i.evaluateJS(r); err != nil {
			if i.cfg.Verbosity >= 2 {
				fmt.Printf("JS challenge required: ref=%s reason=%v\n", ref, err)
			}

			challenge, err := i.generateChallenge(r)
			if err != nil {
				fmt.Printf("failed to generate challenge: %v\n", err)
				addHeaders(w, ref, "challenge")
				i.replyWithPlainBlocked(w, ref, "chl-error")
				return
			}
			i.replyWithChallenge(w, r, i.buildChallengeScript(challenge), ref)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
