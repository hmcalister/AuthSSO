package apiv1

import "github.com/go-chi/chi/v5"

func ApiV1Router() *chi.Mux {
	apiRouter := chi.NewRouter()

	apiRouter.Get("/heartbeat", heartbeat)
	apiRouter.Get("/register", register)

	return apiRouter
}
