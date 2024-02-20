package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog/log"
)

// Catches any unintended panics in the server and recovers, logging the result for future reference
//
// This middleware should be placed first (or at least early) such that all errors are recovered from.
func RecoverWithInternalServerError(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Defer so that any panics are caught when logger is resolved
		defer func() {
			if rec := recover(); rec != nil {
				log.Error().
					Interface("RecoverInformation", rec).
					Bytes("DebugStack", debug.Stack()).
					Msg("ErrorRecoveryDuringZerologMiddleware")

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
