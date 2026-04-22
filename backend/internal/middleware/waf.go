package middleware

import (
	"log"
	"net/http"

	coreruleset "github.com/corazawaf/coraza-coreruleset"
	"github.com/corazawaf/coraza/v3"
	txhttp "github.com/corazawaf/coraza/v3/http"
)

// InitializeWAF initializes and returns a Coraza WAF instance with OWASP CRS
func InitializeWAF() (coraza.WAF, error) {
	waf, err := coraza.NewWAF(
		coraza.NewWAFConfig().
			WithRootFS(coreruleset.FS).
			WithDirectivesFromFile("@coraza.conf-recommended").
			WithDirectivesFromFile("@crs-setup.conf.example").
			WithDirectivesFromFile("@owasp_crs/*.conf"),
	)
	if err != nil {
		return nil, err
	}

	log.Println("✅ Coraza WAF initialized with OWASP Core Rule Set")
	return waf, nil
}

// WrapHandlerWithWAF wraps an HTTP handler with Coraza WAF protection
func WrapHandlerWithWAF(waf coraza.WAF, handler http.Handler) http.Handler {
	return txhttp.WrapHandler(waf, handler)
}
