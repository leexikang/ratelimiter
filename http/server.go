package http

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"ratelimiter"
	"ratelimiter/memory"

	"gopkg.in/yaml.v3"
)

type Rule struct {
	Path           string `yaml:"path"`
	Limit          int    `ymal:"limit"`
	WindowInSecond int    `ymal:"windowinsecond"`
}

type Rules map[string]Rule

type Config struct {
	Port   string `yaml:"port"`
	Remote string `yaml:"remote"`
	Rules  []Rule `yaml:"rules"`
}

func Run() {
	config := readConfig()
	remote, err := url.Parse(config.Remote)
	rateLimiter := memory.New()
	registerRules(config.Rules, rateLimiter)

	if err != nil {
		panic(err)
	}

	handler := func(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Print(r.RequestURI)
			proxy.ServeHTTP(w, r)
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxyHandler := http.HandlerFunc(handler(proxy))
	http.Handle("/", throttleMiddleware(proxyHandler, rateLimiter))
	log.Printf("Listen on server port %s", config.Port)
	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		panic(err)
	}
}

func registerRules(rules []Rule, rateLimiter ratelimiter.RateLimiter) {
	for _, rule := range rules {
		log.Print(rule.Limit)
		log.Print(rule.WindowInSecond)
		rateLimiter.Create(rule.Path, rule.Limit, rule.WindowInSecond)
	}
}

func readConfig() Config {
	config := Config{}
	buf, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		panic(err)
	}

	yaml.Unmarshal(buf, &config)
	return config
}

func throttleMiddleware(next http.Handler, rateLimiter ratelimiter.RateLimiter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := rateLimiter.IsAllowed(r.URL.Path)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
