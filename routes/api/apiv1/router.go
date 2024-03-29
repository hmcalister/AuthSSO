package apiv1

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/hmcalister/AuthSSO/database"
)

const (
	passwordMaxLen = 1024
)

type ApiHandler struct {
	apiRouter          *chi.Mux
	databaseConnection *database.DatabaseManager
	tokenSecretKey     []byte
	tokenAuth          *jwtauth.JWTAuth
}

func (api *ApiHandler) Router() *chi.Mux {
	return api.apiRouter
}

func NewApiRouter(db *database.DatabaseManager, tokenSecretKey []byte) *ApiHandler {
	tokenAuth := jwtauth.New(tokenSigningMethod.Alg(), tokenSecretKey, nil)

	apiRouterData := &ApiHandler{
		apiRouter:          chi.NewRouter(),
		databaseConnection: db,
		tokenSecretKey:     tokenSecretKey,
		tokenAuth:          tokenAuth,
	}

	apiRouterData.apiRouter.Get("/heartbeat", heartbeat)
	apiRouterData.apiRouter.Post("/register", apiRouterData.register)
	apiRouterData.apiRouter.Post("/login", apiRouterData.login)

	apiRouterData.apiRouter.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Get("/private", private)
	})

	return apiRouterData
}
