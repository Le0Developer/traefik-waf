package internal

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	RuleSetEnabled      bool
	RuleSetPath         string
	RequireJS           bool
	AssetPath           string
	RefHeader           string
	ChallengeCookie     string
	ChallengePassage    time.Duration
	ChallengeDifficulty int
	Verbosity           int
	XffCount            int
	FooterName          string
	FooterUrl           string
	WafName             string
}

func NewConfigFromEnv() *Config {
	return &Config{
		RuleSetEnabled:      getEnv("WAF_RULESET_ENABLED", "true") == "true",
		RuleSetPath:         getEnv("WAF_RULESET_PATH", "/rules/*.conf,/rules/*/*.conf,/rules/*/*/*.conf,/rules/*/*/*/*.conf"),
		RequireJS:           getEnv("WAF_REQUIRE_JS", "false") == "true",
		AssetPath:           getEnv("WAF_ASSET_PATH", "/.waf"),
		RefHeader:           getEnv("WAF_REF_HEADER", ""),
		ChallengeCookie:     getEnv("WAF_CHALLENGE_COOKIE", "_wafchlp"),
		ChallengePassage:    getEnvDuration("WAF_CHALLENGE_PASSAGE", "60m"),
		ChallengeDifficulty: getEnvInt("WAF_CHALLENGE_DIFFICULTY", "18"),
		Verbosity:           getEnvInt("WAF_VERBOSITY", "1"),
		XffCount:            getEnvInt("WAF_XFF_COUNT", "-1"),
		FooterName:          getEnv("WAF_FOOTER_NAME", "Web Application Firewall"),
		FooterUrl:           getEnv("WAF_FOOTER_URL", "https://github.com/le0developer/traefik-waf"),
		WafName:             getEnv("WAF_NAME", "Web Application Firewall"),
	}
}

func getEnv(key string, defaults ...string) string {
	value := os.Getenv(key)
	if value == "" {
		if len(defaults) > 0 {
			return defaults[0]
		}
		panic("Environment variable " + key + " is required but not set")
	}

	if strings.HasPrefix(value, "/") {
		if data, err := os.ReadFile(value); err == nil {
			return strings.TrimSpace(string(data))
		}
	}

	return value
}

func getEnvInt(key string, defaults ...string) int {
	value := getEnv(key, defaults...)
	intVal, err := strconv.Atoi(value)
	if err != nil {
		panic("Environment variable " + key + " is not a valid integer: " + err.Error())
	}
	return intVal
}

func getEnvDuration(key string, defaults ...string) time.Duration {
	value := getEnv(key, defaults...)
	dur, err := time.ParseDuration(value)
	if err != nil {
		panic("Environment variable " + key + " is not a valid duration: " + err.Error())
	}
	return dur
}
