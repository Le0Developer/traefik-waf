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

type responseTemplater struct {
	blocked   string
	challenge string
	head      string
}

func (rt *responseTemplater) getBlockedResponse(reference string) string {
	return strings.Replace(strings.Replace(rt.blocked, "<!--HEAD-->", rt.head, 1), "<!--REF-->", html.EscapeString(reference), 1)
}

func (rt *responseTemplater) getChallengeResponse(challenge string, reference string) string {
	return strings.Replace(strings.Replace(strings.Replace(rt.challenge, "<!--HEAD-->", rt.head, 1), "<!--CHALLENGE-->", challenge, 1), "<!--REF-->", html.EscapeString(reference), 1)
}

func newResponseTemplater() *responseTemplater {
	blocked := blockedTemplate
	challenge := challengeTemplate
	head := ""

	if customBlocked, _ := os.ReadFile("/custom/blocked.html"); customBlocked != nil {
		blocked = string(customBlocked)
	}
	if customChallenge, _ := os.ReadFile("/custom/challenge.html"); customChallenge != nil {
		challenge = string(customChallenge)
	}
	if customHead, _ := os.ReadFile("/custom/head.html"); customHead != nil {
		head = string(customHead)
	}

	return &responseTemplater{
		blocked:   blocked,
		challenge: challenge,
		head:      head,
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

func addHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
	// w.Header().Set("Pragma", "no-cache")
}

func (i *Instance) replyWithPlainBlocked(w http.ResponseWriter, r *http.Request, tag byte) {
	reference := i.reference(r)
	addHeaders(w)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte("403 ref=" + reference + " tag=" + string(tag) + "\n"))
}

func (i *Instance) replyWithBlocked(w http.ResponseWriter, r *http.Request) {
	if !acceptsHTML(r) {
		i.replyWithPlainBlocked(w, r, 'B')
		return
	}

	reference := i.reference(r)
	response := i.response.getBlockedResponse(reference)

	addHeaders(w)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte(response))
}

func (i *Instance) replyWithChallenge(w http.ResponseWriter, r *http.Request, challenge string) {
	if !acceptsHTML(r) {
		i.replyWithPlainBlocked(w, r, 'C')
		return
	}

	reference := i.reference(r)
	response := i.response.getChallengeResponse(challenge, reference)

	addHeaders(w)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte(response))
}
