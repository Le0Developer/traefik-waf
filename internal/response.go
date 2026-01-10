package internal

import (
	_ "embed"
	"html"
	"net/http"
	"os"
	"strings"
)

//go:embed blocked.html
var blockedTemplate string

//go:embed challenge.html
var challengeTemplate string

//go:embed jspowobfdata.js
var challengePowLib string

type responseTemplater struct {
	blocked   string
	challenge string
}

func (rt *responseTemplater) getBlockedResponse(reference string) string {
	return strings.Replace(rt.blocked, "<!--REF-->", html.EscapeString(reference), 1)
}

func (rt *responseTemplater) getChallengeResponse(challenge string, reference string) string {
	return strings.Replace(strings.Replace(rt.challenge, "<!--CHALLENGE-->", challenge, 1), "<!--REF-->", html.EscapeString(reference), 1)
}

func newResponseTemplater(cfg *Config) *responseTemplater {
	blocked := blockedTemplate
	challenge := challengeTemplate

	if customBlocked, _ := os.ReadFile("/custom/blocked.html"); customBlocked != nil {
		blocked = string(customBlocked)
	}
	blocked = strings.ReplaceAll(blocked, "{{WAF_NAME}}", cfg.WafName)
	blocked = strings.ReplaceAll(blocked, "{{FOOTER_NAME}}", cfg.FooterName)
	blocked = strings.ReplaceAll(blocked, "{{FOOTER_URL}}", cfg.FooterUrl)

	if customChallenge, _ := os.ReadFile("/custom/challenge.html"); customChallenge != nil {
		challenge = string(customChallenge)
	}
	challenge = strings.ReplaceAll(challenge, "{{WAF_NAME}}", cfg.WafName)
	challenge = strings.ReplaceAll(challenge, "{{FOOTER_NAME}}", cfg.FooterName)
	challenge = strings.ReplaceAll(challenge, "{{FOOTER_URL}}", cfg.FooterUrl)

	challenge = strings.ReplaceAll(challenge, "<!--CHALLENGE-->", challengePowLib+"<!--CHALLENGE-->")

	if customHead, _ := os.ReadFile("/custom/head.html"); customHead != nil {
		head := string(customHead)

		blocked = strings.Replace(blocked, "<!--HEAD-->", head, 1)
		challenge = strings.Replace(challenge, "<!--HEAD-->", head, 1)
	} else {
		blocked = strings.Replace(blocked, "<!--HEAD-->", "", 1)
		challenge = strings.Replace(challenge, "<!--HEAD-->", "", 1)
	}

	return &responseTemplater{
		blocked:   blocked,
		challenge: challenge,
	}
}

func acceptsHTML(r *http.Request) bool {
	acceptHeader := r.Header.Get("Accept")
	for v := range strings.SplitSeq(acceptHeader, ",") {
		trimV := strings.SplitN(strings.TrimSpace(v), ";", 2)[0]
		switch trimV {
		case "text/html", "text/*", "*/*":
			return true
		}
	}
	return false
}

func addHeaders(w http.ResponseWriter, reference string, action string) {
	w.Header().Set("Cache-Control", "no-store")
	// w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Waf-Ref", reference)
	if action != "" {
		w.Header().Set("Waf-Action", action)
	}
}

func (i *Instance) replyWithPlainBlocked(w http.ResponseWriter, reference string, action string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte("403 ref=" + reference + " action=" + action + "\n"))
}

func (i *Instance) replyWithBlocked(w http.ResponseWriter, r *http.Request, reference string) {
	addHeaders(w, reference, "block")
	if !acceptsHTML(r) {
		i.replyWithPlainBlocked(w, reference, "block")
		return
	}

	response := i.response.getBlockedResponse(reference)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte(response))
}

func (i *Instance) replyWithChallenge(w http.ResponseWriter, r *http.Request, challenge string, reference string) {
	addHeaders(w, reference, "challenge")
	if !acceptsHTML(r) {
		i.replyWithPlainBlocked(w, reference, "challenge")
		return
	}

	response := i.response.getChallengeResponse(challenge, reference)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte(response))
}
