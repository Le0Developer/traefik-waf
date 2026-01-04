package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//go:embed challenge.js
var challengeScriptTemplate string

func (i *Instance) evaluateJS(r *http.Request) bool {
	if !i.requiresJS(r) {
		return true
	}

	if err := i.hasValidChallenge(r); err != nil {
		fmt.Println("did not have previous valid challenge:", err)
		return false
	}

	return true
}

func (i *Instance) requiresJS(r *http.Request) bool {
	if i.cfg.RequireJS {
		return true
	}

	return r.Header.Get("X-WAF-Require-JS") != ""
}

func (i *Instance) hasValidChallenge(r *http.Request) error {
	cookie, err := r.Cookie(i.cfg.ChallengeCookie)
	if err != nil {
		return fmt.Errorf("no challenge cookie: %w", err)
	}

	parts := strings.SplitN(cookie.Value, ".", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid challenge cookie format")
	}

	issuance, err := strconv.ParseInt(parts[0], 16, 64)
	if err != nil {
		return fmt.Errorf("invalid challenge issuance: %w", err)
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return fmt.Errorf("invalid challenge signature: %w", err)
	}

	if time.Since(time.Unix(issuance, 0)) > i.cfg.ChallengePassage {
		return fmt.Errorf("challenge expired")
	}

	h := hmac.New(sha256.New, i.challengeSubkeyFor(r))
	h.Write([]byte(parts[0]))
	expectedSignature := h.Sum(nil)
	if !hmac.Equal(signature, expectedSignature) {
		return fmt.Errorf("invalid challenge signature")
	}

	return nil
}

func (i *Instance) generateChallenge(r *http.Request) (string, error) {
	issuance := strconv.FormatInt(time.Now().Unix(), 16)
	h := hmac.New(sha256.New, i.challengeSubkeyFor(r))
	h.Write([]byte(issuance))
	signature := h.Sum(nil)

	signedPassage := issuance + "." + base64.RawURLEncoding.EncodeToString(signature)

	key := make([]byte, 32)
	nonce := make([]byte, 12)
	_, _ = rand.Read(key)
	_, _ = rand.Read(nonce)

	cipher2, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(cipher2)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	encryptedPassage := gcm.Seal(nil, nonce, []byte(signedPassage), nil)

	// zero the first N bits of the key according to difficulty
	for i := range i.cfg.ChallengeDifficulty {
		key[i>>3] &^= 1 << (i & 7)
	}

	challenge := base64.RawStdEncoding.EncodeToString(
		append(
			[]byte{byte(i.cfg.ChallengeDifficulty)},
			append(
				key,
				append(
					nonce,
					encryptedPassage...,
				)...,
			)...,
		),
	)

	return challenge, nil
}

func (i *Instance) buildChallengeScript(challenge string) string {
	script := strings.ReplaceAll(challengeScriptTemplate, "{{CHALLENGE}}", challenge)
	script = strings.ReplaceAll(script, "{{COOKIE_NAME}}", i.cfg.ChallengeCookie)
	script = strings.ReplaceAll(script, "{{PASSAGE_DURATION}}", fmt.Sprintf("%d", int(i.cfg.ChallengePassage.Seconds())))
	return "<script>" + script + "</script>"
}

func (i *Instance) challengeSubkeyFor(r *http.Request) []byte {
	h := hmac.New(sha256.New, i.secret)
	h.Write([]byte(r.URL.Host))
	h.Write([]byte(r.UserAgent()))

	// record static header values
	for _, name := range []string{"sec-ch-ua", "sec-ch-ua-mobile", "sec-ch-ua-platform", "accept-language"} {
		h.Write([]byte(r.Header.Get(name)))
	}

	// record if headers are present
	for _, name := range []string{"sec-fetch-site", "sec-fetch-mode", "sec-fetch-dest", "dnt", "sec-ch-ua", "sec-ch-ua-mobile", "sec-ch-ua-platform", "accept", "accept-encoding", "accept-language", "sec-gpc"} {
		if r.Header.Get(name) != "" {
			h.Write([]byte{1})
		} else {
			h.Write([]byte{0})
		}
	}

	return h.Sum(nil)
}
