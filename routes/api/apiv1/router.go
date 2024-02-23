package apiv1

import (
	"github.com/go-chi/chi/v5"
)

type ApiHandler struct {
	apiRouter *chi.Mux
	// databaseConn
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
