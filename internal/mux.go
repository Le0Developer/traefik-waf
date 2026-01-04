package internal

import (
	"fmt"
	"net/http"
)

func (i *Instance) Mux() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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
			r.URL.Path = r.Header.Get("X-Forwarded-Uri")
		}

		if interruption := i.evaluateRules(*r); interruption != nil {
			fmt.Printf("request blocked: %v\n", interruption)
			i.replyWithBlocked(w, r)
			return
		}

		if !i.evaluateJS(r) {
			challenge, err := i.generateChallenge(r)
			if err != nil {
				fmt.Printf("failed to generate challenge: %v\n", err)
				i.replyWithPlainBlocked(w, r, '5')
				return
			}
			i.replyWithChallenge(w, r, i.buildChallengeScript(challenge))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
