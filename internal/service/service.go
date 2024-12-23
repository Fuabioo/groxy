package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type (
	Service struct {
		httputil.ReverseProxy

		url       url.URL
		endpoints map[string]*EndpointConfiguration
	}
)

func New() (*Service, error) {

	domain := viper.GetString("domain")

	log.Debug("🔗 parsing domain url", "domain", domain)

	if !strings.HasPrefix(domain, "https://") && !strings.HasPrefix(domain, "http://") {
		domain = fmt.Sprintf("https://%s", domain)
	}

	url, err := url.Parse(domain)
	if err != nil {
		return nil, fmt.Errorf("could not parse domain: %w", err)
	}

	domainsConfig := viper.GetStringMap("endpoints")
	endpoints := make(map[string]*EndpointConfiguration, len(domainsConfig))
	for endpoint, stringValue := range domainsConfig {
		log.Debug("📖 parsing endpoint configuration",
			"endpoint", endpoint,
		)

		raw, err := yaml.Marshal(stringValue)
		if err != nil {
			return nil, fmt.Errorf("could not marshal endpoint config: %w", err)
		}

		log.Debug("📦 marshaled endpoint",
			"endpoint", endpoint,
			"config", string(raw),
		)

		var config EndpointConfiguration
		if err := yaml.Unmarshal(raw, &config); err != nil {
			return nil, fmt.Errorf("could not unmarshal endpoint config: %w", err)
		}

		endpoints[endpoint] = &config
		log.Debug("🧐 parsed endpoint",
			"endpoint", endpoint,
			"config", endpoint,
		)

	}

	log.Debug("🦖 initializing handler",
		"scheme", url.Scheme,
		"host", url.Host,
		"endpoints", len(domainsConfig),
	)

	reverseProxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = url.Scheme
			req.URL.Host = url.Host
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: viper.GetBool("insecure"),
			},
		},
	}

	return &Service{
		ReverseProxy: reverseProxy,
		url:          *url,
		endpoints:    endpoints,
	}, nil
}

type responseInterceptor struct {
	http.ResponseWriter
	data       []byte
	statusCode int
}

func NewResponseInterceptor(w http.ResponseWriter) *responseInterceptor {
	return &responseInterceptor{
		ResponseWriter: w,
	}
}

func (r *responseInterceptor) Write(data []byte) (int, error) {
	r.data = append(r.data, data...)
	return r.ResponseWriter.Write(data)
}

func (r *responseInterceptor) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (h *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h == nil {
		http.Error(w, "service not initialized", http.StatusInternalServerError)
		return
	}

	response := NewResponseInterceptor(w)
	w = response

	r.Host = h.url.Host

	defer func() {
		data, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Error("could not dump request", "err", err)
		} else {
			var buffer bytes.Buffer
			buffer.Write([]byte(fmt.Sprintln("HTTP Request:")))
			buffer.Write(bytes.Repeat([]byte{'='}, 80))
			buffer.Write([]byte(fmt.Sprintln("📥 incoming request")))
			buffer.Write(data)
			buffer.Write(bytes.Repeat([]byte{'-'}, 80))
			buffer.Write([]byte(fmt.Sprintln()))
			buffer.Write([]byte(fmt.Sprintln("📤 outgoing response")))
			buffer.Write([]byte(fmt.Sprintf("Status: %d", response.statusCode)))
			buffer.Write([]byte(fmt.Sprintln()))
			for k, v := range w.Header() {
				buffer.Write([]byte(fmt.Sprintf("%s: %s", k, strings.Join(v, ", "))))
				buffer.Write([]byte(fmt.Sprintln()))
			}
			prettyPrintedData := response.data
			if response.Header().Get("Content-Type") == "application/json" {
				var res map[string]interface{}
				if err := json.Unmarshal(response.data, &res); err == nil {
					data, err := json.MarshalIndent(res, "", "  ")
					if err == nil {
						prettyPrintedData = data
					}
				}
			}
			buffer.Write(prettyPrintedData)
			buffer.Write([]byte(fmt.Sprintln()))
			buffer.Write(bytes.Repeat([]byte{'='}, 80))
			log.Debug(buffer.String())
		}
	}()

	// if the request route matches an endpoint configuration, apply it
	if endpoint, ok := h.endpoints[r.URL.Path]; ok {
		endpoint.Lock()

		log.Debug("✍️ applying endpoint configuration",
			"endpoint", r.URL.Path,
			"config", endpoint,
		)

		endpoint.Delay.Take()

		var response *Response
		if err := endpoint.Error.Chance.Take(); err != nil {
			response = endpoint.Error.Response
		} else if err := endpoint.Error.Every.Take(); err != nil {
			response = endpoint.Error.Response
			if response.StatusCode == 0 {
				response.StatusCode = http.StatusBadRequest
			}
		} else if endpoint.MockResponse != nil {
			response = endpoint.MockResponse
		}

		endpoint.Unlock()

		if response != nil {
			log.Debug("🤡 responding with mock response",
				"endpoint", r.URL.Path,
			)
			contentType := response.ContentType
			if contentType == "" {
				contentType = "application/json"
			}
			w.Header().Set("Content-Type", contentType)

			statusCode := http.StatusOK
			if response.StatusCode != 0 {
				statusCode = response.StatusCode
			}
			w.WriteHeader(statusCode)

			if response.Body != "" {
				w.Write([]byte(response.Body))
			}

			// instead of calling the reverse proxy, return here
			return
		}
	}

	h.ReverseProxy.ServeHTTP(w, r)
}
