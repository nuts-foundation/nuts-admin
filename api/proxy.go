package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type proxyRoute struct {
	method       string
	path         string
	compiledPath *regexp.Regexp
}

const proxyPath = "/api/proxy/"

func init() {
	for i := range allowedProxyRoutes {
		allowedProxyRoutes[i].compiledPath = regexp.MustCompile(allowedProxyRoutes[i].path)
	}
}

var allowedProxyRoutes = []proxyRoute{
	// List Discovery Services
	{
		method: http.MethodGet,
		path:   "/internal/discovery/v1",
	},
	// Search VPs on a Discovery Service
	{
		method: http.MethodGet,
		path:   "/internal/discovery/v1/([a-z-A-Z0-9_\\-\\:\\.%]+)",
	},
	// Activate Discovery Services for a DID
	{
		method: http.MethodPost,
		path:   "/internal/discovery/v1/([a-z-A-Z0-9_\\-\\:\\.%]+)/([a-z-A-Z0-0_\\-\\:\\.%]+)",
	},
	// Deactivate Discovery Services for a DID
	{
		method: http.MethodDelete,
		path:   "/internal/discovery/v1/([a-z-A-Z0-9_\\-\\:\\.%]+)/([a-z-A-Z0-0_\\-\\:\\.%]+)",
	},
	// Issue Verifiable Credentials
	{
		method: http.MethodPost,
		path:   "/internal/vcr/v2/issuer/vc",
	},
	// Search for issued Verifiable Credentials
	{
		method: http.MethodGet,
		path:   "/internal/vcr/v2/issuer/vc/search",
	},
	// Load Verifiable Credential into wallet
	{
		method: http.MethodPost,
		path:   "/internal/vcr/v2/holder/([a-z-A-Z0-9_\\-\\:\\.%]+)/vc",
	},
}

// ConfigureProxy configures the proxy middleware for the given Nuts node address.
// It allows the web application to call a curated list of endpoints on the Nuts node.
func ConfigureProxy(e *echo.Echo, nodeAddress *url.URL) {
	e.Use(middleware.ProxyWithConfig(middleware.ProxyConfig{
		Balancer: middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
			{
				URL: nodeAddress,
			},
		}),
		Rewrite: map[string]string{
			"^" + proxyPath + "*": "/$1",
		},
		Skipper: func(c echo.Context) bool {
			proxyURL := c.Request().URL.Path
			if !strings.HasPrefix(proxyURL, proxyPath) {
				// Not a proxy request
				return true
			}
			proxyURL = strings.TrimPrefix(proxyURL, proxyPath)
			proxyURL = "/" + strings.TrimLeft(proxyURL, "/")
			for _, route := range allowedProxyRoutes {
				if c.Request().Method == route.method && route.compiledPath.MatchString(proxyURL) {
					log.Printf("proxying %s %s", c.Request().Method, proxyURL)
					return false
				}
			}
			_ = c.String(http.StatusForbidden, "Proxy route not allowed")
			return true
		},
	}))
}
