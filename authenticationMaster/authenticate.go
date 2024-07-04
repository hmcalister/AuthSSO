package authenticationmaster

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/rs/zerolog/log"
)

type authorizedUserData struct {
	UserID   string
	Username string
}

