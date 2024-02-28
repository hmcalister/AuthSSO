package apiv1

import (
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/rs/zerolog/log"
)

func private(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse JWT!")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occurred during authentication. Please try again."))
		return
	}

	w.Write([]byte(fmt.Sprintf("Private Route, Welcome User %v\n", claims["Subject"])))
}
