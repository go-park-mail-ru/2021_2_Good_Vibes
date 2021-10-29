package csrf

import (
	"crypto/subtle"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/random"
	"net/http"
	"strings"
	"time"
)

type (
	// CSRFConfig defines the config for CSRF middleware.
	CSRFConfig struct {
		Skipper        middleware.Skipper
		TokenLength    uint8         `yaml:"token_length"`
		TokenLookup    string        `yaml:"token_lookup"`
		ContextKey     string        `yaml:"context_key"`
		CookieName     string        `yaml:"cookie_name"`
		CookieDomain   string        `yaml:"cookie_domain"`
		CookiePath     string        `yaml:"cookie_path"`
		CookieMaxAge   int           `yaml:"cookie_max_age"`
		CookieHTTPOnly bool          `yaml:"cookie_http_only"`
		CookieSameSite http.SameSite `yaml:"cookie_same_site"`
	}
	csrfTokenExtractor func(echo.Context) (string, error)
)

var (
	// DefaultCSRFConfig is the default CSRF middleware config.
	DefaultCSRFConfig = CSRFConfig{
		Skipper:        middleware.DefaultSkipper,
		TokenLength:    32,
		TokenLookup:    "header:" + echo.HeaderXCSRFToken,
		ContextKey:     "csrf",
		CookieName:     "_csrf",
		CookieMaxAge:   86400,
		CookieSameSite: http.SameSiteDefaultMode,
	}
)

func CSRFWithConfig() echo.MiddlewareFunc {
	// Initialize
	parts := strings.Split(DefaultCSRFConfig.TokenLookup, ":")
	extractor := csrfTokenFromHeader(parts[1])

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if DefaultCSRFConfig.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			k, err := c.Cookie(DefaultCSRFConfig.CookieName)
			token := ""

			// Generate token
			if err != nil {
				token = random.String(DefaultCSRFConfig.TokenLength)
			} else {
				// Reuse token
				token = k.Value
			}

			switch req.Method {
			case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
			default:
				// Validate token only for requests which are not defined as 'safe' by RFC7231
				if req.RequestURI != "/login" {
					clientToken, err := extractor(c)
					if err != nil {
						return echo.NewHTTPError(http.StatusBadRequest, err.Error())
					}
					if !validateCSRFToken(token, clientToken) {
						return echo.NewHTTPError(http.StatusForbidden, "invalid csrf token")
					}
				}
			}

			cookie := new(http.Cookie)
			cookie.Name = DefaultCSRFConfig.CookieName
			cookie.Value = token
			cookie.Path = DefaultCSRFConfig.CookiePath
			cookie.Domain = DefaultCSRFConfig.CookieDomain
			cookie.SameSite = DefaultCSRFConfig.CookieSameSite

			cookie.Expires = time.Now().Add(time.Duration(DefaultCSRFConfig.CookieMaxAge) * time.Second)
			cookie.Secure = true
			cookie.HttpOnly = DefaultCSRFConfig.CookieHTTPOnly
			c.SetCookie(cookie)

			// Store token in the context
			c.Set(DefaultCSRFConfig.ContextKey, token)

			// Protect clients from caching the response
			c.Response().Header().Add(echo.HeaderVary, echo.HeaderCookie)

			return next(c)
		}
	}
}

// csrfTokenFromForm returns a `csrfTokenExtractor` that extracts token from the
// provided request header.
func csrfTokenFromHeader(header string) csrfTokenExtractor {
	return func(c echo.Context) (string, error) {
		return c.Request().Header.Get(header), nil
	}
}

func validateCSRFToken(token, clientToken string) bool {
	return subtle.ConstantTimeCompare([]byte(token), []byte(clientToken)) == 1
}
