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

func NewApiRouter() *ApiHandler {
	apiRouterData := &ApiHandler{
		apiRouter: chi.NewRouter(),
	}

	apiRouterData.apiRouter.Get("/heartbeat", heartbeat)
	apiRouterData.apiRouter.Get("/register", apiRouterData.register)

	return apiRouterData
}
