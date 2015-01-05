package cors

import (
	"fmt"
	"net/http"
	"strings"
)

// Config for CORS middleware.
type Config struct {
	// Allowed methods
	Methods []string
	// Allowed origins
	Origins []string
	// Allowed headers
	Headers []string
	// Exposed methods
	ExposedHeaders []string
	// Allow credentials
	AllowCredentials bool
	// Max age
	MaxAge float64
}

// IsOriginAllowed determines whether given orgin is allowed
func (c Config) IsOriginAllowed(origin string) bool {
	for i := 0; i < len(c.Origins); i++ {
		if "*" == c.Origins[i] {
			return true
		} else if origin == c.Origins[i] {
			return true
		}
	}
	return false
}

// Config to allow all requests.
var AllowAll = Config{
	Methods:          []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH"},
	Origins:          []string{"*"},
	Headers:          []string{"Content-Type"},
	ExposedHeaders:   []string{"Content-Type"},
	AllowCredentials: true,
	MaxAge:           0,
}

// New returns CORS middleware
func New(config Config) func(http.Handler) http.Handler {
	// prepare header strings
	var methods = strings.Join(config.Methods, ", ")
	var headers = strings.Join(config.Headers, ", ")
	var exposed = strings.Join(config.ExposedHeaders, ", ")
	var allowCredentials = "false"
	if config.AllowCredentials {
		allowCredentials = "true"
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var origin = r.Header.Get("Origin")
			if len(origin) > 0 && config.IsOriginAllowed(origin) {
				w.Header().Add("Access-Control-Allow-Origin", origin)
				w.Header().Add("Access-Control-Allow-Methods", methods)
				w.Header().Add("Access-Control-Allow-Headers", headers)
				w.Header().Add("Access-Control-Expose-Headers", exposed)
				w.Header().Add("Access-Control-Allow-Credentials", allowCredentials)
				if config.MaxAge > 0 {
					w.Header().Add("Access-Control-Max-Age", fmt.Sprintf("%9.f", config.MaxAge))
				}
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
